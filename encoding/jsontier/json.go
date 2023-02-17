package jsontier

import (
	"github.com/Dert12318/Utilities/encoding"
	jsoniter "github.com/json-iterator/go"
)

type (
	implementation struct {
		json jsoniter.API
	}
)

func (i implementation) Marshal(val interface{}) ([]byte, error) {
	return i.json.Marshal(val)
}

func (i implementation) Unmarshal(data []byte, val interface{}) error {
	return i.json.Unmarshal(data, val)
}

func NewEncoding() encoding.Encoding {
	return &implementation{json: jsoniter.ConfigCompatibleWithStandardLibrary}
}
