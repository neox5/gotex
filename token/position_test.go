package token

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func checkPos(t *testing.T, msg string, got, want Position) {
	if got.Filename != want.Filename {
		t.Errorf("%s: got filename = %q; want %q", msg, got.Filename, want.Filename)
	}
	if got.Offset != want.Offset {
		t.Errorf("%s: got offset = %d; want %d", msg, got.Offset, want.Offset)
	}
	if got.Line != want.Line {
		t.Errorf("%s: got line = %d; want %d", msg, got.Line, want.Line)
	}
	if got.Column != want.Column {
		t.Errorf("%s: got column = %d; want %d", msg, got.Column, want.Column)
	}
}

func TestNoPos(t *testing.T) {
	var fset *FileSet
	checkPos(t, "nil NoPos", fset.Position(NoPos), Position{})
	fset = NewFileSet()
	checkPos(t, "fset NoPos", fset.Position(NoPos), Position{})
}

var tests = []struct {
	filename string
	source   []byte // may be nil
	size     int
	lines    []int
}{
	{"document.tex", []byte{}, 0, []int{}},
	{"preamble.tex", []byte("\\documentclass{article}"), 23, []int{0}},
	{"example.tex", []byte("\n\n\n\n\n\n\n\n\n"), 9, []int{0, 1, 2, 3, 4, 5, 6, 7, 8}},
	{"math.tex", nil, 100, []int{0, 5, 10, 20, 30, 70, 71, 72, 80, 85, 90, 99}},
	{"long.tex", nil, 777, []int{0, 80, 100, 120, 130, 180, 267, 455, 500, 567, 620}},
	{"simple.tex", []byte("\\section{Introduction}\n\n\\begin{document}"), 40, []int{0, 23, 24}},
	{"command.tex", []byte("\\section{Introduction}\n\n\\begin{document}\n"), 41, []int{0, 23, 24}},
	{"complex.tex", []byte("\\section{Introduction}\n\n\\begin{document}\n "), 42, []int{0, 23, 24, 35}},
}

func linecol(lines []int, offs int) (int, int) {
	prevLineOffs := 0
	for line, lineOffs := range lines {
		if offs < lineOffs {
			return line, offs - prevLineOffs + 1
		}
		prevLineOffs = lineOffs
	}
	return len(lines), offs - prevLineOffs + 1
}

func verifyPositions(t *testing.T, fset *FileSet, f *File, lines []int) {
	for offs := range f.Size() {
		p := Pos(f.Base() + offs)
		line, col := linecol(lines, offs)
		msg := fmt.Sprintf("%s (offs = %d, p = %d)", f.Name(), offs, p)
		checkPos(t, msg, fset.Position(p), Position{f.Name(), offs, line, col})
	}
}

func TestPositions(t *testing.T) {
	const delta = 7 // a non-zero base offset increment
	fset := NewFileSet()
	for _, test := range tests {
		// verify consistency of test case
		if test.source != nil && len(test.source) != test.size {
			t.Errorf("%s: inconsistent test case: got file size %d; want %d", test.filename, len(test.source), test.size)
		}

		// add file and verify name and size
		f := fset.AddFile(test.filename, fset.Base()+delta, test.size)
		if f.Name() != test.filename {
			t.Errorf("got filename %q; want %q", f.Name(), test.filename)
		}
		if f.Size() != test.size {
			t.Errorf("%s: got file size %d; want %d", f.Name(), f.Size(), test.size)
		}
		if fset.File(Pos(f.Base())) != f {
			t.Errorf("%s: f.Base() was not found in f", f.Name())
		}

		// add lines individually and verify all positions
		for i, offset := range test.lines {
			f.AddLine(offset)
			if f.LineCount() != i+1 {
				t.Errorf("%s, AddLine: got line count %d; want %d", f.Name(), f.LineCount(), i+1)
			}
			// adding the same offset again should be ignored
			f.AddLine(offset)
			if f.LineCount() != i+1 {
				t.Errorf("%s, AddLine: got unchanged line count %d; want %d", f.Name(), f.LineCount(), i+1)
			}
			verifyPositions(t, fset, f, test.lines[0:i+1])
		}
	}
}

func TestFiles(t *testing.T) {
	fset := NewFileSet()
	files := make([]*File, len(tests))
	for i, test := range tests {
		base := fset.Base()
		files[i] = fset.AddFile(test.filename, base, test.size)
	}

	// Test file lookup by position
	for i, f := range files {
		if i > 0 {
			// Check that position at the end of the previous file
			// doesn't return this file
			prevFilePos := Pos(f.Base() - 1)
			prevFile := files[i-1]
			got := fset.File(prevFilePos)
			if got == f {
				t.Errorf("Position %d should belong to file %q (base=%d, size=%d), but got file %q (base=%d, size=%d)",
					prevFilePos, prevFile.Name(), prevFile.Base(), prevFile.Size(),
					f.Name(), f.Base(), f.Size())
			}
			if got != prevFile {
				t.Errorf("Position %d should belong to file %q (base=%d, size=%d), but got %v",
					prevFilePos, prevFile.Name(), prevFile.Base(), prevFile.Size(),
					got)
			}
		}

		// Check position at the start of this file
		startPos := Pos(f.Base())
		got := fset.File(startPos)
		if got != f {
			t.Errorf("Position %d at start of file %q (base=%d) returned wrong file: %v",
				startPos, f.Name(), f.Base(),
				got)
		}

		// Check position in the middle of this file
		midPos := Pos(f.Base() + f.Size()/2)
		got = fset.File(midPos)
		if got != f {
			t.Errorf("Position %d in middle of file %q (base=%d, size=%d) returned wrong file: %v",
				midPos, f.Name(), f.Base(), f.Size(),
				got)
		}

		// Check position at the end of this file
		endPos := Pos(f.Base() + f.Size())
		got = fset.File(endPos)
		if got != f {
			t.Errorf("Position %d at end of file %q (base=%d, size=%d) returned wrong file: %v",
				endPos, f.Name(), f.Base(), f.Size(),
				got)
		}
	}
}

// FileSet.File should return nil if Pos is past the end of the FileSet.
func TestFileSetPastEnd(t *testing.T) {
	fset := NewFileSet()
	for _, test := range tests {
		fset.AddFile(test.filename, fset.Base(), test.size)
	}
	if f := fset.File(Pos(fset.Base())); f != nil {
		t.Errorf("got %v, want nil", f)
	}
}

// Test that concurrent use of FileSet.Position does not trigger a
// race in the FileSet position cache.
func TestFileSetRacePosition(t *testing.T) {
	fset := NewFileSet()
	for i := range 100 {
		fset.AddFile(fmt.Sprintf("tex-file-%d.tex", i), fset.Base(), 1031)
	}
	max := int32(fset.Base())
	var stop sync.WaitGroup
	r := rand.New(rand.NewSource(7))
	for range 2 {
		r := rand.New(rand.NewSource(r.Int63()))
		stop.Add(1)
		go func() {
			for range 1000 {
				fset.Position(Pos(r.Int31n(max)))
			}
			stop.Done()
		}()
	}
	stop.Wait()
}

// Test that concurrent use of File.AddLine and FileSet.Position
// does not trigger a race in the FileSet position cache.
func TestFileSetRaceAddLinePosition(t *testing.T) {
	const N = 1000
	var (
		fset = NewFileSet()
		file = fset.AddFile("concurrent.tex", fset.Base(), N)
		ch   = make(chan int, 2)
	)

	go func() {
		for i := range N {
			file.AddLine(i)
		}
		ch <- 1
	}()

	go func() {
		pos := Pos(file.Base())
		for range N {
			fset.Position(pos)
		}
		ch <- 1
	}()

	<-ch
	<-ch
}

func TestLatexPositions(t *testing.T) {
	src := `\documentclass{article}
\begin{document}
\section{Test}
Text here.
\end{document}`

	fset := NewFileSet()
	basePos := fset.Base()
	f := fset.AddFile("example.tex", basePos, len(src))

	// Line start offsets
	lineOffsets := []int{0, 24, 41, 56, 67}
	for _, offset := range lineOffsets {
		f.AddLine(offset)
	}

	cases := []struct {
		offset int
		pos    Pos
		line   int
		col    int
		char   byte
	}{
		// Line 1: \documentclass{article}
		{0, Pos(basePos), 1, 1, '\\'},        // \
		{1, Pos(basePos + 1), 1, 2, 'd'},     // d
		{2, Pos(basePos + 2), 1, 3, 'o'},     // o
		{3, Pos(basePos + 3), 1, 4, 'c'},     // c
		{4, Pos(basePos + 4), 1, 5, 'u'},     // u
		{5, Pos(basePos + 5), 1, 6, 'm'},     // m
		{6, Pos(basePos + 6), 1, 7, 'e'},     // e
		{7, Pos(basePos + 7), 1, 8, 'n'},     // n
		{8, Pos(basePos + 8), 1, 9, 't'},     // t
		{9, Pos(basePos + 9), 1, 10, 'c'},    // c
		{10, Pos(basePos + 10), 1, 11, 'l'},  // l
		{11, Pos(basePos + 11), 1, 12, 'a'},  // a
		{12, Pos(basePos + 12), 1, 13, 's'},  // s
		{13, Pos(basePos + 13), 1, 14, 's'},  // s
		{14, Pos(basePos + 14), 1, 15, '{'},  // {
		{15, Pos(basePos + 15), 1, 16, 'a'},  // a
		{16, Pos(basePos + 16), 1, 17, 'r'},  // r
		{17, Pos(basePos + 17), 1, 18, 't'},  // t
		{18, Pos(basePos + 18), 1, 19, 'i'},  // i
		{19, Pos(basePos + 19), 1, 20, 'c'},  // c
		{20, Pos(basePos + 20), 1, 21, 'l'},  // l
		{21, Pos(basePos + 21), 1, 22, 'e'},  // e
		{22, Pos(basePos + 22), 1, 23, '}'},  // }
		{23, Pos(basePos + 23), 1, 24, '\n'}, // newline
		// Line 2: \begin{document}
		{24, Pos(basePos + 24), 2, 1, '\\'},  // \
		{25, Pos(basePos + 25), 2, 2, 'b'},   // b
		{26, Pos(basePos + 26), 2, 3, 'e'},   // e
		{27, Pos(basePos + 27), 2, 4, 'g'},   // g
		{28, Pos(basePos + 28), 2, 5, 'i'},   // i
		{29, Pos(basePos + 29), 2, 6, 'n'},   // n
		{30, Pos(basePos + 30), 2, 7, '{'},   // {
		{31, Pos(basePos + 31), 2, 8, 'd'},   // d
		{32, Pos(basePos + 32), 2, 9, 'o'},   // o
		{33, Pos(basePos + 33), 2, 10, 'c'},  // c
		{34, Pos(basePos + 34), 2, 11, 'u'},  // u
		{35, Pos(basePos + 35), 2, 12, 'm'},  // m
		{36, Pos(basePos + 36), 2, 13, 'e'},  // e
		{37, Pos(basePos + 37), 2, 14, 'n'},  // n
		{38, Pos(basePos + 38), 2, 15, 't'},  // t
		{39, Pos(basePos + 39), 2, 16, '}'},  // }
		{40, Pos(basePos + 40), 2, 17, '\n'}, // newline
		// Line 3: \section{Test}
		{41, Pos(basePos + 41), 3, 1, '\\'},  // \
		{42, Pos(basePos + 42), 3, 2, 's'},   // s
		{43, Pos(basePos + 43), 3, 3, 'e'},   // e
		{44, Pos(basePos + 44), 3, 4, 'c'},   // c
		{45, Pos(basePos + 45), 3, 5, 't'},   // t
		{46, Pos(basePos + 46), 3, 6, 'i'},   // i
		{47, Pos(basePos + 47), 3, 7, 'o'},   // o
		{48, Pos(basePos + 48), 3, 8, 'n'},   // n
		{49, Pos(basePos + 49), 3, 9, '{'},   // {
		{50, Pos(basePos + 50), 3, 10, 'T'},  // T
		{51, Pos(basePos + 51), 3, 11, 'e'},  // e
		{52, Pos(basePos + 52), 3, 12, 's'},  // s
		{53, Pos(basePos + 53), 3, 13, 't'},  // t
		{54, Pos(basePos + 54), 3, 14, '}'},  // }
		{55, Pos(basePos + 55), 3, 15, '\n'}, // newline
		// Line 4: Text here.
		{56, Pos(basePos + 56), 4, 1, 'T'},   // T
		{57, Pos(basePos + 57), 4, 2, 'e'},   // e
		{58, Pos(basePos + 58), 4, 3, 'x'},   // x
		{59, Pos(basePos + 59), 4, 4, 't'},   // t
		{60, Pos(basePos + 60), 4, 5, ' '},   // space
		{61, Pos(basePos + 61), 4, 6, 'h'},   // h
		{62, Pos(basePos + 62), 4, 7, 'e'},   // e
		{63, Pos(basePos + 63), 4, 8, 'r'},   // r
		{64, Pos(basePos + 64), 4, 9, 'e'},   // e
		{65, Pos(basePos + 65), 4, 10, '.'},  // .
		{66, Pos(basePos + 66), 4, 11, '\n'}, // newline
		// Line 5: \end{document}
		{67, Pos(basePos + 67), 5, 1, '\\'}, // \
		{68, Pos(basePos + 68), 5, 2, 'e'},  // e
		{69, Pos(basePos + 69), 5, 3, 'n'},  // n
		{70, Pos(basePos + 70), 5, 4, 'd'},  // d
		{71, Pos(basePos + 71), 5, 5, '{'},  // {
		{72, Pos(basePos + 72), 5, 6, 'd'},  // d
		{73, Pos(basePos + 73), 5, 7, 'o'},  // o
		{74, Pos(basePos + 74), 5, 8, 'c'},  // c
		{75, Pos(basePos + 75), 5, 9, 'u'},  // u
		{76, Pos(basePos + 76), 5, 10, 'm'}, // m
		{77, Pos(basePos + 77), 5, 11, 'e'}, // e
		{78, Pos(basePos + 78), 5, 12, 'n'}, // n
		{79, Pos(basePos + 79), 5, 13, 't'}, // t
		{80, Pos(basePos + 80), 5, 14, '}'}, // }
	}

	for _, c := range cases {
		// Test 1: Verify Pos to Position conversion
		pos := fset.Position(c.pos)
		expectedPos := Position{"example.tex", c.offset, c.line, c.col}
		checkPos(t, fmt.Sprintf("Pos %d", c.pos), pos, expectedPos)

		// Test 2: Verify File lookup by Pos
		file := fset.File(c.pos)
		if file != f {
			t.Errorf("Pos %d: got file %v, want %v", c.pos, file, f)
		}

		// Test 3: Verify file.Position(Pos) works too
		filePos := f.Position(c.pos)
		checkPos(t, fmt.Sprintf("file.Position(%d)", c.pos), filePos, expectedPos)

		// Test 4: Verify character at position
		if c.offset < len(src) && src[c.offset] != c.char {
			t.Errorf("offset %d (Pos %d): got %q, want %q",
				c.offset, c.pos, src[c.offset], c.char)
		}
	}

	// Test 5: Verify invalid Pos values
	if fset.File(NoPos) != nil {
		t.Errorf("Expected nil file for NoPos, got %v", fset.File(NoPos))
	}

	// Test 6: Verify out-of-bounds Pos values
	outOfBoundsPos := Pos(basePos + len(src) + 10)
	if fset.File(outOfBoundsPos) != nil {
		t.Errorf("Expected nil file for out-of-bounds Pos %d, got %v",
			outOfBoundsPos, fset.File(outOfBoundsPos))
	}
}
