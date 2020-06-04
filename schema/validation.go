package schema

import "github.com/xeipuuv/gojsonschema"

var schema *gojsonschema.Schema

func init() {

	loader := gojsonschema.NewBytesLoader(gohex)

	if s, err := gojsonschema.NewSchema(loader); err != nil {
		panic(err)
	} else {
		schema = s
	}
}

func Validate(document interface{}) (*gojsonschema.Result, error) {

	loader := gojsonschema.NewGoLoader(document)

	if result, err := schema.Validate(loader); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
