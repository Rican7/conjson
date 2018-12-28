package transformer

import "testing"

const snakeCaseJSON = `
{
    "title": "Example Title",
    "description": "whatever",
    "created_at": "2015-11-17T20:43:31.0463576-05:00",
	"updated_at": "2018-12-24T13:21:15.7883416-07:00",
    "is_active": true,
    "image_url": "https://example.com/image.png",
	"time_to_live": 600
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
	"timeToLive": 600
}
`

func TestConventionalKeys(t *testing.T) {
	// Compile-time functional-interface type check/enforcement "test"
	var trans Transformer = ConventionalKeys()

	if output := trans([]byte(camelCaseJSON), Marshal); string(output) != snakeCaseJSON {
		t.Errorf("Marshal output of %q doesn't match expected %q", output, snakeCaseJSON)
	}

	if output := trans([]byte(snakeCaseJSON), Unmarshal); string(output) != camelCaseJSON {
		t.Errorf("Unmarshal output of %q doesn't match expected %q", output, camelCaseJSON)
	}
}
