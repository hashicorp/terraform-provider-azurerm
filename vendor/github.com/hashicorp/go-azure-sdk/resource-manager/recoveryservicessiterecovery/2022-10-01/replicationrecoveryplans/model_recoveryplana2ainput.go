package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RecoveryPlanProviderSpecificInput = RecoveryPlanA2AInput{}

type RecoveryPlanA2AInput struct {
	PrimaryExtendedLocation  *edgezones.Model `json:"primaryExtendedLocation,omitempty"`
	PrimaryZone              *string          `json:"primaryZone,omitempty"`
	RecoveryExtendedLocation *edgezones.Model `json:"recoveryExtendedLocation,omitempty"`
	RecoveryZone             *string          `json:"recoveryZone,omitempty"`

	// Fields inherited from RecoveryPlanProviderSpecificInput
}

var _ json.Marshaler = RecoveryPlanA2AInput{}

func (s RecoveryPlanA2AInput) MarshalJSON() ([]byte, error) {
	type wrapper RecoveryPlanA2AInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RecoveryPlanA2AInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanA2AInput: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RecoveryPlanA2AInput: %+v", err)
	}

	return encoded, nil
}
