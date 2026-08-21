package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"smtp_honeypot/protocol"
	"strings"
	"time"

	"github.com/google/uuid"
)

var tlsConfig *tls.Config
var handlerSemaphor = make(chan int, 10)

func _sendMessageWithLog(conn net.Conn, logger *SessionLogger, data []byte) {
	conn.Write(data)
	err := logger.RecordMessage(MessageDirection(MESSAGE_DIRECTION_SERVER_TO_CLIENT), data)
	if err != nil {
		fmt.Printf("error recording log: %s\n", err.Error())
	}
}

func HandleConnection(conn net.Conn, connLogger *ConnectionLogger) {
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	id := uuid.New().String()

	logger, err := NewSessionLogger(id)
	if err != nil {
		fmt.Printf("warning: session logging failed to initialize: %s\n", err.Error())
		return
	}

	connection := protocol.SmtpConnection{
		Guid:         id,
		SrcAddr:      conn.LocalAddr().String(),
		StartEpochMs: uint64(time.Now().UnixMilli()),
	}

	defer func() {
		conn.Close()
		logger.Close()
		connection.DurationMs = uint64(time.Now().UnixMilli()) - connection.StartEpochMs
		connLogger.WriteTransaction(connection)

		<-handlerSemaphor
	}()

	_sendMessageWithLog(conn, logger, []byte("220 mail.example.com ESMTP\r\n"))

	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			if os.IsTimeout(err) || err == io.EOF {
				return
			}

			fmt.Printf("error reading message from client: %s\n", err.Error())
			return
		}

		command := bytes.TrimRight(buf, "\r\n")
		logger.RecordMessage(MessageDirection(MESSAGE_DIRECTION_CLIENT_TO_SERVER), buf)

		if strings.Contains(string(command), "QUIT") {
			_sendMessageWithLog(conn, logger, []byte("221 2.0.0 Bye\r\n"))
			return
		}

		if strings.Contains(string(command), "STARTTLS") {
			_sendMessageWithLog(conn, logger, []byte("220 2.0.0 Ready to start TLS\r\n"))

			tlsConn := tls.Server(conn, tlsConfig)

			if err := tlsConn.Handshake(); err != nil {
				conn.Close()
				return
			}

			conn = tlsConn
			reader = bufio.NewReader(conn)
			connection.TLS = true
			connection.TLSInfo = protocol.GetTLSInfo(tlsConn)
			continue
		}

		resp := protocol.Handle(&connection, command)

		if resp != "" {
			_sendMessageWithLog(conn, logger, []byte(resp+"\r\n"))
		}
	}
}

func main() {
	port := os.Getenv("SMTP_HONEYPOT_PORT")
	if port == "" {
		port = "2525"
	}

	certPath := os.Getenv("SMTP_HONEYPOT_CERT_PATH")
	if certPath == "" {
		certPath = "./"
	}

	if !strings.HasSuffix(certPath, "/") {
		certPath = certPath + "/"
	}

	cert, err := tls.LoadX509KeyPair(certPath+"server.crt", certPath+"server.key")
	if err != nil {
		log.Fatalf("error loading certs: %s", err.Error())
	}

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("error listening on port 2525: %s\n", err.Error())
	}

	transactionLogger, err := NewTransactionLogger()
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("warning: client failed to connect: %s\n", err.Error())
			continue
		}

		select {
		case handlerSemaphor <- 0:
			go HandleConnection(conn, transactionLogger)
		default:
			conn.Close()
		}
	}
}
