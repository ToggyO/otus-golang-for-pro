package progress_bar

import (
	"fmt"
	"io"
	"strings"
	"sync/atomic"
	"time"
)

const DEFAULT_REFRESH_RATE = time.Millisecond * 200

type ProgressBar struct {
	cur     int64
	total   int64
	percent int64

	rate  string
	graph string

	finish      chan struct{}
	refreshRate time.Duration
}

func NewProgressBar(total int64) *ProgressBar {
	return &ProgressBar{
		total: total,

		finish:      make(chan struct{}),
		refreshRate: DEFAULT_REFRESH_RATE,
	}
}

func (pb *ProgressBar) SetCurrent(current int64) {
	pb.cur = current
}

func (pb *ProgressBar) Start() {
	if pb.graph == "" {
		pb.graph = "#"
	}

	pb.percent = pb.getPercent(pb.cur, pb.total)
	for i := 0; i < int(pb.percent); i += 2 {
		pb.rate += pb.graph
	}
	go pb.refresher()
}

func (pb *ProgressBar) Add(cur int64) {
	atomic.AddInt64(&pb.cur, cur)
}

func (pb *ProgressBar) Finish() {
	close(pb.finish)
	pb.print()
	fmt.Println()
}

func (pb *ProgressBar) NewProxyReader(reader io.Reader) *Reader {
	return &Reader{
		reader: reader,
		pb:     pb,
	}
}

func (pb *ProgressBar) print() {
	cur := atomic.LoadInt64(&pb.cur)
	total := atomic.LoadInt64(&pb.total)
	last := atomic.LoadInt64(&pb.percent)

	percent := pb.getPercent(cur, total)
	atomic.StoreInt64(&pb.percent, percent)

	if percent != last {
		graphsCount := pb.percent / 2
		pb.rate = strings.Repeat(pb.graph, int(graphsCount))
	}
	fmt.Printf("\r[%-50s]%3d%% %8d/%d", pb.rate, pb.percent, cur, total)
}

func (pb *ProgressBar) getPercent(current, total int64) int64 {
	return int64((float32(current) / float32(total)) * 100)
}

func (pb *ProgressBar) refresher() {
	for {
		select {
		case <-pb.finish:
			return
		case <-time.After(pb.refreshRate):
			pb.print()
		}
	}
}
