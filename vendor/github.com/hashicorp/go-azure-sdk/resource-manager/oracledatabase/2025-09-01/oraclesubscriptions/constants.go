package oraclesubscriptions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddSubscriptionOperationState string

const (
	AddSubscriptionOperationStateFailed    AddSubscriptionOperationState = "Failed"
	AddSubscriptionOperationStateSucceeded AddSubscriptionOperationState = "Succeeded"
	AddSubscriptionOperationStateUpdating  AddSubscriptionOperationState = "Updating"
)

func PossibleValuesForAddSubscriptionOperationState() []string {
	return []string{
		string(AddSubscriptionOperationStateFailed),
		string(AddSubscriptionOperationStateSucceeded),
		string(AddSubscriptionOperationStateUpdating),
	}
}

func (s *AddSubscriptionOperationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAddSubscriptionOperationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAddSubscriptionOperationState(input string) (*AddSubscriptionOperationState, error) {
	vals := map[string]AddSubscriptionOperationState{
		"failed":    AddSubscriptionOperationStateFailed,
		"succeeded": AddSubscriptionOperationStateSucceeded,
		"updating":  AddSubscriptionOperationStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AddSubscriptionOperationState(input)
	return &out, nil
}

type CloudAccountProvisioningState string

const (
	CloudAccountProvisioningStateAvailable    CloudAccountProvisioningState = "Available"
	CloudAccountProvisioningStatePending      CloudAccountProvisioningState = "Pending"
	CloudAccountProvisioningStateProvisioning CloudAccountProvisioningState = "Provisioning"
)

func PossibleValuesForCloudAccountProvisioningState() []string {
	return []string{
		string(CloudAccountProvisioningStateAvailable),
		string(CloudAccountProvisioningStatePending),
		string(CloudAccountProvisioningStateProvisioning),
	}
}

func (s *CloudAccountProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCloudAccountProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCloudAccountProvisioningState(input string) (*CloudAccountProvisioningState, error) {
	vals := map[string]CloudAccountProvisioningState{
		"available":    CloudAccountProvisioningStateAvailable,
		"pending":      CloudAccountProvisioningStatePending,
		"provisioning": CloudAccountProvisioningStateProvisioning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CloudAccountProvisioningState(input)
	return &out, nil
}

type Intent string

const (
	IntentReset  Intent = "Reset"
	IntentRetain Intent = "Retain"
)

func PossibleValuesForIntent() []string {
	return []string{
		string(IntentReset),
		string(IntentRetain),
	}
}

func (s *Intent) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntent(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntent(input string) (*Intent, error) {
	vals := map[string]Intent{
		"reset":  IntentReset,
		"retain": IntentRetain,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Intent(input)
	return &out, nil
}

type OracleSubscriptionProvisioningState string

const (
	OracleSubscriptionProvisioningStateCanceled  OracleSubscriptionProvisioningState = "Canceled"
	OracleSubscriptionProvisioningStateFailed    OracleSubscriptionProvisioningState = "Failed"
	OracleSubscriptionProvisioningStateSucceeded OracleSubscriptionProvisioningState = "Succeeded"
)

func PossibleValuesForOracleSubscriptionProvisioningState() []string {
	return []string{
		string(OracleSubscriptionProvisioningStateCanceled),
		string(OracleSubscriptionProvisioningStateFailed),
		string(OracleSubscriptionProvisioningStateSucceeded),
	}
}

func (s *OracleSubscriptionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOracleSubscriptionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOracleSubscriptionProvisioningState(input string) (*OracleSubscriptionProvisioningState, error) {
	vals := map[string]OracleSubscriptionProvisioningState{
		"canceled":  OracleSubscriptionProvisioningStateCanceled,
		"failed":    OracleSubscriptionProvisioningStateFailed,
		"succeeded": OracleSubscriptionProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OracleSubscriptionProvisioningState(input)
	return &out, nil
}
