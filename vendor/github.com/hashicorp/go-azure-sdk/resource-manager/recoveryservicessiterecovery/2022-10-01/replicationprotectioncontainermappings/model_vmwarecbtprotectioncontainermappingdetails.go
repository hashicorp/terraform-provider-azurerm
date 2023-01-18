package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionContainerMappingProviderSpecificDetails = VMwareCbtProtectionContainerMappingDetails{}

type VMwareCbtProtectionContainerMappingDetails struct {
	KeyVaultId                           *string           `json:"keyVaultId,omitempty"`
	KeyVaultUri                          *string           `json:"keyVaultUri,omitempty"`
	RoleSizeToNicCountMap                *map[string]int64 `json:"roleSizeToNicCountMap,omitempty"`
	ServiceBusConnectionStringSecretName *string           `json:"serviceBusConnectionStringSecretName,omitempty"`
	StorageAccountId                     *string           `json:"storageAccountId,omitempty"`
	StorageAccountSasSecretName          *string           `json:"storageAccountSasSecretName,omitempty"`
	TargetLocation                       *string           `json:"targetLocation,omitempty"`

	// Fields inherited from ProtectionContainerMappingProviderSpecificDetails
}

var _ json.Marshaler = VMwareCbtProtectionContainerMappingDetails{}

func (s VMwareCbtProtectionContainerMappingDetails) MarshalJSON() ([]byte, error) {
	type wrapper VMwareCbtProtectionContainerMappingDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMwareCbtProtectionContainerMappingDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMwareCbtProtectionContainerMappingDetails: %+v", err)
	}
	decoded["instanceType"] = "VMwareCbt"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMwareCbtProtectionContainerMappingDetails: %+v", err)
	}

	return encoded, nil
}
