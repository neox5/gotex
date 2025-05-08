package scanner

import (
	"unicode"
	"unicode/utf8"
)

func isSpaceChar(ch rune) bool   { return ch == ' ' || ch == '\t' }
func lower(ch rune) rune         { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isDigit(ch rune) bool       { return '0' <= ch && ch <= '9' }
func isCommandChar(ch rune) bool { return 'a' <= lower(ch) && lower(ch) <= 'z' }

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}
