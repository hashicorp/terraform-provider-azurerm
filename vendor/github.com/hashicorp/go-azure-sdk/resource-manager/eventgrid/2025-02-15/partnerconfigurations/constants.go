package partnerconfigurations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerConfigurationProvisioningState string

const (
	PartnerConfigurationProvisioningStateCanceled  PartnerConfigurationProvisioningState = "Canceled"
	PartnerConfigurationProvisioningStateCreating  PartnerConfigurationProvisioningState = "Creating"
	PartnerConfigurationProvisioningStateDeleting  PartnerConfigurationProvisioningState = "Deleting"
	PartnerConfigurationProvisioningStateFailed    PartnerConfigurationProvisioningState = "Failed"
	PartnerConfigurationProvisioningStateSucceeded PartnerConfigurationProvisioningState = "Succeeded"
	PartnerConfigurationProvisioningStateUpdating  PartnerConfigurationProvisioningState = "Updating"
)

func PossibleValuesForPartnerConfigurationProvisioningState() []string {
	return []string{
		string(PartnerConfigurationProvisioningStateCanceled),
		string(PartnerConfigurationProvisioningStateCreating),
		string(PartnerConfigurationProvisioningStateDeleting),
		string(PartnerConfigurationProvisioningStateFailed),
		string(PartnerConfigurationProvisioningStateSucceeded),
		string(PartnerConfigurationProvisioningStateUpdating),
	}
}

func (s *PartnerConfigurationProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePartnerConfigurationProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePartnerConfigurationProvisioningState(input string) (*PartnerConfigurationProvisioningState, error) {
	vals := map[string]PartnerConfigurationProvisioningState{
		"canceled":  PartnerConfigurationProvisioningStateCanceled,
		"creating":  PartnerConfigurationProvisioningStateCreating,
		"deleting":  PartnerConfigurationProvisioningStateDeleting,
		"failed":    PartnerConfigurationProvisioningStateFailed,
		"succeeded": PartnerConfigurationProvisioningStateSucceeded,
		"updating":  PartnerConfigurationProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerConfigurationProvisioningState(input)
	return &out, nil
}
