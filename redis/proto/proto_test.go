package proto

import "testing"

type objectSerialisationTest struct {
	in  Object
	out string
}

var objectSerialisationTestData = []objectSerialisationTest{
	{
		bulk{false, "hello"},
		"$5\r\nhello\r\n",
	},
	{
		bulk{false, ""},
		"$0\r\n\r\n",
	},
	{
		bulk{true, ""},
		"$-1\r\n",
	},
	{
		bulk{true, "This stuff is ignored because of .isNil"},
		"$-1\r\n",
	},
	{
		multiBulk{true, nil},
		"*-1\r\n",
	},
	{
		multiBulk{true, []Object{
			integer(42),
			bulk{false, "Ignored too, because of .isNil"},
		}},
		"*-1\r\n",
	},
	{
		multiBulk{false, []Object{}},
		"*0\r\n",
	},
	{
		multiBulk{false, []Object{
			bulk{false, "Hello, "},
			bulk{false, "World!"},
		}},
		"*2\r\n" +
			"$7\r\nHello, \r\n" +
			"$6\r\nWorld!\r\n",
	},
	{
		multiBulk{false, []Object{
			bulk{false, "Hello, "},
			bulk{true, ""},
			bulk{false, "World!"},
		}},
		"*3\r\n" +
			"$7\r\nHello, \r\n" +
			"$-1\r\n" +
			"$6\r\nWorld!\r\n",
	},
	{
		multiBulk{false, []Object{
			bulk{false, "Hello, "},
			integer(-42),
			bulk{false, "World!"},
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
	{
		Command("HSET", "key", "field", "some \r\n unsafe \x00 stuff"),
		"*4\r\n" +
			"$4\r\nHSET\r\n" +
			"$3\r\nkey\r\n" +
			"$5\r\nfield\r\n" +
			"$22\r\nsome \r\n unsafe \x00 stuff\r\n",
	},
}

func TestObjectSerialisation(t *testing.T) {
	for i, test := range objectSerialisationTestData {
		if test.out != test.in.repr() {
			t.Errorf("#%d: Bad result: %#v (expected %#v)",
				i, test.in.repr(), test.out)
		}
	}
}

func assertKind(t *testing.T, o Object, ok ObjectKind) {
	if o.Kind() != ok {
		t.Errorf("Unexpected ObjectKind %v (expected %v)", o.Kind(), ok)
	}
}

func TestObjectKinds(t *testing.T) {
	assertKind(t, integer(-42), Integer)
	assertKind(t, bulk{false, ""}, String)
	assertKind(t, bulk{true, "nil"}, String)
	assertKind(t, multiBulk{false, []Object{integer(0)}}, List)
	assertKind(t, multiBulk{true, nil}, List)
	assertKind(t, status("OK"), Status)
	assertKind(t, error("Some error."), Error)
}
