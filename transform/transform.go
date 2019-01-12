package transform

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

// Bytes TODO
func Bytes(data []byte, direction Direction, transformers ...Transformer) []byte {
	// Make a copy of the source data to make sure that transformers don't
	// modify the source, but instead return a copy of the data as the
	// interface intends
	transformed := make([]byte, len(data))
	copy(transformed, data)

	for _, transformer := range transformers {
		transformed = transformer(transformed, direction)
	}

	return transformed
}

// OnlyForDirection TODO
//
// TODO: Naming?
func OnlyForDirection(only Direction, transform Transformer) Transformer {
	return func(data []byte, direction Direction) []byte {
		if direction == only {
			return transform(data, direction)
		}

		return data
	}
}

// AlwaysAsDirection TODO
//
// TODO: Naming?
func AlwaysAsDirection(always Direction, transform Transformer) Transformer {
	return func(data []byte, direction Direction) []byte {
		return transform(data, always)
	}
}

// ConventionalKeys TODO
func ConventionalKeys() Transformer {
	return func(data []byte, direction Direction) []byte {
		return replaceKeys(
			data,
			func(key []byte) []byte {
				if Unmarshal == direction {
					return snakeCaseToCamelCaseWordBarrier(key)
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

// CamelCaseKeys TODO
func CamelCaseKeys() Transformer {
	return func(data []byte, direction Direction) []byte {
		return replaceKeys(
			data,
			func(key []byte) []byte {
				// Translate hyphenated keys to underscored ("snake_case")
				key = bytes.Replace(key, []byte("-"), []byte("_"), -1)

				key = snakeCaseToCamelCaseWordBarrier(key)

				return append(bytes.ToLower(key[0:1]), key[1:]...)
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

func snakeCaseToCamelCaseWordBarrier(data []byte) []byte {
	return snakeCaseWordBarrierRegex.ReplaceAllFunc(data, func(match []byte) []byte {
		// Remove the underscore from the match, and capitalize the rest
		return append(match[0:1], bytes.ToUpper(match[2:])...)
	})
}
