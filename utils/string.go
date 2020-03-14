package utils

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type buffer struct {
	r         []byte
	runeBytes [utf8.UTFMax]byte
}

func (b *buffer) write(r rune) {
	if r < utf8.RuneSelf {
		b.r = append(b.r, byte(r))
		return
	}
	n := utf8.EncodeRune(b.runeBytes[0:], r)
	b.r = append(b.r, b.runeBytes[0:n]...)
}

func (b *buffer) indent() {
	if len(b.r) > 0 {
		b.r = append(b.r, '_')
	}
}

// ConvertCamelToUnderscore convert camel case string to string with underscore
func ConvertCamelToUnderscore(s string) string {
	b := buffer{
		r: make([]byte, 0, len(s)),
	}
	var m rune
	var w bool
	for _, ch := range s {
		if unicode.IsUpper(ch) {
			if m != 0 {
				if !w {
					b.indent()
					w = true
				}
				b.write(m)
			}
			m = unicode.ToLower(ch)
		} else {
			if m != 0 {
				b.indent()
				b.write(m)
				m = 0
				w = false
			}
			b.write(ch)
		}
	}
	if m != 0 {
		if !w {
			b.indent()
		}
		b.write(m)
	}
	return string(b.r)
}

// SplitParamsToMap for spliting array of string to map
func SplitParamsToMap(param []string) map[string]string {
	res := make(map[string]string)

	for _, v := range param {
		splt := strings.Split(v, "=")
		if len(splt) < 2 {
			continue
		}
		res[splt[0]] = splt[1]
	}

	return res
}

// UcFirst for first to Upper
func UcFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// ConvertUnderscoreToCamel for convert it
func ConvertUnderscoreToCamel(s string) string {
	splitV := strings.Split(s, "_")
	result := ""

	for _, v := range splitV {
		result += UcFirst(v)
	}

	return result
}

// ToLowerFirst for lower all character first
func ToLowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
