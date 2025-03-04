package verifiedpartners

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VerifiedPartnerProvisioningState string

const (
	VerifiedPartnerProvisioningStateCanceled  VerifiedPartnerProvisioningState = "Canceled"
	VerifiedPartnerProvisioningStateCreating  VerifiedPartnerProvisioningState = "Creating"
	VerifiedPartnerProvisioningStateDeleting  VerifiedPartnerProvisioningState = "Deleting"
	VerifiedPartnerProvisioningStateFailed    VerifiedPartnerProvisioningState = "Failed"
	VerifiedPartnerProvisioningStateSucceeded VerifiedPartnerProvisioningState = "Succeeded"
	VerifiedPartnerProvisioningStateUpdating  VerifiedPartnerProvisioningState = "Updating"
)

func PossibleValuesForVerifiedPartnerProvisioningState() []string {
	return []string{
		string(VerifiedPartnerProvisioningStateCanceled),
		string(VerifiedPartnerProvisioningStateCreating),
		string(VerifiedPartnerProvisioningStateDeleting),
		string(VerifiedPartnerProvisioningStateFailed),
		string(VerifiedPartnerProvisioningStateSucceeded),
		string(VerifiedPartnerProvisioningStateUpdating),
	}
}

func (s *VerifiedPartnerProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVerifiedPartnerProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVerifiedPartnerProvisioningState(input string) (*VerifiedPartnerProvisioningState, error) {
	vals := map[string]VerifiedPartnerProvisioningState{
		"canceled":  VerifiedPartnerProvisioningStateCanceled,
		"creating":  VerifiedPartnerProvisioningStateCreating,
		"deleting":  VerifiedPartnerProvisioningStateDeleting,
		"failed":    VerifiedPartnerProvisioningStateFailed,
		"succeeded": VerifiedPartnerProvisioningStateSucceeded,
		"updating":  VerifiedPartnerProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VerifiedPartnerProvisioningState(input)
	return &out, nil
}
