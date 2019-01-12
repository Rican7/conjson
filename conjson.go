package conjson

import (
	"encoding/json"

	"github.com/Rican7/conjson/transform"
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
	transformers []transform.Transformer
}

type unmarshaler struct {
	value        interface{}
	transformers []transform.Transformer
}

type encoder struct {
	inner        *json.Encoder
	transformers []transform.Transformer
}

type decoder struct {
	inner        *json.Decoder
	transformers []transform.Transformer
}

// NewMarshaler TODO
func NewMarshaler(value interface{}, transformers ...transform.Transformer) json.Marshaler {
	return &marshaler{value, transformers}
}

// NewUnmarshaler TODO
func NewUnmarshaler(value interface{}, transformers ...transform.Transformer) json.Unmarshaler {
	return &unmarshaler{value, transformers}
}

// NewEncoder TODO
func NewEncoder(inner *json.Encoder, transformers ...transform.Transformer) Encoder {
	return &encoder{inner, transformers}
}

// NewDecoder TODO
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
