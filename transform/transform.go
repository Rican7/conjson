// Package transform provides mechanisms to transform a data source, with a
// specific focus on JSON data handling.
//
// Copyright © Trevor N. Suarez (Rican7)
package transform

import (
	"bytes"
	"regexp"
	"unicode"
)

// Direction flags the direction of marshaling/unmarshaling (encoding/decoding).
type Direction bool

// Transformer defines a function that transforms a source bytes of data, in a
// given direction, to a result of bytes.
type Transformer func([]byte, Direction) []byte

const (
	// Marshal defines the direction of marshaling or encoding.
	Marshal Direction = false

	// Unmarshal defines the direction of unmarshaling or decoding.
	Unmarshal Direction = true
)

var (
	keyMatchRegex                     = regexp.MustCompile(`\"([^\"]*?)\"\s*?:`)
	camelCaseWordBarrierRegex         = regexp.MustCompile(`([^A-Z])([A-Z])`)
	snakeCaseWordBarrierRegex         = regexp.MustCompile(`(?:[^_])_(.)`)
	repeatedUpperCaseWordBarrierRegex = regexp.MustCompile(`(?:[A-Z])([A-Z]+?)(?:[^A-Z]|$)`)
)

// String satisfies the fmt.Stringer interface to provide a human-readable name
// for a Direction value.
func (d Direction) String() string {
	if Marshal == d {
		return "Marshal"
	}

	return "Unmarshal"
}

// Bytes takes a source bytes of data, a Direction, and a variable number of
// transformers and returns a new byte slice with the result each transformer
// having run on a copy of the original passed data.
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

// OnlyForDirection takes a given direction and a Transformer and returns a
// new Transformer that only executes the given Transformer when the
// transformation direction matches the given direction.
func OnlyForDirection(only Direction, transform Transformer) Transformer {
	return func(data []byte, direction Direction) []byte {
		if direction == only {
			return transform(data, direction)
		}

		return data
	}
}

// AlwaysAsDirection takes a given direction and a Transformer and returns a
// new Transformer that always executes the given Transformer as if the
// transformation direction was the given direction.
func AlwaysAsDirection(always Direction, transform Transformer) Transformer {
	return func(data []byte, direction Direction) []byte {
		return transform(data, always)
	}
}

// ReverseDirection takes a Transformer and returns a new Transformer that
// executes the given Transformer as if the transformation direction was the
// oppopsite of the given direction.
func ReverseDirection(transform Transformer) Transformer {
	return func(data []byte, direction Direction) []byte {
		return transform(data, !direction)
	}
}

// ConventionalKeys returns a Transformer that converts every JSON "key" (or
// JSON object "name") in the transformed data set, depending on the
// transformation direction, based on common JSON data style conventions.
//
// For the "Marshal" direction, JSON keys are converted to `snake_case` style.
// For the "Unmarshal" direction, JSON keys are converted to `camelCase` style.
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

// CamelCaseKeys returns a Transformer that converts every JSON "key" (or JSON
// object "name") in the transformed data set to be in `camelCase` style.
//
// If the passed lowerRepeatedCaps param is `true`, then repeated capital
// letters (such as "URL" or "HTTP") will be converted to typical "Title" case
// (such as "Url" or "Http").
func CamelCaseKeys(lowerRepeatedCaps bool) Transformer {
	return func(data []byte, direction Direction) []byte {
		return replaceKeys(
			data,
			func(key []byte) []byte {
				// Translate hyphenated keys to underscored ("snake_case")
				key = bytes.Replace(key, []byte("-"), []byte("_"), -1)

				// Transform snake-case keys to camel-case keys
				key = snakeCaseToCamelCaseWordBarrier(key)

				if lowerRepeatedCaps {
					// Remove repeated upper-case letters
					key = repeatedUpperCaseWordBarrierRegex.ReplaceAllFunc(key, func(key []byte) []byte {
						// If the last letter in the repeated-upper-case find is lower-case,
						// then we have a successive "word"
						if unicode.IsLower(rune(key[len(key)-1])) {
							// Only lower-case the first "word"
							return append(
								key[0:1],
								append(
									bytes.ToLower(key[1:len(key)-2]),
									key[len(key)-2:]...,
								)...,
							)
						}

						return append(key[0:1], bytes.ToLower(key[1:])...)
					})
				}

				// Lower-case the first letter
				return append(bytes.ToLower(key[0:1]), key[1:]...)
			},
		)
	}
}

// ValidIdentifierKeys returns a Transformer that converts every JSON "key" (or
// JSON object "name") in the transformed data set to a key that's format
// matches the Go specification to be considered a valid "identifier". It does
// this conversion by simply stripping characters from the key/name that would
// otherwise make for an invalid Go identifier, according to specification.
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

// replaceKeys takes a source JSON data and a replaceFunc and runs the given
// `replaceFunc` only on every JSON "key" (or JSON object "name") found in the
// given source JSON data, properly rebuilding the surrounding JSON data/source
// so that the new/replaced "key" is in the resulting data.
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

// snakeCaseToCamelCaseWordBarrier takes a source JSON data and replaces all
// `snake_case` style JSON keys with `camelCase` style JSON keys, based on a
// "word-barrier" regular expression, and returns the resulting bytes.
func snakeCaseToCamelCaseWordBarrier(data []byte) []byte {
	return snakeCaseWordBarrierRegex.ReplaceAllFunc(data, func(match []byte) []byte {
		// Remove the underscore from the match, and capitalize the rest
		return append(match[0:1], bytes.ToUpper(match[2:])...)
	})
}
