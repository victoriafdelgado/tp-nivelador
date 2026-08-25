package protocol

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/domain"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

type MessageType int

const (
	SendMessage = iota
	ReceiveMessage
	DoneSendingBets
	SendBatch
	BatchACK
)

func serializeBet(bet domain.Bet) []byte {
	payload := fmt.Sprintf("%d,%s,%s,%d,%s,%d", bet.AgencyId, bet.Name, bet.Surname, bet.Id, bet.DateOfBirth, bet.Number)
	return []byte(payload)
}

func serializeBets(bets []domain.Bet) []byte {
	var serializedBets []byte
	for i, bet := range bets {
		serializedBets = append(serializedBets, serializeBet(bet)...)
		if i < len(bets)-1 {
			serializedBets = append(serializedBets, '\n')
		}
	}
	return serializedBets
}

func SendBetMessage(bet domain.Bet, socket io.Writer) error {
	payload := serializeBet(bet)
	header := make([]byte, 5)
	header[0] = SendMessage
	binary.BigEndian.PutUint32(header[1:], uint32(len(payload)))
	message := append(header, payload...)
	return safe_socket.SendAll(socket, message)
}

func SendBatchMessage(bets []domain.Bet, socket io.Writer) error {
	payload := serializeBets(bets)
	header := make([]byte, 5)
	header[0] = SendBatch
	binary.BigEndian.PutUint32(header[1:], uint32(len(payload)))
	message := append(header, payload...)
	return safe_socket.SendAll(socket, message)
}

func SendDoneMessage(socket io.Writer) error {
	header := make([]byte, 5)
	header[0] = byte(DoneSendingBets)
	binary.BigEndian.PutUint32(header[1:], 0)
	return safe_socket.SendAll(socket, header)
}

func ReceiveResultMessage(socket io.Reader) (string, error) {
	header, err := safe_socket.RecvAll(socket, 5)
	if err != nil {
		return "", err
	}
	if header[0] != ReceiveMessage {
		return "", fmt.Errorf("Codigo invalido")
	}
	length := binary.BigEndian.Uint32(header[1:])
	payload, err := safe_socket.RecvAll(socket, int(length))
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func ReceiveBatchACK(socket io.Reader) (string, error) {
	header, err := safe_socket.RecvAll(socket, 5)
	if err != nil {
		return "", err
	}
	if header[0] != BatchACK {
		return "", fmt.Errorf("Codigo invalido")
	}
	length := binary.BigEndian.Uint32(header[1:])
	payload, err := safe_socket.RecvAll(socket, int(length))
	if err != nil {
		return "", err
	}
	return string(payload), nil
}
