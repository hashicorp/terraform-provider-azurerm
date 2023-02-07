package monitorsresource

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiftrResourceCategories string

const (
	LiftrResourceCategoriesMonitorLogs LiftrResourceCategories = "MonitorLogs"
	LiftrResourceCategoriesUnknown     LiftrResourceCategories = "Unknown"
)

func PossibleValuesForLiftrResourceCategories() []string {
	return []string{
		string(LiftrResourceCategoriesMonitorLogs),
		string(LiftrResourceCategoriesUnknown),
	}
}

func parseLiftrResourceCategories(input string) (*LiftrResourceCategories, error) {
	vals := map[string]LiftrResourceCategories{
		"monitorlogs": LiftrResourceCategoriesMonitorLogs,
		"unknown":     LiftrResourceCategoriesUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LiftrResourceCategories(input)
	return &out, nil
}

type ManagedIdentityTypes string

const (
	ManagedIdentityTypesSystemAssigned ManagedIdentityTypes = "SystemAssigned"
	ManagedIdentityTypesUserAssigned   ManagedIdentityTypes = "UserAssigned"
)

func PossibleValuesForManagedIdentityTypes() []string {
	return []string{
		string(ManagedIdentityTypesSystemAssigned),
		string(ManagedIdentityTypesUserAssigned),
	}
}

func parseManagedIdentityTypes(input string) (*ManagedIdentityTypes, error) {
	vals := map[string]ManagedIdentityTypes{
		"systemassigned": ManagedIdentityTypesSystemAssigned,
		"userassigned":   ManagedIdentityTypesUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedIdentityTypes(input)
	return &out, nil
}

type MarketplaceSubscriptionStatus string

const (
	MarketplaceSubscriptionStatusActive       MarketplaceSubscriptionStatus = "Active"
	MarketplaceSubscriptionStatusProvisioning MarketplaceSubscriptionStatus = "Provisioning"
	MarketplaceSubscriptionStatusSuspended    MarketplaceSubscriptionStatus = "Suspended"
	MarketplaceSubscriptionStatusUnsubscribed MarketplaceSubscriptionStatus = "Unsubscribed"
)

func PossibleValuesForMarketplaceSubscriptionStatus() []string {
	return []string{
		string(MarketplaceSubscriptionStatusActive),
		string(MarketplaceSubscriptionStatusProvisioning),
		string(MarketplaceSubscriptionStatusSuspended),
		string(MarketplaceSubscriptionStatusUnsubscribed),
	}
}

func parseMarketplaceSubscriptionStatus(input string) (*MarketplaceSubscriptionStatus, error) {
	vals := map[string]MarketplaceSubscriptionStatus{
		"active":       MarketplaceSubscriptionStatusActive,
		"provisioning": MarketplaceSubscriptionStatusProvisioning,
		"suspended":    MarketplaceSubscriptionStatusSuspended,
		"unsubscribed": MarketplaceSubscriptionStatusUnsubscribed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MarketplaceSubscriptionStatus(input)
	return &out, nil
}

type MonitoringStatus string

const (
	MonitoringStatusDisabled MonitoringStatus = "Disabled"
	MonitoringStatusEnabled  MonitoringStatus = "Enabled"
)

func PossibleValuesForMonitoringStatus() []string {
	return []string{
		string(MonitoringStatusDisabled),
		string(MonitoringStatusEnabled),
	}
}

func parseMonitoringStatus(input string) (*MonitoringStatus, error) {
	vals := map[string]MonitoringStatus{
		"disabled": MonitoringStatusDisabled,
		"enabled":  MonitoringStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MonitoringStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleted      ProvisioningState = "Deleted"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateNotSpecified ProvisioningState = "NotSpecified"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateNotSpecified),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"creating":     ProvisioningStateCreating,
		"deleted":      ProvisioningStateDeleted,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"notspecified": ProvisioningStateNotSpecified,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
