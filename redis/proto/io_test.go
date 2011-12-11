package proto

import (
	"bytes"
	"testing"
)

func TestWriteObjects(t *testing.T) {
	// Running tests on the same data as ./proto_test.go does with
	// TestObjectSerialisation
	for i, test := range objectSerialisationTestData {
		buf := bytes.NewBuffer([]byte{})
		w := &ObjectWriter{w: buf}
		err := w.WriteObject(test.in)
		if err != nil {
			t.Fatalf("#%d: Write failed: %s", i, err)
		}
		got := buf.String()
		if test.out != got {
			t.Fatalf("#%d: Wrong data written: %#v (expected %#v)",
				i, got, test.out)
		}
	}
}
