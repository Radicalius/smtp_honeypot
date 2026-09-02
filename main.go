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
		SrcAddr:      conn.RemoteAddr().String(),
		DstAddr:      conn.LocalAddr().String(),
		StartEpochMs: uint64(time.Now().UnixMilli()),
	}

	defer func() {
		conn.Close()
		logger.Close()
		connection.DurationMs = uint64(time.Now().UnixMilli()) - connection.StartEpochMs
		connLogger.WriteTransaction(connection)

		<-handlerSemaphor
	}()

	conn = NewBufferedTLSConn(conn)
	reader := bufio.NewReader(conn)

	if isTls, err := conn.(*BufferedTLSConn).TLSCheck(); isTls {
		conn, reader, err = TlsUpgrade(conn, &connection)
		if err != nil {
			fmt.Printf("error upgrading tls connection: %s\n", err.Error())
			return
		}
	}

	_sendMessageWithLog(conn, logger, []byte("220 mail.example.com ESMTP\r\n"))

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

			conn, reader, err = TlsUpgrade(conn, &connection)
			if err != nil {
				fmt.Printf("error upgrading tls connection on starttls: %s\n", err.Error())
				return
			}

			continue
		}

		resp := protocol.Handle(&connection, command)

		if resp != "" {
			_sendMessageWithLog(conn, logger, []byte(resp+"\r\n"))
		}

		if resp != "" {
			connection.Commands = append(connection.Commands, protocol.SmtpCommand{
				Command:  string(command),
				Response: resp,
			})
		}
	}
}

func Listen(port string, connLogger *ConnectionLogger) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("error listening on port 2525: %s\n", err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("warning: client failed to connect: %s\n", err.Error())
			continue
		}

		select {
		case handlerSemaphor <- 0:
			go HandleConnection(conn, connLogger)
		default:
			conn.Close()
		}
	}
}

func main() {
	ports := os.Getenv("SMTP_HONEYPOT_PORT")
	if ports == "" {
		ports = "2525"
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

	connectionLogger, err := NewConnectionLogger()
	if err != nil {
		log.Fatal(err)
	}

	allPorts := strings.Split(ports, ",")
	if len(allPorts) > 1 {
		for _, port := range allPorts[1:] {
			go Listen(port, connectionLogger)
		}
	}

	Listen(allPorts[0], connectionLogger)
}
