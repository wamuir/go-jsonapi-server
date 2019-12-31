package schema

import "github.com/xeipuuv/gojsonschema"

var schema *gojsonschema.Schema = setSchema(gohex)

func setSchema(bytes []byte) *gojsonschema.Schema {

	loader := gojsonschema.NewBytesLoader(bytes)

	s, err := gojsonschema.NewSchema(loader)
	if err != nil {
		panic(err)
	}

	return s
}

func Validate(document interface{}) (*gojsonschema.Result, error) {

	loader := gojsonschema.NewGoLoader(document)

	result, err := schema.Validate(loader)
	if err != nil {
		return result, err
	}

	return result, nil
}
