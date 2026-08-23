package client

import (
	"bufio"
	"net"
	"os"
	"time"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/domain"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol"
)

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 200

type ClientConfig struct {
	ServerHost string
	ServerPort string
	AgencyId   string
	InputFile  string
	OutputFile string
}

type Client struct {
	conn   net.Conn
	config ClientConfig
}

func NewClient(config ClientConfig) (*Client, error) {
	conn, err := connectToServer(config.ServerHost, config.ServerPort)
	if err != nil {
		logger.Warn("connect-to-server", logger.Fail)
		return nil, err
	}

	client := &Client{conn: conn, config: config}
	return client, nil
}

func connectToServer(host, port string) (net.Conn, error) {
	const action = "connect-to-server"
	var err error
	var conn net.Conn

	logger.Info(action, logger.InProgress)
	for i := range CONNECTION_ATTEMPTS_MAX {
		conn, err = net.Dial("tcp", host+":"+port)
		if err != nil {
			logger.Warn(action, logger.Fail, "attempt", i)
			time.Sleep(CONNECTION_ATTEMPS_DELAY_MS * time.Millisecond)
			continue
		}

		logger.Info(action, logger.Success)
		break
	}

	return conn, err
}

func (client *Client) Run() error {
	defer client.conn.Close()

	inputFile, err := os.Open(client.config.InputFile)
	if err != nil {
		logger.Error("open-input-file", logger.Fail, "input-file", client.config.InputFile)
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(client.config.OutputFile)
	if err != nil {
		logger.Error("crate-output-file", logger.Fail, "output-file", client.config.OutputFile)
		return err
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		lineWithAgency := client.config.AgencyId + "," + line

		bet, err := domain.ParseBetFromString(lineWithAgency)
		if err != nil {
			logger.Error("parser-error", logger.Fail)
			return err
		}

		if err := protocol.SendBetMessage(bet, client.conn); err != nil {
			logger.Error("send-error", logger.Fail)
			return err
		}

	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if err := protocol.SendDoneMessage(client.conn); err != nil {
		logger.Error("send-error", logger.Fail)
		return err
	}

	response, err := protocol.ReceiveResultMessage(client.conn)
	if err != nil {
		logger.Error("recv-response", logger.Fail)
		return err
	}

	if _, err := writer.WriteString(string(response) + "\n"); err != nil {
		logger.Error("write-response", logger.Fail)
		return err
	}

	return nil
}
