import socket
import logger
import protocol
import domain
import threading
import signal
from lottery import Lottery

class Server:
    def __init__(self, server_host: str, server_port: int, lottery: Lottery, agency_quorum_min: int) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.lottery = lottery
        self.agency_quorum_min = agency_quorum_min
        self.shutting_down = False
        self.server_socket = None
        self.condition = threading.Condition()
        self.storage_lock = threading.Lock()
        self.finished_agencies = set()
        self.client_threads = []
        self.client_sockets = set()
        self.client_sockets_lock = threading.Lock()

    def _handle_client(self, client_socket):
        action = "handle-client"
        message_amount = 0
        agency_id = None

        with self.client_sockets_lock:
            self.client_sockets.add(client_socket)

        try:
            logger.info(action, logger.LogResult.in_progress)
            msg_type, payload = protocol.agency_id_announcement(client_socket)
            
            if msg_type != protocol.RECIEVE_AGENCY_ID:
                logger.error(action, logger.LogResult.fail, "agency id not received correctly", message_amount)
                return

            agency_id = int(payload)

            while True:
                msg_type, payload = protocol.recieve_bet_chunk(client_socket)

                if msg_type == protocol.RECIEVE_BET_CHUNK:
                    message_amount = self.receive_chunk(payload, message_amount, agency_id)
                    protocol.send_ack(client_socket, "OK")

                if msg_type == protocol.CLIENT_DONE:
                    self.handle_lottery(client_socket, agency_id)
                    logger.info(action, logger.LogResult.success, "messages-amount", message_amount)
                    return
                
                if msg_type is None:
                    logger.info(action, logger.LogResult.success, "messages-amount", message_amount)
                    return
                
        except Exception as e:
            logger.error(action, logger.LogResult.fail, "messages-amount", message_amount)
            protocol.send_ack(client_socket, "ERROR")
            raise e
        finally:
            with self.client_sockets_lock:
                self.client_sockets.remove(client_socket)
            client_socket.close()

    def receive_chunk(self,payload, message_amount, agency_id):
        bets = domain.strings_to_bets(payload, agency_id)

        with self.storage_lock:
            self.lottery.store_bets(bets)
        
        return message_amount + 1

    def handle_lottery(self, client_socket, agency_id):
        with self.condition:
            if agency_id is not None:
                self.finished_agencies.add(agency_id)
            while len(self.finished_agencies) < self.agency_quorum_min and not self.shutting_down: 
                self.condition.wait()
            self.condition.notify_all()

        if self.shutting_down:
            return
        
        with self.storage_lock:
            bets = list(self.lottery.load_bets())

        winner_strings = []
        for b in bets:
            if self.lottery.has_won(b) and int(b.agency_id) == agency_id:
                winner_strings.append(domain.bet_to_string(b))
        
        payload = "\n".join(winner_strings) 
        protocol.send_result_message(client_socket, payload)

    def _handle_sigterm(self, signum, frame):
        logger.info("sigterm", logger.LogResult.in_progress)
        self.shutting_down = True
        self.server_socket.close()
        with self.condition:
            self.condition.notify_all()
        with self.client_sockets_lock:
            for client_socket in self.client_sockets:
                client_socket.close()
        for t in self.client_threads:
            t.join()
    
    def run(self):
        action = "accept-connection"
        signal.signal(signal.SIGTERM, self._handle_sigterm)

        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()
            self.server_socket = server_socket
            while self.shutting_down is False:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except Exception as e:
                    if self.shutting_down:
                        break
                    logger.error(action, logger.LogResult.fail)
                    raise e
                logger.info(action, logger.LogResult.success)

                t = threading.Thread(target=self._handle_client, args=(client_socket,), daemon=True)
                t.start()
                self.client_threads.append(t)  