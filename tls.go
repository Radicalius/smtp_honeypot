package main

import (
	"bufio"
	"crypto/tls"
	"net"
	"os"
	"smtp_honeypot/protocol"
	"strconv"
	"time"
)

type BufferedTLSConn struct {
	conn       net.Conn
	byteBuffer []byte
	isTls      bool
}

func NewBufferedTLSConn(conn net.Conn) *BufferedTLSConn {
	return &BufferedTLSConn{
		conn: conn,
	}
}

func (b *BufferedTLSConn) TLSCheck() (bool, error) {
	var err error
	var immTlsWindow int64
	if immTlsWindow, err = strconv.ParseInt(os.Getenv("SMTP_HONEYPOT_IMMEDIATE_TLS_WINDOW"), 10, 64); err != nil {
		immTlsWindow = 100
	}

	b.conn.SetReadDeadline(time.Now().Add(time.Duration(immTlsWindow) * time.Millisecond))

	b.byteBuffer = make([]byte, 1)
	_, err = b.conn.Read(b.byteBuffer)
	if err != nil {
		return false, err
	}

	return b.byteBuffer[0] == '\x16', nil
}

func (b *BufferedTLSConn) Read(p []byte) (n int, err error) {
	if b.byteBuffer != nil && b.byteBuffer[0] != 0 {
		p[0] = b.byteBuffer[0]
		n, err := b.conn.Read(p[1:])
		if err != nil {
			return 0, err
		}

		b.byteBuffer = nil
		return n + 1, err
	}

	return b.conn.Read(p)
}

func (b *BufferedTLSConn) Write(p []byte) (n int, err error) {
	return b.conn.Write(p)
}

func (b *BufferedTLSConn) Close() error {
	return b.conn.Close()
}

func (b *BufferedTLSConn) LocalAddr() net.Addr {
	return b.conn.LocalAddr()
}

func (b *BufferedTLSConn) RemoteAddr() net.Addr {
	return b.conn.RemoteAddr()
}

func (b *BufferedTLSConn) SetDeadline(t time.Time) error {
	return b.conn.SetDeadline(t)
}

func (b *BufferedTLSConn) SetReadDeadline(t time.Time) error {
	return b.conn.SetReadDeadline(t)
}

func (b *BufferedTLSConn) SetWriteDeadline(t time.Time) error {
	return b.conn.SetWriteDeadline(t)
}

func TlsUpgrade(conn net.Conn, connection *protocol.SmtpConnection) (net.Conn, *bufio.Reader, error) {
	tlsConn := tls.Server(conn, tlsConfig)

	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		return conn, nil, err
	}

	reader := bufio.NewReader(tlsConn)
	connection.TLS = true
	connection.TLSInfo = protocol.GetTLSInfo(tlsConn)

	return tlsConn, reader, nil
}
