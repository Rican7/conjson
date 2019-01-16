package conjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Rican7/conjson"
	"github.com/Rican7/conjson/transform"
)

type exampleModel struct {
	Title         string
	Description   string
	ImageURL      string
	ReferredByURL string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

const (
	dateTimeFormat = time.RFC3339

	inceptionDateTime = "2015-11-17T20:43:31-05:00"
	packageDateTime   = "2018-12-24T13:21:15-07:00"

	marshalPrefix = ""
	marshalIndent = "    "
)

var (
	inceptionTime, _ = time.Parse(dateTimeFormat, inceptionDateTime)
	packageTime, _   = time.Parse(dateTimeFormat, packageDateTime)
)

func Example() {
	model := exampleModel{
		Title:         "Example Title",
		Description:   "This is a description.",
		ImageURL:      "https://example.com/image.png",
		ReferredByURL: "https://example.com/referrer/index.html",
		IsActive:      true,
		CreatedAt:     inceptionTime,
		UpdatedAt:     packageTime,
	}

	marshaler := conjson.NewMarshaler(model, transform.ConventionalKeys())

	encoded, _ := json.MarshalIndent(marshaler, marshalPrefix, marshalIndent)

	fmt.Println(string(encoded))
	// Output:
	// {
	//     "title": "Example Title",
	//     "description": "This is a description.",
	//     "image_url": "https://example.com/image.png",
	//     "referred_by_url": "https://example.com/referrer/index.html",
	//     "is_active": true,
	//     "created_at": "2015-11-17T20:43:31-05:00",
	//     "updated_at": "2018-12-24T13:21:15-07:00"
	// }
}

func Example_unmarshal() {
	sampleJSON := `
	{
	    "title": "Example Title",
	    "description": "This is a description.",
	    "image_url": "https://example.com/image.png",
	    "referred_by_url": "https://example.com/referrer/index.html",
	    "is_active": true,
	    "created_at": "2015-11-17T20:43:31-05:00",
	    "updated_at": "2018-12-24T13:21:15-07:00"
	}
	`

	var model exampleModel

	json.Unmarshal(
		[]byte(sampleJSON),
		conjson.NewUnmarshaler(&model, transform.ConventionalKeys()),
	)

	// Print the "raw" model JSON to show result
	rawJSON, _ := json.MarshalIndent(model, marshalPrefix, marshalIndent)
	fmt.Println(string(rawJSON))
	// Output:
	// {
	//     "Title": "Example Title",
	//     "Description": "This is a description.",
	//     "ImageURL": "https://example.com/image.png",
	//     "ReferredByURL": "https://example.com/referrer/index.html",
	//     "IsActive": true,
	//     "CreatedAt": "2015-11-17T20:43:31-05:00",
	//     "UpdatedAt": "2018-12-24T13:21:15-07:00"
	// }
}

func ExampleEncoder() {
	model := exampleModel{
		Title:         "Example Title",
		Description:   "This is a description.",
		ImageURL:      "https://example.com/image.png",
		ReferredByURL: "https://example.com/referrer/index.html",
		IsActive:      true,
		CreatedAt:     inceptionTime,
		UpdatedAt:     packageTime,
	}

	jsonEncoder := json.NewEncoder(os.Stdout)
	jsonEncoder.SetIndent(marshalPrefix, marshalIndent)

	conjson.NewEncoder(jsonEncoder, transform.CamelCaseKeys()).Encode(model)

	// Output:
	// {
	//     "title": "Example Title",
	//     "description": "This is a description.",
	//     "imageUrl": "https://example.com/image.png",
	//     "referredByUrl": "https://example.com/referrer/index.html",
	//     "isActive": true,
	//     "createdAt": "2015-11-17T20:43:31-05:00",
	//     "updatedAt": "2018-12-24T13:21:15-07:00"
	// }
}

func ExampleDecoder() {
	sampleJSON := `
	{
	    "$title--": "Example Title",
	    "$description--": "This is a description.",
	    "$image_url--": "https://example.com/image.png",
	    "$referred_by_url--": "https://example.com/referrer/index.html",
	    "$is_active--": true,
	    "created_at--": "2015-11-17T20:43:31-05:00",
	    "updated_at--": "2018-12-24T13:21:15-07:00"
	}
	`

	var model exampleModel

	decoder := conjson.NewDecoder(
		json.NewDecoder(bytes.NewBufferString(sampleJSON)),
		transform.ConventionalKeys(),
		transform.ValidIdentifierKeys(),
	)

	decoder.Decode(&model)

	// Print the "raw" model JSON to show result
	rawJSON, _ := json.MarshalIndent(model, marshalPrefix, marshalIndent)
	fmt.Println(string(rawJSON))
	// Output:
	// {
	//     "Title": "Example Title",
	//     "Description": "This is a description.",
	//     "ImageURL": "https://example.com/image.png",
	//     "ReferredByURL": "https://example.com/referrer/index.html",
	//     "IsActive": true,
	//     "CreatedAt": "2015-11-17T20:43:31-05:00",
	//     "UpdatedAt": "2018-12-24T13:21:15-07:00"
	// }
}
