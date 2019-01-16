# conjson

[![Build Status](https://travis-ci.com/Rican7/conjson.svg?branch=master)](https://travis-ci.com/Rican7/conjson)
[![Coverage Status](https://coveralls.io/repos/github/Rican7/conjson/badge.svg)](https://coveralls.io/github/Rican7/conjson)
[![Go Report Card](https://goreportcard.com/badge/github.com/Rican7/conjson)](https://goreportcard.com/report/github.com/Rican7/conjson)
[![GoDoc](https://godoc.org/github.com/Rican7/conjson?status.svg)](https://godoc.org/github.com/Rican7/conjson)
[![Latest Stable Version](https://img.shields.io/github/release/Rican7/conjson.svg?style=flat)](https://github.com/Rican7/conjson/releases)

**conjson** - (conventional, consistent, conformative) JSON

A simple, functional, no-tags-required mechanism to handle and transform JSON representations of values, consistently.


## Project Status

This project is currently in "pre-release". While the code is heavily tested, the API may change.
Vendor or "lock" this dependency if you plan on using it.


## History

For the curious, this project was born from [a 3+ year old (at the time of creation) "Gist"](https://gist.github.com/Rican7/39a3dc10c1499384ca91).

Both that Gist, and eventually this larger project, were inspired by a desire to more easily work with APIs that
accepted/returned JSON that had "snake_case"-style object keys.

Basically, I wanted a way to Marshal and Unmarshal Go structures without having to add "tags" to each and every field of
each and every structure. That Gist solved that problem for me, and now this library can do the same but with more
power and flexibility.


## Examples

### Marshal a Go structure into "conventional" style JSON

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
//     "created_at": "2015-11-17T20:43:31-05:00",
//     "updated_at": "2018-12-24T13:21:15-07:00"
// }
```

### Unmarshal "conventional" style JSON into a Go structure

```go
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
```

### Encode a Go structure into "camelCase" style JSON

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
//     "createdAt": "2015-11-17T20:43:31-05:00",
//     "updatedAt": "2018-12-24T13:21:15-07:00"
// }
```

### Decode JSON with atypical keys into a Go structure

```go
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
```
