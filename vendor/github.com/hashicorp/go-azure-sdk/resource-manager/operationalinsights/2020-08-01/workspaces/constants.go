package workspaces

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataIngestionStatus string

const (
	DataIngestionStatusApproachingQuota      DataIngestionStatus = "ApproachingQuota"
	DataIngestionStatusForceOff              DataIngestionStatus = "ForceOff"
	DataIngestionStatusForceOn               DataIngestionStatus = "ForceOn"
	DataIngestionStatusOverQuota             DataIngestionStatus = "OverQuota"
	DataIngestionStatusRespectQuota          DataIngestionStatus = "RespectQuota"
	DataIngestionStatusSubscriptionSuspended DataIngestionStatus = "SubscriptionSuspended"
)

func PossibleValuesForDataIngestionStatus() []string {
	return []string{
		string(DataIngestionStatusApproachingQuota),
		string(DataIngestionStatusForceOff),
		string(DataIngestionStatusForceOn),
		string(DataIngestionStatusOverQuota),
		string(DataIngestionStatusRespectQuota),
		string(DataIngestionStatusSubscriptionSuspended),
	}
}

func parseDataIngestionStatus(input string) (*DataIngestionStatus, error) {
	vals := map[string]DataIngestionStatus{
		"approachingquota":      DataIngestionStatusApproachingQuota,
		"forceoff":              DataIngestionStatusForceOff,
		"forceon":               DataIngestionStatusForceOn,
		"overquota":             DataIngestionStatusOverQuota,
		"respectquota":          DataIngestionStatusRespectQuota,
		"subscriptionsuspended": DataIngestionStatusSubscriptionSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataIngestionStatus(input)
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

type PurgeState string

const (
	PurgeStateCompleted PurgeState = "completed"
	PurgeStatePending   PurgeState = "pending"
)

func PossibleValuesForPurgeState() []string {
	return []string{
		string(PurgeStateCompleted),
		string(PurgeStatePending),
	}
}

func parsePurgeState(input string) (*PurgeState, error) {
	vals := map[string]PurgeState{
		"completed": PurgeStateCompleted,
		"pending":   PurgeStatePending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PurgeState(input)
	return &out, nil
}

type SearchSortEnum string

const (
	SearchSortEnumAsc  SearchSortEnum = "asc"
	SearchSortEnumDesc SearchSortEnum = "desc"
)

func PossibleValuesForSearchSortEnum() []string {
	return []string{
		string(SearchSortEnumAsc),
		string(SearchSortEnumDesc),
	}
}

func parseSearchSortEnum(input string) (*SearchSortEnum, error) {
	vals := map[string]SearchSortEnum{
		"asc":  SearchSortEnumAsc,
		"desc": SearchSortEnumDesc,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SearchSortEnum(input)
	return &out, nil
}

type WorkspaceEntityStatus string

const (
	WorkspaceEntityStatusCanceled            WorkspaceEntityStatus = "Canceled"
	WorkspaceEntityStatusCreating            WorkspaceEntityStatus = "Creating"
	WorkspaceEntityStatusDeleting            WorkspaceEntityStatus = "Deleting"
	WorkspaceEntityStatusFailed              WorkspaceEntityStatus = "Failed"
	WorkspaceEntityStatusProvisioningAccount WorkspaceEntityStatus = "ProvisioningAccount"
	WorkspaceEntityStatusSucceeded           WorkspaceEntityStatus = "Succeeded"
	WorkspaceEntityStatusUpdating            WorkspaceEntityStatus = "Updating"
)

func PossibleValuesForWorkspaceEntityStatus() []string {
	return []string{
		string(WorkspaceEntityStatusCanceled),
		string(WorkspaceEntityStatusCreating),
		string(WorkspaceEntityStatusDeleting),
		string(WorkspaceEntityStatusFailed),
		string(WorkspaceEntityStatusProvisioningAccount),
		string(WorkspaceEntityStatusSucceeded),
		string(WorkspaceEntityStatusUpdating),
	}
}

func parseWorkspaceEntityStatus(input string) (*WorkspaceEntityStatus, error) {
	vals := map[string]WorkspaceEntityStatus{
		"canceled":            WorkspaceEntityStatusCanceled,
		"creating":            WorkspaceEntityStatusCreating,
		"deleting":            WorkspaceEntityStatusDeleting,
		"failed":              WorkspaceEntityStatusFailed,
		"provisioningaccount": WorkspaceEntityStatusProvisioningAccount,
		"succeeded":           WorkspaceEntityStatusSucceeded,
		"updating":            WorkspaceEntityStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkspaceEntityStatus(input)
	return &out, nil
}

type WorkspaceSkuNameEnum string

const (
	WorkspaceSkuNameEnumCapacityReservation  WorkspaceSkuNameEnum = "CapacityReservation"
	WorkspaceSkuNameEnumFree                 WorkspaceSkuNameEnum = "Free"
	WorkspaceSkuNameEnumLACluster            WorkspaceSkuNameEnum = "LACluster"
	WorkspaceSkuNameEnumPerGBTwoZeroOneEight WorkspaceSkuNameEnum = "PerGB2018"
	WorkspaceSkuNameEnumPerNode              WorkspaceSkuNameEnum = "PerNode"
	WorkspaceSkuNameEnumPremium              WorkspaceSkuNameEnum = "Premium"
	WorkspaceSkuNameEnumStandalone           WorkspaceSkuNameEnum = "Standalone"
	WorkspaceSkuNameEnumStandard             WorkspaceSkuNameEnum = "Standard"
)

func PossibleValuesForWorkspaceSkuNameEnum() []string {
	return []string{
		string(WorkspaceSkuNameEnumCapacityReservation),
		string(WorkspaceSkuNameEnumFree),
		string(WorkspaceSkuNameEnumLACluster),
		string(WorkspaceSkuNameEnumPerGBTwoZeroOneEight),
		string(WorkspaceSkuNameEnumPerNode),
		string(WorkspaceSkuNameEnumPremium),
		string(WorkspaceSkuNameEnumStandalone),
		string(WorkspaceSkuNameEnumStandard),
	}
}

func parseWorkspaceSkuNameEnum(input string) (*WorkspaceSkuNameEnum, error) {
	vals := map[string]WorkspaceSkuNameEnum{
		"capacityreservation": WorkspaceSkuNameEnumCapacityReservation,
		"free":                WorkspaceSkuNameEnumFree,
		"lacluster":           WorkspaceSkuNameEnumLACluster,
		"pergb2018":           WorkspaceSkuNameEnumPerGBTwoZeroOneEight,
		"pernode":             WorkspaceSkuNameEnumPerNode,
		"premium":             WorkspaceSkuNameEnumPremium,
		"standalone":          WorkspaceSkuNameEnumStandalone,
		"standard":            WorkspaceSkuNameEnumStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkspaceSkuNameEnum(input)
	return &out, nil
}
