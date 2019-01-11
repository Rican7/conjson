package transformer

import (
	"bytes"
	"regexp"
	"unicode"
)

// Direction TODO
type Direction bool

// Transformer TODO
type Transformer func([]byte, Direction) []byte

const (
	// Marshal TODO
	Marshal Direction = false

	// Unmarshal TODO
	Unmarshal Direction = true
)

var (
	keyMatchRegex             = regexp.MustCompile(`\"(.+?)\"\s*?:`)
	camelCaseWordBarrierRegex = regexp.MustCompile(`(.)([A-Z])`)
	snakeCaseWordBarrierRegex = regexp.MustCompile(`(?:[^_])_(.)`)
)

// String satisfies the fmt.Stringer interface to provide a human-readable name
// for a Direction value
func (d Direction) String() string {
	if Marshal == d {
		return "Marshal"
	}

	return "Unmarshal"
}

// ConventionalKeys TODO
func ConventionalKeys() Transformer {
	return func(data []byte, direction Direction) []byte {
		return replaceKeys(
			data,
			func(key []byte) []byte {
				if Unmarshal == direction {
					return snakeCaseWordBarrierRegex.ReplaceAllFunc(key, func(match []byte) []byte {
						// Remove the underscore from the match, and capitalize the rest
						return append(match[0:1], bytes.ToUpper(match[2:])...)
					})
				}

				return bytes.ToLower(
					camelCaseWordBarrierRegex.ReplaceAll(
						key,
						[]byte("${1}_${2}"),
					),
				)
			},
		)
	}
}

// ValidIdentifierKeys TODO
//
// https://golang.org/ref/spec#Identifiers
func ValidIdentifierKeys() Transformer {
	return func(data []byte, direction Direction) []byte {
		return replaceKeys(
			data,
			func(key []byte) []byte {
				key = bytes.TrimLeftFunc(key, func(r rune) bool {
					return !unicode.IsLetter(r)
				})

				fields := bytes.FieldsFunc(key, func(r rune) bool {
					return !unicode.In(r, unicode.Letter, unicode.Digit)
				})

				return bytes.Join(fields, nil)
			},
		)
	}
}

func replaceKeys(data []byte, replaceFunc func(key []byte) []byte) []byte {
	return keyMatchRegex.ReplaceAllFunc(
		data,
		func(match []byte) []byte {
			key := replaceFunc(keyMatchRegex.FindSubmatch(match)[1])
			indexes := keyMatchRegex.FindSubmatchIndex(match)

			return append([]byte{match[indexes[0]]}, append(key, match[indexes[3]:indexes[1]]...)...)
		},
	)
}
