package encodings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobInput interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawJobInputImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalJobInputImplementation(input []byte) (JobInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling JobInput into map[string]interface: %+v", err)
	}

	value, ok := temp["@odata.type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.JobInputAsset") {
		var out JobInputAsset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JobInputAsset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.JobInputClip") {
		var out JobInputClip
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JobInputClip: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.JobInputHttp") {
		var out JobInputHTTP
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JobInputHTTP: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.JobInputSequence") {
		var out JobInputSequence
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JobInputSequence: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.JobInputs") {
		var out JobInputs
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JobInputs: %+v", err)
		}
		return out, nil
	}

	out := RawJobInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
