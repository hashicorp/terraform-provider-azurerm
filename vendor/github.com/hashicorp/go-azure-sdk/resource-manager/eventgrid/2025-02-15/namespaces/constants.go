package namespaces

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomainIdentityType string

const (
	CustomDomainIdentityTypeSystemAssigned CustomDomainIdentityType = "SystemAssigned"
	CustomDomainIdentityTypeUserAssigned   CustomDomainIdentityType = "UserAssigned"
)

func PossibleValuesForCustomDomainIdentityType() []string {
	return []string{
		string(CustomDomainIdentityTypeSystemAssigned),
		string(CustomDomainIdentityTypeUserAssigned),
	}
}

func (s *CustomDomainIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomDomainIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomDomainIdentityType(input string) (*CustomDomainIdentityType, error) {
	vals := map[string]CustomDomainIdentityType{
		"systemassigned": CustomDomainIdentityTypeSystemAssigned,
		"userassigned":   CustomDomainIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomDomainIdentityType(input)
	return &out, nil
}

type CustomDomainValidationState string

const (
	CustomDomainValidationStateApproved                 CustomDomainValidationState = "Approved"
	CustomDomainValidationStateErrorRetrievingDnsRecord CustomDomainValidationState = "ErrorRetrievingDnsRecord"
	CustomDomainValidationStatePending                  CustomDomainValidationState = "Pending"
)

func PossibleValuesForCustomDomainValidationState() []string {
	return []string{
		string(CustomDomainValidationStateApproved),
		string(CustomDomainValidationStateErrorRetrievingDnsRecord),
		string(CustomDomainValidationStatePending),
	}
}

func (s *CustomDomainValidationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomDomainValidationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomDomainValidationState(input string) (*CustomDomainValidationState, error) {
	vals := map[string]CustomDomainValidationState{
		"approved":                 CustomDomainValidationStateApproved,
		"errorretrievingdnsrecord": CustomDomainValidationStateErrorRetrievingDnsRecord,
		"pending":                  CustomDomainValidationStatePending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomDomainValidationState(input)
	return &out, nil
}

type IPActionType string

const (
	IPActionTypeAllow IPActionType = "Allow"
)

func PossibleValuesForIPActionType() []string {
	return []string{
		string(IPActionTypeAllow),
	}
}

func (s *IPActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPActionType(input string) (*IPActionType, error) {
	vals := map[string]IPActionType{
		"allow": IPActionTypeAllow,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPActionType(input)
	return &out, nil
}

type NamespaceProvisioningState string

const (
	NamespaceProvisioningStateCanceled      NamespaceProvisioningState = "Canceled"
	NamespaceProvisioningStateCreateFailed  NamespaceProvisioningState = "CreateFailed"
	NamespaceProvisioningStateCreating      NamespaceProvisioningState = "Creating"
	NamespaceProvisioningStateDeleteFailed  NamespaceProvisioningState = "DeleteFailed"
	NamespaceProvisioningStateDeleted       NamespaceProvisioningState = "Deleted"
	NamespaceProvisioningStateDeleting      NamespaceProvisioningState = "Deleting"
	NamespaceProvisioningStateFailed        NamespaceProvisioningState = "Failed"
	NamespaceProvisioningStateSucceeded     NamespaceProvisioningState = "Succeeded"
	NamespaceProvisioningStateUpdatedFailed NamespaceProvisioningState = "UpdatedFailed"
	NamespaceProvisioningStateUpdating      NamespaceProvisioningState = "Updating"
)

func PossibleValuesForNamespaceProvisioningState() []string {
	return []string{
		string(NamespaceProvisioningStateCanceled),
		string(NamespaceProvisioningStateCreateFailed),
		string(NamespaceProvisioningStateCreating),
		string(NamespaceProvisioningStateDeleteFailed),
		string(NamespaceProvisioningStateDeleted),
		string(NamespaceProvisioningStateDeleting),
		string(NamespaceProvisioningStateFailed),
		string(NamespaceProvisioningStateSucceeded),
		string(NamespaceProvisioningStateUpdatedFailed),
		string(NamespaceProvisioningStateUpdating),
	}
}

func (s *NamespaceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNamespaceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNamespaceProvisioningState(input string) (*NamespaceProvisioningState, error) {
	vals := map[string]NamespaceProvisioningState{
		"canceled":      NamespaceProvisioningStateCanceled,
		"createfailed":  NamespaceProvisioningStateCreateFailed,
		"creating":      NamespaceProvisioningStateCreating,
		"deletefailed":  NamespaceProvisioningStateDeleteFailed,
		"deleted":       NamespaceProvisioningStateDeleted,
		"deleting":      NamespaceProvisioningStateDeleting,
		"failed":        NamespaceProvisioningStateFailed,
		"succeeded":     NamespaceProvisioningStateSucceeded,
		"updatedfailed": NamespaceProvisioningStateUpdatedFailed,
		"updating":      NamespaceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NamespaceProvisioningState(input)
	return &out, nil
}

type PersistedConnectionStatus string

const (
	PersistedConnectionStatusApproved     PersistedConnectionStatus = "Approved"
	PersistedConnectionStatusDisconnected PersistedConnectionStatus = "Disconnected"
	PersistedConnectionStatusPending      PersistedConnectionStatus = "Pending"
	PersistedConnectionStatusRejected     PersistedConnectionStatus = "Rejected"
)

func PossibleValuesForPersistedConnectionStatus() []string {
	return []string{
		string(PersistedConnectionStatusApproved),
		string(PersistedConnectionStatusDisconnected),
		string(PersistedConnectionStatusPending),
		string(PersistedConnectionStatusRejected),
	}
}

func (s *PersistedConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePersistedConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePersistedConnectionStatus(input string) (*PersistedConnectionStatus, error) {
	vals := map[string]PersistedConnectionStatus{
		"approved":     PersistedConnectionStatusApproved,
		"disconnected": PersistedConnectionStatusDisconnected,
		"pending":      PersistedConnectionStatusPending,
		"rejected":     PersistedConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PersistedConnectionStatus(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func (s *PublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type ResourceProvisioningState string

const (
	ResourceProvisioningStateCanceled  ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateCreating  ResourceProvisioningState = "Creating"
	ResourceProvisioningStateDeleting  ResourceProvisioningState = "Deleting"
	ResourceProvisioningStateFailed    ResourceProvisioningState = "Failed"
	ResourceProvisioningStateSucceeded ResourceProvisioningState = "Succeeded"
	ResourceProvisioningStateUpdating  ResourceProvisioningState = "Updating"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateCreating),
		string(ResourceProvisioningStateDeleting),
		string(ResourceProvisioningStateFailed),
		string(ResourceProvisioningStateSucceeded),
		string(ResourceProvisioningStateUpdating),
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
		"canceled":  ResourceProvisioningStateCanceled,
		"creating":  ResourceProvisioningStateCreating,
		"deleting":  ResourceProvisioningStateDeleting,
		"failed":    ResourceProvisioningStateFailed,
		"succeeded": ResourceProvisioningStateSucceeded,
		"updating":  ResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}

type RoutingIdentityType string

const (
	RoutingIdentityTypeNone           RoutingIdentityType = "None"
	RoutingIdentityTypeSystemAssigned RoutingIdentityType = "SystemAssigned"
	RoutingIdentityTypeUserAssigned   RoutingIdentityType = "UserAssigned"
)

func PossibleValuesForRoutingIdentityType() []string {
	return []string{
		string(RoutingIdentityTypeNone),
		string(RoutingIdentityTypeSystemAssigned),
		string(RoutingIdentityTypeUserAssigned),
	}
}

func (s *RoutingIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoutingIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoutingIdentityType(input string) (*RoutingIdentityType, error) {
	vals := map[string]RoutingIdentityType{
		"none":           RoutingIdentityTypeNone,
		"systemassigned": RoutingIdentityTypeSystemAssigned,
		"userassigned":   RoutingIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoutingIdentityType(input)
	return &out, nil
}

type SkuName string

const (
	SkuNameStandard SkuName = "Standard"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameStandard),
	}
}

func (s *SkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"standard": SkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}

type StaticRoutingEnrichmentType string

const (
	StaticRoutingEnrichmentTypeString StaticRoutingEnrichmentType = "String"
)

func PossibleValuesForStaticRoutingEnrichmentType() []string {
	return []string{
		string(StaticRoutingEnrichmentTypeString),
	}
}

func (s *StaticRoutingEnrichmentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStaticRoutingEnrichmentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStaticRoutingEnrichmentType(input string) (*StaticRoutingEnrichmentType, error) {
	vals := map[string]StaticRoutingEnrichmentType{
		"string": StaticRoutingEnrichmentTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StaticRoutingEnrichmentType(input)
	return &out, nil
}

type TlsVersion string

const (
	TlsVersionOnePointOne  TlsVersion = "1.1"
	TlsVersionOnePointTwo  TlsVersion = "1.2"
	TlsVersionOnePointZero TlsVersion = "1.0"
)

func PossibleValuesForTlsVersion() []string {
	return []string{
		string(TlsVersionOnePointOne),
		string(TlsVersionOnePointTwo),
		string(TlsVersionOnePointZero),
	}
}

func (s *TlsVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTlsVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTlsVersion(input string) (*TlsVersion, error) {
	vals := map[string]TlsVersion{
		"1.1": TlsVersionOnePointOne,
		"1.2": TlsVersionOnePointTwo,
		"1.0": TlsVersionOnePointZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TlsVersion(input)
	return &out, nil
}

type TopicSpacesConfigurationState string

const (
	TopicSpacesConfigurationStateDisabled TopicSpacesConfigurationState = "Disabled"
	TopicSpacesConfigurationStateEnabled  TopicSpacesConfigurationState = "Enabled"
)

func PossibleValuesForTopicSpacesConfigurationState() []string {
	return []string{
		string(TopicSpacesConfigurationStateDisabled),
		string(TopicSpacesConfigurationStateEnabled),
	}
}

func (s *TopicSpacesConfigurationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTopicSpacesConfigurationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTopicSpacesConfigurationState(input string) (*TopicSpacesConfigurationState, error) {
	vals := map[string]TopicSpacesConfigurationState{
		"disabled": TopicSpacesConfigurationStateDisabled,
		"enabled":  TopicSpacesConfigurationStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TopicSpacesConfigurationState(input)
	return &out, nil
}
