package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateApplianceForReplicationProtectedItemProviderSpecificInput interface {
}

func unmarshalUpdateApplianceForReplicationProtectedItemProviderSpecificInputImplementation(input []byte) (UpdateApplianceForReplicationProtectedItemProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling UpdateApplianceForReplicationProtectedItemProviderSpecificInput into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmUpdateApplianceForReplicationProtectedItemInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmUpdateApplianceForReplicationProtectedItemInput: %+v", err)
		}
		return out, nil
	}

	type RawUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
