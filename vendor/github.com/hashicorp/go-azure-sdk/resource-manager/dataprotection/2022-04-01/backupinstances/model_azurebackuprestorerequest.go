package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBackupRestoreRequest interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAzureBackupRestoreRequestImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalAzureBackupRestoreRequestImplementation(input []byte) (AzureBackupRestoreRequest, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureBackupRestoreRequest into map[string]interface: %+v", err)
	}

	value, ok := temp["objectType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureBackupRecoveryPointBasedRestoreRequest") {
		var out AzureBackupRecoveryPointBasedRestoreRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBackupRecoveryPointBasedRestoreRequest: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBackupRecoveryTimeBasedRestoreRequest") {
		var out AzureBackupRecoveryTimeBasedRestoreRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBackupRecoveryTimeBasedRestoreRequest: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBackupRestoreWithRehydrationRequest") {
		var out AzureBackupRestoreWithRehydrationRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBackupRestoreWithRehydrationRequest: %+v", err)
		}
		return out, nil
	}

	out := RawAzureBackupRestoreRequestImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
