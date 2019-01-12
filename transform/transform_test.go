package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/parser"
	"testing"
)

func TestDirection_String(t *testing.T) {
	for _, testData := range []struct {
		direction      Direction
		expectedString string
	}{
		{Marshal, "Marshal"},
		{Unmarshal, "Unmarshal"},
	} {
		if str := testData.direction.String(); testData.expectedString != str {
			t.Errorf("%q doesn't match expected %q", str, testData.expectedString)
		}
	}
}

func TestOnlyForDirection(t *testing.T) {
	mockDataReturn := []byte("mock data")

	mockTransformer := func(data []byte, direction Direction) []byte {
		return mockDataReturn
	}

	// Compile-time functional-interface type check/enforcement "test"
	var marshalTrans Transformer = OnlyForDirection(Marshal, mockTransformer)
	var unmarshalTrans Transformer = OnlyForDirection(Unmarshal, mockTransformer)

	for _, testData := range []struct {
		trans        Transformer
		direction    Direction
		callExpected bool
	}{
		{marshalTrans, Marshal, true},
		{marshalTrans, Unmarshal, false},
		{unmarshalTrans, Marshal, false},
		{unmarshalTrans, Unmarshal, true},
	} {
		if output := testData.trans([]byte(""), testData.direction); testData.callExpected && !bytes.Equal(output, mockDataReturn) {
			t.Errorf("%s output of %q doesn't match expected %q", testData.direction, output, mockDataReturn)
		}
	}
}

func TestAlwaysAsDirection(t *testing.T) {
	mockTransformer := func(data []byte, direction Direction) []byte {
		return []byte(direction.String())
	}

	// Compile-time functional-interface type check/enforcement "test"
	var marshalTrans Transformer = AlwaysAsDirection(Marshal, mockTransformer)
	var unmarshalTrans Transformer = AlwaysAsDirection(Unmarshal, mockTransformer)

	for _, testData := range []struct {
		trans             Transformer
		expectedDirection Direction
		direction         Direction
	}{
		{marshalTrans, Marshal, Marshal},
		{marshalTrans, Marshal, Unmarshal},
		{unmarshalTrans, Unmarshal, Marshal},
		{unmarshalTrans, Unmarshal, Unmarshal},
	} {
		expectedDirectionBytes := []byte(testData.expectedDirection.String())

		if output := testData.trans([]byte(""), testData.direction); !bytes.Equal(output, expectedDirectionBytes) {
			t.Errorf("%s output of %q doesn't match expected %q", testData.direction, output, expectedDirectionBytes)
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

	// Compile-time functional-interface type check/enforcement "test"
	var trans Transformer = ConventionalKeys()

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
		"UpperCamelCaseKey": "an UpperCamelCase key"
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
		"upperCamelCaseKey": "an UpperCamelCase key"
	}
	`

	// Compile-time functional-interface type check/enforcement "test"
	var trans Transformer = CamelCaseKeys()

	for _, testData := range []struct {
		jsonBytes string
		direction Direction
	}{
		{originalJSON, Marshal},
		{originalJSON, Unmarshal},
		{camelCaseJSON, Marshal},
		{camelCaseJSON, Unmarshal},
	} {
		if output := trans([]byte(testData.jsonBytes), testData.direction); string(output) != camelCaseJSON {
			t.Errorf("%s output of %s doesn't match expected %s", testData.direction, output, camelCaseJSON)
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

	// Compile-time functional-interface type check/enforcement "test"
	var trans Transformer = ValidIdentifierKeys()

	for _, testData := range []struct {
		jsonBytes string
		direction Direction
	}{
		{invalidKeyJSON, Marshal},
		{invalidKeyJSON, Unmarshal},
		{validKeyJSON, Marshal},
		{validKeyJSON, Unmarshal},
	} {
		output := trans([]byte(testData.jsonBytes), testData.direction)

		if string(output) != validKeyJSON {
			t.Errorf("%s output of %s doesn't match expected %s", testData.direction, output, validKeyJSON)
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
