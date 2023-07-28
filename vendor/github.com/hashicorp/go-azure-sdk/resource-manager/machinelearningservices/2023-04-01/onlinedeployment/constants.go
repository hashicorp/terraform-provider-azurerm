package onlinedeployment

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerType string

const (
	ContainerTypeInferenceServer    ContainerType = "InferenceServer"
	ContainerTypeStorageInitializer ContainerType = "StorageInitializer"
)

func PossibleValuesForContainerType() []string {
	return []string{
		string(ContainerTypeInferenceServer),
		string(ContainerTypeStorageInitializer),
	}
}

func (s *ContainerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerType(input string) (*ContainerType, error) {
	vals := map[string]ContainerType{
		"inferenceserver":    ContainerTypeInferenceServer,
		"storageinitializer": ContainerTypeStorageInitializer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerType(input)
	return &out, nil
}

type DeploymentProvisioningState string

const (
	DeploymentProvisioningStateCanceled  DeploymentProvisioningState = "Canceled"
	DeploymentProvisioningStateCreating  DeploymentProvisioningState = "Creating"
	DeploymentProvisioningStateDeleting  DeploymentProvisioningState = "Deleting"
	DeploymentProvisioningStateFailed    DeploymentProvisioningState = "Failed"
	DeploymentProvisioningStateScaling   DeploymentProvisioningState = "Scaling"
	DeploymentProvisioningStateSucceeded DeploymentProvisioningState = "Succeeded"
	DeploymentProvisioningStateUpdating  DeploymentProvisioningState = "Updating"
)

func PossibleValuesForDeploymentProvisioningState() []string {
	return []string{
		string(DeploymentProvisioningStateCanceled),
		string(DeploymentProvisioningStateCreating),
		string(DeploymentProvisioningStateDeleting),
		string(DeploymentProvisioningStateFailed),
		string(DeploymentProvisioningStateScaling),
		string(DeploymentProvisioningStateSucceeded),
		string(DeploymentProvisioningStateUpdating),
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
		"canceled":  DeploymentProvisioningStateCanceled,
		"creating":  DeploymentProvisioningStateCreating,
		"deleting":  DeploymentProvisioningStateDeleting,
		"failed":    DeploymentProvisioningStateFailed,
		"scaling":   DeploymentProvisioningStateScaling,
		"succeeded": DeploymentProvisioningStateSucceeded,
		"updating":  DeploymentProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentProvisioningState(input)
	return &out, nil
}

type EgressPublicNetworkAccessType string

const (
	EgressPublicNetworkAccessTypeDisabled EgressPublicNetworkAccessType = "Disabled"
	EgressPublicNetworkAccessTypeEnabled  EgressPublicNetworkAccessType = "Enabled"
)

func PossibleValuesForEgressPublicNetworkAccessType() []string {
	return []string{
		string(EgressPublicNetworkAccessTypeDisabled),
		string(EgressPublicNetworkAccessTypeEnabled),
	}
}

func (s *EgressPublicNetworkAccessType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEgressPublicNetworkAccessType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEgressPublicNetworkAccessType(input string) (*EgressPublicNetworkAccessType, error) {
	vals := map[string]EgressPublicNetworkAccessType{
		"disabled": EgressPublicNetworkAccessTypeDisabled,
		"enabled":  EgressPublicNetworkAccessTypeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EgressPublicNetworkAccessType(input)
	return &out, nil
}

type EndpointComputeType string

const (
	EndpointComputeTypeAzureMLCompute EndpointComputeType = "AzureMLCompute"
	EndpointComputeTypeKubernetes     EndpointComputeType = "Kubernetes"
	EndpointComputeTypeManaged        EndpointComputeType = "Managed"
)

func PossibleValuesForEndpointComputeType() []string {
	return []string{
		string(EndpointComputeTypeAzureMLCompute),
		string(EndpointComputeTypeKubernetes),
		string(EndpointComputeTypeManaged),
	}
}

func (s *EndpointComputeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEndpointComputeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEndpointComputeType(input string) (*EndpointComputeType, error) {
	vals := map[string]EndpointComputeType{
		"azuremlcompute": EndpointComputeTypeAzureMLCompute,
		"kubernetes":     EndpointComputeTypeKubernetes,
		"managed":        EndpointComputeTypeManaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointComputeType(input)
	return &out, nil
}

type ScaleType string

const (
	ScaleTypeDefault           ScaleType = "Default"
	ScaleTypeTargetUtilization ScaleType = "TargetUtilization"
)

func PossibleValuesForScaleType() []string {
	return []string{
		string(ScaleTypeDefault),
		string(ScaleTypeTargetUtilization),
	}
}

func (s *ScaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScaleType(input string) (*ScaleType, error) {
	vals := map[string]ScaleType{
		"default":           ScaleTypeDefault,
		"targetutilization": ScaleTypeTargetUtilization,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScaleType(input)
	return &out, nil
}

type SkuScaleType string

const (
	SkuScaleTypeAutomatic SkuScaleType = "Automatic"
	SkuScaleTypeManual    SkuScaleType = "Manual"
	SkuScaleTypeNone      SkuScaleType = "None"
)

func PossibleValuesForSkuScaleType() []string {
	return []string{
		string(SkuScaleTypeAutomatic),
		string(SkuScaleTypeManual),
		string(SkuScaleTypeNone),
	}
}

func (s *SkuScaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuScaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuScaleType(input string) (*SkuScaleType, error) {
	vals := map[string]SkuScaleType{
		"automatic": SkuScaleTypeAutomatic,
		"manual":    SkuScaleTypeManual,
		"none":      SkuScaleTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuScaleType(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierFree     SkuTier = "Free"
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
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
		"basic":    SkuTierBasic,
		"free":     SkuTierFree,
		"premium":  SkuTierPremium,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}
