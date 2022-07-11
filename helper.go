package api

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var jsonEscapeHTML = jsoniter.Config{
	EscapeHTML:                    false,
	SortMapKeys:                   true,
	ValidateJsonRawMessage:        false,
	ObjectFieldMustBeSimpleString: true,
}.Froze()
