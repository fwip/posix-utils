package text

import (
	"bytes"
	"fmt"
	"io"
)

// PieceTable stores data as a series of pieces
type PieceTable struct {
	orig   io.ReaderAt
	append *bytes.Buffer
	pieces []piece // TODO: This should probably be a better data structure
}

// NewPieceTable creates a new PieceTable, ready to use
func NewPieceTable(r io.ReaderAt, length int) PieceTable {
	b := bytes.NewBuffer(nil)
	return PieceTable{
		orig:   r,
		append: b,
		pieces: []piece{{length: length, offset: 0, reader: r}},
	}
}

// piece is a single piece of edit.
type piece struct {
	length int
	offset int64
	reader io.ReaderAt
}

// Read will always read exactly p.length bytes
func (p piece) Read(bytes []byte) (int, error) {
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
	length := 0
	for _, p := range pt.pieces {
		length += p.length
	}

	out := make([]byte, length, length)
	idx := 0
	for _, p := range pt.pieces {
		n, err := p.Read(out[idx : idx+p.length])
		if err != nil {
			panic(err)
		}
		idx += n
	}
	return string(out)
}

// Insert adds text
func (pt *PieceTable) Insert(text []byte, at int) {
	idx := 0
	for i, p := range pt.pieces {
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

			// Add pieces into the array.
			pt.pieces = append(pt.pieces, piece{}, piece{})
			copy(pt.pieces[i+2:], pt.pieces[i:])
			pt.pieces[i] = beforePiece
			pt.pieces[i+1] = newPiece
			pt.pieces[i+2] = afterPiece
			return
		}
	}
	panic(fmt.Errorf("at (%d) too big, max %d", at, idx))
}

func (pt *PieceTable) Delete(length int, at int) {
	end := 0
	for i, p := range pt.pieces {
		end += p.length
		if end >= at {
			start := end - p.length
			pt.pieces[i].length = at - start - 1
			if end > at+length {
				newPiece := piece{length: end - at - length, offset: p.offset + int64(at-start+length), reader: p.reader}
				// Add new Piece
				pt.pieces = append(pt.pieces, piece{})
				copy(pt.pieces[i+1:], pt.pieces[i:])
				pt.pieces[i+1] = newPiece
			}

			for _, p := range pt.pieces {
				fmt.Println(p)
			}
			return
		}
	}
	panic(fmt.Errorf("at (%d) too big, max %d", at, end))
}

// TODO: PieceTable.Read()
