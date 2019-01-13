package conjson_test

import (
	"encoding/json"
	"fmt"
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

	inceptionDateTime = "2015-11-17T20:43:31.0463576-05:00"
	packageDateTime   = "2018-12-24T13:21:15.7883416-07:00"

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
	//     "created_at": "2015-11-17T20:43:31.0463576-05:00",
	//     "updated_at": "2018-12-24T13:21:15.7883416-07:00"
	// }
}
