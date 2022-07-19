package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ToggyO/otus-golang-for-pro/hw07_file_copying/progressbar"
)

const (
	tmpFilePattern     = "hw07_dd_tmp"
	initialPermissions = 0o644
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	stat, err := os.Stat(fromPath)
	if err != nil || !stat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	sourceFileSize := stat.Size()
	if sourceFileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	limit = alignLimit(sourceFileSize, offset, limit)

	src, err := os.Open(fromPath)
	defer handleFileClose(src)
	if err != nil {
		return err
	}

	_, err = src.Seek(offset, 0)
	if err != nil {
		return err
	}

	progressBar := progressbar.NewProgressBar(limit)
	progressBar.Start()
	reader := progressBar.NewProxyReader(src)

	k := filepath.Dir(toPath)
	tmp, err := os.CreateTemp(k, tmpFilePattern)
	if err != nil {
		return err
	}

	_, err = io.CopyN(tmp, reader, limit)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}

	if err = tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}

	if err = os.Chmod(tmp.Name(), os.FileMode(initialPermissions)); err != nil {
		os.Remove(tmp.Name())
		return err
	}

	if err = os.Rename(tmp.Name(), toPath); err != nil {
		os.Remove(tmp.Name())
		return err
	}

	progressBar.Finish()
	return nil
}

func alignLimit(fileSize, offset, limit int64) int64 {
	if limit == 0 {
		return fileSize
	}

	if limit+offset > fileSize {
		return fileSize - offset
	}

	return limit
}

func handleFileClose(src *os.File) {
	if err := src.Close(); err != nil {
		fmt.Println(err)
	}
}
