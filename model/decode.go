package model

import (
	"encoding/json"
	"github.com/wamuir/go-jsonapi-server/schema"
	"io"
	"net/http"
)

func Decode(body io.Reader) (*Document, *ErrorObject) {

	var document *Document

	// Decode JSON document
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&document)
	if err != nil {
		errObj := MakeError(http.StatusBadRequest)
		errObj.Code = "e6f91b"
		return document, errObj
	}

	// Decode Data member
	data, errObj := decodeDataMbr(document.Data)
	if errObj != nil {
		return document, errObj
	}
	document.Data = data

	// Validate
	result, err := schema.Validate(document)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "52f5cf"
		errObj.Title = "Encountered internal error while validating request body against JSON:API schema"
		errObj.Detail = err.Error()
		return document, errObj
	}

	if !result.Valid() {
		errObj := MakeError(http.StatusBadRequest)
		errObj.Code = "d96685"
		errObj.Title = "Request body failed to validate against JSON:API schema"
		// errObj.Detail = <-- need to add this... pointer?
		return document, errObj
	}

	return document, nil
}

func decodeDataMbr(i interface{}) (interface{}, *ErrorObject) {

	j, err := json.Marshal(i)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "e36d46"
		errObj.Detail = err.Error()
		return i, errObj
	}

	var resource Resource
	err = json.Unmarshal(j, &resource)
	if err == nil {
		return resource, nil
	}

	var collection Collection
	err = json.Unmarshal(j, &collection)
	if err == nil {
		return collection, nil
	}

	errObj := MakeError(http.StatusBadRequest)
	errObj.Code = "885bbb"
	return nil, errObj
}
