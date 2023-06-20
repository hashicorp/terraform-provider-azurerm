package contentkeypolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPolicyRestriction interface {
}

func unmarshalContentKeyPolicyRestrictionImplementation(input []byte) (ContentKeyPolicyRestriction, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyRestriction into map[string]interface: %+v", err)
	}

	value, ok := temp["@odata.type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.ContentKeyPolicyOpenRestriction") {
		var out ContentKeyPolicyOpenRestriction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentKeyPolicyOpenRestriction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.ContentKeyPolicyTokenRestriction") {
		var out ContentKeyPolicyTokenRestriction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentKeyPolicyTokenRestriction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.ContentKeyPolicyUnknownRestriction") {
		var out ContentKeyPolicyUnknownRestriction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentKeyPolicyUnknownRestriction: %+v", err)
		}
		return out, nil
	}

	type RawContentKeyPolicyRestrictionImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawContentKeyPolicyRestrictionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
