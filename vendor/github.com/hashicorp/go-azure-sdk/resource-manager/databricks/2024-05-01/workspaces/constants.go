package workspaces

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomaticClusterUpdateValue string

const (
	AutomaticClusterUpdateValueDisabled AutomaticClusterUpdateValue = "Disabled"
	AutomaticClusterUpdateValueEnabled  AutomaticClusterUpdateValue = "Enabled"
)

func PossibleValuesForAutomaticClusterUpdateValue() []string {
	return []string{
		string(AutomaticClusterUpdateValueDisabled),
		string(AutomaticClusterUpdateValueEnabled),
	}
}

func (s *AutomaticClusterUpdateValue) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomaticClusterUpdateValue(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomaticClusterUpdateValue(input string) (*AutomaticClusterUpdateValue, error) {
	vals := map[string]AutomaticClusterUpdateValue{
		"disabled": AutomaticClusterUpdateValueDisabled,
		"enabled":  AutomaticClusterUpdateValueEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomaticClusterUpdateValue(input)
	return &out, nil
}

type ComplianceSecurityProfileValue string

const (
	ComplianceSecurityProfileValueDisabled ComplianceSecurityProfileValue = "Disabled"
	ComplianceSecurityProfileValueEnabled  ComplianceSecurityProfileValue = "Enabled"
)

func PossibleValuesForComplianceSecurityProfileValue() []string {
	return []string{
		string(ComplianceSecurityProfileValueDisabled),
		string(ComplianceSecurityProfileValueEnabled),
	}
}

func (s *ComplianceSecurityProfileValue) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComplianceSecurityProfileValue(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComplianceSecurityProfileValue(input string) (*ComplianceSecurityProfileValue, error) {
	vals := map[string]ComplianceSecurityProfileValue{
		"disabled": ComplianceSecurityProfileValueDisabled,
		"enabled":  ComplianceSecurityProfileValueEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComplianceSecurityProfileValue(input)
	return &out, nil
}

type ComplianceStandard string

const (
	ComplianceStandardHIPAA  ComplianceStandard = "HIPAA"
	ComplianceStandardNONE   ComplianceStandard = "NONE"
	ComplianceStandardPCIDSS ComplianceStandard = "PCI_DSS"
)

func PossibleValuesForComplianceStandard() []string {
	return []string{
		string(ComplianceStandardHIPAA),
		string(ComplianceStandardNONE),
		string(ComplianceStandardPCIDSS),
	}
}

func (s *ComplianceStandard) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComplianceStandard(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComplianceStandard(input string) (*ComplianceStandard, error) {
	vals := map[string]ComplianceStandard{
		"hipaa":   ComplianceStandardHIPAA,
		"none":    ComplianceStandardNONE,
		"pci_dss": ComplianceStandardPCIDSS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComplianceStandard(input)
	return &out, nil
}

type CustomParameterType string

const (
	CustomParameterTypeBool   CustomParameterType = "Bool"
	CustomParameterTypeObject CustomParameterType = "Object"
	CustomParameterTypeString CustomParameterType = "String"
)

func PossibleValuesForCustomParameterType() []string {
	return []string{
		string(CustomParameterTypeBool),
		string(CustomParameterTypeObject),
		string(CustomParameterTypeString),
	}
}

func (s *CustomParameterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomParameterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomParameterType(input string) (*CustomParameterType, error) {
	vals := map[string]CustomParameterType{
		"bool":   CustomParameterTypeBool,
		"object": CustomParameterTypeObject,
		"string": CustomParameterTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomParameterType(input)
	return &out, nil
}

type DefaultStorageFirewall string

const (
	DefaultStorageFirewallDisabled DefaultStorageFirewall = "Disabled"
	DefaultStorageFirewallEnabled  DefaultStorageFirewall = "Enabled"
)

func PossibleValuesForDefaultStorageFirewall() []string {
	return []string{
		string(DefaultStorageFirewallDisabled),
		string(DefaultStorageFirewallEnabled),
	}
}

func (s *DefaultStorageFirewall) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDefaultStorageFirewall(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDefaultStorageFirewall(input string) (*DefaultStorageFirewall, error) {
	vals := map[string]DefaultStorageFirewall{
		"disabled": DefaultStorageFirewallDisabled,
		"enabled":  DefaultStorageFirewallEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DefaultStorageFirewall(input)
	return &out, nil
}

type EncryptionKeySource string

const (
	EncryptionKeySourceMicrosoftPointKeyvault EncryptionKeySource = "Microsoft.Keyvault"
)

func PossibleValuesForEncryptionKeySource() []string {
	return []string{
		string(EncryptionKeySourceMicrosoftPointKeyvault),
	}
}

func (s *EncryptionKeySource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionKeySource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionKeySource(input string) (*EncryptionKeySource, error) {
	vals := map[string]EncryptionKeySource{
		"microsoft.keyvault": EncryptionKeySourceMicrosoftPointKeyvault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionKeySource(input)
	return &out, nil
}

type EnhancedSecurityMonitoringValue string

const (
	EnhancedSecurityMonitoringValueDisabled EnhancedSecurityMonitoringValue = "Disabled"
	EnhancedSecurityMonitoringValueEnabled  EnhancedSecurityMonitoringValue = "Enabled"
)

func PossibleValuesForEnhancedSecurityMonitoringValue() []string {
	return []string{
		string(EnhancedSecurityMonitoringValueDisabled),
		string(EnhancedSecurityMonitoringValueEnabled),
	}
}

func (s *EnhancedSecurityMonitoringValue) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnhancedSecurityMonitoringValue(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnhancedSecurityMonitoringValue(input string) (*EnhancedSecurityMonitoringValue, error) {
	vals := map[string]EnhancedSecurityMonitoringValue{
		"disabled": EnhancedSecurityMonitoringValueDisabled,
		"enabled":  EnhancedSecurityMonitoringValueEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnhancedSecurityMonitoringValue(input)
	return &out, nil
}

type IdentityType string

const (
	IdentityTypeSystemAssigned IdentityType = "SystemAssigned"
	IdentityTypeUserAssigned   IdentityType = "UserAssigned"
)

func PossibleValuesForIdentityType() []string {
	return []string{
		string(IdentityTypeSystemAssigned),
		string(IdentityTypeUserAssigned),
	}
}

func (s *IdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIdentityType(input string) (*IdentityType, error) {
	vals := map[string]IdentityType{
		"systemassigned": IdentityTypeSystemAssigned,
		"userassigned":   IdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentityType(input)
	return &out, nil
}

type InitialType string

const (
	InitialTypeHiveMetastore InitialType = "HiveMetastore"
	InitialTypeUnityCatalog  InitialType = "UnityCatalog"
)

func PossibleValuesForInitialType() []string {
	return []string{
		string(InitialTypeHiveMetastore),
		string(InitialTypeUnityCatalog),
	}
}

func (s *InitialType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInitialType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInitialType(input string) (*InitialType, error) {
	vals := map[string]InitialType{
		"hivemetastore": InitialTypeHiveMetastore,
		"unitycatalog":  InitialTypeUnityCatalog,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InitialType(input)
	return &out, nil
}

type KeySource string

const (
	KeySourceDefault                KeySource = "Default"
	KeySourceMicrosoftPointKeyvault KeySource = "Microsoft.Keyvault"
)

func PossibleValuesForKeySource() []string {
	return []string{
		string(KeySourceDefault),
		string(KeySourceMicrosoftPointKeyvault),
	}
}

func (s *KeySource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeySource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeySource(input string) (*KeySource, error) {
	vals := map[string]KeySource{
		"default":            KeySourceDefault,
		"microsoft.keyvault": KeySourceMicrosoftPointKeyvault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeySource(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
	PrivateEndpointConnectionProvisioningStateUpdating  PrivateEndpointConnectionProvisioningState = "Updating"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
		string(PrivateEndpointConnectionProvisioningStateUpdating),
	}
}

func (s *PrivateEndpointConnectionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointConnectionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":  PrivateEndpointConnectionProvisioningStateCreating,
		"deleting":  PrivateEndpointConnectionProvisioningStateDeleting,
		"failed":    PrivateEndpointConnectionProvisioningStateFailed,
		"succeeded": PrivateEndpointConnectionProvisioningStateSucceeded,
		"updating":  PrivateEndpointConnectionProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type PrivateLinkServiceConnectionStatus string

const (
	PrivateLinkServiceConnectionStatusApproved     PrivateLinkServiceConnectionStatus = "Approved"
	PrivateLinkServiceConnectionStatusDisconnected PrivateLinkServiceConnectionStatus = "Disconnected"
	PrivateLinkServiceConnectionStatusPending      PrivateLinkServiceConnectionStatus = "Pending"
	PrivateLinkServiceConnectionStatusRejected     PrivateLinkServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateLinkServiceConnectionStatus() []string {
	return []string{
		string(PrivateLinkServiceConnectionStatusApproved),
		string(PrivateLinkServiceConnectionStatusDisconnected),
		string(PrivateLinkServiceConnectionStatusPending),
		string(PrivateLinkServiceConnectionStatusRejected),
	}
}

func (s *PrivateLinkServiceConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateLinkServiceConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateLinkServiceConnectionStatus(input string) (*PrivateLinkServiceConnectionStatus, error) {
	vals := map[string]PrivateLinkServiceConnectionStatus{
		"approved":     PrivateLinkServiceConnectionStatusApproved,
		"disconnected": PrivateLinkServiceConnectionStatusDisconnected,
		"pending":      PrivateLinkServiceConnectionStatusPending,
		"rejected":     PrivateLinkServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateLinkServiceConnectionStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreated   ProvisioningState = "Created"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleted   ProvisioningState = "Deleted"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateReady     ProvisioningState = "Ready"
	ProvisioningStateRunning   ProvisioningState = "Running"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreated),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateReady),
		string(ProvisioningStateRunning),
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
		"accepted":  ProvisioningStateAccepted,
		"canceled":  ProvisioningStateCanceled,
		"created":   ProvisioningStateCreated,
		"creating":  ProvisioningStateCreating,
		"deleted":   ProvisioningStateDeleted,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"ready":     ProvisioningStateReady,
		"running":   ProvisioningStateRunning,
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

type RequiredNsgRules string

const (
	RequiredNsgRulesAllRules               RequiredNsgRules = "AllRules"
	RequiredNsgRulesNoAzureDatabricksRules RequiredNsgRules = "NoAzureDatabricksRules"
	RequiredNsgRulesNoAzureServiceRules    RequiredNsgRules = "NoAzureServiceRules"
)

func PossibleValuesForRequiredNsgRules() []string {
	return []string{
		string(RequiredNsgRulesAllRules),
		string(RequiredNsgRulesNoAzureDatabricksRules),
		string(RequiredNsgRulesNoAzureServiceRules),
	}
}

func (s *RequiredNsgRules) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRequiredNsgRules(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRequiredNsgRules(input string) (*RequiredNsgRules, error) {
	vals := map[string]RequiredNsgRules{
		"allrules":               RequiredNsgRulesAllRules,
		"noazuredatabricksrules": RequiredNsgRulesNoAzureDatabricksRules,
		"noazureservicerules":    RequiredNsgRulesNoAzureServiceRules,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequiredNsgRules(input)
	return &out, nil
}
