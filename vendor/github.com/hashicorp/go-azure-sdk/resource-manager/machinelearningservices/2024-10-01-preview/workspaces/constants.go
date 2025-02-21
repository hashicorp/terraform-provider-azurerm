package workspaces

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DefaultActionType string

const (
	DefaultActionTypeAllow DefaultActionType = "Allow"
	DefaultActionTypeDeny  DefaultActionType = "Deny"
)

func PossibleValuesForDefaultActionType() []string {
	return []string{
		string(DefaultActionTypeAllow),
		string(DefaultActionTypeDeny),
	}
}

func (s *DefaultActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDefaultActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDefaultActionType(input string) (*DefaultActionType, error) {
	vals := map[string]DefaultActionType{
		"allow": DefaultActionTypeAllow,
		"deny":  DefaultActionTypeDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DefaultActionType(input)
	return &out, nil
}

type EncryptionStatus string

const (
	EncryptionStatusDisabled EncryptionStatus = "Disabled"
	EncryptionStatusEnabled  EncryptionStatus = "Enabled"
)

func PossibleValuesForEncryptionStatus() []string {
	return []string{
		string(EncryptionStatusDisabled),
		string(EncryptionStatusEnabled),
	}
}

func (s *EncryptionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionStatus(input string) (*EncryptionStatus, error) {
	vals := map[string]EncryptionStatus{
		"disabled": EncryptionStatusDisabled,
		"enabled":  EncryptionStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionStatus(input)
	return &out, nil
}

type EndpointServiceConnectionStatus string

const (
	EndpointServiceConnectionStatusApproved     EndpointServiceConnectionStatus = "Approved"
	EndpointServiceConnectionStatusDisconnected EndpointServiceConnectionStatus = "Disconnected"
	EndpointServiceConnectionStatusPending      EndpointServiceConnectionStatus = "Pending"
	EndpointServiceConnectionStatusRejected     EndpointServiceConnectionStatus = "Rejected"
	EndpointServiceConnectionStatusTimeout      EndpointServiceConnectionStatus = "Timeout"
)

func PossibleValuesForEndpointServiceConnectionStatus() []string {
	return []string{
		string(EndpointServiceConnectionStatusApproved),
		string(EndpointServiceConnectionStatusDisconnected),
		string(EndpointServiceConnectionStatusPending),
		string(EndpointServiceConnectionStatusRejected),
		string(EndpointServiceConnectionStatusTimeout),
	}
}

func (s *EndpointServiceConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEndpointServiceConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEndpointServiceConnectionStatus(input string) (*EndpointServiceConnectionStatus, error) {
	vals := map[string]EndpointServiceConnectionStatus{
		"approved":     EndpointServiceConnectionStatusApproved,
		"disconnected": EndpointServiceConnectionStatusDisconnected,
		"pending":      EndpointServiceConnectionStatusPending,
		"rejected":     EndpointServiceConnectionStatusRejected,
		"timeout":      EndpointServiceConnectionStatusTimeout,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointServiceConnectionStatus(input)
	return &out, nil
}

type FirewallSku string

const (
	FirewallSkuBasic    FirewallSku = "Basic"
	FirewallSkuStandard FirewallSku = "Standard"
)

func PossibleValuesForFirewallSku() []string {
	return []string{
		string(FirewallSkuBasic),
		string(FirewallSkuStandard),
	}
}

func (s *FirewallSku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallSku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallSku(input string) (*FirewallSku, error) {
	vals := map[string]FirewallSku{
		"basic":    FirewallSkuBasic,
		"standard": FirewallSkuStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallSku(input)
	return &out, nil
}

type IsolationMode string

const (
	IsolationModeAllowInternetOutbound     IsolationMode = "AllowInternetOutbound"
	IsolationModeAllowOnlyApprovedOutbound IsolationMode = "AllowOnlyApprovedOutbound"
	IsolationModeDisabled                  IsolationMode = "Disabled"
)

func PossibleValuesForIsolationMode() []string {
	return []string{
		string(IsolationModeAllowInternetOutbound),
		string(IsolationModeAllowOnlyApprovedOutbound),
		string(IsolationModeDisabled),
	}
}

func (s *IsolationMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIsolationMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIsolationMode(input string) (*IsolationMode, error) {
	vals := map[string]IsolationMode{
		"allowinternetoutbound":     IsolationModeAllowInternetOutbound,
		"allowonlyapprovedoutbound": IsolationModeAllowOnlyApprovedOutbound,
		"disabled":                  IsolationModeDisabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IsolationMode(input)
	return &out, nil
}

type ManagedNetworkStatus string

const (
	ManagedNetworkStatusActive   ManagedNetworkStatus = "Active"
	ManagedNetworkStatusInactive ManagedNetworkStatus = "Inactive"
)

func PossibleValuesForManagedNetworkStatus() []string {
	return []string{
		string(ManagedNetworkStatusActive),
		string(ManagedNetworkStatusInactive),
	}
}

func (s *ManagedNetworkStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedNetworkStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedNetworkStatus(input string) (*ManagedNetworkStatus, error) {
	vals := map[string]ManagedNetworkStatus{
		"active":   ManagedNetworkStatusActive,
		"inactive": ManagedNetworkStatusInactive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedNetworkStatus(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
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
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUnknown   ProvisioningState = "Unknown"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnknown),
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
		"unknown":   ProvisioningStateUnknown,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicNetworkAccessType string

const (
	PublicNetworkAccessTypeDisabled PublicNetworkAccessType = "Disabled"
	PublicNetworkAccessTypeEnabled  PublicNetworkAccessType = "Enabled"
)

func PossibleValuesForPublicNetworkAccessType() []string {
	return []string{
		string(PublicNetworkAccessTypeDisabled),
		string(PublicNetworkAccessTypeEnabled),
	}
}

func (s *PublicNetworkAccessType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccessType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccessType(input string) (*PublicNetworkAccessType, error) {
	vals := map[string]PublicNetworkAccessType{
		"disabled": PublicNetworkAccessTypeDisabled,
		"enabled":  PublicNetworkAccessTypeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccessType(input)
	return &out, nil
}

type RuleAction string

const (
	RuleActionAllow RuleAction = "Allow"
	RuleActionDeny  RuleAction = "Deny"
)

func PossibleValuesForRuleAction() []string {
	return []string{
		string(RuleActionAllow),
		string(RuleActionDeny),
	}
}

func (s *RuleAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRuleAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRuleAction(input string) (*RuleAction, error) {
	vals := map[string]RuleAction{
		"allow": RuleActionAllow,
		"deny":  RuleActionDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleAction(input)
	return &out, nil
}

type RuleCategory string

const (
	RuleCategoryDependency  RuleCategory = "Dependency"
	RuleCategoryRecommended RuleCategory = "Recommended"
	RuleCategoryRequired    RuleCategory = "Required"
	RuleCategoryUserDefined RuleCategory = "UserDefined"
)

func PossibleValuesForRuleCategory() []string {
	return []string{
		string(RuleCategoryDependency),
		string(RuleCategoryRecommended),
		string(RuleCategoryRequired),
		string(RuleCategoryUserDefined),
	}
}

func (s *RuleCategory) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRuleCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRuleCategory(input string) (*RuleCategory, error) {
	vals := map[string]RuleCategory{
		"dependency":  RuleCategoryDependency,
		"recommended": RuleCategoryRecommended,
		"required":    RuleCategoryRequired,
		"userdefined": RuleCategoryUserDefined,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleCategory(input)
	return &out, nil
}

type RuleStatus string

const (
	RuleStatusActive   RuleStatus = "Active"
	RuleStatusInactive RuleStatus = "Inactive"
)

func PossibleValuesForRuleStatus() []string {
	return []string{
		string(RuleStatusActive),
		string(RuleStatusInactive),
	}
}

func (s *RuleStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRuleStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRuleStatus(input string) (*RuleStatus, error) {
	vals := map[string]RuleStatus{
		"active":   RuleStatusActive,
		"inactive": RuleStatusInactive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleStatus(input)
	return &out, nil
}

type RuleType string

const (
	RuleTypeFQDN            RuleType = "FQDN"
	RuleTypePrivateEndpoint RuleType = "PrivateEndpoint"
	RuleTypeServiceTag      RuleType = "ServiceTag"
)

func PossibleValuesForRuleType() []string {
	return []string{
		string(RuleTypeFQDN),
		string(RuleTypePrivateEndpoint),
		string(RuleTypeServiceTag),
	}
}

func (s *RuleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRuleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRuleType(input string) (*RuleType, error) {
	vals := map[string]RuleType{
		"fqdn":            RuleTypeFQDN,
		"privateendpoint": RuleTypePrivateEndpoint,
		"servicetag":      RuleTypeServiceTag,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleType(input)
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

type SystemDatastoresAuthMode string

const (
	SystemDatastoresAuthModeAccessKey         SystemDatastoresAuthMode = "AccessKey"
	SystemDatastoresAuthModeIdentity          SystemDatastoresAuthMode = "Identity"
	SystemDatastoresAuthModeUserDelegationSAS SystemDatastoresAuthMode = "UserDelegationSAS"
)

func PossibleValuesForSystemDatastoresAuthMode() []string {
	return []string{
		string(SystemDatastoresAuthModeAccessKey),
		string(SystemDatastoresAuthModeIdentity),
		string(SystemDatastoresAuthModeUserDelegationSAS),
	}
}

func (s *SystemDatastoresAuthMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSystemDatastoresAuthMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSystemDatastoresAuthMode(input string) (*SystemDatastoresAuthMode, error) {
	vals := map[string]SystemDatastoresAuthMode{
		"accesskey":         SystemDatastoresAuthModeAccessKey,
		"identity":          SystemDatastoresAuthModeIdentity,
		"userdelegationsas": SystemDatastoresAuthModeUserDelegationSAS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SystemDatastoresAuthMode(input)
	return &out, nil
}
