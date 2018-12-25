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
					return bytes.ToUpper(
						snakeCaseWordBarrierRegex.ReplaceAll(
							match,
							[]byte("${1}"),
						),
					)
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
