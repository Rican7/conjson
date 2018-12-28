package transformer

import (
	"bytes"
	"regexp"
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
	keyMatchRegex             = regexp.MustCompile(`\"(\w+)\":`)
	camelCaseWordBarrierRegex = regexp.MustCompile(`([a-z])([A-Z])`)
	snakeCaseWordBarrierRegex = regexp.MustCompile(`_(\w)`)
)

// ConventionalKeys TODO
func ConventionalKeys() Transformer {
	return func(data []byte, direction Direction) []byte {
		return keyMatchRegex.ReplaceAllFunc(
			data,
			func(match []byte) []byte {
				if Unmarshal == direction {
					return snakeCaseWordBarrierRegex.ReplaceAllFunc(match, func(match []byte) []byte {
						// Remove the first byte (the underscore) from the match, and capitalize the rest
						return bytes.ToUpper(match[1:])
					})
				}

				return bytes.ToLower(
					camelCaseWordBarrierRegex.ReplaceAll(
						match,
						[]byte("${1}_${2}"),
					),
				)
			},
		)
	}
}
