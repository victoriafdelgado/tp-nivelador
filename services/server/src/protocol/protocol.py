import safe_socket
import socket
import struct

def recieve_bet_message(socket):
    header_bytes = safe_socket.recv_all(socket, 5)
    if len(header_bytes) == 0:
        return None, None 
    msg_type = header_bytes[0]
    size = struct.unpack(">I", header_bytes[1:5])[0]
    payload_bytes = safe_socket.recv_all(socket, size)
    payload = payload_bytes.decode('utf-8')
    return msg_type, payload

def send_result_message(socket, result):
    payload = result.encode('utf-8')
    msg_type = 1
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
    msg_type = 4
    header = struct.pack(">BI", msg_type, len(payload))
    message = header + payload
    return safe_socket.send_all(socket, message)