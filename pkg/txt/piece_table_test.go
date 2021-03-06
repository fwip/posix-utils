package txt

import (
	"bytes"
	"testing"
)

func expectEqual(expected, actual string, t *testing.T) {
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestPieceTableSimpleInsert(t *testing.T) {
	orig := bytes.NewReader([]byte("abcdefghi"))

	pt := NewPieceTable(orig, orig.Len())

	pt.Insert([]byte("123"), 3)

	out := pt.String()
	expected := "abc123defghi"

	if out != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}

	pt.Insert([]byte("!"), 0)
	expectEqual("!abc123defghi", pt.String(), t)

}

func TestPieceTableSimpleDelete(t *testing.T) {
	orig := bytes.NewReader([]byte("123456789"))

	pt := NewPieceTable(orig, orig.Len())

	pt.Delete(1, 0)
	expectEqual("23456789", pt.String(), t)

	pt.Delete(1, 4)
	expectEqual("2345789", pt.String(), t)

	pt.Delete(2, 3)
	expectEqual("23489", pt.String(), t)
}

func TestNaiveTableSimpleInsert(t *testing.T) {
	nt := naiveTable([]byte("abcdefghi"))

	nt.Insert([]byte("123"), 3)
	expectEqual("abc123defghi", nt.String(), t)

	nt.Insert([]byte("!"), 0)
	expectEqual("!abc123defghi", nt.String(), t)
}

func TestNaiveTableSimpleDelete(t *testing.T) {

	nt := naiveTable([]byte("123456789"))

	nt.Delete(1, 0)
	expectEqual("23456789", nt.String(), t)

	nt.Delete(1, 4)
	expectEqual("2345789", nt.String(), t)

	nt.Delete(2, 3)
	expectEqual("23489", nt.String(), t)
}
