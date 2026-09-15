package datasources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataDeletionDetectionPolicy = SoftDeleteColumnDeletionDetectionPolicy{}

type SoftDeleteColumnDeletionDetectionPolicy struct {
	SoftDeleteColumnName  *string `json:"softDeleteColumnName,omitempty"`
	SoftDeleteMarkerValue *string `json:"softDeleteMarkerValue,omitempty"`

	// Fields inherited from DataDeletionDetectionPolicy

	OdataType string `json:"@odata.type"`
}

func (s SoftDeleteColumnDeletionDetectionPolicy) DataDeletionDetectionPolicy() BaseDataDeletionDetectionPolicyImpl {
	return BaseDataDeletionDetectionPolicyImpl{
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = SoftDeleteColumnDeletionDetectionPolicy{}

func (s SoftDeleteColumnDeletionDetectionPolicy) MarshalJSON() ([]byte, error) {
	type wrapper SoftDeleteColumnDeletionDetectionPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SoftDeleteColumnDeletionDetectionPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SoftDeleteColumnDeletionDetectionPolicy: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.SoftDeleteColumnDeletionDetectionPolicy"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SoftDeleteColumnDeletionDetectionPolicy: %+v", err)
	}

	return encoded, nil
}
