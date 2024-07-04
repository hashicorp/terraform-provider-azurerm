package clusters

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingType string

const (
	BillingTypeCluster    BillingType = "Cluster"
	BillingTypeWorkspaces BillingType = "Workspaces"
)

func PossibleValuesForBillingType() []string {
	return []string{
		string(BillingTypeCluster),
		string(BillingTypeWorkspaces),
	}
}

func (s *BillingType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBillingType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBillingType(input string) (*BillingType, error) {
	vals := map[string]BillingType{
		"cluster":    BillingTypeCluster,
		"workspaces": BillingTypeWorkspaces,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BillingType(input)
	return &out, nil
}

type Capacity int64

const (
	CapacityFiveHundred      Capacity = 500
	CapacityFiveThousand     Capacity = 5000
	CapacityFiveZeroThousand Capacity = 50000
	CapacityFourHundred      Capacity = 400
	CapacityOneHundred       Capacity = 100
	CapacityOneThousand      Capacity = 1000
	CapacityOneZeroThousand  Capacity = 10000
	CapacityThreeHundred     Capacity = 300
	CapacityTwoFiveThousand  Capacity = 25000
	CapacityTwoHundred       Capacity = 200
	CapacityTwoThousand      Capacity = 2000
)

func PossibleValuesForCapacity() []int64 {
	return []int64{
		int64(CapacityFiveHundred),
		int64(CapacityFiveThousand),
		int64(CapacityFiveZeroThousand),
		int64(CapacityFourHundred),
		int64(CapacityOneHundred),
		int64(CapacityOneThousand),
		int64(CapacityOneZeroThousand),
		int64(CapacityThreeHundred),
		int64(CapacityTwoFiveThousand),
		int64(CapacityTwoHundred),
		int64(CapacityTwoThousand),
	}
}

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

func (s *ClusterEntityStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterEntityStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *ClusterSkuNameEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterSkuNameEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
