package cacertificates

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CaCertificateProvisioningState string

const (
	CaCertificateProvisioningStateCanceled  CaCertificateProvisioningState = "Canceled"
	CaCertificateProvisioningStateCreating  CaCertificateProvisioningState = "Creating"
	CaCertificateProvisioningStateDeleted   CaCertificateProvisioningState = "Deleted"
	CaCertificateProvisioningStateDeleting  CaCertificateProvisioningState = "Deleting"
	CaCertificateProvisioningStateFailed    CaCertificateProvisioningState = "Failed"
	CaCertificateProvisioningStateSucceeded CaCertificateProvisioningState = "Succeeded"
	CaCertificateProvisioningStateUpdating  CaCertificateProvisioningState = "Updating"
)

func PossibleValuesForCaCertificateProvisioningState() []string {
	return []string{
		string(CaCertificateProvisioningStateCanceled),
		string(CaCertificateProvisioningStateCreating),
		string(CaCertificateProvisioningStateDeleted),
		string(CaCertificateProvisioningStateDeleting),
		string(CaCertificateProvisioningStateFailed),
		string(CaCertificateProvisioningStateSucceeded),
		string(CaCertificateProvisioningStateUpdating),
	}
}

func (s *CaCertificateProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCaCertificateProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCaCertificateProvisioningState(input string) (*CaCertificateProvisioningState, error) {
	vals := map[string]CaCertificateProvisioningState{
		"canceled":  CaCertificateProvisioningStateCanceled,
		"creating":  CaCertificateProvisioningStateCreating,
		"deleted":   CaCertificateProvisioningStateDeleted,
		"deleting":  CaCertificateProvisioningStateDeleting,
		"failed":    CaCertificateProvisioningStateFailed,
		"succeeded": CaCertificateProvisioningStateSucceeded,
		"updating":  CaCertificateProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CaCertificateProvisioningState(input)
	return &out, nil
}
