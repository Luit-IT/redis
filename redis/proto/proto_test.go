package proto

import (
	"bufio"
	"strings"
	"testing"
)

type objectSerialisationTest struct {
	in  Object
	out string
}

var objectSerialisationTestData = []objectSerialisationTest{
	{
		bulk{false, []byte("hello")},
		"$5\r\nhello\r\n",
	},
	{
		bulk{false, []byte{}},
		"$0\r\n\r\n",
	},
	{
		bulk{true, nil},
		"$-1\r\n",
	},
	{
		bulk{true, []byte("This stuff is ignored because of .isNil")},
		"$-1\r\n",
	},
	{
		multiBulk{true, nil},
		"*-1\r\n",
	},
	{
		multiBulk{true, []Object{
			integer(42),
			bulk{false, []byte("Ignored too, because of .isNil")},
		}},
		"*-1\r\n",
	},
	{
		multiBulk{false, []Object{}},
		"*0\r\n",
	},
	{
		multiBulk{false, []Object{
			bulk{false, []byte("Hello, ")},
			bulk{false, []byte("World!")},
		}},
		"*2\r\n" +
			"$7\r\nHello, \r\n" +
			"$6\r\nWorld!\r\n",
	},
	{
		multiBulk{false, []Object{
			bulk{false, []byte("Hello, ")},
			bulk{true, nil},
			bulk{false, []byte("World!")},
		}},
		"*3\r\n" +
			"$7\r\nHello, \r\n" +
			"$-1\r\n" +
			"$6\r\nWorld!\r\n",
	},
	{
		multiBulk{false, []Object{
			bulk{false, []byte("Hello, ")},
			integer(-42),
			bulk{false, []byte("World!")},
			error("Whoa!"),
		}},
		"*4\r\n" +
			"$7\r\nHello, \r\n" +
			":-42\r\n" +
			"$6\r\nWorld!\r\n" +
			"-Whoa!\r\n",
	},
	{
		status("OK"),
		"+OK\r\n",
	},
	{
		error("Something something"),
		"-Something something\r\n",
	},
}

func TestObjectSerialisation(t *testing.T) {
	for i, test := range objectSerialisationTestData {
		if test.out != test.in.string() {
			t.Errorf("#%d: Bad result: %#v (expected %#v)",
				i, test.in.string(), test.out)
		}
	}
}

func TestReadInteger(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(":-32\r\n"))
	o, err := readObject(r)
	if err != nil {
		t.Errorf("readObject() failed: %s", err.String())
		return
	}
	if o.Type() != Integer {
		t.Errorf("Unexpected ObjectType: %v (expected %v)", o.Type(), Integer)
		return
	}
	i, err := o.Integer()
	if err != nil {
		t.Errorf("Object.Integer() failed: %s", err.String())
		return
	}
	if int64(i) != int64(-32) {
		t.Errorf("Bad readInteger() result: %d (expected -32)", i)
	}
}
