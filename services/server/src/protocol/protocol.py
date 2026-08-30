import safe_socket
import socket
import struct

SEND_RESULT_MESSAGE = 0
CLIENT_DONE = 1
RECIEVE_BET_CHUNK = 2
SEND_BATCH_ACK = 3

def send_result_message(socket, result):
    payload = result.encode('utf-8')
    msg_type = SEND_RESULT_MESSAGE
    header = struct.pack(">BI", msg_type, len(payload))
    message = header + payload
    return safe_socket.send_all(socket, message)

def recieve_bet_chunk(socket):
    header_bytes = safe_socket.recv_all(socket, 5)
    if len(header_bytes) == 0:
        return None, None
    msg_type = header_bytes[0]
    size = struct.unpack(">I", header_bytes[1:5])[0]

    if size == 0:
        return msg_type, []
    
    payload_bytes = safe_socket.recv_all(socket, size)
    payload = payload_bytes.decode('utf-8')

    bet_lines = [line for line in payload.split('\n') if line.strip()]
    return msg_type, bet_lines
    
def send_ack(socket, result):
    payload = result.encode('utf-8')
    msg_type = SEND_BATCH_ACK
    header = struct.pack(">BI", msg_type, len(payload))
    message = header + payload
    return safe_socket.send_all(socket, message)