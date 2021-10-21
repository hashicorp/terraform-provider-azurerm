package videoanalyzer

import (
	"encoding/json"
	"fmt"
)

type AuthenticationBase interface {
}

func unmarshalAuthenticationBase(body []byte) (AuthenticationBase, error) {
	type intermediateType struct {
		Type string `json:"@type"`
	}
	var intermediate intermediateType
	if err := json.Unmarshal(body, &intermediate); err != nil {
		return nil, fmt.Errorf("decoding: %+v", err)
	}

	switch intermediate.Type {

	case "#Microsoft.VideoAnalyzer.JwtAuthentication":
		{
			var out JwtAuthentication
			if err := json.Unmarshal(body, &out); err != nil {
				return nil, fmt.Errorf("unmarshalling into %q: %+v", "JwtAuthentication", err)
			}
			return &out, nil
		}

	}

	return nil, fmt.Errorf("unknown value for Type: %q", intermediate.Type)
}
