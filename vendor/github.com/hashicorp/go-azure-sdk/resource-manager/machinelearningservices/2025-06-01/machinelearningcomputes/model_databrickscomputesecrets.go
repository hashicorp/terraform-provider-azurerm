package machinelearningcomputes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ComputeSecrets = DatabricksComputeSecrets{}

type DatabricksComputeSecrets struct {
	DatabricksAccessToken *string `json:"databricksAccessToken,omitempty"`

	// Fields inherited from ComputeSecrets

	ComputeType ComputeType `json:"computeType"`
}

func (s DatabricksComputeSecrets) ComputeSecrets() BaseComputeSecretsImpl {
	return BaseComputeSecretsImpl{
		ComputeType: s.ComputeType,
	}
}

var _ json.Marshaler = DatabricksComputeSecrets{}

func (s DatabricksComputeSecrets) MarshalJSON() ([]byte, error) {
	type wrapper DatabricksComputeSecrets
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DatabricksComputeSecrets: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DatabricksComputeSecrets: %+v", err)
	}

	decoded["computeType"] = "Databricks"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DatabricksComputeSecrets: %+v", err)
	}

	return encoded, nil
}
