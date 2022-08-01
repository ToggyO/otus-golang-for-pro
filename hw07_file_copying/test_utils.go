package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	inputFilePath  = "./testdata/input.txt"
	testDir        = "./tmp_test"
	outputFileName = "out_test.txt"
)

func setup(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(testDir); os.IsExist(err) {
		return
	}

	const perm uint32 = 0o755
	err := os.Mkdir(testDir, os.FileMode(perm))
	if err != nil {
		t.Fatal(err)
	}
}

func cleanUp(t *testing.T, dir string) {
	t.Helper()

	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
}

func getRandHash() int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000))
	return n.Int64()
}

func getOutputFilePath() string {
	return filepath.Join(testDir, fmt.Sprintf("%d_%s", getRandHash(), outputFileName))
}

func execCp(t *testing.T, offset, limit int64) (os.FileInfo, os.FileInfo) {
	t.Helper()

	inputFileStat, err := os.Stat(inputFilePath)
	require.NoError(t, err)

	outputFilePath := getOutputFilePath()
	err = Copy(inputFilePath, outputFilePath, offset, limit)
	require.NoError(t, err)

	outputFileStat, err := os.Stat(outputFilePath)
	require.NoError(t, err)

	return inputFileStat, outputFileStat
}
