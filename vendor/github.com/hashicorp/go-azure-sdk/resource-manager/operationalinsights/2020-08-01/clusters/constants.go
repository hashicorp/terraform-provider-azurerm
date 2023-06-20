package clusters

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterEntityStatus string

const (
	ClusterEntityStatusCanceled            ClusterEntityStatus = "Canceled"
	ClusterEntityStatusCreating            ClusterEntityStatus = "Creating"
	ClusterEntityStatusDeleting            ClusterEntityStatus = "Deleting"
	ClusterEntityStatusFailed              ClusterEntityStatus = "Failed"
	ClusterEntityStatusProvisioningAccount ClusterEntityStatus = "ProvisioningAccount"
	ClusterEntityStatusSucceeded           ClusterEntityStatus = "Succeeded"
	ClusterEntityStatusUpdating            ClusterEntityStatus = "Updating"
)

func PossibleValuesForClusterEntityStatus() []string {
	return []string{
		string(ClusterEntityStatusCanceled),
		string(ClusterEntityStatusCreating),
		string(ClusterEntityStatusDeleting),
		string(ClusterEntityStatusFailed),
		string(ClusterEntityStatusProvisioningAccount),
		string(ClusterEntityStatusSucceeded),
		string(ClusterEntityStatusUpdating),
	}
}

func parseClusterEntityStatus(input string) (*ClusterEntityStatus, error) {
	vals := map[string]ClusterEntityStatus{
		"canceled":            ClusterEntityStatusCanceled,
		"creating":            ClusterEntityStatusCreating,
		"deleting":            ClusterEntityStatusDeleting,
		"failed":              ClusterEntityStatusFailed,
		"provisioningaccount": ClusterEntityStatusProvisioningAccount,
		"succeeded":           ClusterEntityStatusSucceeded,
		"updating":            ClusterEntityStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterEntityStatus(input)
	return &out, nil
}

type ClusterSkuNameEnum string

const (
	ClusterSkuNameEnumCapacityReservation ClusterSkuNameEnum = "CapacityReservation"
)

func PossibleValuesForClusterSkuNameEnum() []string {
	return []string{
		string(ClusterSkuNameEnumCapacityReservation),
	}
}

func parseClusterSkuNameEnum(input string) (*ClusterSkuNameEnum, error) {
	vals := map[string]ClusterSkuNameEnum{
		"capacityreservation": ClusterSkuNameEnumCapacityReservation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterSkuNameEnum(input)
	return &out, nil
}
