package conjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/Rican7/conjson/transform"
)

var (
	// Compile time interface assertion
	_ json.Marshaler   = (*marshaler)(nil)
	_ json.Unmarshaler = (*unmarshaler)(nil)
	_ Encoder          = (*encoder)(nil)
	_ Decoder          = (*decoder)(nil)
)

// Mocks

var noOpTransformer transform.Transformer = func(data []byte, direction transform.Direction) []byte {
	return data
}

type value bool
type errorMarshaler bool

func (v *value) MarshalJSON() ([]byte, error) {
	*v = true

	return json.Marshal((bool)(*v))
}

func (v *value) UnmarshalJSON(data []byte) error {
	*v = true

	return json.Unmarshal(data, (*bool)(v))
}

func (em *errorMarshaler) MarshalJSON() ([]byte, error) {
	if *em {
		return nil, errors.New("expected error")
	}

	return json.Marshal((bool)(*em))
}

func mockTransformer(timesRan *int, directionRan *transform.Direction) transform.Transformer {
	return func(data []byte, direction transform.Direction) []byte {
		*timesRan++
		*directionRan = direction

		return data
	}
}

// Tests

func TestNewMarshaler(t *testing.T) {
	if _, isJSONMarshaler := NewMarshaler(nil).(json.Marshaler); !isJSONMarshaler {
		t.Error("NewMarshaler didn't return a type compatible with json.Marshaler")
	}

	if _, isJSONMarshaler := NewMarshaler(nil, nil).(json.Marshaler); !isJSONMarshaler {
		t.Error("NewMarshaler didn't return a type compatible with json.Marshaler")
	}

	if _, isJSONMarshaler := NewMarshaler(nil, noOpTransformer).(json.Marshaler); !isJSONMarshaler {
		t.Error("NewMarshaler didn't return a type compatible with json.Marshaler")
	}

	if _, isJSONMarshaler := NewMarshaler(nil, noOpTransformer, noOpTransformer).(json.Marshaler); !isJSONMarshaler {
		t.Error("NewMarshaler didn't return a type compatible with json.Marshaler")
	}

	if _, isJSONMarshaler := NewMarshaler(time.Now()).(json.Marshaler); !isJSONMarshaler {
		t.Error("NewMarshaler didn't return a type compatible with json.Marshaler")
	}

	if _, isJSONMarshaler := NewMarshaler(time.Now(), nil).(json.Marshaler); !isJSONMarshaler {
		t.Error("NewMarshaler didn't return a type compatible with json.Marshaler")
	}

	if _, isJSONMarshaler := NewMarshaler(time.Now(), noOpTransformer).(json.Marshaler); !isJSONMarshaler {
		t.Error("NewMarshaler didn't return a type compatible with json.Marshaler")
	}

	if _, isJSONMarshaler := NewMarshaler(time.Now(), noOpTransformer, noOpTransformer).(json.Marshaler); !isJSONMarshaler {
		t.Error("NewMarshaler didn't return a type compatible with json.Marshaler")
	}
}

func TestNewUnmarshaler(t *testing.T) {
	if _, isJSONUnmarshaler := NewUnmarshaler(nil).(json.Unmarshaler); !isJSONUnmarshaler {
		t.Error("NewUnmarshaler didn't return a type compatible with json.Unmarshaler")
	}

	if _, isJSONUnmarshaler := NewUnmarshaler(nil, nil).(json.Unmarshaler); !isJSONUnmarshaler {
		t.Error("NewUnmarshaler didn't return a type compatible with json.Unmarshaler")
	}

	if _, isJSONUnmarshaler := NewUnmarshaler(nil, noOpTransformer).(json.Unmarshaler); !isJSONUnmarshaler {
		t.Error("NewUnmarshaler didn't return a type compatible with json.Unmarshaler")
	}

	if _, isJSONUnmarshaler := NewUnmarshaler(nil, noOpTransformer, noOpTransformer).(json.Unmarshaler); !isJSONUnmarshaler {
		t.Error("NewUnmarshaler didn't return a type compatible with json.Unmarshaler")
	}

	if _, isJSONUnmarshaler := NewUnmarshaler(time.Now()).(json.Unmarshaler); !isJSONUnmarshaler {
		t.Error("NewUnmarshaler didn't return a type compatible with json.Unmarshaler")
	}

	if _, isJSONUnmarshaler := NewUnmarshaler(time.Now(), nil).(json.Unmarshaler); !isJSONUnmarshaler {
		t.Error("NewUnmarshaler didn't return a type compatible with json.Unmarshaler")
	}

	if _, isJSONUnmarshaler := NewUnmarshaler(time.Now(), noOpTransformer).(json.Unmarshaler); !isJSONUnmarshaler {
		t.Error("NewUnmarshaler didn't return a type compatible with json.Unmarshaler")
	}

	if _, isJSONUnmarshaler := NewUnmarshaler(time.Now(), noOpTransformer, noOpTransformer).(json.Unmarshaler); !isJSONUnmarshaler {
		t.Error("NewUnmarshaler didn't return a type compatible with json.Unmarshaler")
	}
}

func TestNewEncoder(t *testing.T) {
	testJSONEncoder := json.NewEncoder(ioutil.Discard)

	if _, isEncoder := NewEncoder(nil).(Encoder); !isEncoder {
		t.Error("NewEncoder didn't return a type compatible with Encoder")
	}

	if _, isEncoder := NewEncoder(nil, nil).(Encoder); !isEncoder {
		t.Error("NewEncoder didn't return a type compatible with Encoder")
	}

	if _, isEncoder := NewEncoder(nil, noOpTransformer).(Encoder); !isEncoder {
		t.Error("NewEncoder didn't return a type compatible with Encoder")
	}

	if _, isEncoder := NewEncoder(nil, noOpTransformer, noOpTransformer).(Encoder); !isEncoder {
		t.Error("NewEncoder didn't return a type compatible with Encoder")
	}

	if _, isEncoder := NewEncoder(testJSONEncoder).(Encoder); !isEncoder {
		t.Error("NewEncoder didn't return a type compatible with Encoder")
	}

	if _, isEncoder := NewEncoder(testJSONEncoder, nil).(Encoder); !isEncoder {
		t.Error("NewEncoder didn't return a type compatible with Encoder")
	}

	if _, isEncoder := NewEncoder(testJSONEncoder, noOpTransformer).(Encoder); !isEncoder {
		t.Error("NewEncoder didn't return a type compatible with Encoder")
	}

	if _, isEncoder := NewEncoder(testJSONEncoder, noOpTransformer, noOpTransformer).(Encoder); !isEncoder {
		t.Error("NewEncoder didn't return a type compatible with Encoder")
	}
}

func TestNewDecoder(t *testing.T) {
	testJSONDecoder := json.NewDecoder(&bytes.Buffer{})

	if _, isDecoder := NewDecoder(nil).(Decoder); !isDecoder {
		t.Error("NewDecoder didn't return a type compatible with Decoder")
	}

	if _, isDecoder := NewDecoder(nil, nil).(Decoder); !isDecoder {
		t.Error("NewDecoder didn't return a type compatible with Decoder")
	}

	if _, isDecoder := NewDecoder(nil, noOpTransformer).(Decoder); !isDecoder {
		t.Error("NewDecoder didn't return a type compatible with Decoder")
	}

	if _, isDecoder := NewDecoder(nil, noOpTransformer, noOpTransformer).(Decoder); !isDecoder {
		t.Error("NewDecoder didn't return a type compatible with Decoder")
	}

	if _, isDecoder := NewDecoder(testJSONDecoder).(Decoder); !isDecoder {
		t.Error("NewDecoder didn't return a type compatible with Decoder")
	}

	if _, isDecoder := NewDecoder(testJSONDecoder, nil).(Decoder); !isDecoder {
		t.Error("NewDecoder didn't return a type compatible with Decoder")
	}

	if _, isDecoder := NewDecoder(testJSONDecoder, noOpTransformer).(Decoder); !isDecoder {
		t.Error("NewDecoder didn't return a type compatible with Decoder")
	}

	if _, isDecoder := NewDecoder(testJSONDecoder, noOpTransformer, noOpTransformer).(Decoder); !isDecoder {
		t.Error("NewDecoder didn't return a type compatible with Decoder")
	}
}

func TestMarshaler_MarshalJSON(t *testing.T) {
	var val value

	if _, err := NewMarshaler(nil).MarshalJSON(); nil != err {
		t.Errorf("Unexpected error (%T) %q", err, err)
	}

	val = false
	if _, err := NewMarshaler(&val).MarshalJSON(); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Marshal wasn't called on the inner value")
		}
	}

	val, timesRan, directionRan := false, 0, transform.Unmarshal
	if _, err := NewMarshaler(&val, mockTransformer(&timesRan, &directionRan)).MarshalJSON(); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Marshal wasn't called on the inner value")
		}

		if 1 != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `1`", timesRan)
		}

		if transform.Marshal != directionRan {
			t.Error("directionRan isn't the expected transform.Marshal")
		}
	}

	val, timesRan, directionRan = false, 0, transform.Unmarshal
	transformers := []transform.Transformer{
		mockTransformer(&timesRan, &directionRan),
		mockTransformer(&timesRan, &directionRan),
	}
	if _, err := NewMarshaler(&val, transformers...).MarshalJSON(); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Marshal wasn't called on the inner value")
		}

		if len(transformers) != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `%d`", timesRan, len(transformers))
		}

		if transform.Marshal != directionRan {
			t.Error("directionRan isn't the expected transform.Marshal")
		}
	}

	shouldError, timesRan := errorMarshaler(true), 0
	if _, err := NewMarshaler(&shouldError, mockTransformer(&timesRan, &directionRan)).MarshalJSON(); true {
		if nil == err {
			t.Error("Expected error was nil")
		}

		if 0 != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `0`", timesRan)
		}
	}
}

func TestUnmarshaler_UnmarshalJSON(t *testing.T) {
	testJSONBytes := []byte(`true`)
	var val value

	val = false
	if err := NewUnmarshaler(&val).UnmarshalJSON(testJSONBytes); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Unmarshal wasn't called on the inner value")
		}
	}

	val, timesRan, directionRan := false, 0, transform.Marshal
	if err := NewUnmarshaler(&val, mockTransformer(&timesRan, &directionRan)).UnmarshalJSON(testJSONBytes); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Unmarshal wasn't called on the inner value")
		}

		if 1 != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `1`", timesRan)
		}

		if transform.Unmarshal != directionRan {
			t.Error("directionRan isn't the expected transform.Ununmarshal")
		}
	}

	val, timesRan, directionRan = false, 0, transform.Marshal
	transformers := []transform.Transformer{
		mockTransformer(&timesRan, &directionRan),
		mockTransformer(&timesRan, &directionRan),
	}
	if err := NewUnmarshaler(&val, transformers...).UnmarshalJSON(testJSONBytes); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Unmarshal wasn't called on the inner value")
		}

		if len(transformers) != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `%d`", timesRan, len(transformers))
		}

		if transform.Unmarshal != directionRan {
			t.Error("directionRan isn't the expected transform.Ununmarshal")
		}
	}

	timesRan, directionRan = 0, transform.Marshal
	if err := NewUnmarshaler(false, mockTransformer(&timesRan, &directionRan)).UnmarshalJSON(testJSONBytes); true {
		if nil == err {
			t.Error("Expected error was nil")
		}

		if 1 != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `1`", timesRan)
		}

		if transform.Unmarshal != directionRan {
			t.Error("directionRan isn't the expected transform.Ununmarshal")
		}
	}
}

func TestEncoder_Encode(t *testing.T) {
	var val value
	var buf bytes.Buffer

	if err := NewEncoder(json.NewEncoder(&buf)).Encode(nil); nil != err {
		t.Errorf("Unexpected error (%T) %q", err, err)
	}

	val, buf = false, bytes.Buffer{}
	if err := NewEncoder(json.NewEncoder(&buf)).Encode(&val); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Marshal wasn't called on the inner value")
		}
	}

	val, buf, timesRan, directionRan := false, bytes.Buffer{}, 0, transform.Unmarshal
	if err := NewEncoder(json.NewEncoder(&buf), mockTransformer(&timesRan, &directionRan)).Encode(&val); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Marshal wasn't called on the inner value")
		}

		if 1 != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `1`", timesRan)
		}

		if transform.Marshal != directionRan {
			t.Error("directionRan isn't the expected transform.Marshal")
		}
	}

	val, buf, timesRan, directionRan = false, bytes.Buffer{}, 0, transform.Unmarshal
	transformers := []transform.Transformer{
		mockTransformer(&timesRan, &directionRan),
		mockTransformer(&timesRan, &directionRan),
	}
	if err := NewEncoder(json.NewEncoder(&buf), transformers...).Encode(&val); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Marshal wasn't called on the inner value")
		}

		if len(transformers) != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `%d`", timesRan, len(transformers))
		}

		if transform.Marshal != directionRan {
			t.Error("directionRan isn't the expected transform.Marshal")
		}
	}

	shouldError, buf, timesRan := errorMarshaler(true), bytes.Buffer{}, 0
	if err := NewEncoder(json.NewEncoder(&buf), mockTransformer(&timesRan, &directionRan)).Encode(&shouldError); true {
		if nil == err {
			t.Error("Expected error was nil")
		}

		if 0 != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `0`", timesRan)
		}
	}
}

func TestDecoder_UnmarshalJSON(t *testing.T) {
	testJSONBytes := []byte(`true`)
	var val value
	var buf *bytes.Buffer

	val, buf = false, bytes.NewBuffer(testJSONBytes)
	if err := NewDecoder(json.NewDecoder(buf)).Decode(&val); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Unmarshal wasn't called on the inner value")
		}
	}

	val, buf, timesRan, directionRan := false, bytes.NewBuffer(testJSONBytes), 0, transform.Marshal
	if err := NewDecoder(json.NewDecoder(buf), mockTransformer(&timesRan, &directionRan)).Decode(&val); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Unmarshal wasn't called on the inner value")
		}

		if 1 != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `1`", timesRan)
		}

		if transform.Unmarshal != directionRan {
			t.Error("directionRan isn't the expected transform.Ununmarshal")
		}
	}

	val, buf, timesRan, directionRan = false, bytes.NewBuffer(testJSONBytes), 0, transform.Marshal
	transformers := []transform.Transformer{
		mockTransformer(&timesRan, &directionRan),
		mockTransformer(&timesRan, &directionRan),
	}
	if err := NewDecoder(json.NewDecoder(buf), transformers...).Decode(&val); true {
		if nil != err {
			t.Errorf("Unexpected error (%T) %q", err, err)
		}

		if !val {
			t.Error("val is false, so json.Unmarshal wasn't called on the inner value")
		}

		if len(transformers) != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `%d`", timesRan, len(transformers))
		}

		if transform.Unmarshal != directionRan {
			t.Error("directionRan isn't the expected transform.Ununmarshal")
		}
	}

	buf, timesRan, directionRan = bytes.NewBuffer(testJSONBytes), 0, transform.Marshal
	if err := NewDecoder(json.NewDecoder(buf), mockTransformer(&timesRan, &directionRan)).Decode(false); true {
		if nil == err {
			t.Error("Expected error was nil")
		}

		if 1 != timesRan {
			t.Errorf("timesRan was `%d`, when expected to be `1`", timesRan)
		}

		if transform.Unmarshal != directionRan {
			t.Error("directionRan isn't the expected transform.Ununmarshal")
		}
	}
}
