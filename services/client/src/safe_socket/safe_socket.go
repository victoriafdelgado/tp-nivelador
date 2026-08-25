package safe_socket

import "io"

func SendAll(socket io.Writer, bytes []byte) error {
	bytesSent := 0
	for bytesSent < len(bytes) {
		n, err := socket.Write(bytes[bytesSent:])
		if n > 0 {
			bytesSent += n
		}
		if err != nil {
			return err
		}
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
		if bytesReceived == size {
			return buff, nil
		}
		if err != nil {
			if err == io.EOF {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}
	}
	return buff, nil
}
