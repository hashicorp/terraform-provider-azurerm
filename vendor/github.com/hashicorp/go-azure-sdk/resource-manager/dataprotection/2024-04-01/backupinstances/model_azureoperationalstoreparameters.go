package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataStoreParameters = AzureOperationalStoreParameters{}

type AzureOperationalStoreParameters struct {
	ResourceGroupId *string `json:"resourceGroupId,omitempty"`

	// Fields inherited from DataStoreParameters

	DataStoreType DataStoreTypes `json:"dataStoreType"`
	ObjectType    string         `json:"objectType"`
}

func (s AzureOperationalStoreParameters) DataStoreParameters() BaseDataStoreParametersImpl {
	return BaseDataStoreParametersImpl{
		DataStoreType: s.DataStoreType,
		ObjectType:    s.ObjectType,
	}
}

var _ json.Marshaler = AzureOperationalStoreParameters{}

func (s AzureOperationalStoreParameters) MarshalJSON() ([]byte, error) {
	type wrapper AzureOperationalStoreParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureOperationalStoreParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureOperationalStoreParameters: %+v", err)
	}

	decoded["objectType"] = "AzureOperationalStoreParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureOperationalStoreParameters: %+v", err)
	}

	return encoded, nil
}
