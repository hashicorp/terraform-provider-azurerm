// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/servers"
)

type StorageTiers struct {
	DefaultTier      servers.AzureManagedDiskPerformanceTiers
	ValidTiers       *[]string
	PossibleTiersInt *[]int
}

// Creates a map of valid StorageTiers based on the storage_mb for the PostgreSQL Flexible Server
func InitializeFlexibleServerStorageTierDefaults() map[int]StorageTiers {
	storageTiersMappings := map[int]StorageTiers{
		32768: {servers.AzureManagedDiskPerformanceTiersPFour, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPFour),
			string(servers.AzureManagedDiskPerformanceTiersPSix),
			string(servers.AzureManagedDiskPerformanceTiersPOneZero),
			string(servers.AzureManagedDiskPerformanceTiersPOneFive),
			string(servers.AzureManagedDiskPerformanceTiersPTwoZero),
			string(servers.AzureManagedDiskPerformanceTiersPThreeZero),
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{4, 6, 10, 15, 20, 30, 40, 50}},
		65536: {servers.AzureManagedDiskPerformanceTiersPSix, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPSix),
			string(servers.AzureManagedDiskPerformanceTiersPOneZero),
			string(servers.AzureManagedDiskPerformanceTiersPOneFive),
			string(servers.AzureManagedDiskPerformanceTiersPTwoZero),
			string(servers.AzureManagedDiskPerformanceTiersPThreeZero),
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{6, 10, 15, 20, 30, 40, 50}},
		131072: {servers.AzureManagedDiskPerformanceTiersPOneZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPOneZero),
			string(servers.AzureManagedDiskPerformanceTiersPOneFive),
			string(servers.AzureManagedDiskPerformanceTiersPTwoZero),
			string(servers.AzureManagedDiskPerformanceTiersPThreeZero),
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{10, 15, 20, 30, 40, 50}},
		262144: {servers.AzureManagedDiskPerformanceTiersPOneFive, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPOneFive),
			string(servers.AzureManagedDiskPerformanceTiersPTwoZero),
			string(servers.AzureManagedDiskPerformanceTiersPThreeZero),
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{15, 20, 30, 40, 50}},
		524288: {servers.AzureManagedDiskPerformanceTiersPTwoZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPTwoZero),
			string(servers.AzureManagedDiskPerformanceTiersPThreeZero),
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{20, 30, 40, 50}},
		1048576: {servers.AzureManagedDiskPerformanceTiersPThreeZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPThreeZero),
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{30, 40, 50}},
		2097152: {servers.AzureManagedDiskPerformanceTiersPFourZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPFourZero),
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{40, 50}},
		4193280: {servers.AzureManagedDiskPerformanceTiersPFiveZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{50}},
		4194304: {servers.AzureManagedDiskPerformanceTiersPFiveZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPFiveZero),
		}, &[]int{50}},
		8388608: {servers.AzureManagedDiskPerformanceTiersPSixZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPSixZero),
			string(servers.AzureManagedDiskPerformanceTiersPSevenZero),
			string(servers.AzureManagedDiskPerformanceTiersPEightZero),
		}, &[]int{60, 70, 80}},
		16777216: {servers.AzureManagedDiskPerformanceTiersPSevenZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPSevenZero),
			string(servers.AzureManagedDiskPerformanceTiersPEightZero),
		}, &[]int{70, 80}},
		33553408: {servers.AzureManagedDiskPerformanceTiersPEightZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTiersPEightZero),
		}, &[]int{80}},
	}

	return storageTiersMappings
}
