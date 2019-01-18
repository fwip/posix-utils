package txt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// PieceTable stores data as a series of pieces
type PieceTable struct {
	orig   io.ReaderAt
	append *bytes.Buffer
	head   *piece
}

// NewPieceTable creates a new PieceTable, ready to use
func NewPieceTable(r io.ReaderAt, length int) PieceTable {
	b := bytes.NewBuffer(nil)
	return PieceTable{
		orig:   r,
		append: b,
		head:   &piece{length: length, offset: 0, reader: r, next: nil},
	}
}

// piece is a single piece of edit.
type piece struct {
	length int
	offset int64
	reader io.ReaderAt
	next   *piece
}

func (p piece) String() string {
	if p.length == 0 {
		return ""
	}
	b := make([]byte, p.length)
	p.Read(b)
	return fmt.Sprintf("%d at %d: %s  -> (%s)", p.length, p.offset, string(b), p.next)
}

// Read will always read exactly p.length bytes
func (p piece) Read(bytes []byte) (int, error) {
	if p.length == 0 {
		return 0, nil
	}
	if len(bytes) < p.length {
		return 0, fmt.Errorf("%d bytes available, %d needed", len(bytes), p.length)
	}
	n := 0
	for n < p.length {
		i, err := p.reader.ReadAt(bytes[n:p.length], p.offset+int64(n))
		if err != nil {
			return 0, err
		}
		n += i
	}
	if n != p.length {
		return 0, fmt.Errorf("Meant to read %d bytes, read %d instead", p.length, n)
	}
	return n, nil
}

// String combines all of the pieces of the PieceTable
func (pt PieceTable) String() string {
	//pt.clean()
	length := 0
	for p := pt.head; p != nil; p = p.next {
		length += p.length
	}

	out := make([]byte, length, length)
	idx := 0
	for p := pt.head; p != nil; p = p.next {
		if p.length == 0 {
			continue
		}
		if p.length < 0 {
			panic("Length is < 0")
		}
		n, err := p.Read(out[idx : idx+p.length])
		if err != nil {
			panic(err)
		}
		idx += n
	}
	return string(out)
}

// TODO: not used
func (p *piece) append(pieces ...*piece) {
	tmp := p.next
	for _, p2 := range pieces {
		p.next = p2
		p = p2
	}
	p.next = tmp
}

// Insert adds text
func (pt *PieceTable) Insert(text []byte, at int) {
	idx := 0
	var prev *piece
	for p := pt.head; p != nil; p = p.next {
		idx += p.length
		if idx >= at {
			leftSize := at - (idx - p.length)
			// Split into two pieces and add a new one between them
			beforePiece := piece{length: leftSize, offset: p.offset, reader: p.reader}
			afterPiece := piece{length: p.length - leftSize, offset: p.offset + int64(leftSize), reader: p.reader}
			// Insert in this piece
			buflen := pt.append.Len()
			pt.append.Write(text)
			newPiece := piece{length: len(text), offset: int64(buflen), reader: bytes.NewReader(pt.append.Bytes())}

			// Add pieces into the list.
			if prev == nil {
				pt.head = &beforePiece
			} else {
				prev.next = &beforePiece
			}
			beforePiece.next = &newPiece
			newPiece.next = &afterPiece
			afterPiece.next = p.next
			return
		}
		prev = p
	}
	panic(fmt.Errorf("at (%d) too big, max %d", at, idx))
}

func (pt *PieceTable) Delete(length int, at int) {
	end := 0
	var first *piece
	for p := pt.head; p != nil; p = p.next {
		start := end
		end += p.length
		if end >= at && first == nil {
			first = p
			p.length = at - start
			if end >= at+length {
				newPiece := &piece{
					next:   p.next,
					length: end - (at + length),
					offset: p.offset + int64(at-start+length),
					reader: p.reader,
				}
				p.next = newPiece
				return
			}
		}
		if first != nil && end >= at+length {
			first.next = p
			trimLen := at + length - start
			p.length -= trimLen
			p.offset += int64(at + length - start)
			return
		}
	}
	panic(fmt.Errorf("at (%d) too big, max %d", at, end))
}

func (pt *PieceTable) clean() {
	// Remove 0-length nodes
	for pt.head != nil && pt.head.length <= 0 {
		pt.head = pt.head.next
	}
	var prev *piece
	for p := pt.head; p != nil; p = p.next {
		if p.length <= 0 {
			prev.next = p.next
		}
		prev = p
	}

}

// TODO: PieceTable.Read()

// naiveTable impelements Insert/Delete in a naive way, for ease of testing
type naiveTable []byte

func (nt *naiveTable) String() string {
	return string(*nt)
}
func (nt *naiveTable) Insert(text []byte, at int) {
	old := []byte(*nt)
	new := make([]byte, len(old)+len(text))
	copy(new[:at], old[:at])
	copy(new[at:], text)
	copy(new[at+len(text):], old[at:])
	*nt = naiveTable(new)
}
func (nt *naiveTable) Delete(length int, at int) {
	*nt = append((*nt)[:at], (*nt)[at+length:]...)
}

func Fuzz(data []byte) int {
	s := bufio.NewScanner(bytes.NewReader(data))
	s.Scan()
	if s.Err() != nil {
		return 0
	}
	orig := s.Text()
	pt := NewPieceTable(strings.NewReader(orig), len(orig))
	nt := naiveTable(orig)

	for s.Scan() {
		cmd := s.Text()
		var typ int
		var at int
		var txt string
		n, err := fmt.Sscanf(cmd, "%d %d %s", &typ, &at, &txt)
		if n != 3 || err != nil || at < 0 || typ == 0 {
			return 0
		}
		// Make sure 'at' isn't too big
		if at > len(nt) {
			at = len(nt)
		}
		if typ > 0 {
			// insert
			nt.Insert([]byte(txt), at)
			pt.Insert([]byte(txt), at)
		} else {
			length, err := strconv.Atoi(txt)
			if err != nil || length < 0 || at+length >= len(nt) {
				return 0
			}
			nt.Delete(length, at)
			pt.Delete(length, at)
		}
		if nt.String() != pt.String() {
			panic("problem")
		}
	}
	if s.Err() != nil {
		return 0
	}

	return 0
}
