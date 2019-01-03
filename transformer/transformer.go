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
	keyMatchRegex             = regexp.MustCompile(`\"(.+?)\"\s*?:`)
	camelCaseWordBarrierRegex = regexp.MustCompile(`(.)([A-Z])`)
	snakeCaseWordBarrierRegex = regexp.MustCompile(`(?:[^_])_(.)`)
)

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
