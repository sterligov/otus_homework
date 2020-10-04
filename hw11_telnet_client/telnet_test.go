package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("timeout", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("1ns")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			err = client.Connect()
			err = errors.Unwrap(err)

			require.IsType(t, &net.OpError{}, err)
			require.True(t, err.(*net.OpError).Timeout(), "is timeout error")
		}()

		wg.Wait()
	})

	t.Run("close nil connection", func(t *testing.T) {
		client := NewTelnetClient("127.0.0.1", time.Second, nil, nil)
		err := client.Close()

		require.EqualError(t, err, ErrCloseNilConnection.Error())
	})

	t.Run("bad connection hostport", func(t *testing.T) {
		client := NewTelnetClient("bad_hostport", time.Second, nil, nil)
		err := client.Connect()

		require.Error(t, err)
	})

	t.Run("send to close connection", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			require.NoError(t, l.Close())

			in.WriteString("hello\n")
			err = client.Send()
			require.Error(t, err)
		}()

		wg.Wait()
	})

	t.Run("receive error", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), nil)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			require.NoError(t, l.Close())

			in.WriteString("hello\n")
			err = client.Receive()
			require.Error(t, err)
		}()

		wg.Wait()
	})
}
