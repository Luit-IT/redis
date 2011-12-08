package proto

import "testing"

type bulkTypeTest struct {
	in  Bulk
	out string
}

var bulkTypeTestData = []bulkTypeTest{
	{
		Bulk{false, []byte("hello")},
		"$5\r\nhello\r\n",
	},
	{
		Bulk{false, []byte{}},
		"$0\r\n\r\n",
	},
	{
		Bulk{true, []byte{}},
		"$-1\r\n",
	},
	{
		Bulk{true, []byte("This stuff is ignored because of .Nil")},
		"$-1\r\n",
	},
}

func TestBulkType(t *testing.T) {
	for i, test := range bulkTypeTestData {
		if test.out != test.in.String() {
			t.Errorf("#%d: Bad result: %#v (expected %#v)",
				i, test.in.String(), test.out)
		}
	}
}
