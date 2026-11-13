package datasources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataChangeDetectionPolicy = HighWaterMarkChangeDetectionPolicy{}

type HighWaterMarkChangeDetectionPolicy struct {
	HighWaterMarkColumnName string `json:"highWaterMarkColumnName"`

	// Fields inherited from DataChangeDetectionPolicy

	OdataType string `json:"@odata.type"`
}

func (s HighWaterMarkChangeDetectionPolicy) DataChangeDetectionPolicy() BaseDataChangeDetectionPolicyImpl {
	return BaseDataChangeDetectionPolicyImpl{
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = HighWaterMarkChangeDetectionPolicy{}

func (s HighWaterMarkChangeDetectionPolicy) MarshalJSON() ([]byte, error) {
	type wrapper HighWaterMarkChangeDetectionPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HighWaterMarkChangeDetectionPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HighWaterMarkChangeDetectionPolicy: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.HighWaterMarkChangeDetectionPolicy"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HighWaterMarkChangeDetectionPolicy: %+v", err)
	}

	return encoded, nil
}
