package rsh

import (
	"bufio"
	"io"
	"net"
)

type Transport interface {
	io.ReadWriter

	Send(data []byte) error
	Receive() ([]byte, error)
	Close() error
}

func NewTCPServer(addr string) (Transport, error) {
	tcp := &tcpTransport{}
	err := tcp.sListen(addr)
	if err != nil {
		return nil, err
	}

	return tcp, nil
}

func NewTCPClient(addr string) (Transport, error) {
	tcp := &tcpTransport{}
	err := tcp.cConnect(addr)

	if err != nil {
		return nil, err
	}

	return tcp, nil
}

type tcpTransport struct {
	conn net.Conn
}

func (t *tcpTransport) sListen(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	conn, err := ln.Accept()
	if err != nil {
		return err
	}

	t.conn = conn
	return nil
}

func (t *tcpTransport) cConnect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	t.conn = conn
	return nil
}

func (t *tcpTransport) Receive() ([]byte, error) {
	message, err := bufio.NewReader(t.conn).ReadString('\n')
	if err != nil {
		return nil, err
	}
	return []byte(message), nil
}

func (t *tcpTransport) Send(message []byte) error {
	_, err := t.conn.Write(append(message, byte('\n')))
	if err != nil {
		return err
	}

	return nil
}

func (t *tcpTransport) Close() error {
	return t.conn.Close()
}

func (t *tcpTransport) Read(p []byte) (int, error) {
	return t.conn.Read(p)
}

func (t *tcpTransport) Write(p []byte) (int, error) {
	return t.conn.Write(p)
}
