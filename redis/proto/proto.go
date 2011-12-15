package proto

import (
	"os"
	"strconv"
)

type integer int64
type bulk struct {
	isNil bool
	bulk  string
}
type multiBulk struct {
	isNil     bool
	multiBulk []Object
}
type status string
type error string

// string methods (raw repr.)

func (i integer) repr() string {
	return ":" + strconv.Itoa64(int64(i)) + "\r\n"
}

func (b bulk) repr() string {
	if b.isNil {
		return "$-1\r\n"
	}
	return "$" + strconv.Itoa64(int64(len(b.bulk))) + "\r\n" +
		b.bulk + "\r\n"
}

func (m multiBulk) repr() string {
	if m.isNil {
		return "*-1\r\n"
	}
	if len(m.multiBulk) == 0 {
		return "*0\r\n"
	}
	s := "*" + strconv.Itoa64(int64(len(m.multiBulk))) + "\r\n"
	for _, c := range m.multiBulk {
		s += c.repr()
	}
	return s
}

func (s status) repr() string {
	return "+" + string(s) + "\r\n"
}

func (e error) repr() string {
	return "-" + string(e) + "\r\n"
}

type ObjectKind int

const (
	Integer ObjectKind = iota // Integer reply
	String                    // Bulk reply
	List                      // Multi-Bulk reply or a Command (request)
	Status                    // Status reply
	Error                     // Error reply
)

// A Redis protocol object, being either a request (i.e. command) or a
// response.
type Object interface {
	repr() string               // Method to get wire format representation
	Kind() ObjectKind           // Method to get the ObjectKind of Object
	Integer() (int64, os.Error) // Method to get the integer value of Object
	String() (string, os.Error) // Method to get the string value of Object
	List() ([]Object, os.Error) // Method to get the list of Objects out of Object with ObjectKind List
	Nil() bool                  // Method to see if Object with ObjectKind String or List is nil
}

// func (o Object) Kind() methods

func (i integer) Kind() ObjectKind   { return Integer }
func (b bulk) Kind() ObjectKind      { return String }
func (m multiBulk) Kind() ObjectKind { return List }
func (s status) Kind() ObjectKind    { return Status }
func (e error) Kind() ObjectKind     { return Error }

// func (o Object) Nil() methods

func (i integer) Nil() bool   { return false } // Shouldn't be called on this type
func (b bulk) Nil() bool      { return b.isNil }
func (m multiBulk) Nil() bool { return m.isNil }
func (s status) Nil() bool    { return false } // Shouldn't be called on this type
func (e error) Nil() bool     { return false } // Shouldn't be called on this type

// func (o Object) Integer() methods

func (i integer) Integer() (int64, os.Error) { return int64(i), nil }
func (b bulk) Integer() (int64, os.Error) {
	if b.isNil {
		return 0, ErrIsNilBulk
	}
	return strconv.Atoi64(b.bulk)
}
func (m multiBulk) Integer() (int64, os.Error) { return 0, ErrNotIntegerableObject }
func (s status) Integer() (int64, os.Error)    { return 0, ErrNotIntegerableObject }
func (e error) Integer() (int64, os.Error)     { return 0, ErrNotIntegerableObject }

// func (o Object) String() methods

func (i integer) String() (string, os.Error) { return strconv.Itoa64(int64(i)), nil }
func (b bulk) String() (string, os.Error) {
	if b.isNil {
		return "", ErrIsNilBulk
	}
	return string(b.bulk), nil
}
func (m multiBulk) String() (string, os.Error) { return "", ErrNotStringableObject }
func (s status) String() (string, os.Error)    { return string(s), nil }
func (e error) String() (string, os.Error)     { return string(e), nil }

// func (o Object) List() methods

func (i integer) List() ([]Object, os.Error) { return nil, ErrNotListObject }
func (b bulk) List() ([]Object, os.Error)    { return nil, ErrNotListObject }
func (m multiBulk) List() ([]Object, os.Error) {
	if m.isNil {
		return nil, ErrIsNilMultiBulk
	}
	return m.multiBulk, nil
}
func (s status) List() ([]Object, os.Error) { return nil, ErrNotListObject }
func (e error) List() ([]Object, os.Error)  { return nil, ErrNotListObject }

// Create a request object.
func Command(args ...string) Object {
	o := multiBulk{false, make([]Object, len(args))}
	for x := range args {
		o.multiBulk[x] = bulk{false, args[x]}
	}
	return o
}

// Some simple checks to make sure the Object interface is satisfied
var (
	_ Object = integer(0)
	_ Object = bulk{false, ""}
	_ Object = multiBulk{false, []Object{integer(0)}}
	_ Object = status("")
	_ Object = error("")
)
