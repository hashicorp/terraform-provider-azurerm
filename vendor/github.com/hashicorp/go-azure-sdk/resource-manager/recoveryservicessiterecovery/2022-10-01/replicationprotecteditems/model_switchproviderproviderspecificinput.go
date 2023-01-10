package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SwitchProviderProviderSpecificInput interface {
}

func unmarshalSwitchProviderProviderSpecificInputImplementation(input []byte) (SwitchProviderProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SwitchProviderProviderSpecificInput into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out InMageAzureV2SwitchProviderProviderInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageAzureV2SwitchProviderProviderInput: %+v", err)
		}
		return out, nil
	}

	type RawSwitchProviderProviderSpecificInputImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSwitchProviderProviderSpecificInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
