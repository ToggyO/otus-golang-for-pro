package progress_bar

import "io"

type Reader struct {
	reader io.Reader
	pb     *ProgressBar
}

func (r *Reader) Read(bytes []byte) (int, error) {
	n, err := r.reader.Read(bytes)
	r.pb.Add(int64(n))
	return n, err
}

func (r *Reader) Close() error {
	r.pb.Finish()
	if closer, ok := r.reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
