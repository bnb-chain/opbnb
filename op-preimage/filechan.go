package preimage

import (
	"io"
	"os"
)

// FileChannel is a unidirectional channel for file I/O
type FileChannel interface {
	io.ReadWriteCloser
	// Reader returns the file that is used for reading.
	Reader() *os.File
	// Writer returns the file that is used for writing.
	Writer() *os.File
}

type ReadWritePair struct {
	r *os.File
	w *os.File
}

// NewReadWritePair creates a new FileChannel that uses the given files
func NewReadWritePair(r *os.File, w *os.File) *ReadWritePair {
	return &ReadWritePair{r: r, w: w}
}

func (rw *ReadWritePair) Read(p []byte) (int, error) {
	return rw.r.Read(p)
}

func (rw *ReadWritePair) Write(p []byte) (int, error) {
	return rw.w.Write(p)
}

func (rw *ReadWritePair) Reader() *os.File {
	return rw.r
}

func (rw *ReadWritePair) Writer() *os.File {
	return rw.w
}

func (rw *ReadWritePair) Close() error {
	if err := rw.r.Close(); err != nil {
		return err
	}
	return rw.w.Close()
}

// CreateBidirectionalChannel creates a pair of FileChannels that are connected to each other.
func CreateBidirectionalChannel() (FileChannel, FileChannel, error) {
	ar, bw, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	br, aw, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return NewReadWritePair(ar, aw), NewReadWritePair(br, bw), nil
}

const (
	HClientRFd = 3
	HClientWFd = 4
	PClientRFd = 5
	PClientWFd = 6
)

func ClientHinterChannel() *ReadWritePair {
	r := os.NewFile(HClientRFd, "preimage-hint-read")
	w := os.NewFile(HClientWFd, "preimage-hint-write")
	return NewReadWritePair(r, w)
}

// ClientPreimageChannel returns a FileChannel for the preimage oracle in a detached context
func ClientPreimageChannel() *ReadWritePair {
	r := os.NewFile(PClientRFd, "preimage-oracle-read")
	w := os.NewFile(PClientWFd, "preimage-oracle-write")
	return NewReadWritePair(r, w)
}
