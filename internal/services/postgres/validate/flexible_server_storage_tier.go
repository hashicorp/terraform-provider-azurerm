// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/servers"
)

type StorageTiers struct {
	DefaultTier      string
	ValidTiers       *[]string
	PossibleTiersInt *[]int
}

// Creates a map of valid StorageTiers based on the storage_gb for the PostgreSQL Flexible Server
func InitializeFlexibleServerStorageTierDefaults() map[int]StorageTiers {

	storageTiersMappings := map[int]StorageTiers{
		32768: {string(servers.AzureManagedDiskPerformanceTiersPFour), &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPSix),
			string(servers.AzureManagedDiskPerformanceTiersPOneZero),
			string(servers.AzureManagedDiskPerformanceTiersPOneFive),
			string(servers.AzureManagedDiskPerformanceTiersPTwoZero),
			string(servers.AzureManagedDiskPerformanceTiersPThreeZero),
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{6, 10, 15, 20, 30, 40, 50}},
		65536: {string(servers.AzureManagedDiskPerformanceTiersPSix), &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPOneZero),
			string(servers.AzureManagedDiskPerformanceTiersPOneFive),
			string(servers.AzureManagedDiskPerformanceTiersPTwoZero),
			string(servers.AzureManagedDiskPerformanceTiersPThreeZero),
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{10, 15, 20, 30, 40, 50}},
		131072: {string(servers.AzureManagedDiskPerformanceTiersPOneZero), &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPOneFive),
			string(servers.AzureManagedDiskPerformanceTiersPTwoZero),
			string(servers.AzureManagedDiskPerformanceTiersPThreeZero),
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{15, 20, 30, 40, 50}},
		262144: {string(servers.AzureManagedDiskPerformanceTiersPOneFive), &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPTwoZero),
			string(servers.AzureManagedDiskPerformanceTiersPThreeZero),
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{20, 30, 40, 50}},
		524288: {string(servers.AzureManagedDiskPerformanceTiersPTwoZero), &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPThreeZero),
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{30, 40, 50}},
		1048576: {string(servers.AzureManagedDiskPerformanceTiersPThreeZero), &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{40, 50}},
		2097152: {string(servers.AzureManagedDiskPerformanceTiersPFourZero), &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{50}},
		4194304: {string(servers.AzureManagedDiskPerformanceTiersPFiveZero), nil, nil},
		8388608: {string(servers.AzureManagedDiskPerformanceTiersPSixZero), &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPSevenZero),
			string(servers.AzureManagedDiskPerformanceTiersPEightZero),
		}, &[]int{70, 80}},
		16777216: {string(servers.AzureManagedDiskPerformanceTiersPSevenZero), &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPEightZero),
		}, &[]int{80}},
		33553408: {string(servers.AzureManagedDiskPerformanceTiersPEightZero), nil, nil},
	}

	return storageTiersMappings
}

func StorageTierNameToDefaultStorageMb(tier string) *int {
	var result int

	switch tier {
	case string(servers.AzureManagedDiskPerformanceTiersPFour):
		result = 32768
	case string(servers.AzureManagedDiskPerformanceTiersPSix):
		result = 65536
	case string(servers.AzureManagedDiskPerformanceTiersPOneZero):
		result = 131072
	case string(servers.AzureManagedDiskPerformanceTiersPOneFive):
		result = 262144
	case string(servers.AzureManagedDiskPerformanceTiersPTwoZero):
		result = 524288
	case string(servers.AzureManagedDiskPerformanceTiersPThreeZero):
		result = 1048576
	case string(servers.AzureManagedDiskPerformanceTiersPFourZero):
		result = 2097152
	case string(servers.AzureManagedDiskPerformanceTiersPFiveZero):
		result = 4194304
	case string(servers.AzureManagedDiskPerformanceTiersPSixZero):
		result = 8388608
	case string(servers.AzureManagedDiskPerformanceTiersPSevenZero):
		result = 16777216
	case string(servers.AzureManagedDiskPerformanceTiersPEightZero):
		result = 33553408
	default:
		return nil
	}

	return &result
}
