package machinelearningcomputes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ComputeSecrets = AksComputeSecrets{}

type AksComputeSecrets struct {
	AdminKubeConfig     *string `json:"adminKubeConfig,omitempty"`
	ImagePullSecretName *string `json:"imagePullSecretName,omitempty"`
	UserKubeConfig      *string `json:"userKubeConfig,omitempty"`

	// Fields inherited from ComputeSecrets
}

var _ json.Marshaler = AksComputeSecrets{}

func (s AksComputeSecrets) MarshalJSON() ([]byte, error) {
	type wrapper AksComputeSecrets
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AksComputeSecrets: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AksComputeSecrets: %+v", err)
	}
	decoded["computeType"] = "AKS"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AksComputeSecrets: %+v", err)
	}

	return encoded, nil
}
