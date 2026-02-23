package containerinstance

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFileShareAccessTier string

const (
	AzureFileShareAccessTierCool                 AzureFileShareAccessTier = "Cool"
	AzureFileShareAccessTierHot                  AzureFileShareAccessTier = "Hot"
	AzureFileShareAccessTierPremium              AzureFileShareAccessTier = "Premium"
	AzureFileShareAccessTierTransactionOptimized AzureFileShareAccessTier = "TransactionOptimized"
)

func PossibleValuesForAzureFileShareAccessTier() []string {
	return []string{
		string(AzureFileShareAccessTierCool),
		string(AzureFileShareAccessTierHot),
		string(AzureFileShareAccessTierPremium),
		string(AzureFileShareAccessTierTransactionOptimized),
	}
}

func (s *AzureFileShareAccessTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureFileShareAccessTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureFileShareAccessTier(input string) (*AzureFileShareAccessTier, error) {
	vals := map[string]AzureFileShareAccessTier{
		"cool":                 AzureFileShareAccessTierCool,
		"hot":                  AzureFileShareAccessTierHot,
		"premium":              AzureFileShareAccessTierPremium,
		"transactionoptimized": AzureFileShareAccessTierTransactionOptimized,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFileShareAccessTier(input)
	return &out, nil
}

type AzureFileShareAccessType string

const (
	AzureFileShareAccessTypeExclusive AzureFileShareAccessType = "Exclusive"
	AzureFileShareAccessTypeShared    AzureFileShareAccessType = "Shared"
)

func PossibleValuesForAzureFileShareAccessType() []string {
	return []string{
		string(AzureFileShareAccessTypeExclusive),
		string(AzureFileShareAccessTypeShared),
	}
}

func (s *AzureFileShareAccessType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureFileShareAccessType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureFileShareAccessType(input string) (*AzureFileShareAccessType, error) {
	vals := map[string]AzureFileShareAccessType{
		"exclusive": AzureFileShareAccessTypeExclusive,
		"shared":    AzureFileShareAccessTypeShared,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFileShareAccessType(input)
	return &out, nil
}

type ContainerGroupIPAddressType string

const (
	ContainerGroupIPAddressTypePrivate ContainerGroupIPAddressType = "Private"
	ContainerGroupIPAddressTypePublic  ContainerGroupIPAddressType = "Public"
)

func PossibleValuesForContainerGroupIPAddressType() []string {
	return []string{
		string(ContainerGroupIPAddressTypePrivate),
		string(ContainerGroupIPAddressTypePublic),
	}
}

func (s *ContainerGroupIPAddressType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerGroupIPAddressType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerGroupIPAddressType(input string) (*ContainerGroupIPAddressType, error) {
	vals := map[string]ContainerGroupIPAddressType{
		"private": ContainerGroupIPAddressTypePrivate,
		"public":  ContainerGroupIPAddressTypePublic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerGroupIPAddressType(input)
	return &out, nil
}

type ContainerGroupNetworkProtocol string

const (
	ContainerGroupNetworkProtocolTCP ContainerGroupNetworkProtocol = "TCP"
	ContainerGroupNetworkProtocolUDP ContainerGroupNetworkProtocol = "UDP"
)

func PossibleValuesForContainerGroupNetworkProtocol() []string {
	return []string{
		string(ContainerGroupNetworkProtocolTCP),
		string(ContainerGroupNetworkProtocolUDP),
	}
}

func (s *ContainerGroupNetworkProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerGroupNetworkProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerGroupNetworkProtocol(input string) (*ContainerGroupNetworkProtocol, error) {
	vals := map[string]ContainerGroupNetworkProtocol{
		"tcp": ContainerGroupNetworkProtocolTCP,
		"udp": ContainerGroupNetworkProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerGroupNetworkProtocol(input)
	return &out, nil
}

type ContainerGroupPriority string

const (
	ContainerGroupPriorityRegular ContainerGroupPriority = "Regular"
	ContainerGroupPrioritySpot    ContainerGroupPriority = "Spot"
)

func PossibleValuesForContainerGroupPriority() []string {
	return []string{
		string(ContainerGroupPriorityRegular),
		string(ContainerGroupPrioritySpot),
	}
}

func (s *ContainerGroupPriority) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerGroupPriority(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerGroupPriority(input string) (*ContainerGroupPriority, error) {
	vals := map[string]ContainerGroupPriority{
		"regular": ContainerGroupPriorityRegular,
		"spot":    ContainerGroupPrioritySpot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerGroupPriority(input)
	return &out, nil
}

type ContainerGroupProvisioningState string

const (
	ContainerGroupProvisioningStateAccepted       ContainerGroupProvisioningState = "Accepted"
	ContainerGroupProvisioningStateCanceled       ContainerGroupProvisioningState = "Canceled"
	ContainerGroupProvisioningStateCreating       ContainerGroupProvisioningState = "Creating"
	ContainerGroupProvisioningStateDeleting       ContainerGroupProvisioningState = "Deleting"
	ContainerGroupProvisioningStateFailed         ContainerGroupProvisioningState = "Failed"
	ContainerGroupProvisioningStateNotAccessible  ContainerGroupProvisioningState = "NotAccessible"
	ContainerGroupProvisioningStateNotSpecified   ContainerGroupProvisioningState = "NotSpecified"
	ContainerGroupProvisioningStatePending        ContainerGroupProvisioningState = "Pending"
	ContainerGroupProvisioningStatePreProvisioned ContainerGroupProvisioningState = "PreProvisioned"
	ContainerGroupProvisioningStateRepairing      ContainerGroupProvisioningState = "Repairing"
	ContainerGroupProvisioningStateSucceeded      ContainerGroupProvisioningState = "Succeeded"
	ContainerGroupProvisioningStateUnhealthy      ContainerGroupProvisioningState = "Unhealthy"
	ContainerGroupProvisioningStateUpdating       ContainerGroupProvisioningState = "Updating"
)

func PossibleValuesForContainerGroupProvisioningState() []string {
	return []string{
		string(ContainerGroupProvisioningStateAccepted),
		string(ContainerGroupProvisioningStateCanceled),
		string(ContainerGroupProvisioningStateCreating),
		string(ContainerGroupProvisioningStateDeleting),
		string(ContainerGroupProvisioningStateFailed),
		string(ContainerGroupProvisioningStateNotAccessible),
		string(ContainerGroupProvisioningStateNotSpecified),
		string(ContainerGroupProvisioningStatePending),
		string(ContainerGroupProvisioningStatePreProvisioned),
		string(ContainerGroupProvisioningStateRepairing),
		string(ContainerGroupProvisioningStateSucceeded),
		string(ContainerGroupProvisioningStateUnhealthy),
		string(ContainerGroupProvisioningStateUpdating),
	}
}

func (s *ContainerGroupProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerGroupProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerGroupProvisioningState(input string) (*ContainerGroupProvisioningState, error) {
	vals := map[string]ContainerGroupProvisioningState{
		"accepted":       ContainerGroupProvisioningStateAccepted,
		"canceled":       ContainerGroupProvisioningStateCanceled,
		"creating":       ContainerGroupProvisioningStateCreating,
		"deleting":       ContainerGroupProvisioningStateDeleting,
		"failed":         ContainerGroupProvisioningStateFailed,
		"notaccessible":  ContainerGroupProvisioningStateNotAccessible,
		"notspecified":   ContainerGroupProvisioningStateNotSpecified,
		"pending":        ContainerGroupProvisioningStatePending,
		"preprovisioned": ContainerGroupProvisioningStatePreProvisioned,
		"repairing":      ContainerGroupProvisioningStateRepairing,
		"succeeded":      ContainerGroupProvisioningStateSucceeded,
		"unhealthy":      ContainerGroupProvisioningStateUnhealthy,
		"updating":       ContainerGroupProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerGroupProvisioningState(input)
	return &out, nil
}

type ContainerGroupRestartPolicy string

const (
	ContainerGroupRestartPolicyAlways    ContainerGroupRestartPolicy = "Always"
	ContainerGroupRestartPolicyNever     ContainerGroupRestartPolicy = "Never"
	ContainerGroupRestartPolicyOnFailure ContainerGroupRestartPolicy = "OnFailure"
)

func PossibleValuesForContainerGroupRestartPolicy() []string {
	return []string{
		string(ContainerGroupRestartPolicyAlways),
		string(ContainerGroupRestartPolicyNever),
		string(ContainerGroupRestartPolicyOnFailure),
	}
}

func (s *ContainerGroupRestartPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerGroupRestartPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerGroupRestartPolicy(input string) (*ContainerGroupRestartPolicy, error) {
	vals := map[string]ContainerGroupRestartPolicy{
		"always":    ContainerGroupRestartPolicyAlways,
		"never":     ContainerGroupRestartPolicyNever,
		"onfailure": ContainerGroupRestartPolicyOnFailure,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerGroupRestartPolicy(input)
	return &out, nil
}

type ContainerGroupSku string

const (
	ContainerGroupSkuConfidential ContainerGroupSku = "Confidential"
	ContainerGroupSkuDedicated    ContainerGroupSku = "Dedicated"
	ContainerGroupSkuNotSpecified ContainerGroupSku = "NotSpecified"
	ContainerGroupSkuStandard     ContainerGroupSku = "Standard"
)

func PossibleValuesForContainerGroupSku() []string {
	return []string{
		string(ContainerGroupSkuConfidential),
		string(ContainerGroupSkuDedicated),
		string(ContainerGroupSkuNotSpecified),
		string(ContainerGroupSkuStandard),
	}
}

func (s *ContainerGroupSku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerGroupSku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerGroupSku(input string) (*ContainerGroupSku, error) {
	vals := map[string]ContainerGroupSku{
		"confidential": ContainerGroupSkuConfidential,
		"dedicated":    ContainerGroupSkuDedicated,
		"notspecified": ContainerGroupSkuNotSpecified,
		"standard":     ContainerGroupSkuStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerGroupSku(input)
	return &out, nil
}

type ContainerNetworkProtocol string

const (
	ContainerNetworkProtocolTCP ContainerNetworkProtocol = "TCP"
	ContainerNetworkProtocolUDP ContainerNetworkProtocol = "UDP"
)

func PossibleValuesForContainerNetworkProtocol() []string {
	return []string{
		string(ContainerNetworkProtocolTCP),
		string(ContainerNetworkProtocolUDP),
	}
}

func (s *ContainerNetworkProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerNetworkProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerNetworkProtocol(input string) (*ContainerNetworkProtocol, error) {
	vals := map[string]ContainerNetworkProtocol{
		"tcp": ContainerNetworkProtocolTCP,
		"udp": ContainerNetworkProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerNetworkProtocol(input)
	return &out, nil
}

type DnsNameLabelReusePolicy string

const (
	DnsNameLabelReusePolicyNoreuse            DnsNameLabelReusePolicy = "Noreuse"
	DnsNameLabelReusePolicyResourceGroupReuse DnsNameLabelReusePolicy = "ResourceGroupReuse"
	DnsNameLabelReusePolicySubscriptionReuse  DnsNameLabelReusePolicy = "SubscriptionReuse"
	DnsNameLabelReusePolicyTenantReuse        DnsNameLabelReusePolicy = "TenantReuse"
	DnsNameLabelReusePolicyUnsecure           DnsNameLabelReusePolicy = "Unsecure"
)

func PossibleValuesForDnsNameLabelReusePolicy() []string {
	return []string{
		string(DnsNameLabelReusePolicyNoreuse),
		string(DnsNameLabelReusePolicyResourceGroupReuse),
		string(DnsNameLabelReusePolicySubscriptionReuse),
		string(DnsNameLabelReusePolicyTenantReuse),
		string(DnsNameLabelReusePolicyUnsecure),
	}
}

func (s *DnsNameLabelReusePolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDnsNameLabelReusePolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDnsNameLabelReusePolicy(input string) (*DnsNameLabelReusePolicy, error) {
	vals := map[string]DnsNameLabelReusePolicy{
		"noreuse":            DnsNameLabelReusePolicyNoreuse,
		"resourcegroupreuse": DnsNameLabelReusePolicyResourceGroupReuse,
		"subscriptionreuse":  DnsNameLabelReusePolicySubscriptionReuse,
		"tenantreuse":        DnsNameLabelReusePolicyTenantReuse,
		"unsecure":           DnsNameLabelReusePolicyUnsecure,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DnsNameLabelReusePolicy(input)
	return &out, nil
}

type GpuSku string

const (
	GpuSkuKEightZero  GpuSku = "K80"
	GpuSkuPOneHundred GpuSku = "P100"
	GpuSkuVOneHundred GpuSku = "V100"
)

func PossibleValuesForGpuSku() []string {
	return []string{
		string(GpuSkuKEightZero),
		string(GpuSkuPOneHundred),
		string(GpuSkuVOneHundred),
	}
}

func (s *GpuSku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGpuSku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGpuSku(input string) (*GpuSku, error) {
	vals := map[string]GpuSku{
		"k80":  GpuSkuKEightZero,
		"p100": GpuSkuPOneHundred,
		"v100": GpuSkuVOneHundred,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GpuSku(input)
	return &out, nil
}

type IdentityAccessLevel string

const (
	IdentityAccessLevelAll    IdentityAccessLevel = "All"
	IdentityAccessLevelSystem IdentityAccessLevel = "System"
	IdentityAccessLevelUser   IdentityAccessLevel = "User"
)

func PossibleValuesForIdentityAccessLevel() []string {
	return []string{
		string(IdentityAccessLevelAll),
		string(IdentityAccessLevelSystem),
		string(IdentityAccessLevelUser),
	}
}

func (s *IdentityAccessLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIdentityAccessLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIdentityAccessLevel(input string) (*IdentityAccessLevel, error) {
	vals := map[string]IdentityAccessLevel{
		"all":    IdentityAccessLevelAll,
		"system": IdentityAccessLevelSystem,
		"user":   IdentityAccessLevelUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentityAccessLevel(input)
	return &out, nil
}

type LogAnalyticsLogType string

const (
	LogAnalyticsLogTypeContainerInsights     LogAnalyticsLogType = "ContainerInsights"
	LogAnalyticsLogTypeContainerInstanceLogs LogAnalyticsLogType = "ContainerInstanceLogs"
)

func PossibleValuesForLogAnalyticsLogType() []string {
	return []string{
		string(LogAnalyticsLogTypeContainerInsights),
		string(LogAnalyticsLogTypeContainerInstanceLogs),
	}
}

func (s *LogAnalyticsLogType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLogAnalyticsLogType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLogAnalyticsLogType(input string) (*LogAnalyticsLogType, error) {
	vals := map[string]LogAnalyticsLogType{
		"containerinsights":     LogAnalyticsLogTypeContainerInsights,
		"containerinstancelogs": LogAnalyticsLogTypeContainerInstanceLogs,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LogAnalyticsLogType(input)
	return &out, nil
}

type NGroupProvisioningState string

const (
	NGroupProvisioningStateCanceled  NGroupProvisioningState = "Canceled"
	NGroupProvisioningStateCreating  NGroupProvisioningState = "Creating"
	NGroupProvisioningStateDeleting  NGroupProvisioningState = "Deleting"
	NGroupProvisioningStateFailed    NGroupProvisioningState = "Failed"
	NGroupProvisioningStateMigrating NGroupProvisioningState = "Migrating"
	NGroupProvisioningStateSucceeded NGroupProvisioningState = "Succeeded"
	NGroupProvisioningStateUpdating  NGroupProvisioningState = "Updating"
)

func PossibleValuesForNGroupProvisioningState() []string {
	return []string{
		string(NGroupProvisioningStateCanceled),
		string(NGroupProvisioningStateCreating),
		string(NGroupProvisioningStateDeleting),
		string(NGroupProvisioningStateFailed),
		string(NGroupProvisioningStateMigrating),
		string(NGroupProvisioningStateSucceeded),
		string(NGroupProvisioningStateUpdating),
	}
}

func (s *NGroupProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNGroupProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNGroupProvisioningState(input string) (*NGroupProvisioningState, error) {
	vals := map[string]NGroupProvisioningState{
		"canceled":  NGroupProvisioningStateCanceled,
		"creating":  NGroupProvisioningStateCreating,
		"deleting":  NGroupProvisioningStateDeleting,
		"failed":    NGroupProvisioningStateFailed,
		"migrating": NGroupProvisioningStateMigrating,
		"succeeded": NGroupProvisioningStateSucceeded,
		"updating":  NGroupProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NGroupProvisioningState(input)
	return &out, nil
}

type NGroupUpdateMode string

const (
	NGroupUpdateModeManual  NGroupUpdateMode = "Manual"
	NGroupUpdateModeRolling NGroupUpdateMode = "Rolling"
)

func PossibleValuesForNGroupUpdateMode() []string {
	return []string{
		string(NGroupUpdateModeManual),
		string(NGroupUpdateModeRolling),
	}
}

func (s *NGroupUpdateMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNGroupUpdateMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNGroupUpdateMode(input string) (*NGroupUpdateMode, error) {
	vals := map[string]NGroupUpdateMode{
		"manual":  NGroupUpdateModeManual,
		"rolling": NGroupUpdateModeRolling,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NGroupUpdateMode(input)
	return &out, nil
}

type OperatingSystemTypes string

const (
	OperatingSystemTypesLinux   OperatingSystemTypes = "Linux"
	OperatingSystemTypesWindows OperatingSystemTypes = "Windows"
)

func PossibleValuesForOperatingSystemTypes() []string {
	return []string{
		string(OperatingSystemTypesLinux),
		string(OperatingSystemTypesWindows),
	}
}

func (s *OperatingSystemTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperatingSystemTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperatingSystemTypes(input string) (*OperatingSystemTypes, error) {
	vals := map[string]OperatingSystemTypes{
		"linux":   OperatingSystemTypesLinux,
		"windows": OperatingSystemTypesWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatingSystemTypes(input)
	return &out, nil
}

type Priority string

const (
	PriorityRegular Priority = "Regular"
	PrioritySpot    Priority = "Spot"
)

func PossibleValuesForPriority() []string {
	return []string{
		string(PriorityRegular),
		string(PrioritySpot),
	}
}

func (s *Priority) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePriority(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePriority(input string) (*Priority, error) {
	vals := map[string]Priority{
		"regular": PriorityRegular,
		"spot":    PrioritySpot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Priority(input)
	return &out, nil
}

type Scheme string

const (
	SchemeHTTP  Scheme = "http"
	SchemeHTTPS Scheme = "https"
)

func PossibleValuesForScheme() []string {
	return []string{
		string(SchemeHTTP),
		string(SchemeHTTPS),
	}
}

func (s *Scheme) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheme(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheme(input string) (*Scheme, error) {
	vals := map[string]Scheme{
		"http":  SchemeHTTP,
		"https": SchemeHTTPS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Scheme(input)
	return &out, nil
}
