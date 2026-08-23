import safe_socket
import socket
import struct

def recieve_bet_message(sock):
    header_bytes = safe_socket.recv_all(sock, 5)
    if len(header_bytes) == 0:
        return None, None 
    msg_type = header_bytes[0]
    size = struct.unpack(">I", header_bytes[1:5])[0]
    payload_bytes = safe_socket.recv_all(sock, size)
    payload = payload_bytes.decode('utf-8')
    return msg_type, payload

def send_result_message(socket: socket.socket, result):
    payload = result.encode('utf-8')
    msg_type = 1
    header = struct.pack(">BI", msg_type, len(payload))
    message = header + payload
    return safe_socket.send_all(socket, message)