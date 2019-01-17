package text

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

}

func TestPieceTableSimpleDelete(t *testing.T) {
	orig := bytes.NewReader([]byte("abcdefghi"))

	pt := NewPieceTable(orig, orig.Len())

	pt.Delete(3, 4)
	expectEqual("abchi", pt.String(), t)

	pt.Delete(2, 2)
	expectEqual("abi", pt.String(), t)
}
