package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

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
	f, err := os.Create(fmt.Sprintf("data/sessions/%s.bin", id))
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
