# conjson

[![Build Status](https://travis-ci.org/Rican7/conjson.svg?branch=master)](https://travis-ci.org/Rican7/conjson)
[![Coverage Status](https://coveralls.io/repos/github/Rican7/conjson/badge.svg)](https://coveralls.io/github/Rican7/conjson)
[![Go Report Card](https://goreportcard.com/badge/Rican7/conjson)](http://goreportcard.com/report/Rican7/conjson)
[![GoDoc](https://godoc.org/github.com/Rican7/conjson?status.svg)](https://godoc.org/github.com/Rican7/conjson)
[![Latest Stable Version](https://img.shields.io/github/release/Rican7/conjson.svg?style=flat)](https://github.com/Rican7/conjson/releases)

**conjson** - (conventional, consistent, conformative) JSON

A simple, functional, no-tags-required mechanism to handle and transform JSON representations of values, consistently.


## Project Status

This project is currently in "pre-release". While the code is heavily tested, the API may change.
Vendor or "lock" this dependency if you plan on using it.


## Examples

### Marshal

```go
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
```

### Unmarshal

```go
sampleJSON := `
{
	"title": "Example Title",
	"description": "This is a description.",
	"image_url": "https://example.com/image.png",
	"referred_by_url": "https://example.com/referrer/index.html",
	"is_active": true,
	"created_at": "2015-11-17T20:43:31.0463576-05:00",
	"updated_at": "2018-12-24T13:21:15.7883416-07:00"
}
`

var model exampleModel

json.Unmarshal(
	[]byte(sampleJSON),
	conjson.NewUnmarshaler(&model, transform.ConventionalKeys()),
)

fmt.Println(model.ReferredByURL)
// Output: https://example.com/referrer/index.html
```

### Encode

```go
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
//     "createdAt": "2015-11-17T20:43:31.0463576-05:00",
//     "updatedAt": "2018-12-24T13:21:15.7883416-07:00"
// }
```

### Decode

```go
sampleJSON := `
{
	"$title--": "Example Title",
	"$description--": "This is a description.",
	"$image_url--": "https://example.com/image.png",
	"$referred_by_url--": "https://example.com/referrer/index.html",
	"$is_active--": true,
	"$created_at--": "2015-11-17T20:43:31.0463576-05:00",
	"updated_at--": "2018-12-24T13:21:15.7883416-07:00"
}
`

var model exampleModel

decoder := conjson.NewDecoder(
	json.NewDecoder(bytes.NewBufferString(sampleJSON)),
	transform.ConventionalKeys(),
	transform.ValidIdentifierKeys(),
)

decoder.Decode(&model)

fmt.Println(model.Title)
fmt.Println(model.Description)
fmt.Println(model.ImageURL)
// Output:
// Example Title
// This is a description.
// https://example.com/image.png
```
