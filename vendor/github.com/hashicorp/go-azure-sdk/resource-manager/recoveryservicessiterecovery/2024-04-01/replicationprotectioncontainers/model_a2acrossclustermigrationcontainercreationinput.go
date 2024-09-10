package replicationprotectioncontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificContainerCreationInput = A2ACrossClusterMigrationContainerCreationInput{}

type A2ACrossClusterMigrationContainerCreationInput struct {

	// Fields inherited from ReplicationProviderSpecificContainerCreationInput
}

var _ json.Marshaler = A2ACrossClusterMigrationContainerCreationInput{}

func (s A2ACrossClusterMigrationContainerCreationInput) MarshalJSON() ([]byte, error) {
	type wrapper A2ACrossClusterMigrationContainerCreationInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2ACrossClusterMigrationContainerCreationInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2ACrossClusterMigrationContainerCreationInput: %+v", err)
	}
	decoded["instanceType"] = "A2ACrossClusterMigration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2ACrossClusterMigrationContainerCreationInput: %+v", err)
	}

	return encoded, nil
}
