package mediatype

// Code in this helpers file is borrowed from the google "headers" library

// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

import "strings"

// Octet types from RFC 2616.
var octetTypes [256]octetType

type octetType byte

const (
	isToken octetType = 1 << iota
	isSpace
)

func init() {
	// OCTET      = <any 8-bit sequence of data>
	// CHAR       = <any US-ASCII character (octets 0 - 127)>
	// CTL        = <any US-ASCII control character (octets 0 - 31) and DEL (127)>
	// CR         = <US-ASCII CR, carriage return (13)>
	// LF         = <US-ASCII LF, linefeed (10)>
	// SP         = <US-ASCII SP, space (32)>
	// HT         = <US-ASCII HT, horizontal-tab (9)>
	// <">        = <US-ASCII double-quote mark (34)>
	// CRLF       = CR LF
	// LWS        = [CRLF] 1*( SP | HT )
	// TEXT       = <any OCTET except CTLs, but including LWS>
	// separators = "(" | ")" | "<" | ">" | "@" | "," | ";" | ":" | "\" | <">
	//              | "/" | "[" | "]" | "?" | "=" | "{" | "}" | SP | HT
	// token      = 1*<any CHAR except CTLs or separators>
	// qdtext     = <any TEXT except <">>

	for c := 0; c < 256; c++ {
		var t octetType
		isCtl := c <= 31 || c == 127
		isChar := 0 <= c && c <= 127
		isSeparator := strings.IndexRune(" \t\"(),/:;<=>?@[]\\{}", rune(c)) >= 0
		if strings.IndexRune(" \t\r\n", rune(c)) >= 0 {
			t |= isSpace
		}
		if isChar && !isCtl && !isSeparator {
			t |= isToken
		}
		octetTypes[c] = t
	}
}

// ParseList parses a comma separated list of values. Commas are ignored in
// quoted strings. Quoted values are not unescaped or unquoted. Whitespace is
// trimmed.
func parseList(s string) []string {
	var result []string
	begin := 0
	end := 0
	escape := false
	quote := false
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch {
		case escape:
			escape = false
			end = i + 1
		case quote:
			switch b {
			case '\\':
				escape = true
			case '"':
				quote = false
			}
			end = i + 1
		case b == '"':
			quote = true
			end = i + 1
		case octetTypes[b]&isSpace != 0:
			if begin == end {
				begin = i + 1
				end = begin
			}
		case b == ',':
			if begin < end {
				result = append(result, s[begin:end])
			}
			begin = i + 1
			end = begin
		default:
			end = i + 1
		}
	}
	if begin < end {
		result = append(result, s[begin:end])
	}
	return result
}

// parseValueAndParams parses a comma separated list of values with optional
// semicolon separated name-value pairs. Content-Type and Content-Disposition
// headers are in this format.
func parseValueAndParams(s string) (value string, params map[string]string) {
	params = make(map[string]string)
	value, s = expectTokenSlash(s)
	if value == "" {
		return
	}
	value = strings.ToLower(value)
	s = skipSpace(s)
	for strings.HasPrefix(s, ";") {
		var pkey string
		pkey, s = expectToken(skipSpace(s[1:]))
		if pkey == "" {
			return
		}
		if !strings.HasPrefix(s, "=") {
			return
		}
		var pvalue string
		pvalue, s = expectTokenOrQuoted(s[1:])
		if pvalue == "" {
			return
		}
		pkey = strings.ToLower(pkey)
		params[pkey] = pvalue
		s = skipSpace(s)
	}
	return
}
func skipSpace(s string) (rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isSpace == 0 {
			break
		}
	}
	return s[i:]
}

func expectToken(s string) (token, rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isToken == 0 {
			break
		}
	}
	return s[:i], s[i:]
}

func expectTokenSlash(s string) (token, rest string) {
	i := 0
	for ; i < len(s); i++ {
		b := s[i]
		if (octetTypes[b]&isToken == 0) && b != '/' {
			break
		}
	}
	return s[:i], s[i:]
}

func expectTokenOrQuoted(s string) (value string, rest string) {
	if !strings.HasPrefix(s, "\"") {
		v, r := expectToken(s)
		return v, r
	}
	s = s[1:]
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			return s[:i], s[i+1:]
		case '\\':
			p := make([]byte, len(s)-1)
			j := copy(p, s[:i])
			escape := true
			for i = i + 1; i < len(s); i++ {
				b := s[i]
				switch {
				case escape:
					escape = false
					p[j] = b
					j++
				case b == '\\':
					escape = true
				case b == '"':
					return string(p[:j]), s[i+1:]
				default:
					p[j] = b
					j++
				}
			}
			return "", ""
		}
	}
	return "", ""
}
