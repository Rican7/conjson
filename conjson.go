package conjson

import (
	"encoding/json"

	"github.com/Rican7/conjson/transformer"
)

// Encoder TODO
type Encoder interface {
	Encode(interface{}) error
}

// Decoder TODO
type Decoder interface {
	Decode(interface{}) error
}

type marshaler struct {
	value        interface{}
	transformers []transformer.Transformer
}

type unmarshaler struct {
	value        interface{}
	transformers []transformer.Transformer
}

type encoder struct {
	inner        *json.Encoder
	transformers []transformer.Transformer
}

type decoder struct {
	inner        *json.Decoder
	transformers []transformer.Transformer
}

// NewMarshaler TODO
func NewMarshaler(value interface{}, transformers ...transformer.Transformer) json.Marshaler {
	return &marshaler{value, transformers}
}

// NewUnmarshaler TODO
func NewUnmarshaler(value interface{}, transformers ...transformer.Transformer) json.Unmarshaler {
	return &unmarshaler{value, transformers}
}

// NewEncoder TODO
func NewEncoder(inner *json.Encoder, transformers ...transformer.Transformer) Encoder {
	return &encoder{inner, transformers}
}

// NewDecoder TODO
func NewDecoder(inner *json.Decoder, transformers ...transformer.Transformer) Decoder {
	return &decoder{inner, transformers}
}

func (m *marshaler) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(m.value)

	if nil == err {
		marshalled = transform(marshalled, transformer.Marshal, m.transformers...)
	}

	return marshalled, err
}

func (um *unmarshaler) UnmarshalJSON(data []byte) error {
	data = transform(data, transformer.Unmarshal, um.transformers...)

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

func transform(data []byte, direction transformer.Direction, transformers ...transformer.Transformer) []byte {
	for _, transformer := range transformers {
		data = transformer(data, direction)
	}

	return data
}
