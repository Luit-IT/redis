package proto

import (
	"io"
	"os"
	"strconv"
)

type Reader struct {
	r io.Reader
}

// Wrap an io.Reader to be able to read Objects from it.
func NewReader(r io.Reader) *Reader { return &Reader{r: r} }

func (r *Reader) readByte() (c byte, err os.Error) {
	br, ok := r.r.(io.ByteReader)
	if ok {
		return br.ReadByte()
	}
	p := make([]byte, 1)
	n, err := r.r.Read(p)
	if n == 0 && err != nil {
		return 0, err
	}
	// os.EOF could have occurred, but doesn't matter, we've got 1
	// character, and Read will return (0, os.EOF) on next read anyway.
	return p[0], nil
}

func (r *Reader) readInt64() (int64, os.Error) {
	// TODO: Need I be more accurate? 20 would be a tighter fit (2 ^ 63 is
	// a 19 (base 10) character long number, +1 for detecting the '\r')
	buf := make([]byte, 32)
	for x := range buf {
		c, err := r.readByte()
		if err != nil {
			return 0, err
		}
		if c == '\r' {
			nc, err := r.readByte()
			if err != nil {
				return 0, err
			}
			if nc != '\n' {
				return 0, ErrProtocolError
			}
			return strconv.Atoi64(string(buf[:x-1]))
		}
		buf[x] = c
	}
	return 0, ErrInt64ReadSize
}

func (r *Reader) readLine() (string, os.Error) {
	buf := make([]byte, 0)
	for {
		c, err := r.readByte()
		if err != nil {
			return "", err
		}
		if c == '\r' {
			nc, err := r.readByte()
			if err != nil {
				return "", err
			}
			if nc != '\n' {
				return "", ErrProtocolError
			}
			return string(buf), nil
		}
		buf = append(buf, c)
	}
	panic("THIS ... SENTENCE ... IS ... FALSE! " +
		"*dontthinkaboutitdontthinkaboutitdontthinkaboutit*")
}

// Read an Object from the underlying io.Reader
func (r *Reader) Read() (Object, os.Error) {
	c, err := r.readByte()
	if err != nil {
		return nil, err
	}
	switch c {
	case ':':
		i, err := r.readInt64()
		if err != nil {
			return nil, err
		}
		return integer(i), nil
	case '$':
		l, err := r.readInt64()
		if err != nil {
			return nil, err
		}
		if l == -1 {
			return bulk{true, ""}, nil
		}
		buf := make([]byte, l+2)
		n, err := r.r.Read(buf)
		if n == len(buf) {
			if buf[n-2] != '\r' || buf[n-1] != '\n' {
				return nil, ErrProtocolError
			}
			return bulk{false, string(buf[:n-2])}, nil
		}
		if err != nil {
			panic("func (r *io.Reader) Read(p []byte) (n int, err os.Error) call didn't fill p, but didn't return err either")
		}
		return nil, err
	case '*':
		l, err := r.readInt64()
		if err != nil {
			return nil, err
		}
		if l == -1 {
			return multiBulk{true, nil}, nil
		}
		buf := make([]Object, l)
		for x := range buf {
			o, err := r.Read()
			if err != nil {
				return nil, err
			}
			buf[x] = o
		}
		return multiBulk{false, buf}, nil
	case '+':
		s, err := r.readLine()
		if err != nil {
			return nil, err
		}
		return status(s), nil
	case '-':
		s, err := r.readLine()
		if err != nil {
			return nil, err
		}
		return error(s), nil
	}
	return nil, ErrProtocolError
}

type Writer struct {
	w io.Writer
}

// Wrap an io.Writer to be able to write Objects to it.
func NewWriter(w io.Writer) *Writer { return &Writer{w: w} }

// Write the wire representation of Object to the underlying io.Writer
func (w *Writer) Write(o Object) os.Error {
	s := o.repr()
	n, err := w.w.Write([]byte(s))
	if n == len(s) {
		// Full write is done. Any error that did occur will come back
		// with the next Write(), so this one is fine.
		return nil
	}
	return err
}
