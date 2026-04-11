package deployments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentModelVersionUpgradeOption string

const (
	DeploymentModelVersionUpgradeOptionNoAutoUpgrade                  DeploymentModelVersionUpgradeOption = "NoAutoUpgrade"
	DeploymentModelVersionUpgradeOptionOnceCurrentVersionExpired      DeploymentModelVersionUpgradeOption = "OnceCurrentVersionExpired"
	DeploymentModelVersionUpgradeOptionOnceNewDefaultVersionAvailable DeploymentModelVersionUpgradeOption = "OnceNewDefaultVersionAvailable"
)

func PossibleValuesForDeploymentModelVersionUpgradeOption() []string {
	return []string{
		string(DeploymentModelVersionUpgradeOptionNoAutoUpgrade),
		string(DeploymentModelVersionUpgradeOptionOnceCurrentVersionExpired),
		string(DeploymentModelVersionUpgradeOptionOnceNewDefaultVersionAvailable),
	}
}

func (s *DeploymentModelVersionUpgradeOption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentModelVersionUpgradeOption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentModelVersionUpgradeOption(input string) (*DeploymentModelVersionUpgradeOption, error) {
	vals := map[string]DeploymentModelVersionUpgradeOption{
		"noautoupgrade":                  DeploymentModelVersionUpgradeOptionNoAutoUpgrade,
		"oncecurrentversionexpired":      DeploymentModelVersionUpgradeOptionOnceCurrentVersionExpired,
		"oncenewdefaultversionavailable": DeploymentModelVersionUpgradeOptionOnceNewDefaultVersionAvailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentModelVersionUpgradeOption(input)
	return &out, nil
}

type DeploymentProvisioningState string

const (
	DeploymentProvisioningStateAccepted  DeploymentProvisioningState = "Accepted"
	DeploymentProvisioningStateCanceled  DeploymentProvisioningState = "Canceled"
	DeploymentProvisioningStateCreating  DeploymentProvisioningState = "Creating"
	DeploymentProvisioningStateDeleting  DeploymentProvisioningState = "Deleting"
	DeploymentProvisioningStateDisabled  DeploymentProvisioningState = "Disabled"
	DeploymentProvisioningStateFailed    DeploymentProvisioningState = "Failed"
	DeploymentProvisioningStateMoving    DeploymentProvisioningState = "Moving"
	DeploymentProvisioningStateSucceeded DeploymentProvisioningState = "Succeeded"
)

func PossibleValuesForDeploymentProvisioningState() []string {
	return []string{
		string(DeploymentProvisioningStateAccepted),
		string(DeploymentProvisioningStateCanceled),
		string(DeploymentProvisioningStateCreating),
		string(DeploymentProvisioningStateDeleting),
		string(DeploymentProvisioningStateDisabled),
		string(DeploymentProvisioningStateFailed),
		string(DeploymentProvisioningStateMoving),
		string(DeploymentProvisioningStateSucceeded),
	}
}

func (s *DeploymentProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentProvisioningState(input string) (*DeploymentProvisioningState, error) {
	vals := map[string]DeploymentProvisioningState{
		"accepted":  DeploymentProvisioningStateAccepted,
		"canceled":  DeploymentProvisioningStateCanceled,
		"creating":  DeploymentProvisioningStateCreating,
		"deleting":  DeploymentProvisioningStateDeleting,
		"disabled":  DeploymentProvisioningStateDisabled,
		"failed":    DeploymentProvisioningStateFailed,
		"moving":    DeploymentProvisioningStateMoving,
		"succeeded": DeploymentProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentProvisioningState(input)
	return &out, nil
}

type DeploymentScaleType string

const (
	DeploymentScaleTypeManual   DeploymentScaleType = "Manual"
	DeploymentScaleTypeStandard DeploymentScaleType = "Standard"
)

func PossibleValuesForDeploymentScaleType() []string {
	return []string{
		string(DeploymentScaleTypeManual),
		string(DeploymentScaleTypeStandard),
	}
}

func (s *DeploymentScaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentScaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentScaleType(input string) (*DeploymentScaleType, error) {
	vals := map[string]DeploymentScaleType{
		"manual":   DeploymentScaleTypeManual,
		"standard": DeploymentScaleTypeStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentScaleType(input)
	return &out, nil
}

type DeploymentState string

const (
	DeploymentStatePaused  DeploymentState = "Paused"
	DeploymentStateRunning DeploymentState = "Running"
)

func PossibleValuesForDeploymentState() []string {
	return []string{
		string(DeploymentStatePaused),
		string(DeploymentStateRunning),
	}
}

func (s *DeploymentState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentState(input string) (*DeploymentState, error) {
	vals := map[string]DeploymentState{
		"paused":  DeploymentStatePaused,
		"running": DeploymentStateRunning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentState(input)
	return &out, nil
}

type RoutingMode string

const (
	RoutingModeAccuracy RoutingMode = "accuracy"
	RoutingModeBalanced RoutingMode = "balanced"
	RoutingModeCost     RoutingMode = "cost"
)

func PossibleValuesForRoutingMode() []string {
	return []string{
		string(RoutingModeAccuracy),
		string(RoutingModeBalanced),
		string(RoutingModeCost),
	}
}

func (s *RoutingMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoutingMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoutingMode(input string) (*RoutingMode, error) {
	vals := map[string]RoutingMode{
		"accuracy": RoutingModeAccuracy,
		"balanced": RoutingModeBalanced,
		"cost":     RoutingModeCost,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoutingMode(input)
	return &out, nil
}

type ServiceTier string

const (
	ServiceTierDefault  ServiceTier = "Default"
	ServiceTierPriority ServiceTier = "Priority"
)

func PossibleValuesForServiceTier() []string {
	return []string{
		string(ServiceTierDefault),
		string(ServiceTierPriority),
	}
}

func (s *ServiceTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServiceTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServiceTier(input string) (*ServiceTier, error) {
	vals := map[string]ServiceTier{
		"default":  ServiceTierDefault,
		"priority": ServiceTierPriority,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceTier(input)
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
