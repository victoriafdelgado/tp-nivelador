import socket

# TODO: Complete with a short-read/short-write tolerant implementation

def recv_all(socket: socket.socket, size):
    buffer = []
    bytesRecieved = 0
    while bytesRecieved < size:
        buff = socket.recv(size - bytesRecieved)   
        if not buff:
            if bytesRecieved == 0:
                return b""
            raise ConnectionError("Socket connection failed")
        buffer.append(buff)
        bytesRecieved += len(buff)
    return b''.join(buffer)

def send_all(socket: socket.socket, bytes):
    bytesSent = 0

    while bytesSent < len(bytes):
        n = socket.send(bytes[bytesSent:])

        if n is None:
            raise ConnectionError("Socket connection failed")

        bytesSent += n

    return bytesSent

