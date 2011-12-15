package proto

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

func readByte(r io.Reader) (c byte, err os.Error) {
	if br, ok := r.(io.ByteReader); ok {
		return br.ReadByte()
	}
	p := make([]byte, 1)
	n, err := r.Read(p)
	if n == 0 && err != nil {
		return 0, err
	}
	// os.EOF could have occurred, but doesn't matter, we've got 1
	// character, and Read will return (0, os.EOF) on next read anyway.
	return p[0], nil
}

func readInt64(r io.Reader) (int64, os.Error) {
	// TODO: Need I be more accurate? 20 would be a tighter fit (2 ^ 63 is
	// a 19 (base 10) character long number, +1 for detecting the '\r')
	buf := make([]byte, 32)
	for x := range buf {
		c, err := readByte(r)
		if err != nil {
			return 0, err
		}
		if c == '\r' {
			nc, err := readByte(r)
			if err != nil {
				return 0, err
			}
			if nc != '\n' {
				return 0, ErrProtocolError
			}
			return strconv.Atoi64(string(buf[:x]))
		}
		buf[x] = c
	}
	return 0, ErrInt64ReadSize
}

func readLine(r io.Reader) (string, os.Error) {
	// Try to use a bufio.Reader
	if br, ok := r.(*bufio.Reader); ok {
		wholeLine := ""
		isPrefix := true
		for isPrefix {
			var line []byte
			var err os.Error
			line, isPrefix, err = br.ReadLine()
			if err != nil {
				return "", err
			}
			wholeLine += string(line)
		}
		return wholeLine, nil
	}
	// No bufio.Reader, so we have to read byte by byte
	buf := make([]byte, 0)
	for {
		c, err := readByte(r)
		if err != nil {
			return "", err
		}
		if c == '\r' {
			nc, err := readByte(r)
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

// Read an Object from the io.Reader
func ReadObject(r io.Reader) (Object, os.Error) {
	c, err := readByte(r)
	if err != nil {
		return nil, err
	}
	switch c {
	case ':':
		i, err := readInt64(r)
		if err != nil {
			return nil, err
		}
		return integer(i), nil
	case '$':
		l, err := readInt64(r)
		if err != nil {
			return nil, err
		}
		if l == -1 {
			return bulk{true, ""}, nil
		}
		buf := make([]byte, l+2)
		n, err := r.Read(buf)
		if n == len(buf) {
			if buf[n-2] != '\r' || buf[n-1] != '\n' {
				return nil, ErrProtocolError
			}
			return bulk{false, string(buf[:n-2])}, nil
		}
		if err != nil {
			panic("func (r *io.Reader) Read(p []byte) (n int, " +
				"err os.Error) call didn't fill p, but " +
				"didn't return err either")
		}
		return nil, err
	case '*':
		l, err := readInt64(r)
		if err != nil {
			return nil, err
		}
		if l == -1 {
			return multiBulk{true, nil}, nil
		}
		buf := make([]Object, l)
		for x := range buf {
			o, err := ReadObject(r)
			if err != nil {
				return nil, err
			}
			buf[x] = o
		}
		return multiBulk{false, buf}, nil
	case '+':
		s, err := readLine(r)
		if err != nil {
			return nil, err
		}
		return status(s), nil
	case '-':
		s, err := readLine(r)
		if err != nil {
			return nil, err
		}
		return error(s), nil
	}
	return nil, ErrProtocolError
}

// Write an Object to the io.Writer
func WriteObject(w io.Writer, o Object) os.Error {
	s := o.repr()
	n, err := w.Write([]byte(s))
	if n == len(s) {
		// Full write is done. Any error that did occur will come back
		// with the next WriteObject() call, so this one is fine.
		return nil
	}
	return err
}
