package testlines

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type TestLinePurpose string

const (
	TestLinePurposeAutomated TestLinePurpose = "Automated"
	TestLinePurposeManual    TestLinePurpose = "Manual"
)

func PossibleValuesForTestLinePurpose() []string {
	return []string{
		string(TestLinePurposeAutomated),
		string(TestLinePurposeManual),
	}
}

func parseTestLinePurpose(input string) (*TestLinePurpose, error) {
	vals := map[string]TestLinePurpose{
		"automated": TestLinePurposeAutomated,
		"manual":    TestLinePurposeManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TestLinePurpose(input)
	return &out, nil
}
