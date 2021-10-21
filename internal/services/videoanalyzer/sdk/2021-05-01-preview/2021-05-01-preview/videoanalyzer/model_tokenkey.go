package videoanalyzer

import (
	"encoding/json"
	"fmt"
)

type TokenKey interface {
}

func unmarshalTokenKey(body []byte) (TokenKey, error) {
	type intermediateType struct {
		Type string `json:"@type"`
	}
	var intermediate intermediateType
	if err := json.Unmarshal(body, &intermediate); err != nil {
		return nil, fmt.Errorf("decoding: %+v", err)
	}

	switch intermediate.Type {

	case "#Microsoft.VideoAnalyzer.EccTokenKey":
		{
			var out EccTokenKey
			if err := json.Unmarshal(body, &out); err != nil {
				return nil, fmt.Errorf("unmarshalling into %q: %+v", "EccTokenKey", err)
			}
			return &out, nil
		}

	case "#Microsoft.VideoAnalyzer.RsaTokenKey":
		{
			var out RsaTokenKey
			if err := json.Unmarshal(body, &out); err != nil {
				return nil, fmt.Errorf("unmarshalling into %q: %+v", "RsaTokenKey", err)
			}
			return &out, nil
		}

	}

	return nil, fmt.Errorf("unknown value for Type: %q", intermediate.Type)
}
