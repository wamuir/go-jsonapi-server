package schema

import (
	_ "embed"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed jsonapi-schema.modified.json
var b []byte

var schema *gojsonschema.Schema

func init() {
	loader := gojsonschema.NewBytesLoader(b)
	s, err := gojsonschema.NewSchema(loader)
	if err != nil {
		panic(err)
	}
	schema = s
}

func Validate(document interface{}) (*gojsonschema.Result, error) {
	loader := gojsonschema.NewGoLoader(document)
	return schema.Validate(loader)
}
