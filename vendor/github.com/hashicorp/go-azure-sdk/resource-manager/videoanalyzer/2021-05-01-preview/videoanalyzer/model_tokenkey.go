package videoanalyzer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TokenKey interface {
}

func unmarshalTokenKeyImplementation(input []byte) (TokenKey, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TokenKey into map[string]interface: %+v", err)
	}

	value, ok := temp["@type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "#Microsoft.VideoAnalyzer.EccTokenKey") {
		var out EccTokenKey
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EccTokenKey: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.VideoAnalyzer.RsaTokenKey") {
		var out RsaTokenKey
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RsaTokenKey: %+v", err)
		}
		return out, nil
	}

	type RawTokenKeyImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawTokenKeyImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
