package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificContainerMappingInput = VMwareCbtContainerMappingInput{}

type VMwareCbtContainerMappingInput struct {
	KeyVaultId                           *string `json:"keyVaultId,omitempty"`
	KeyVaultUri                          *string `json:"keyVaultUri,omitempty"`
	ServiceBusConnectionStringSecretName *string `json:"serviceBusConnectionStringSecretName,omitempty"`
	StorageAccountId                     string  `json:"storageAccountId"`
	StorageAccountSasSecretName          *string `json:"storageAccountSasSecretName,omitempty"`
	TargetLocation                       string  `json:"targetLocation"`

	// Fields inherited from ReplicationProviderSpecificContainerMappingInput

	InstanceType string `json:"instanceType"`
}

func (s VMwareCbtContainerMappingInput) ReplicationProviderSpecificContainerMappingInput() BaseReplicationProviderSpecificContainerMappingInputImpl {
	return BaseReplicationProviderSpecificContainerMappingInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = VMwareCbtContainerMappingInput{}

func (s VMwareCbtContainerMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper VMwareCbtContainerMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMwareCbtContainerMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMwareCbtContainerMappingInput: %+v", err)
	}

	decoded["instanceType"] = "VMwareCbt"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMwareCbtContainerMappingInput: %+v", err)
	}

	return encoded, nil
}
