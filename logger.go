package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"smtp_honeypot/protocol"
	"time"
)

var logRootPath = os.Getenv("SMTP_HONEYPOT_LOG_ROOT")

const (
	MESSAGE_DIRECTION_CLIENT_TO_SERVER uint8 = 1
	MESSAGE_DIRECTION_SERVER_TO_CLIENT uint8 = 2
)

type MessageDirection uint8

type SessionFrame struct {
	EpochTimestampMs uint64
	Direction        MessageDirection
	Data             []byte
}

func (f SessionFrame) writeFrame(w io.Writer) error {
	length := uint32(len(f.Data))

	if err := binary.Write(w, binary.BigEndian, f.EpochTimestampMs); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, f.Direction); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}

	_, err := w.Write(f.Data)
	return err
}

type SessionLogger struct {
	guid string
	fd   *os.File
}

func NewSessionLogger(id string) (*SessionLogger, error) {
	f, err := os.Create(fmt.Sprintf(logRootPath+"/sessions/%s.bin", id))
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %s", err.Error())
	}

	return &SessionLogger{
		guid: id,
		fd:   f,
	}, nil
}

func (s *SessionLogger) RecordMessage(dir MessageDirection, data []byte) error {
	frame := SessionFrame{
		EpochTimestampMs: uint64(time.Now().UnixMilli()),
		Direction:        dir,
		Data:             data,
	}
	return frame.writeFrame(s.fd)
}

func (s *SessionLogger) Close() {
	s.fd.Close()
}

type ConnectionLogger struct {
	fd *os.File
	ds string
}

func NewConnectionLogger() (*ConnectionLogger, error) {
	dateStamp := time.Now().Format("2006-01-02")
	f, err := os.Create(fmt.Sprintf(logRootPath+"/transactions/%s.jsonl", dateStamp))
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %s", err.Error())
	}

	return &ConnectionLogger{
		fd: f,
		ds: dateStamp,
	}, nil
}

func (t *ConnectionLogger) WriteTransaction(connection protocol.SmtpConnection) error {
	dateStamp := time.Now().Format("2006-01-02")
	if dateStamp != t.ds {
		f, err := os.Create(fmt.Sprintf(logRootPath+"/transactions/%s.jsonl", dateStamp))
		if err != nil {
			return fmt.Errorf("error opening log file: %s", err.Error())
		}

		t.fd.Close()
		t.fd = f
		t.ds = dateStamp
	}

	data, err := json.Marshal(connection)
	if err != nil {
		return fmt.Errorf("error marshalling transaction: %s\n", err.Error())
	}

	_, err = t.fd.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("error writing transaction: %s\n", err.Error())
	}

	return nil
}

func (c *ConnectionLogger) PublishAxiomTransaction(connection protocol.SmtpConnection) error {
	axiomUrl := os.Getenv("SMTP_HONEYPOT_AXIOM_URL")
	if axiomUrl == "" {
		return fmt.Errorf("SMTP_HONEYPOT_AXIOM_URL not set")
	}

	axiomToken := os.Getenv("SMTP_HONEYPOT_AXIOM_TOKEN")
	if axiomToken == "" {
		return fmt.Errorf("SMTP_HONEYPOT_AXIOM_TOKEN not set")
	}

	data, err := json.Marshal(connection)
	if err != nil {
		return fmt.Errorf("error marshalling transaction: %s\n", err.Error())
	}

	req, err := http.NewRequest("POST", axiomUrl, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer xaat-8bfa0c1a-67f6-44db-a21c-77337a636c88")

	_, err = http.DefaultClient.Do(req)
	return err
}
