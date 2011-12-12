package proto

import (
	"bufio"
	"io"
	"net"
	"net/textproto"
	"os"
)

type Conn struct {
	objectReader
	objectWriter
	textproto.Pipeline
	conn io.ReadWriteCloser
}

func NewConn(conn io.ReadWriteCloser) *Conn {
	return &Conn{
		objectReader: objectReader{r: bufio.NewReader(conn)},
		objectWriter: objectWriter{w: bufio.NewWriter(conn)},
		conn:         conn,
	}
}

func (c *Conn) Close() os.Error {
	return c.conn.Close()
}

func Dial(network, addr string) (*Conn, os.Error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return NewConn(c), nil
}

func (c *Conn) WriteRequest(o Object) (uint, os.Error) {
	id := c.Next()
	c.StartRequest(id)
	defer c.EndRequest(id)
	err := c.objectWriter.writeObject(o)
	if err != nil {
		return 0, err
	}
	if bw, ok := c.w.(*bufio.Writer); ok {
		go bw.Flush()
	}
	return id, nil

}

func (c *Conn) ReadResponse(id uint) (Object, os.Error) {
	c.StartResponse(id)
	defer c.EndResponse(id)
	return c.readObject()
}

func (c *Conn) Command(args ...string) (Object, os.Error) {
	id, err := c.WriteRequest(Command(args...))
	if err != nil {
		return nil, err
	}
	return c.ReadResponse(id)
}
