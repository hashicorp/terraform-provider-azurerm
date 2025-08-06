package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EnableProtectionProviderSpecificInput = A2ACrossClusterMigrationEnableProtectionInput{}

type A2ACrossClusterMigrationEnableProtectionInput struct {
	FabricObjectId      *string `json:"fabricObjectId,omitempty"`
	RecoveryContainerId *string `json:"recoveryContainerId,omitempty"`

	// Fields inherited from EnableProtectionProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s A2ACrossClusterMigrationEnableProtectionInput) EnableProtectionProviderSpecificInput() BaseEnableProtectionProviderSpecificInputImpl {
	return BaseEnableProtectionProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = A2ACrossClusterMigrationEnableProtectionInput{}

func (s A2ACrossClusterMigrationEnableProtectionInput) MarshalJSON() ([]byte, error) {
	type wrapper A2ACrossClusterMigrationEnableProtectionInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2ACrossClusterMigrationEnableProtectionInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2ACrossClusterMigrationEnableProtectionInput: %+v", err)
	}

	decoded["instanceType"] = "A2ACrossClusterMigration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2ACrossClusterMigrationEnableProtectionInput: %+v", err)
	}

	return encoded, nil
}
