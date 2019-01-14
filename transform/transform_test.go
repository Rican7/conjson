package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/parser"
	"testing"
)

func TestDirection_String(t *testing.T) {
	for _, testCase := range []struct {
		direction      Direction
		expectedString string
	}{
		{Marshal, "Marshal"},
		{Unmarshal, "Unmarshal"},
	} {
		if str := testCase.direction.String(); testCase.expectedString != str {
			t.Errorf("%q doesn't match expected %q", str, testCase.expectedString)
		}
	}
}

func TestBytes(t *testing.T) {
	testData := []byte("test-data")

	transformers := []Transformer{
		func(data []byte, direction Direction) []byte {
			return bytes.ToUpper(data)
		},
		func(data []byte, direction Direction) []byte {
			return bytes.Replace(data, []byte("-"), []byte(" "), -1)
		},
		func(data []byte, direction Direction) []byte {
			prefix := fmt.Sprintf("%s: ", direction.String())

			return append([]byte(prefix), data...)
		},
	}

	for _, testCase := range []struct {
		transformers   []Transformer
		direction      Direction
		expectedOutput []byte
	}{
		{transformers[0:1], Marshal, []byte("TEST-DATA")},
		{transformers[0:1], Unmarshal, []byte("TEST-DATA")},
		{transformers[1:2], Marshal, []byte("test data")},
		{transformers[1:2], Unmarshal, []byte("test data")},
		{transformers[2:3], Marshal, []byte("Marshal: test-data")},
		{transformers[2:3], Unmarshal, []byte("Unmarshal: test-data")},
		{transformers, Marshal, []byte("Marshal: TEST DATA")},
		{transformers, Unmarshal, []byte("Unmarshal: TEST DATA")},
	} {
		if output := Bytes(testData, testCase.direction, testCase.transformers...); !bytes.Equal(testCase.expectedOutput, output) {
			t.Errorf("%s output of %q doesn't match expected %q", testCase.direction, output, testCase.expectedOutput)
		}
	}
}

// TestBytes_AssertDoesntModifyPassedBytes tests that our passed bytes aren't
// modified in-place, but rather that we receive a copy
func TestBytes_AssertDoesntModifyPassedBytes(t *testing.T) {
	testData := []byte("test-data")
	const badByte = '!'

	transformerThatModifiesBytesInPlace := func(data []byte, direction Direction) []byte {
		data[0] = badByte

		return data
	}

	// If our output matches our input, then we know we've modified our passed data
	if output := Bytes(testData, Marshal, transformerThatModifiesBytesInPlace); bytes.Equal(testData, output) {
		t.Error("Input data was modified!")
	}
}

func TestOnlyForDirection(t *testing.T) {
	mockDataReturn := []byte("mock data")

	mockTransformer := func(data []byte, direction Direction) []byte {
		return mockDataReturn
	}

	marshalTrans := OnlyForDirection(Marshal, mockTransformer)
	unmarshalTrans := OnlyForDirection(Unmarshal, mockTransformer)

	for _, testCase := range []struct {
		trans        Transformer
		direction    Direction
		callExpected bool
	}{
		{marshalTrans, Marshal, true},
		{marshalTrans, Unmarshal, false},
		{unmarshalTrans, Marshal, false},
		{unmarshalTrans, Unmarshal, true},
	} {
		if output := testCase.trans([]byte(""), testCase.direction); testCase.callExpected && !bytes.Equal(output, mockDataReturn) {
			t.Errorf("%s output of %q doesn't match expected %q", testCase.direction, output, mockDataReturn)
		}
	}
}

func TestAlwaysAsDirection(t *testing.T) {
	mockTransformer := func(data []byte, direction Direction) []byte {
		return []byte(direction.String())
	}

	marshalTrans := AlwaysAsDirection(Marshal, mockTransformer)
	unmarshalTrans := AlwaysAsDirection(Unmarshal, mockTransformer)

	for _, testCase := range []struct {
		trans             Transformer
		expectedDirection Direction
		direction         Direction
	}{
		{marshalTrans, Marshal, Marshal},
		{marshalTrans, Marshal, Unmarshal},
		{unmarshalTrans, Unmarshal, Marshal},
		{unmarshalTrans, Unmarshal, Unmarshal},
	} {
		expectedDirectionBytes := []byte(testCase.expectedDirection.String())

		if output := testCase.trans([]byte(""), testCase.direction); !bytes.Equal(output, expectedDirectionBytes) {
			t.Errorf("%s output of %q doesn't match expected %q", testCase.direction, output, expectedDirectionBytes)
		}
	}
}

func TestReverseDirection(t *testing.T) {
	trans := ReverseDirection(
		func(data []byte, direction Direction) []byte {
			return []byte(direction.String())
		},
	)

	for _, testCase := range []struct {
		direction         Direction
		expectedDirection Direction
	}{
		{Marshal, Unmarshal},
		{Unmarshal, Marshal},
	} {
		expectedDirectionBytes := []byte(testCase.expectedDirection.String())

		if output := trans([]byte(""), testCase.direction); !bytes.Equal(output, expectedDirectionBytes) {
			t.Errorf("%s output of %q doesn't match expected %q", testCase.direction, output, expectedDirectionBytes)
		}
	}
}

func TestConventionalKeys(t *testing.T) {
	const snakeCaseJSON = `
	{
		"title": "Example Title",
		"description": "whatever",
		"created_at": "2015-11-17T20:43:31.0463576-05:00",
		"updated_at": "2018-12-24T13:21:15.7883416-07:00",
		"is_active": true,
		"image_url": "https://example.com/image.png",
		"time_to_live": 600,
		"$weird_key" : "with-colon-spacing-before-value",
		"plan9_key" : "just an example with a numeral",
		"__metadata_key": "You see this in OData 2.0",
		"__metadata": "You see this in OData 2.0"
	}
	`

	const camelCaseJSON = `
	{
		"title": "Example Title",
		"description": "whatever",
		"createdAt": "2015-11-17T20:43:31.0463576-05:00",
		"updatedAt": "2018-12-24T13:21:15.7883416-07:00",
		"isActive": true,
		"imageUrl": "https://example.com/image.png",
		"timeToLive": 600,
		"$weirdKey" : "with-colon-spacing-before-value",
		"plan9Key" : "just an example with a numeral",
		"__metadataKey": "You see this in OData 2.0",
		"__metadata": "You see this in OData 2.0"
	}
	`

	trans := ConventionalKeys()

	if output := trans([]byte(camelCaseJSON), Marshal); string(output) != snakeCaseJSON {
		t.Errorf("Marshal output of %s doesn't match expected %s", output, snakeCaseJSON)
	}

	if output := trans([]byte(snakeCaseJSON), Unmarshal); string(output) != camelCaseJSON {
		t.Errorf("Unmarshal output of %s doesn't match expected %s", output, camelCaseJSON)
	}
}

func TestCamelCaseKeys(t *testing.T) {
	const originalJSON = `
	{
		"title": "Example Title",
		"description": "whatever",
		"created_at": "2015-11-17T20:43:31.0463576-05:00",
		"updated_at": "2018-12-24T13:21:15.7883416-07:00",
		"is_active": true,
		"image_url": "https://example.com/image.png",
		"time_to_live": 600,
		"$weird_key" : "with-colon-spacing-before-value",
		"plan9_key" : "just an example with a numeral",
		"__metadata_key": "You see this in OData 2.0",
		"__metadata": "You see this in OData 2.0",
		"Capitalized": "a capitalized key",
		"a-hyphenated-key": "a hyphenated key",
		"camelCaseKey": "a camelCase key",
		"UpperCamelCaseKey": "an UpperCamelCase key",
		"a-strange_mixedAndConfused_CaseStyle": "an UpperCamelCase key",
		"aKeyWithRepeatedCapitalLettersLikeURL": "a key with repeated capitals"
	}
	`

	const camelCaseJSON = `
	{
		"title": "Example Title",
		"description": "whatever",
		"createdAt": "2015-11-17T20:43:31.0463576-05:00",
		"updatedAt": "2018-12-24T13:21:15.7883416-07:00",
		"isActive": true,
		"imageUrl": "https://example.com/image.png",
		"timeToLive": 600,
		"$weirdKey" : "with-colon-spacing-before-value",
		"plan9Key" : "just an example with a numeral",
		"__metadataKey": "You see this in OData 2.0",
		"__metadata": "You see this in OData 2.0",
		"capitalized": "a capitalized key",
		"aHyphenatedKey": "a hyphenated key",
		"camelCaseKey": "a camelCase key",
		"upperCamelCaseKey": "an UpperCamelCase key",
		"aStrangeMixedAndConfusedCaseStyle": "an UpperCamelCase key",
		"aKeyWithRepeatedCapitalLettersLikeUrl": "a key with repeated capitals"
	}
	`

	trans := CamelCaseKeys()

	for _, testCase := range []struct {
		jsonBytes string
		direction Direction
	}{
		{originalJSON, Marshal},
		{originalJSON, Unmarshal},
		{camelCaseJSON, Marshal},
		{camelCaseJSON, Unmarshal},
	} {
		if output := trans([]byte(testCase.jsonBytes), testCase.direction); string(output) != camelCaseJSON {
			t.Errorf("%s output of %s doesn't match expected %s", testCase.direction, output, camelCaseJSON)
		}
	}
}

func TestValidIdentifierKeys(t *testing.T) {
	const invalidKeyJSON = `
	{
		"title": "Example Title",
		"description": "whatever",
		"created_at": "2015-11-17T20:43:31.0463576-05:00",
		"updated_at": "2018-12-24T13:21:15.7883416-07:00",
		"is_active": true,
		"image_url": "https://example.com/image.png",
		"time_to_live": 600,
		"$weird_key" : "with-colon-spacing-before-value",
		"plan9_key" : "just an example with a numeral",
		"__metadata_key": "You see this in OData 2.0",
		"__metadata": "You see this in OData 2.0",
		"Capitalized": "a capitalized key",
		"a-hyphenated-key": "a hyphenated key",
		"3^%!@#*identifier-pls&": "Weird key that isn't a valid Go identifier"
	}
	`

	const validKeyJSON = `
	{
		"title": "Example Title",
		"description": "whatever",
		"createdat": "2015-11-17T20:43:31.0463576-05:00",
		"updatedat": "2018-12-24T13:21:15.7883416-07:00",
		"isactive": true,
		"imageurl": "https://example.com/image.png",
		"timetolive": 600,
		"weirdkey" : "with-colon-spacing-before-value",
		"plan9key" : "just an example with a numeral",
		"metadatakey": "You see this in OData 2.0",
		"metadata": "You see this in OData 2.0",
		"Capitalized": "a capitalized key",
		"ahyphenatedkey": "a hyphenated key",
		"identifierpls": "Weird key that isn't a valid Go identifier"
	}
	`

	trans := ValidIdentifierKeys()

	for _, testCase := range []struct {
		jsonBytes string
		direction Direction
	}{
		{invalidKeyJSON, Marshal},
		{invalidKeyJSON, Unmarshal},
		{validKeyJSON, Marshal},
		{validKeyJSON, Unmarshal},
	} {
		output := trans([]byte(testCase.jsonBytes), testCase.direction)

		if string(output) != validKeyJSON {
			t.Errorf("%s output of %s doesn't match expected %s", testCase.direction, output, validKeyJSON)
		}

		var testMapForKeys map[string]json.RawMessage

		json.Unmarshal(output, &testMapForKeys)

		for key := range testMapForKeys {
			keyAsGoIdentifierExpression := fmt.Sprintf("%s == 0", key)

			if _, err := parser.ParseExpr(keyAsGoIdentifierExpression); nil != err {
				t.Errorf("Transformed key %q is not a valid Go identifier", key)
			}
		}
	}
}
