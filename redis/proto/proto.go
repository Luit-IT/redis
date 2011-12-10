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
	// Integer reply
	Integer ObjectKind = iota
	// Bulk reply
	String
	// Multi-Bulk reply, or a Command
	List
	// Status reply
	Status
	// Error reply
	Error
)

// A Redis protocol object, being either a request (i.e. command) or a
// response.
type Object interface {
	repr() string // Method to get wire format representation
}

// Find out what type of Object you're dealing with.
func Kind(o Object) ObjectKind {
	switch o.(type) {
	case integer:
		return Integer
	case bulk:
		return String
	case multiBulk:
		return List
	case status:
		return Status
	case error:
		return Error
	}
	panic("Object isn't a proper Redis object")
}

// Create a request object.
func Command(args ...string) Object {
	o := multiBulk{false, make([]Object, len(args))}
	for x := range args {
		o.multiBulk[x] = bulk{false, args[x]}
	}
	return o
}


// Convert object to int64. Returns ErrNotIntegerableObject error for List,
// Status and Error Kinds of object. Integer objects are int64 behind the
// scenes already. Bulk objects (String-Kind object) are run through
// strconv.Atoi64 and return accordingly. With nil bulk this returns an
// ErrIsNilBulk error.
func ObjectInteger(o Object) (int64, os.Error) {
	switch v := o.(type) {
	case integer:
		return int64(v), nil
	case bulk:
		if v.isNil {
			return 0, ErrIsNilBulk
		}
		return strconv.Atoi64(v.bulk)
	}
	return 0, ErrNotIntegerableObject
}

// Convert objects to string. Doesn't work on List-Kind objects. Integer
// objects are converted through strconv.Itoa64. Nil bulk objects return an
// ErrIsNilBulk error.
func ObjectString(o Object) (string, os.Error) {
	switch v := o.(type) {
	case integer:
		return strconv.Itoa64(int64(v)), nil
	case bulk:
		if v.isNil {
			return "", ErrIsNilBulk
		}
		return v.bulk, nil
	case status:
	case error:
		return string(v), nil
	}
	return "", ErrNotStringableObject
}

// Convert object to []Object. Only works on List-Kind objects.
func ObjectList(o Object) ([]Object, os.Error) {
	switch v := o.(type) {
	case multiBulk:
		return v.multiBulk, nil
	}
	return nil, ErrNotListObject
}

// Check if object is Nil
func IsNil(o Object) bool {
	switch v := o.(type) {
	case bulk:
		return v.isNil
	case multiBulk:
		return v.isNil
	}
	// TODO: should this be false, or error?
	return false
}

// Some simple checks to make sure the Object interface is satisfied
var (
	_ Object = integer(0)
	_ Object = bulk{false, ""}
	_ Object = multiBulk{false, []Object{integer(0)}}
	_ Object = status("")
	_ Object = error("")
)
