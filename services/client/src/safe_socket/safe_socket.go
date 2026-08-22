package safe_socket

import "io"

func SendAll(socket io.Writer, bytes []byte) error {
	bytesSent := 0
	for bytesSent < len(bytes) {
		n, err := socket.Write(bytes[bytesSent:])
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrUnexpectedEOF
		}
		bytesSent += n
	}
	return nil
}

func RecvAll(socket io.Reader, size int) ([]byte, error) {
	buff := make([]byte, size)
	bytesReceived := 0

	for bytesReceived < size {
		n, err := socket.Read(buff[bytesReceived:])
		if n > 0 {
			bytesReceived += n
		}
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
	}
	return buff, nil
}
