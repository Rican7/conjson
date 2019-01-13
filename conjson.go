package conjson

import (
	"encoding/json"

	"github.com/Rican7/conjson/transform"
)

// Encoder is an interface defining a simple JSON encoder, with an interface
// compatible with `encoding/json.Encoder`.
type Encoder interface {
	// Encode takes a value and encodes the JSON representation of that value to
	// the underyling/inner encoder.
	//
	// See the documentation for both `encoding/json.Encoder` and
	// `encoding/json.Marshal` for more details about the workings of the
	// underlying encoder.
	Encode(interface{}) error
}

// Decoder is an interface defining a simple JSON decoder, with an interface
// compatible with `encoding/json.Decoder`.
type Decoder interface {
	// Decode reads the next JSON represented value from the underyling/inner
	// decoder and stores the decoded result in the pointed to passed value.
	//
	// See the documentation for both `encoding/json.Decoder` and
	// `encoding/json.Unmarshal` for more details about the workings of the
	// underlying decoder.
	Decode(interface{}) error
}

// marshaler is a structure that wraps a value and a list of transformers to
// enable JSON marshaling with output transformations.
type marshaler struct {
	value        interface{}
	transformers []transform.Transformer
}

// unmarshaler is a structure that wraps a value and a list of transformers to
// enable JSON unmarshaling with input transformations.
type unmarshaler struct {
	value        interface{}
	transformers []transform.Transformer
}

// encoder is a structure that wraps an `encoding/json.Encoder` and a list
// of transformers to enable JSON encoding with output transformations.
type encoder struct {
	inner        *json.Encoder
	transformers []transform.Transformer
}

// decoder is a structure that wraps an `encoding/json.Decoder` and a list
// of transformers to enable JSON decoding with input transformations.
type decoder struct {
	inner        *json.Decoder
	transformers []transform.Transformer
}

// NewMarshaler takes a value and a variable number of `transform.Transformer`s
// and returns an `encoding/json.Marshaler` that runs the given transformers
// upon JSON marshaling.
//
// See the documentation for both `encoding/json.Marshaler` and
// `encoding/json.Marshal` for more details about JSON marshaling.
func NewMarshaler(value interface{}, transformers ...transform.Transformer) json.Marshaler {
	return &marshaler{value, transformers}
}

// NewUnmarshaler takes a pointer value and a variable number of
// `transform.Transformer`s and returns an `encoding/json.Unmarshaler` that runs
// the given transformers upon JSON unmarshaling.
//
// See the documentation for both `encoding/json.Unmarshaler` and
// `encoding/json.Unmarshal` for more details about JSON unmarshaling.
func NewUnmarshaler(value interface{}, transformers ...transform.Transformer) json.Unmarshaler {
	return &unmarshaler{value, transformers}
}

// NewEncoder takes an `encoding/json.Encoder` and a variable number of
// `transform.Transformer`s and returns an `Encoder` that runs the given
// transformers upon JSON encoding.
//
// See the documentation for both `encoding/json.Encoder` and
// `encoding/json.Marshal` for more details about the passed inner encoder.
func NewEncoder(inner *json.Encoder, transformers ...transform.Transformer) Encoder {
	return &encoder{inner, transformers}
}

// NewDecoder takes an `encoding/json.Decoder` and a variable number of
// `transform.Transformer`s and returns an `Decoder` that runs the given
// transformers upon JSON encoding.
//
// See the documentation for both `encoding/json.Decoder` and
// `encoding/json.Unmarshal` for more details about the passed inner decoder.
func NewDecoder(inner *json.Decoder, transformers ...transform.Transformer) Decoder {
	return &decoder{inner, transformers}
}

func (m *marshaler) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(m.value)

	if nil == err {
		marshalled = transform.Bytes(marshalled, transform.Marshal, m.transformers...)
	}

	return marshalled, err
}

func (um *unmarshaler) UnmarshalJSON(data []byte) error {
	data = transform.Bytes(data, transform.Unmarshal, um.transformers...)

	return json.Unmarshal(data, um.value)
}

func (e *encoder) Encode(value interface{}) error {
	return e.inner.Encode(
		NewMarshaler(value, e.transformers...),
	)
}

func (e *decoder) Decode(value interface{}) error {
	return e.inner.Decode(
		NewUnmarshaler(value, e.transformers...),
	)
}
