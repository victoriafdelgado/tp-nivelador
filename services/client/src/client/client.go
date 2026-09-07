package client

import (
	"bufio"
	"context"
	"net"
	"os"
	"strconv"
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
	BatchSize  string
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

func (client *Client) Run(ctx context.Context) error {
	defer client.conn.Close()

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			client.conn.Close()
		case <-done:
		}
	}()

	inputFile, err := os.Open(client.config.InputFile)
	if err != nil {
		logger.Error("open-input-file", logger.Fail, "input-file", client.config.InputFile)
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(client.config.OutputFile)
	if err != nil {
		logger.Error("create-output-file", logger.Fail, "output-file", client.config.OutputFile)
		return err
	}
	defer outputFile.Close()

	batchSize, err := strconv.Atoi(client.config.BatchSize)
	if err != nil {
		logger.Error("parse-batch-size", logger.Fail)
		return err
	}

	if err := client.announceAgency(ctx); err != nil {
		return err
	}

	scanner := bufio.NewScanner(inputFile)
	if err := client.sendBets(ctx, scanner, batchSize); err != nil {
		return err
	}

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()
	if err := client.receiveAndSaveWinners(ctx, writer); err != nil {
		return err
	}

	logger.Info("client-run", logger.Success, "agency", client.config.AgencyId)
	return nil
}

func (client *Client) announceAgency(ctx context.Context) error {
	if err := protocol.SendAnnounceAgencyMessage(client.config.AgencyId, client.conn); err != nil {
		logger.Error("send-agency", logger.Fail)
		if ctx.Err() != nil {
			return nil
		}
		return err
	}
	return nil
}

func (client *Client) sendBets(ctx context.Context, scanner *bufio.Scanner, batchSize int) error {
	batch := make([]domain.Bet, 0, batchSize)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			client.conn.Close()
			return nil
		default:
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		bet, err := domain.ParseBetFromString(line)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			logger.Error("parser-error", logger.Fail)
			return err
		}

		batch = append(batch, bet)

		if len(batch) == batchSize {
			if err := client.sendBatchAndReceiveACK(ctx, batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		client.conn.Close()
		return nil
	default:
	}

	if len(batch) > 0 {
		if err := client.sendBatchAndReceiveACK(ctx, batch); err != nil {
			return err
		}
	}
	return nil
}

func (client *Client) sendBatchAndReceiveACK(ctx context.Context, batch []domain.Bet) error {
	if err := protocol.SendBatchMessage(batch, client.conn); err != nil {
		logger.Error("send-error", logger.Fail)
		if ctx.Err() != nil {
			return nil
		}
		return err
	}

	ack, err := protocol.ReceiveBatchACK(client.conn)
	if err != nil || ack != "OK" {
		if ctx.Err() != nil {
			return nil
		}
		return err
	}
	return nil
}

func (client *Client) receiveAndSaveWinners(ctx context.Context, writer *bufio.Writer) error {
	if err := protocol.SendDoneMessage(client.conn); err != nil {
		logger.Error("send-error", logger.Fail)
		if ctx.Err() != nil {
			return nil
		}
		return err
	}

	response, err := protocol.ReceiveResultMessage(client.conn)
	if err != nil {
		logger.Error("recv-response", logger.Fail)
		if ctx.Err() != nil {
			return nil
		}
		return err
	}

	if _, err := writer.WriteString(string(response) + "\n"); err != nil {
		logger.Error("write-response", logger.Fail)
		if ctx.Err() != nil {
			return nil
		}
		return err
	}
	writer.Flush()
	return nil
}
