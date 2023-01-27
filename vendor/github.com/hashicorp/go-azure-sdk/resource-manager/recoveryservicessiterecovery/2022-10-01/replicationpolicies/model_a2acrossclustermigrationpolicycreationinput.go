package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificInput = A2ACrossClusterMigrationPolicyCreationInput{}

type A2ACrossClusterMigrationPolicyCreationInput struct {

	// Fields inherited from PolicyProviderSpecificInput
}

var _ json.Marshaler = A2ACrossClusterMigrationPolicyCreationInput{}

func (s A2ACrossClusterMigrationPolicyCreationInput) MarshalJSON() ([]byte, error) {
	type wrapper A2ACrossClusterMigrationPolicyCreationInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2ACrossClusterMigrationPolicyCreationInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2ACrossClusterMigrationPolicyCreationInput: %+v", err)
	}
	decoded["instanceType"] = "A2ACrossClusterMigration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2ACrossClusterMigrationPolicyCreationInput: %+v", err)
	}

	return encoded, nil
}
