package msgraph

import (
	"encoding/json"

	"github.com/hashicorp/go-uuid"
)

func MarshalDocs(docs [][]byte) ([]byte, error) {
	out := make(map[string]interface{})
	for _, d := range docs {
		var o map[string]interface{}
		err := json.Unmarshal(d, &o)
		if err != nil {
			return d, err
		}
		for k, v := range o {
			out[k] = v
		}
	}
	return json.Marshal(out)
}

func ValidateId(id *string) bool {
	if id == nil || *id == "" {
		return false
	}
	if _, err := uuid.ParseUUID(*id); err != nil {
		return false
	}
	return true
}
