package provisionedclusterinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddonPhase string

const (
	AddonPhaseDeleting                             AddonPhase = "deleting"
	AddonPhaseFailed                               AddonPhase = "failed"
	AddonPhasePending                              AddonPhase = "pending"
	AddonPhaseProvisioned                          AddonPhase = "provisioned"
	AddonPhaseProvisioning                         AddonPhase = "provisioning"
	AddonPhaseProvisioningHelmChartInstalled       AddonPhase = "provisioning {HelmChartInstalled}"
	AddonPhaseProvisioningMSICertificateDownloaded AddonPhase = "provisioning {MSICertificateDownloaded}"
	AddonPhaseUpgrading                            AddonPhase = "upgrading"
)

func PossibleValuesForAddonPhase() []string {
	return []string{
		string(AddonPhaseDeleting),
		string(AddonPhaseFailed),
		string(AddonPhasePending),
		string(AddonPhaseProvisioned),
		string(AddonPhaseProvisioning),
		string(AddonPhaseProvisioningHelmChartInstalled),
		string(AddonPhaseProvisioningMSICertificateDownloaded),
		string(AddonPhaseUpgrading),
	}
}

func (s *AddonPhase) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAddonPhase(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAddonPhase(input string) (*AddonPhase, error) {
	vals := map[string]AddonPhase{
		"deleting":                          AddonPhaseDeleting,
		"failed":                            AddonPhaseFailed,
		"pending":                           AddonPhasePending,
		"provisioned":                       AddonPhaseProvisioned,
		"provisioning":                      AddonPhaseProvisioning,
		"provisioning {helmchartinstalled}": AddonPhaseProvisioningHelmChartInstalled,
		"provisioning {msicertificatedownloaded}": AddonPhaseProvisioningMSICertificateDownloaded,
		"upgrading": AddonPhaseUpgrading,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AddonPhase(input)
	return &out, nil
}

type AzureHybridBenefit string

const (
	AzureHybridBenefitFalse         AzureHybridBenefit = "False"
	AzureHybridBenefitNotApplicable AzureHybridBenefit = "NotApplicable"
	AzureHybridBenefitTrue          AzureHybridBenefit = "True"
)

func PossibleValuesForAzureHybridBenefit() []string {
	return []string{
		string(AzureHybridBenefitFalse),
		string(AzureHybridBenefitNotApplicable),
		string(AzureHybridBenefitTrue),
	}
}

func (s *AzureHybridBenefit) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureHybridBenefit(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureHybridBenefit(input string) (*AzureHybridBenefit, error) {
	vals := map[string]AzureHybridBenefit{
		"false":         AzureHybridBenefitFalse,
		"notapplicable": AzureHybridBenefitNotApplicable,
		"true":          AzureHybridBenefitTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureHybridBenefit(input)
	return &out, nil
}

type Expander string

const (
	ExpanderLeastNegativewaste Expander = "least-waste"
	ExpanderMostNegativepods   Expander = "most-pods"
	ExpanderPriority           Expander = "priority"
	ExpanderRandom             Expander = "random"
)

func PossibleValuesForExpander() []string {
	return []string{
		string(ExpanderLeastNegativewaste),
		string(ExpanderMostNegativepods),
		string(ExpanderPriority),
		string(ExpanderRandom),
	}
}

func (s *Expander) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpander(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpander(input string) (*Expander, error) {
	vals := map[string]Expander{
		"least-waste": ExpanderLeastNegativewaste,
		"most-pods":   ExpanderMostNegativepods,
		"priority":    ExpanderPriority,
		"random":      ExpanderRandom,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Expander(input)
	return &out, nil
}

type ExtendedLocationTypes string

const (
	ExtendedLocationTypesCustomLocation ExtendedLocationTypes = "CustomLocation"
)

func PossibleValuesForExtendedLocationTypes() []string {
	return []string{
		string(ExtendedLocationTypesCustomLocation),
	}
}

func (s *ExtendedLocationTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtendedLocationTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtendedLocationTypes(input string) (*ExtendedLocationTypes, error) {
	vals := map[string]ExtendedLocationTypes{
		"customlocation": ExtendedLocationTypesCustomLocation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtendedLocationTypes(input)
	return &out, nil
}

type NetworkPolicy string

const (
	NetworkPolicyCalico NetworkPolicy = "calico"
)

func PossibleValuesForNetworkPolicy() []string {
	return []string{
		string(NetworkPolicyCalico),
	}
}

func (s *NetworkPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkPolicy(input string) (*NetworkPolicy, error) {
	vals := map[string]NetworkPolicy{
		"calico": NetworkPolicyCalico,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkPolicy(input)
	return &out, nil
}

type OSSKU string

const (
	OSSKUCBLMariner            OSSKU = "CBLMariner"
	OSSKUWindowsTwoZeroOneNine OSSKU = "Windows2019"
	OSSKUWindowsTwoZeroTwoTwo  OSSKU = "Windows2022"
)

func PossibleValuesForOSSKU() []string {
	return []string{
		string(OSSKUCBLMariner),
		string(OSSKUWindowsTwoZeroOneNine),
		string(OSSKUWindowsTwoZeroTwoTwo),
	}
}

func (s *OSSKU) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOSSKU(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOSSKU(input string) (*OSSKU, error) {
	vals := map[string]OSSKU{
		"cblmariner":  OSSKUCBLMariner,
		"windows2019": OSSKUWindowsTwoZeroOneNine,
		"windows2022": OSSKUWindowsTwoZeroTwoTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OSSKU(input)
	return &out, nil
}

type OsType string

const (
	OsTypeLinux   OsType = "Linux"
	OsTypeWindows OsType = "Windows"
)

func PossibleValuesForOsType() []string {
	return []string{
		string(OsTypeLinux),
		string(OsTypeWindows),
	}
}

func (s *OsType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOsType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOsType(input string) (*OsType, error) {
	vals := map[string]OsType{
		"linux":   OsTypeLinux,
		"windows": OsTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OsType(input)
	return &out, nil
}

type ResourceProvisioningState string

const (
	ResourceProvisioningStateAccepted  ResourceProvisioningState = "Accepted"
	ResourceProvisioningStateCanceled  ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateCreating  ResourceProvisioningState = "Creating"
	ResourceProvisioningStateDeleting  ResourceProvisioningState = "Deleting"
	ResourceProvisioningStateFailed    ResourceProvisioningState = "Failed"
	ResourceProvisioningStatePending   ResourceProvisioningState = "Pending"
	ResourceProvisioningStateSucceeded ResourceProvisioningState = "Succeeded"
	ResourceProvisioningStateUpdating  ResourceProvisioningState = "Updating"
	ResourceProvisioningStateUpgrading ResourceProvisioningState = "Upgrading"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateAccepted),
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateCreating),
		string(ResourceProvisioningStateDeleting),
		string(ResourceProvisioningStateFailed),
		string(ResourceProvisioningStatePending),
		string(ResourceProvisioningStateSucceeded),
		string(ResourceProvisioningStateUpdating),
		string(ResourceProvisioningStateUpgrading),
	}
}

func (s *ResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceProvisioningState(input string) (*ResourceProvisioningState, error) {
	vals := map[string]ResourceProvisioningState{
		"accepted":  ResourceProvisioningStateAccepted,
		"canceled":  ResourceProvisioningStateCanceled,
		"creating":  ResourceProvisioningStateCreating,
		"deleting":  ResourceProvisioningStateDeleting,
		"failed":    ResourceProvisioningStateFailed,
		"pending":   ResourceProvisioningStatePending,
		"succeeded": ResourceProvisioningStateSucceeded,
		"updating":  ResourceProvisioningStateUpdating,
		"upgrading": ResourceProvisioningStateUpgrading,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}
