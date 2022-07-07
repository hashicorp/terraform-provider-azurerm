package policyinsights

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceDiscoveryMode string

const (
	ResourceDiscoveryModeExistingNonCompliant ResourceDiscoveryMode = "ExistingNonCompliant"
	ResourceDiscoveryModeReEvaluateCompliance ResourceDiscoveryMode = "ReEvaluateCompliance"
)

func PossibleValuesForResourceDiscoveryMode() []string {
	return []string{
		string(ResourceDiscoveryModeExistingNonCompliant),
		string(ResourceDiscoveryModeReEvaluateCompliance),
	}
}

func parseResourceDiscoveryMode(input string) (*ResourceDiscoveryMode, error) {
	vals := map[string]ResourceDiscoveryMode{
		"existingnoncompliant": ResourceDiscoveryModeExistingNonCompliant,
		"reevaluatecompliance": ResourceDiscoveryModeReEvaluateCompliance,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceDiscoveryMode(input)
	return &out, nil
}
