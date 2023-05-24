package cosmosdb

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupPolicy interface {
}

func unmarshalBackupPolicyImplementation(input []byte) (BackupPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BackupPolicy into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Continuous") {
		var out ContinuousModeBackupPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContinuousModeBackupPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Periodic") {
		var out PeriodicModeBackupPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PeriodicModeBackupPolicy: %+v", err)
		}
		return out, nil
	}

	type RawBackupPolicyImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawBackupPolicyImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
