// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
)

type StorageTiers struct {
	DefaultTier      servers.AzureManagedDiskPerformanceTier
	ValidTiers       *[]string
	PossibleTiersInt *[]int
}

// Creates a map of valid StorageTiers based on the storage_mb for the PostgreSQL Flexible Server
func InitializeFlexibleServerStorageTierDefaults() map[int]StorageTiers {
	storageTiersMappings := map[int]StorageTiers{
		32768: {servers.AzureManagedDiskPerformanceTierPFour, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPFour),
			string(servers.AzureManagedDiskPerformanceTierPSix),
			string(servers.AzureManagedDiskPerformanceTierPOneZero),
			string(servers.AzureManagedDiskPerformanceTierPOneFive),
			string(servers.AzureManagedDiskPerformanceTierPTwoZero),
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{4, 6, 10, 15, 20, 30, 40, 50}},
		65536: {servers.AzureManagedDiskPerformanceTierPSix, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPSix),
			string(servers.AzureManagedDiskPerformanceTierPOneZero),
			string(servers.AzureManagedDiskPerformanceTierPOneFive),
			string(servers.AzureManagedDiskPerformanceTierPTwoZero),
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{6, 10, 15, 20, 30, 40, 50}},
		131072: {servers.AzureManagedDiskPerformanceTierPOneZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPOneZero),
			string(servers.AzureManagedDiskPerformanceTierPOneFive),
			string(servers.AzureManagedDiskPerformanceTierPTwoZero),
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{10, 15, 20, 30, 40, 50}},
		262144: {servers.AzureManagedDiskPerformanceTierPOneFive, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPOneFive),
			string(servers.AzureManagedDiskPerformanceTierPTwoZero),
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{15, 20, 30, 40, 50}},
		524288: {servers.AzureManagedDiskPerformanceTierPTwoZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPTwoZero),
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{20, 30, 40, 50}},
		1048576: {servers.AzureManagedDiskPerformanceTierPThreeZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{30, 40, 50}},
		2097152: {servers.AzureManagedDiskPerformanceTierPFourZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{40, 50}},
		4193280: {servers.AzureManagedDiskPerformanceTierPFiveZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{50}},
		4194304: {servers.AzureManagedDiskPerformanceTierPFiveZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{50}},
		8388608: {servers.AzureManagedDiskPerformanceTierPSixZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPSixZero),
			string(servers.AzureManagedDiskPerformanceTierPSevenZero),
			string(servers.AzureManagedDiskPerformanceTierPEightZero),
		}, &[]int{60, 70, 80}},
		16777216: {servers.AzureManagedDiskPerformanceTierPSevenZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPSevenZero),
			string(servers.AzureManagedDiskPerformanceTierPEightZero),
		}, &[]int{70, 80}},
		33553408: {servers.AzureManagedDiskPerformanceTierPEightZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPEightZero),
		}, &[]int{80}},
	}

	return storageTiersMappings
}
