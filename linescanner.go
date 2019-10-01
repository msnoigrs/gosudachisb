package gosudachisb

import (
	"bufio"
	"io"
)

type normalizer struct {
	r        io.Reader
	lastChar byte
}

func newNormalizer(r io.Reader) *normalizer {
	return &normalizer{r: r}
}

func (norm *normalizer) Read(p []byte) (n int, err error) {
	n, err = norm.r.Read(p)
	for i := 0; i < n; i++ {
		switch {
		case p[i] == '\n' && norm.lastChar == '\r':
			copy(p[i:n], p[i+1:])
			norm.lastChar = p[i]
			n--
			i--
		case p[i] == '\r':
			norm.lastChar = p[i]
			p[i] = '\n'
		default:
			norm.lastChar = p[i]
		}
	}
	return
}

type LineScanner struct {
	r         *bufio.Reader
	line      []byte
	rawBuffer []byte
	err       error
}

func NewLineScanner(r io.Reader) *LineScanner {
	return &LineScanner{r: bufio.NewReader(newNormalizer(r))}
}

func (s *LineScanner) Bytes() []byte {
	return s.line
}

func (s *LineScanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

func (s *LineScanner) Scan() bool {
	s.line, s.err = s.r.ReadSlice('\n')
	if s.err == bufio.ErrBufferFull {
		s.rawBuffer = append(s.rawBuffer[:0], s.line...)
		for s.err == bufio.ErrBufferFull {
			s.line, s.err = s.r.ReadSlice('\n')
			s.rawBuffer = append(s.rawBuffer, s.line...)
		}
		s.line = s.rawBuffer
	}
	if s.err == io.EOF {
		s.err = nil
		if len(s.line) > 0 {
			return true
		} else {
			return false
		}
	}
	if s.err != nil {
		return false
	}
	s.line = s.line[:len(s.line)-1]
	return true
}

func (s *LineScanner) Text() string {
	return string(s.line)
}
