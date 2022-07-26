package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

var ErrorInvalidEnvVarName = errors.New("invalid env variable name")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntry {
		// Recursive directory scanning
		if entry.IsDir() {
			readDirResult, err := ReadDir(entry.Name())
			if err != nil {
				return env, err
			}
			mergeIntoLeftMap(env, readDirResult)
		}

		fileName := entry.Name()
		if strings.Contains(fileName, "=") {
			return nil, ErrorInvalidEnvVarName
		}

		f, err := os.Open(path.Join(dir, fileName))
		if err != nil {
			return nil, err
		}

		err = handleEnv(env, f, fileName)
		if err != nil {
			return nil, err
		}

		err = f.Close()
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}

func handleEnv(envMap Environment, file *os.File, fileName string) error {
	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	envVar := strings.TrimRight(string(bytes.ReplaceAll(line, []byte{0x00}, []byte{'\n'})), "\t\n ")

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	envMap[fileName] = EnvValue{Value: envVar, NeedRemove: stat.Size() == 0}
	return nil
}

func mergeIntoLeftMap(m1, m2 Environment) {
	for k, v := range m2 {
		m1[k] = v
	}
}
