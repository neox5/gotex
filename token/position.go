package token

import (
	"cmp"
	"fmt"
	"slices"
	"sync"
)

// Position describes a source postion inside a file.
type Position struct {
	Filename string // file name
	Offset   int    // offset = Pos - file.base
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (byte count)
}

// String returns a string in the form file:line:column.
func (p Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.Filename, p.Line, p.Column)
}

// NoPos represents the initialization value for Pos.
const NoPos Pos = 0

// Pos is the absolute Position in the file set.
type Pos int

// File is a handle for a file in a [FileSet].
type File struct {
	name string // file name
	base int    // Pos starting point for this file
	size int    // file size; this gives Pos range of [base... base+size]

	// lines is protected by mutex
	mutex sync.Mutex
	lines []int // lines contain the offset of the first character for each line (lines[0] always 0)
}

// Name returns the file name of file f.
func (f *File) Name() string {
	return f.name
}

// Base returns the base offset of file f.
func (f *File) Base() int {
	return f.base
}

// Size returns the size of file f.
func (f *File) Size() int {
	return f.size
}

// LineCount returns the number of lines in file f.
func (f *File) LineCount() int {
	f.mutex.Lock()
	n := len(f.lines)
	f.mutex.Unlock()
	return n
}

func (f *File) AddLine(offset int) {
	f.mutex.Lock()
	if i := len(f.lines); (i == 0 || f.lines[i-1] < offset) && offset < f.size {
		f.lines = append(f.lines, offset)
	}
	f.mutex.Unlock() // manual unlocking without defer, due to performance costs
}

// Position returns the [Position] value for the given file postion p.
func (f *File) Position(p Pos) (pos Position) {
	if p != NoPos {
		if p >= Pos(f.base) && p < Pos(f.base+f.size) {
			pos = f.position(p)
		}
	}
	return
}

func (f *File) line(offset int) (line int) {
	for i, o := range f.lines {
		if offset < o {
			line = i
			break
		}
	}
	return
}

// column returns the column for a given offset and line number.
func (f *File) column(offset, line int) int {
	return offset - line
}

func (f *File) position(p Pos) Position {
	o := int(p) - f.base
	l := f.line(o)
	c := f.column(o, l)

	return Position{
		Filename: f.name,
		Offset:   o,
		Line:     l,
		Column:   c,
	}
}

type FileSet struct {
	mutex sync.RWMutex // protects the file set
	base  int          // base offset for the next file
	files []*File      // list of files in the order added to the set
}

// NewFileSet creates a new file set.
func NewFileSet() *FileSet {
	return &FileSet{
		base: 1, // 0 == NoPos
	}
}

// Base returns the minimum base which has to be provided to
// [FileSet.AddFile] when adding the next file.
func (s *FileSet) Base() int {
	s.mutex.RLock()
	b := s.base
	s.mutex.RUnlock()
	return b
}

func (s *FileSet) AddFile(filename string, base, size int) *File {
	// Allocate f outside mutex
	f := &File{name: filename, size: size, lines: []int{0}}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if base < s.base {
		panic(fmt.Sprintf("invalid base %d (should be >= %d)", base, s.base))
	}
	f.base = base
	if size < 0 {
		panic(fmt.Sprintf("invalid size %d (should be >= 0)", size))
	}
	base += size + 1 // +1 because EOF also has a position
	if base < 0 {
		panic("token.Pos offset overflow (> 2G of source code in file set)")
	}
	// add the file to the file FileSet
	s.base = base
	s.files = append(s.files, f)
	return f
}

// searchFiles returns the index of the File whose base is ≤ x.
// Assumes 'a' is sorted by base. Out-of-bounds is not checked.
func searchFiles(a []*File, x int) int {
	i, found := slices.BinarySearchFunc(a, x, func(a *File, x int) int {
		return cmp.Compare(a.base, x)
	})
	if !found {
		// If x isn't an exact base match, step back to the previous file
		i--
	}
	return i
}

func (s *FileSet) file(p Pos) *File {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if i := searchFiles(s.files, int(p)); i >= 0 {
		f := s.files[i]
		// check for out-of-bounds
		if int(p) <= f.base+f.size {
			return f
		}
	}
	return nil
}

// File returns the file that contains the position p.
// If no file is found (e.g. p == [NoPos] or p == out-of-bounds),
// the result is nil.
func (s *FileSet) File(p Pos) (f *File) {
	if p != NoPos {
		f = s.file(p)
	}
	return f
}

// Position converts a [Pos] p in the file set into a Position value.
func (s *FileSet) Position(p Pos) (pos Position) {
	if p != NoPos {
		if f := s.file(p); f != nil {
			return f.position(p)
		}
	}
	return
}
