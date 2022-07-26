package main

import (
	"os"
	"testing"
)

var unitTestDir = "./unittestdata"

func setupEnvReaderTests(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(unitTestDir); os.IsExist(err) {
		return
	}

	const perm = 0o755
	err := os.Mkdir(unitTestDir, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func cleanUpEnvReaderTests(t *testing.T) {
	t.Helper()

	if err := os.RemoveAll(unitTestDir); err != nil {
		t.Fatal(err)
	}
}
