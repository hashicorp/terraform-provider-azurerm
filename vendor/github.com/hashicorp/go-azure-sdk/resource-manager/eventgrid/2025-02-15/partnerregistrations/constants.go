package partnerregistrations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerRegistrationProvisioningState string

const (
	PartnerRegistrationProvisioningStateCanceled  PartnerRegistrationProvisioningState = "Canceled"
	PartnerRegistrationProvisioningStateCreating  PartnerRegistrationProvisioningState = "Creating"
	PartnerRegistrationProvisioningStateDeleting  PartnerRegistrationProvisioningState = "Deleting"
	PartnerRegistrationProvisioningStateFailed    PartnerRegistrationProvisioningState = "Failed"
	PartnerRegistrationProvisioningStateSucceeded PartnerRegistrationProvisioningState = "Succeeded"
	PartnerRegistrationProvisioningStateUpdating  PartnerRegistrationProvisioningState = "Updating"
)

func PossibleValuesForPartnerRegistrationProvisioningState() []string {
	return []string{
		string(PartnerRegistrationProvisioningStateCanceled),
		string(PartnerRegistrationProvisioningStateCreating),
		string(PartnerRegistrationProvisioningStateDeleting),
		string(PartnerRegistrationProvisioningStateFailed),
		string(PartnerRegistrationProvisioningStateSucceeded),
		string(PartnerRegistrationProvisioningStateUpdating),
	}
}

func (s *PartnerRegistrationProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePartnerRegistrationProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePartnerRegistrationProvisioningState(input string) (*PartnerRegistrationProvisioningState, error) {
	vals := map[string]PartnerRegistrationProvisioningState{
		"canceled":  PartnerRegistrationProvisioningStateCanceled,
		"creating":  PartnerRegistrationProvisioningStateCreating,
		"deleting":  PartnerRegistrationProvisioningStateDeleting,
		"failed":    PartnerRegistrationProvisioningStateFailed,
		"succeeded": PartnerRegistrationProvisioningStateSucceeded,
		"updating":  PartnerRegistrationProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerRegistrationProvisioningState(input)
	return &out, nil
}
