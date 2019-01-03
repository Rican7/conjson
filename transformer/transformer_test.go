package transformer

import (
	"testing"
)

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
