package recoverypoint

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AzureBackupRecoveryPoint = AzureBackupDiscreteRecoveryPoint{}

type AzureBackupDiscreteRecoveryPoint struct {
	FriendlyName                   *string                          `json:"friendlyName,omitempty"`
	PolicyName                     *string                          `json:"policyName,omitempty"`
	PolicyVersion                  *string                          `json:"policyVersion,omitempty"`
	RecoveryPointDataStoresDetails *[]RecoveryPointDataStoreDetails `json:"recoveryPointDataStoresDetails,omitempty"`
	RecoveryPointId                *string                          `json:"recoveryPointId,omitempty"`
	RecoveryPointTime              string                           `json:"recoveryPointTime"`
	RecoveryPointType              *string                          `json:"recoveryPointType,omitempty"`
	RetentionTagName               *string                          `json:"retentionTagName,omitempty"`
	RetentionTagVersion            *string                          `json:"retentionTagVersion,omitempty"`

	// Fields inherited from AzureBackupRecoveryPoint
}

var _ json.Marshaler = AzureBackupDiscreteRecoveryPoint{}

func (s AzureBackupDiscreteRecoveryPoint) MarshalJSON() ([]byte, error) {
	type wrapper AzureBackupDiscreteRecoveryPoint
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureBackupDiscreteRecoveryPoint: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureBackupDiscreteRecoveryPoint: %+v", err)
	}
	decoded["objectType"] = "AzureBackupDiscreteRecoveryPoint"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureBackupDiscreteRecoveryPoint: %+v", err)
	}

	return encoded, nil
}
