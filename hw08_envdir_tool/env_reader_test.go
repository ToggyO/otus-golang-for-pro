package main

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	setupEnvReaderTests(t)
	defer cleanUpEnvReaderTests(t)

	wg := sync.WaitGroup{}

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "reading success",
			testFunc: func(t *testing.T) {
				t.Helper()
				wg.Add(1)
				defer wg.Done()
				readingSuccess(t)
			},
		},
		{
			name: "clear env value test",
			testFunc: func(t *testing.T) {
				t.Helper()
				wg.Add(1)
				defer wg.Done()
				clearEnvValue(t)
			},
		},
		{
			name: "invalid env var name",
			testFunc: func(t *testing.T) {
				t.Helper()
				wg.Add(1)
				defer wg.Done()
				invalidEnvVarName(t)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.testFunc)
	}

	wg.Wait()
}

func readingSuccess(t *testing.T) {
	t.Helper()

	expected := Environment{
		"BAR":   {"bar", false},
		"UNSET": {"", true},
		"EMPTY": {"", false},
		"FOO":   {"   foo\nwith new line", false},
		"HELLO": {"\"hello\"", false},
	}

	actual, err := ReadDir("testdata/env")
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func clearEnvValue(t *testing.T) {
	t.Helper()

	tempDir, err := os.MkdirTemp(unitTestDir, "clear_env_var")
	require.NoError(t, err)

	dbUser, err := os.CreateTemp(tempDir, "DB_USER")
	require.NoError(t, err)

	_, err = dbUser.Write([]byte("root\t"))
	require.NoError(t, err)

	err = dbUser.Close()
	require.NoError(t, err)

	port, err := os.CreateTemp(tempDir, "PORT")
	require.NoError(t, err)

	_, err = port.Write([]byte("      69\u000069"))
	require.NoError(t, err)

	err = port.Close()
	require.NoError(t, err)

	actual, err := ReadDir(tempDir)
	require.NoError(t, err)
	require.Equal(t, "root",
		actual[filepath.Base(dbUser.Name())].Value,
		"tabs/spaces at the end of the line must be removed")
	require.Equal(t, "      69\n69",
		actual[filepath.Base(port.Name())].Value,
		"null-terminal symbols must be replaced by \\n")

	_ = os.RemoveAll(tempDir)
}

func invalidEnvVarName(t *testing.T) {
	t.Helper()

	tempDir, err := os.MkdirTemp(unitTestDir, "invalid_env_var_name")
	require.NoError(t, err)

	f, err := os.CreateTemp(tempDir, "HOST=")
	require.NoError(t, err)

	err = f.Close()
	require.NoError(t, err)

	_, err = ReadDir(tempDir)
	result := errors.Is(err, ErrorInvalidEnvVarName)
	require.True(t, result, "env var name must not contain `=`")

	_ = os.RemoveAll(tempDir)
}
