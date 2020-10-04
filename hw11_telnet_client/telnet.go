package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrCloseNilConnection = errors.New("attempt to close nil connection")

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

func NewTelnetClient(hostport string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		hostport: hostport,
		timeout:  timeout,
		in:       in,
		out:      out,
	}
}

type Telnet struct {
	hostport string
	timeout  time.Duration
	conn     net.Conn
	in       io.ReadCloser
	out      io.Writer
}

func (t *Telnet) Connect() error {
	conn, err := net.DialTimeout("tcp", t.hostport, t.timeout)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	t.conn = conn

	return nil
}

func (t *Telnet) Close() error {
	if t.conn == nil {
		return ErrCloseNilConnection
	}
	if err := t.conn.Close(); err != nil {
		return fmt.Errorf("close failed: %w", err)
	}

	return nil
}

func (t *Telnet) Send() error {
	if _, err := io.Copy(t.conn, t.in); err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	return nil
}

func (t *Telnet) Receive() error {
	if _, err := io.Copy(t.out, t.conn); err != nil {
		return fmt.Errorf("receive failed: %w", err)
	}

	return nil
}
