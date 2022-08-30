package servicelinker

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetServiceBase interface {
}

func unmarshalTargetServiceBaseImplementation(input []byte) (TargetServiceBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TargetServiceBase into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureResource") {
		var out AzureResource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureResource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConfluentBootstrapServer") {
		var out ConfluentBootstrapServer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConfluentBootstrapServer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConfluentSchemaRegistry") {
		var out ConfluentSchemaRegistry
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConfluentSchemaRegistry: %+v", err)
		}
		return out, nil
	}

	type RawTargetServiceBaseImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawTargetServiceBaseImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
