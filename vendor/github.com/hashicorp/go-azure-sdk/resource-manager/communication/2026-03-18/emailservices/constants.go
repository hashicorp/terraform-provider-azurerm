package emailservices

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EmailServicesProvisioningState string

const (
	EmailServicesProvisioningStateCanceled  EmailServicesProvisioningState = "Canceled"
	EmailServicesProvisioningStateCreating  EmailServicesProvisioningState = "Creating"
	EmailServicesProvisioningStateDeleting  EmailServicesProvisioningState = "Deleting"
	EmailServicesProvisioningStateFailed    EmailServicesProvisioningState = "Failed"
	EmailServicesProvisioningStateMoving    EmailServicesProvisioningState = "Moving"
	EmailServicesProvisioningStateRunning   EmailServicesProvisioningState = "Running"
	EmailServicesProvisioningStateSucceeded EmailServicesProvisioningState = "Succeeded"
	EmailServicesProvisioningStateUnknown   EmailServicesProvisioningState = "Unknown"
	EmailServicesProvisioningStateUpdating  EmailServicesProvisioningState = "Updating"
)

func PossibleValuesForEmailServicesProvisioningState() []string {
	return []string{
		string(EmailServicesProvisioningStateCanceled),
		string(EmailServicesProvisioningStateCreating),
		string(EmailServicesProvisioningStateDeleting),
		string(EmailServicesProvisioningStateFailed),
		string(EmailServicesProvisioningStateMoving),
		string(EmailServicesProvisioningStateRunning),
		string(EmailServicesProvisioningStateSucceeded),
		string(EmailServicesProvisioningStateUnknown),
		string(EmailServicesProvisioningStateUpdating),
	}
}

func (s *EmailServicesProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEmailServicesProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEmailServicesProvisioningState(input string) (*EmailServicesProvisioningState, error) {
	vals := map[string]EmailServicesProvisioningState{
		"canceled":  EmailServicesProvisioningStateCanceled,
		"creating":  EmailServicesProvisioningStateCreating,
		"deleting":  EmailServicesProvisioningStateDeleting,
		"failed":    EmailServicesProvisioningStateFailed,
		"moving":    EmailServicesProvisioningStateMoving,
		"running":   EmailServicesProvisioningStateRunning,
		"succeeded": EmailServicesProvisioningStateSucceeded,
		"unknown":   EmailServicesProvisioningStateUnknown,
		"updating":  EmailServicesProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EmailServicesProvisioningState(input)
	return &out, nil
}
