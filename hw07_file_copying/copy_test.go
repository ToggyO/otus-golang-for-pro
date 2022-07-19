package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	setup(t)
	defer cleanUp(t, testDir)

	wg := sync.WaitGroup{}

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "copy full",
			test: func(t *testing.T) {
				t.Helper()
				wg.Add(1)
				defer wg.Done()
				copyFull(t)
			},
		},
		{
			name: "provide non regular file and get ErrUnsupportedFile",
			test: func(t *testing.T) {
				t.Helper()
				wg.Add(1)
				defer wg.Done()
				nonRegularFile(t)
			},
		},
		{
			name: "get ErrOffsetExceedsFileSize error",
			test: func(t *testing.T) {
				t.Helper()
				wg.Add(1)
				defer wg.Done()
				offsetExceedsError(t)
			},
		},
		{
			name: "offset 0 limit 2048",
			test: func(t *testing.T) {
				t.Helper()
				wg.Add(1)
				defer wg.Done()
				zeroOffsetLimit2048(t)
			},
		},
		{
			name: "offset 10 limit 100",
			test: func(t *testing.T) {
				t.Helper()
				wg.Add(1)
				defer wg.Done()
				offset10Limit100(t)
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, testCase.test)
	}

	wg.Wait()
}

func copyFull(t *testing.T) {
	t.Helper()

	inputFileStat, outputFileStat := execCp(t, 0, 0)
	require.Equal(t, inputFileStat.Size(), outputFileStat.Size())
}

func nonRegularFile(t *testing.T) {
	t.Helper()

	err := Copy("/dev/urandom", getOutputFilePath(), 0, 0)
	result := errors.Is(err, ErrUnsupportedFile)
	require.True(t, result)
}

func offsetExceedsError(t *testing.T) {
	t.Helper()

	err := Copy(inputFilePath, getOutputFilePath(), 8000, 0)
	result := errors.Is(err, ErrOffsetExceedsFileSize)
	require.True(t, result)
}

func zeroOffsetLimit2048(t *testing.T) {
	t.Helper()

	var lim int64 = 2048
	_, outputFileStat := execCp(t, 0, lim)
	require.Equal(t, lim, outputFileStat.Size())
}

func offset10Limit100(t *testing.T) {
	t.Helper()

	var off int64 = 10
	var lim int64 = 100
	_, outputFileStat := execCp(t, off, lim)

	file, err := os.Open(fmt.Sprintf("%s/%s", testDir, outputFileStat.Name()))
	require.NoError(t, err)

	r := bufio.NewReader(file)
	c, _, err := r.ReadRune()
	require.NoError(t, err)

	require.Equal(t, "t", string(c))
}
