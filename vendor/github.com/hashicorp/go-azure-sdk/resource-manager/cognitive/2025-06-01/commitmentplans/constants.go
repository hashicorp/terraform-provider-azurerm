package commitmentplans

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommitmentPlanProvisioningState string

const (
	CommitmentPlanProvisioningStateAccepted  CommitmentPlanProvisioningState = "Accepted"
	CommitmentPlanProvisioningStateCanceled  CommitmentPlanProvisioningState = "Canceled"
	CommitmentPlanProvisioningStateCreating  CommitmentPlanProvisioningState = "Creating"
	CommitmentPlanProvisioningStateDeleting  CommitmentPlanProvisioningState = "Deleting"
	CommitmentPlanProvisioningStateFailed    CommitmentPlanProvisioningState = "Failed"
	CommitmentPlanProvisioningStateMoving    CommitmentPlanProvisioningState = "Moving"
	CommitmentPlanProvisioningStateSucceeded CommitmentPlanProvisioningState = "Succeeded"
)

func PossibleValuesForCommitmentPlanProvisioningState() []string {
	return []string{
		string(CommitmentPlanProvisioningStateAccepted),
		string(CommitmentPlanProvisioningStateCanceled),
		string(CommitmentPlanProvisioningStateCreating),
		string(CommitmentPlanProvisioningStateDeleting),
		string(CommitmentPlanProvisioningStateFailed),
		string(CommitmentPlanProvisioningStateMoving),
		string(CommitmentPlanProvisioningStateSucceeded),
	}
}

func (s *CommitmentPlanProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCommitmentPlanProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCommitmentPlanProvisioningState(input string) (*CommitmentPlanProvisioningState, error) {
	vals := map[string]CommitmentPlanProvisioningState{
		"accepted":  CommitmentPlanProvisioningStateAccepted,
		"canceled":  CommitmentPlanProvisioningStateCanceled,
		"creating":  CommitmentPlanProvisioningStateCreating,
		"deleting":  CommitmentPlanProvisioningStateDeleting,
		"failed":    CommitmentPlanProvisioningStateFailed,
		"moving":    CommitmentPlanProvisioningStateMoving,
		"succeeded": CommitmentPlanProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CommitmentPlanProvisioningState(input)
	return &out, nil
}

type HostingModel string

const (
	HostingModelConnectedContainer    HostingModel = "ConnectedContainer"
	HostingModelDisconnectedContainer HostingModel = "DisconnectedContainer"
	HostingModelProvisionedWeb        HostingModel = "ProvisionedWeb"
	HostingModelWeb                   HostingModel = "Web"
)

func PossibleValuesForHostingModel() []string {
	return []string{
		string(HostingModelConnectedContainer),
		string(HostingModelDisconnectedContainer),
		string(HostingModelProvisionedWeb),
		string(HostingModelWeb),
	}
}

func (s *HostingModel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHostingModel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHostingModel(input string) (*HostingModel, error) {
	vals := map[string]HostingModel{
		"connectedcontainer":    HostingModelConnectedContainer,
		"disconnectedcontainer": HostingModelDisconnectedContainer,
		"provisionedweb":        HostingModelProvisionedWeb,
		"web":                   HostingModelWeb,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HostingModel(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic      SkuTier = "Basic"
	SkuTierEnterprise SkuTier = "Enterprise"
	SkuTierFree       SkuTier = "Free"
	SkuTierPremium    SkuTier = "Premium"
	SkuTierStandard   SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierEnterprise),
		string(SkuTierFree),
		string(SkuTierPremium),
		string(SkuTierStandard),
	}
}

func (s *SkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":      SkuTierBasic,
		"enterprise": SkuTierEnterprise,
		"free":       SkuTierFree,
		"premium":    SkuTierPremium,
		"standard":   SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}
