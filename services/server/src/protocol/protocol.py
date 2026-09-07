import code

import safe_socket

SEND_RESULT_MESSAGE = 0
CLIENT_DONE = 1
RECIEVE_BET_CHUNK = 2
SEND_BATCH_ACK = 3
RECIEVE_AGENCY_ID = 4

FIXED_HEADER_SIZE = 5 # 1 byte para el tipo de mensaje + 4 bytes para el tamaño del payload

def _build_packet(msg_type, payload):
    header = msg_type.to_bytes(1, byteorder='big') + len(payload).to_bytes(4, byteorder='big')
    return header + payload

def _receive_packet(socket):
    header_bytes = safe_socket.recv_all(socket, FIXED_HEADER_SIZE)
    if len(header_bytes) == 0:
        return None, None
    msg_type = header_bytes[0]
    size = int.from_bytes(header_bytes[1:5], byteorder='big')
    payload_bytes = safe_socket.recv_all(socket, size)
    return msg_type, payload_bytes
    
def send_result_message(socket, result):
    message = _build_packet(SEND_RESULT_MESSAGE, result.encode('utf-8'))
    return safe_socket.send_all(socket, message)

def recieve_bet_chunk(socket):
    msg_type, payload_bytes = _receive_packet(socket)
    if msg_type is None:
        return None, None
    if not payload_bytes:
        return msg_type, []
    payload = payload_bytes.decode('utf-8')
    bet_lines = [line for line in payload.split('\n') if line.strip()]
    return msg_type, bet_lines
    
def send_ack(socket, result):
    message = _build_packet(SEND_BATCH_ACK, result.encode('utf-8'))
    return safe_socket.send_all(socket, message)

def agency_id_announcement(socket):
    msg_type, agency_id = _receive_packet(socket)
    if msg_type is None or not agency_id:
        return None, None
    return msg_type, agency_id.decode('utf-8')