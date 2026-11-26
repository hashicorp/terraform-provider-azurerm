package datasources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataChangeDetectionPolicy = SqlIntegratedChangeTrackingPolicy{}

type SqlIntegratedChangeTrackingPolicy struct {

	// Fields inherited from DataChangeDetectionPolicy

	OdataType string `json:"@odata.type"`
}

func (s SqlIntegratedChangeTrackingPolicy) DataChangeDetectionPolicy() BaseDataChangeDetectionPolicyImpl {
	return BaseDataChangeDetectionPolicyImpl{
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = SqlIntegratedChangeTrackingPolicy{}

func (s SqlIntegratedChangeTrackingPolicy) MarshalJSON() ([]byte, error) {
	type wrapper SqlIntegratedChangeTrackingPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SqlIntegratedChangeTrackingPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SqlIntegratedChangeTrackingPolicy: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.SqlIntegratedChangeTrackingPolicy"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SqlIntegratedChangeTrackingPolicy: %+v", err)
	}

	return encoded, nil
}
