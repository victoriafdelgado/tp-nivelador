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
        self.storage_lock = threading.Lock()
        self.barrier = threading.Barrier(agency_quorum_min)
        self.shutting_down = False
        self.server_socket = None
        self.client_threads = []
        
    def _handle_client(self, client_socket):
        action = "handle-client"
        message_amount = 0
        agency_id = None
        try:
            logger.info(action, logger.LogResult.in_progress)
            while True:
                msg_type, payload = protocol.recieve_bet_chunk(client_socket)
         
                if msg_type == 2:
                    bets = domain.strings_to_bets(payload)
                    self.lottery.store_bets(bets)
                    if agency_id is None:
                        agency_id = bets[0].agency_id
                    message_amount +=1
                if msg_type == 1:
                    self.barrier.wait()
                    with self.storage_lock:
                        bets = self.lottery.load_bets()
                    if not bets:
                        protocol.send_ack(client_socket, "Error")
                    protocol.send_ack(client_socket, "OK")
                    winner_strings = []
                    for b in bets:
                        if self.lottery.has_won(b) and int(b.agency_id) == agency_id:
                            winner_strings.append(domain.bet_to_string(b))
                    
                    payload = "\n".join(winner_strings) 
                    protocol.send_result_message(client_socket, payload)
                    return

                if msg_type is None:
                    logger.info(action, logger.LogResult.success, "messages-amount", message_amount)
                    return
        except Exception as e:
            logger.error(action, logger.LogResult.fail, "messages-amount", message_amount)
            raise e

    def _handle_sigterm(self, signum, frame):
        self.shutting_down = True
        self.server_socket.close()
        self.barrier.abort()
    
    def run(self):
        action = "accept-connection"
        signal.signal(signal.SIGTERM, self._handle_sigterm)

        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()
            self.server_socket = server_socket
            while True:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except Exception as e:
                    logger.error(action, logger.LogResult.fail)
                    if self.shutting_down:
                        break
                    raise e
                logger.info(action, logger.LogResult.success)

                t = threading.Thread(target=self._handle_client, args=(client_socket,))
                t.start()
                self.client_threads.append(t)

        for t in self.client_threads:
            t.join()
           


