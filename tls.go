package main

import (
	"bufio"
	"crypto/tls"
	"net"
	"smtp_honeypot/protocol"
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
	b.byteBuffer = make([]byte, 1)
	b.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, err := b.conn.Read(b.byteBuffer)
	if err != nil {
		return false, err
	}

	return b.byteBuffer[0] == '\x16', nil
}

func (b *BufferedTLSConn) Read(p []byte) (n int, err error) {
	if b.byteBuffer != nil {
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
