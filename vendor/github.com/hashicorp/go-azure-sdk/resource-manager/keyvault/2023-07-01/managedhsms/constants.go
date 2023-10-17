package managedhsms

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionsRequired string

const (
	ActionsRequiredNone ActionsRequired = "None"
)

func PossibleValuesForActionsRequired() []string {
	return []string{
		string(ActionsRequiredNone),
	}
}

func parseActionsRequired(input string) (*ActionsRequired, error) {
	vals := map[string]ActionsRequired{
		"none": ActionsRequiredNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionsRequired(input)
	return &out, nil
}

type ActivationStatus string

const (
	ActivationStatusActive       ActivationStatus = "Active"
	ActivationStatusFailed       ActivationStatus = "Failed"
	ActivationStatusNotActivated ActivationStatus = "NotActivated"
	ActivationStatusUnknown      ActivationStatus = "Unknown"
)

func PossibleValuesForActivationStatus() []string {
	return []string{
		string(ActivationStatusActive),
		string(ActivationStatusFailed),
		string(ActivationStatusNotActivated),
		string(ActivationStatusUnknown),
	}
}

func parseActivationStatus(input string) (*ActivationStatus, error) {
	vals := map[string]ActivationStatus{
		"active":       ActivationStatusActive,
		"failed":       ActivationStatusFailed,
		"notactivated": ActivationStatusNotActivated,
		"unknown":      ActivationStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActivationStatus(input)
	return &out, nil
}

type CreateMode string

const (
	CreateModeDefault CreateMode = "default"
	CreateModeRecover CreateMode = "recover"
)

func PossibleValuesForCreateMode() []string {
	return []string{
		string(CreateModeDefault),
		string(CreateModeRecover),
	}
}

func parseCreateMode(input string) (*CreateMode, error) {
	vals := map[string]CreateMode{
		"default": CreateModeDefault,
		"recover": CreateModeRecover,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateMode(input)
	return &out, nil
}

type GeoReplicationRegionProvisioningState string

const (
	GeoReplicationRegionProvisioningStateCleanup         GeoReplicationRegionProvisioningState = "Cleanup"
	GeoReplicationRegionProvisioningStateDeleting        GeoReplicationRegionProvisioningState = "Deleting"
	GeoReplicationRegionProvisioningStateFailed          GeoReplicationRegionProvisioningState = "Failed"
	GeoReplicationRegionProvisioningStatePreprovisioning GeoReplicationRegionProvisioningState = "Preprovisioning"
	GeoReplicationRegionProvisioningStateProvisioning    GeoReplicationRegionProvisioningState = "Provisioning"
	GeoReplicationRegionProvisioningStateSucceeded       GeoReplicationRegionProvisioningState = "Succeeded"
)

func PossibleValuesForGeoReplicationRegionProvisioningState() []string {
	return []string{
		string(GeoReplicationRegionProvisioningStateCleanup),
		string(GeoReplicationRegionProvisioningStateDeleting),
		string(GeoReplicationRegionProvisioningStateFailed),
		string(GeoReplicationRegionProvisioningStatePreprovisioning),
		string(GeoReplicationRegionProvisioningStateProvisioning),
		string(GeoReplicationRegionProvisioningStateSucceeded),
	}
}

func parseGeoReplicationRegionProvisioningState(input string) (*GeoReplicationRegionProvisioningState, error) {
	vals := map[string]GeoReplicationRegionProvisioningState{
		"cleanup":         GeoReplicationRegionProvisioningStateCleanup,
		"deleting":        GeoReplicationRegionProvisioningStateDeleting,
		"failed":          GeoReplicationRegionProvisioningStateFailed,
		"preprovisioning": GeoReplicationRegionProvisioningStatePreprovisioning,
		"provisioning":    GeoReplicationRegionProvisioningStateProvisioning,
		"succeeded":       GeoReplicationRegionProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GeoReplicationRegionProvisioningState(input)
	return &out, nil
}

type ManagedHsmSkuFamily string

const (
	ManagedHsmSkuFamilyB ManagedHsmSkuFamily = "B"
)

func PossibleValuesForManagedHsmSkuFamily() []string {
	return []string{
		string(ManagedHsmSkuFamilyB),
	}
}

func parseManagedHsmSkuFamily(input string) (*ManagedHsmSkuFamily, error) {
	vals := map[string]ManagedHsmSkuFamily{
		"b": ManagedHsmSkuFamilyB,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedHsmSkuFamily(input)
	return &out, nil
}

type ManagedHsmSkuName string

const (
	ManagedHsmSkuNameCustomBSix      ManagedHsmSkuName = "Custom_B6"
	ManagedHsmSkuNameCustomBThreeTwo ManagedHsmSkuName = "Custom_B32"
	ManagedHsmSkuNameStandardBOne    ManagedHsmSkuName = "Standard_B1"
)

func PossibleValuesForManagedHsmSkuName() []string {
	return []string{
		string(ManagedHsmSkuNameCustomBSix),
		string(ManagedHsmSkuNameCustomBThreeTwo),
		string(ManagedHsmSkuNameStandardBOne),
	}
}

func parseManagedHsmSkuName(input string) (*ManagedHsmSkuName, error) {
	vals := map[string]ManagedHsmSkuName{
		"custom_b6":   ManagedHsmSkuNameCustomBSix,
		"custom_b32":  ManagedHsmSkuNameCustomBThreeTwo,
		"standard_b1": ManagedHsmSkuNameStandardBOne,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedHsmSkuName(input)
	return &out, nil
}

type NetworkRuleAction string

const (
	NetworkRuleActionAllow NetworkRuleAction = "Allow"
	NetworkRuleActionDeny  NetworkRuleAction = "Deny"
)

func PossibleValuesForNetworkRuleAction() []string {
	return []string{
		string(NetworkRuleActionAllow),
		string(NetworkRuleActionDeny),
	}
}

func parseNetworkRuleAction(input string) (*NetworkRuleAction, error) {
	vals := map[string]NetworkRuleAction{
		"allow": NetworkRuleActionAllow,
		"deny":  NetworkRuleActionDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkRuleAction(input)
	return &out, nil
}

type NetworkRuleBypassOptions string

const (
	NetworkRuleBypassOptionsAzureServices NetworkRuleBypassOptions = "AzureServices"
	NetworkRuleBypassOptionsNone          NetworkRuleBypassOptions = "None"
)

func PossibleValuesForNetworkRuleBypassOptions() []string {
	return []string{
		string(NetworkRuleBypassOptionsAzureServices),
		string(NetworkRuleBypassOptionsNone),
	}
}

func parseNetworkRuleBypassOptions(input string) (*NetworkRuleBypassOptions, error) {
	vals := map[string]NetworkRuleBypassOptions{
		"azureservices": NetworkRuleBypassOptionsAzureServices,
		"none":          NetworkRuleBypassOptionsNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkRuleBypassOptions(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating     PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting     PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateDisconnected PrivateEndpointConnectionProvisioningState = "Disconnected"
	PrivateEndpointConnectionProvisioningStateFailed       PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded    PrivateEndpointConnectionProvisioningState = "Succeeded"
	PrivateEndpointConnectionProvisioningStateUpdating     PrivateEndpointConnectionProvisioningState = "Updating"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateDisconnected),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
		string(PrivateEndpointConnectionProvisioningStateUpdating),
	}
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":     PrivateEndpointConnectionProvisioningStateCreating,
		"deleting":     PrivateEndpointConnectionProvisioningStateDeleting,
		"disconnected": PrivateEndpointConnectionProvisioningStateDisconnected,
		"failed":       PrivateEndpointConnectionProvisioningStateFailed,
		"succeeded":    PrivateEndpointConnectionProvisioningStateSucceeded,
		"updating":     PrivateEndpointConnectionProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type PrivateEndpointServiceConnectionStatus string

const (
	PrivateEndpointServiceConnectionStatusApproved     PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusDisconnected PrivateEndpointServiceConnectionStatus = "Disconnected"
	PrivateEndpointServiceConnectionStatusPending      PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected     PrivateEndpointServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateEndpointServiceConnectionStatus() []string {
	return []string{
		string(PrivateEndpointServiceConnectionStatusApproved),
		string(PrivateEndpointServiceConnectionStatusDisconnected),
		string(PrivateEndpointServiceConnectionStatusPending),
		string(PrivateEndpointServiceConnectionStatusRejected),
	}
}

func parsePrivateEndpointServiceConnectionStatus(input string) (*PrivateEndpointServiceConnectionStatus, error) {
	vals := map[string]PrivateEndpointServiceConnectionStatus{
		"approved":     PrivateEndpointServiceConnectionStatusApproved,
		"disconnected": PrivateEndpointServiceConnectionStatusDisconnected,
		"pending":      PrivateEndpointServiceConnectionStatusPending,
		"rejected":     PrivateEndpointServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointServiceConnectionStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateActivated             ProvisioningState = "Activated"
	ProvisioningStateDeleting              ProvisioningState = "Deleting"
	ProvisioningStateFailed                ProvisioningState = "Failed"
	ProvisioningStateProvisioning          ProvisioningState = "Provisioning"
	ProvisioningStateRestoring             ProvisioningState = "Restoring"
	ProvisioningStateSecurityDomainRestore ProvisioningState = "SecurityDomainRestore"
	ProvisioningStateSucceeded             ProvisioningState = "Succeeded"
	ProvisioningStateUpdating              ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateActivated),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateRestoring),
		string(ProvisioningStateSecurityDomainRestore),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"activated":             ProvisioningStateActivated,
		"deleting":              ProvisioningStateDeleting,
		"failed":                ProvisioningStateFailed,
		"provisioning":          ProvisioningStateProvisioning,
		"restoring":             ProvisioningStateRestoring,
		"securitydomainrestore": ProvisioningStateSecurityDomainRestore,
		"succeeded":             ProvisioningStateSucceeded,
		"updating":              ProvisioningStateUpdating,
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

type Reason string

const (
	ReasonAccountNameInvalid Reason = "AccountNameInvalid"
	ReasonAlreadyExists      Reason = "AlreadyExists"
)

func PossibleValuesForReason() []string {
	return []string{
		string(ReasonAccountNameInvalid),
		string(ReasonAlreadyExists),
	}
}

func parseReason(input string) (*Reason, error) {
	vals := map[string]Reason{
		"accountnameinvalid": ReasonAccountNameInvalid,
		"alreadyexists":      ReasonAlreadyExists,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Reason(input)
	return &out, nil
}
