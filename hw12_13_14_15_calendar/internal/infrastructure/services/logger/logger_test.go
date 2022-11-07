package logger

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger("WarN", true, &buf)
	require.NotNil(t, logger)

	t.Run("is different levels are printed", func(t *testing.T) {
		reader := bufio.NewReader(&buf)

		logger.Info("Information")
		output, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		message := string(output)
		require.Empty(t, message)

		buf.Reset()

		logger.Warn("Warning")
		output, err = reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		message = string(output)
		require.Contains(t, message, "Warning")
	})
}
