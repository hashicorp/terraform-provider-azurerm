package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DisableProtectionProviderSpecificInput interface {
}

func unmarshalDisableProtectionProviderSpecificInputImplementation(input []byte) (DisableProtectionProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DisableProtectionProviderSpecificInput into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "InMage") {
		var out InMageDisableProtectionProviderSpecificInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageDisableProtectionProviderSpecificInput: %+v", err)
		}
		return out, nil
	}

	type RawDisableProtectionProviderSpecificInputImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawDisableProtectionProviderSpecificInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
