package proto

import (
	"os"
	"strconv"
)

type ObjectType int

const (
	Integer ObjectType = iota
	String
	List
	Status
	Error
)

type Object interface {
	Type() ObjectType
	Nil() bool
	Integer() (int64, os.Error)
	String() (string, os.Error) // Also for ObjectType Status and Error
	List() ([]Object, os.Error)

	string() string // Internal method to get raw representation
}

type integer int64
type bulk struct {
	isNil bool
	bulk  []byte
}
type multiBulk struct {
	isNil     bool
	multiBulk []Object
}
type status string
type error string

// Type methods

func (i integer) Type() ObjectType   { return Integer }
func (b bulk) Type() ObjectType      { return String }
func (m multiBulk) Type() ObjectType { return List }
func (s status) Type() ObjectType    { return Status }
func (e error) Type() ObjectType     { return Error }

// Nil methods

func (i integer) Nil() bool   { return false } // Shouldn't be called on this type
func (b bulk) Nil() bool      { return b.isNil }
func (m multiBulk) Nil() bool { return m.isNil }
func (s status) Nil() bool    { return false } // Shouldn't be called on this type
func (e error) Nil() bool     { return false } // Shouldn't be called on this type

// Integer methods

func (i integer) Integer() (int64, os.Error) { return int64(i), nil }
func (b bulk) Integer() (int64, os.Error) {
	if b.isNil {
		return 0, os.NewError("Not an Integer")
	}
	i, err := strconv.Atoi64(string(b.bulk))
	if err != nil {
		return 0, os.NewError("Not an Integer")
	}
	return i, nil
}
func (m multiBulk) Integer() (int64, os.Error) { return 0, os.NewError("Not an Integer") }
func (s status) Integer() (int64, os.Error)    { return 0, os.NewError("Not an Integer") }
func (e error) Integer() (int64, os.Error)     { return 0, os.NewError("Not an Integer") }

// String methods

func (i integer) String() (string, os.Error) { return strconv.Itoa64(int64(i)), nil }
func (b bulk) String() (string, os.Error) {
	if b.Nil() {
		return "", os.NewError("String Object is Nil")
	}
	return string(b.bulk), nil
}
func (m multiBulk) String() (string, os.Error) { return "", os.NewError("Can't String() a List Object") }
func (s status) String() (string, os.Error)    { return string(s), nil }
func (e error) String() (string, os.Error)     { return string(e), nil }

// List methods

func (i integer) List() ([]Object, os.Error) { return nil, os.NewError("Not a List") }
func (b bulk) List() ([]Object, os.Error)    { return nil, os.NewError("Not a List") }
func (m multiBulk) List() ([]Object, os.Error) {
	if m.Nil() {
		return nil, os.NewError("List Object is Nil")
	}
	return m.multiBulk, nil
}
func (s status) List() ([]Object, os.Error) { return nil, os.NewError("Not a List") }
func (e error) List() ([]Object, os.Error)  { return nil, os.NewError("Not a List") }

// string methods (raw repr.)

func (i integer) string() string {
	return ":" + strconv.Itoa64(int64(i)) + "\r\n"
}

func (b bulk) string() string {
	if b.isNil {
		return "$-1\r\n"
	}
	return "$" + strconv.Itoa64(int64(len(b.bulk))) + "\r\n" +
		string(b.bulk) + "\r\n"
}

func (m multiBulk) string() string {
	if m.isNil {
		return "*-1\r\n"
	}
	if len(m.multiBulk) == 0 {
		return "*0\r\n"
	}
	s := "*" + strconv.Itoa64(int64(len(m.multiBulk))) + "\r\n"
	for _, c := range m.multiBulk {
		s += c.string()
	}
	return s
}

func (s status) string() string {
	return "+" + string(s) + "\r\n"
}

func (e error) string() string {
	return "-" + string(e) + "\r\n"
}

// Some simple checks to make sure the Object interface is satisfied
var (
	_ Object = integer(0)
	_ Object = bulk{false, []byte{}}
	_ Object = multiBulk{false, []Object{integer(0)}}
	_ Object = status("")
	_ Object = error("")
)
