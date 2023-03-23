package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BaseBackupPolicy interface {
}

func unmarshalBaseBackupPolicyImplementation(input []byte) (BaseBackupPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BaseBackupPolicy into map[string]interface: %+v", err)
	}

	value, ok := temp["objectType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "BackupPolicy") {
		var out BackupPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BackupPolicy: %+v", err)
		}
		return out, nil
	}

	type RawBaseBackupPolicyImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawBaseBackupPolicyImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
