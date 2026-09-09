package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/domain"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

const (
	ReceiveMessage byte = iota
	DoneSendingBets
	SendBatch
	BatchACK
	AnnounceAgency
)

const Delimiter = ','
const FixedHeaderSize = 5

func serializeBet(bet domain.Bet) []byte {
	var buf []byte
	buf = append(buf, []byte(bet.Name)...)
	buf = append(buf, Delimiter)
	buf = append(buf, []byte(bet.Surname)...)
	buf = append(buf, Delimiter)
	buf = append(buf, []byte(strconv.Itoa(bet.Id))...)
	buf = append(buf, Delimiter)
	buf = append(buf, []byte(bet.DateOfBirth)...)
	buf = append(buf, Delimiter)
	buf = append(buf, []byte(strconv.Itoa(bet.Number))...)
	return buf
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

func SendAnnounceAgencyMessage(agencyId string, socket io.Writer) error {
	payload := []byte(agencyId)
	header := make([]byte, FixedHeaderSize)
	header[0] = AnnounceAgency
	binary.BigEndian.PutUint32(header[1:], uint32(len(payload)))
	message := append(header, payload...)
	return safe_socket.SendAll(socket, message)
}

func SendBatchMessage(bets []domain.Bet, socket io.Writer) error {
	payload := serializeBets(bets)
	header := make([]byte, FixedHeaderSize)
	header[0] = SendBatch
	binary.BigEndian.PutUint32(header[1:], uint32(len(payload)))
	message := append(header, payload...)
	return safe_socket.SendAll(socket, message)
}

func SendDoneMessage(socket io.Writer) error {
	header := make([]byte, FixedHeaderSize)
	header[0] = DoneSendingBets
	binary.BigEndian.PutUint32(header[1:], 0)
	return safe_socket.SendAll(socket, header)
}

func ReceiveResultMessage(socket io.Reader) (string, error) {
	header, err := safe_socket.RecvAll(socket, FixedHeaderSize)
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
	header, err := safe_socket.RecvAll(socket, FixedHeaderSize)
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
