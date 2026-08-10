package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"smtp_honeypot/protocol"
	"strings"
)

var tlsConfig *tls.Config

func HandleConnection(conn net.Conn) {
	var trans protocol.SmtpTransaction

	conn.Write([]byte("220 mail.example.com ESMTP\r\n"))

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("error reading message from client: %s\n", err.Error())
			return
		}

		command := bytes.TrimSuffix(buf[:n], []byte("\r\n"))
		if string(command) == "QUIT" {
			conn.Write([]byte("221 2.0.0 Bye\r\n"))
			conn.Close()

			fmt.Printf("%v\n", trans)

			return
		}

		if string(command) == "STARTTLS" {
			conn.Write([]byte("220 2.0.0 Ready to start TLS\r\n"))

			tlsConn := tls.Server(conn, tlsConfig)

			if err := tlsConn.Handshake(); err != nil {
				conn.Close()
				return
			}

			conn = tlsConn
			continue
		}

		resp := protocol.Handle(&trans, command)

		conn.Write([]byte(resp + "\r\n"))
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("warning: client failed to connect: %s\n", err.Error())
			continue
		}

		go HandleConnection(conn)
	}
}
