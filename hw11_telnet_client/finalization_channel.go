package main

import "sync"

type FinalizationChannel struct {
	C    chan struct{}
	once sync.Once
}

func NewFinalizationChannel() *FinalizationChannel { //nolint:deadcode
	return &FinalizationChannel{C: make(chan struct{})}
}

func NewBufferedFinalizationChannel(bufferSize int) *FinalizationChannel {
	return &FinalizationChannel{C: make(chan struct{}, bufferSize)}
}

func (fc *FinalizationChannel) SafeClose() {
	fc.once.Do(func() {
		close(fc.C)
	})
}
