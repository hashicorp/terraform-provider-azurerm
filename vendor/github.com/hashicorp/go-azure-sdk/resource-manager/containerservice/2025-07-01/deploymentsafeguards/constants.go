package deploymentsafeguards

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentSafeguardsLevel string

const (
	DeploymentSafeguardsLevelEnforce DeploymentSafeguardsLevel = "Enforce"
	DeploymentSafeguardsLevelWarn    DeploymentSafeguardsLevel = "Warn"
)

func PossibleValuesForDeploymentSafeguardsLevel() []string {
	return []string{
		string(DeploymentSafeguardsLevelEnforce),
		string(DeploymentSafeguardsLevelWarn),
	}
}

func (s *DeploymentSafeguardsLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentSafeguardsLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentSafeguardsLevel(input string) (*DeploymentSafeguardsLevel, error) {
	vals := map[string]DeploymentSafeguardsLevel{
		"enforce": DeploymentSafeguardsLevelEnforce,
		"warn":    DeploymentSafeguardsLevelWarn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentSafeguardsLevel(input)
	return &out, nil
}

type PodSecurityStandardsLevel string

const (
	PodSecurityStandardsLevelBaseline   PodSecurityStandardsLevel = "Baseline"
	PodSecurityStandardsLevelPrivileged PodSecurityStandardsLevel = "Privileged"
	PodSecurityStandardsLevelRestricted PodSecurityStandardsLevel = "Restricted"
)

func PossibleValuesForPodSecurityStandardsLevel() []string {
	return []string{
		string(PodSecurityStandardsLevelBaseline),
		string(PodSecurityStandardsLevelPrivileged),
		string(PodSecurityStandardsLevelRestricted),
	}
}

func (s *PodSecurityStandardsLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePodSecurityStandardsLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePodSecurityStandardsLevel(input string) (*PodSecurityStandardsLevel, error) {
	vals := map[string]PodSecurityStandardsLevel{
		"baseline":   PodSecurityStandardsLevelBaseline,
		"privileged": PodSecurityStandardsLevelPrivileged,
		"restricted": PodSecurityStandardsLevelRestricted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PodSecurityStandardsLevel(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
