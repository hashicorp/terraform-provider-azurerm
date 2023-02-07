package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RecoveryPlanProviderSpecificDetails = RecoveryPlanA2ADetails{}

type RecoveryPlanA2ADetails struct {
	PrimaryExtendedLocation  *edgezones.Model `json:"primaryExtendedLocation,omitempty"`
	PrimaryZone              *string          `json:"primaryZone,omitempty"`
	RecoveryExtendedLocation *edgezones.Model `json:"recoveryExtendedLocation,omitempty"`
	RecoveryZone             *string          `json:"recoveryZone,omitempty"`

	// Fields inherited from RecoveryPlanProviderSpecificDetails
}

var _ json.Marshaler = RecoveryPlanA2ADetails{}

func (s RecoveryPlanA2ADetails) MarshalJSON() ([]byte, error) {
	type wrapper RecoveryPlanA2ADetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RecoveryPlanA2ADetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanA2ADetails: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RecoveryPlanA2ADetails: %+v", err)
	}

	return encoded, nil
}
