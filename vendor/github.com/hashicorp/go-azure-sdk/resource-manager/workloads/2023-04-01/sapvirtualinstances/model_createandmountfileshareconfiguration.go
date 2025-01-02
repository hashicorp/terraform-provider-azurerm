package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FileShareConfiguration = CreateAndMountFileShareConfiguration{}

type CreateAndMountFileShareConfiguration struct {
	ResourceGroup      *string `json:"resourceGroup,omitempty"`
	StorageAccountName *string `json:"storageAccountName,omitempty"`

	// Fields inherited from FileShareConfiguration

	ConfigurationType ConfigurationType `json:"configurationType"`
}

func (s CreateAndMountFileShareConfiguration) FileShareConfiguration() BaseFileShareConfigurationImpl {
	return BaseFileShareConfigurationImpl{
		ConfigurationType: s.ConfigurationType,
	}
}

var _ json.Marshaler = CreateAndMountFileShareConfiguration{}

func (s CreateAndMountFileShareConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper CreateAndMountFileShareConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CreateAndMountFileShareConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CreateAndMountFileShareConfiguration: %+v", err)
	}

	decoded["configurationType"] = "CreateAndMount"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CreateAndMountFileShareConfiguration: %+v", err)
	}

	return encoded, nil
}
