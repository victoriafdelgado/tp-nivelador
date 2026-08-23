import socket
import logger
import protocol
import domain
from lottery.lottery import Lottery

class Server:
    def __init__(self, server_host: str, server_port: int) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.lottery = Lottery("/tmp/bets.csv")

    def _handle_client(self, client_socket):
        action = "handle-client"
        message_amount = 0
        try:
            logger.info(action, logger.LogResult.in_progress)
            while True:
                msg_type, payload = protocol.recieve_bet_message(client_socket)
                if msg_type == 0:
                    bet = domain.string_to_bet(payload)
                    self.lottery.store_bets([bet])
                    message_amount += 1
                    
                if msg_type == 2:
                    bets = self.lottery.load_bets()
                    
                    winner_strings = []
                    for b in bets:
                        if self.lottery.has_won(b):
                            winner_strings.append(domain.bet_to_string(b))
                    
                    payload = "\n".join(winner_strings) 
                    protocol.send_result_message(client_socket, payload)

                if msg_type is None:
                    logger.info(action, logger.LogResult.success, "messages-amount", message_amount)
                    return
        except Exception as e:
            logger.error(action, logger.LogResult.fail, "messages-amount", message_amount)
            raise e

    def run(self):
        action = "accept-connection"
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()
            while True:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except Exception as e:
                    logger.error(action, logger.LogResult.fail)
                    raise e
                logger.info(action, logger.LogResult.success)

                self._handle_client(client_socket)
