package recoverypoint

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBackupRecoveryPoint interface {
}

func unmarshalAzureBackupRecoveryPointImplementation(input []byte) (AzureBackupRecoveryPoint, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureBackupRecoveryPoint into map[string]interface: %+v", err)
	}

	value, ok := temp["objectType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureBackupDiscreteRecoveryPoint") {
		var out AzureBackupDiscreteRecoveryPoint
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBackupDiscreteRecoveryPoint: %+v", err)
		}
		return out, nil
	}

	type RawAzureBackupRecoveryPointImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawAzureBackupRecoveryPointImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
