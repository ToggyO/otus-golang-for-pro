package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l := getListener(t)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout := getDuration(t)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err := client.Send()
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

	t.Run("call Connect() multiple times", func(t *testing.T) {
		l := getListener(t)
		defer func() { require.NoError(t, l.Close()) }()

		timeout := getDuration(t)
		client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(&bytes.Buffer{}), &bytes.Buffer{})

		err := client.Connect()
		require.NoError(t, err)

		err = client.Connect()
		require.ErrorIs(t, err, ErrMultipleConnections)
	})

	t.Run("nil in or out channels", func(t *testing.T) {
		l := getListener(t)
		defer func() { require.NoError(t, l.Close()) }()

		timeout := getDuration(t)
		client := NewTelnetClient(l.Addr().String(), timeout, nil, nil)

		err := client.Connect()
		require.ErrorIs(t, err, ErrNotNilMessenger)
	})
}

func getListener(t *testing.T) net.Listener {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)

	return l
}

func getDuration(t *testing.T) time.Duration {
	t.Helper()

	timeout, err := time.ParseDuration("10s")
	require.NoError(t, err)

	return timeout
}
