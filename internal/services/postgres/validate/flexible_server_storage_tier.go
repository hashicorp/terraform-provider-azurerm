// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"math"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
)

type StorageTiers struct {
	DefaultTier      servers.AzureManagedDiskPerformanceTier
	ValidTiers       *[]string
	PossibleTiersInt *[]int
}

// Creates a map of valid StorageTiers based on the storage_mb for the PostgreSQL Flexible Server
func InitializeFlexibleServerStorageTierDefaults() map[int]StorageTiers {
	return map[int]StorageTiers{
		int(math.Exp2(15)): {servers.AzureManagedDiskPerformanceTierPFour, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPFour),
			string(servers.AzureManagedDiskPerformanceTierPSix),
			string(servers.AzureManagedDiskPerformanceTierPOneZero),
			string(servers.AzureManagedDiskPerformanceTierPOneFive),
			string(servers.AzureManagedDiskPerformanceTierPTwoZero),
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{4, 6, 10, 15, 20, 30, 40, 50}},
		int(math.Exp2(16)): {servers.AzureManagedDiskPerformanceTierPSix, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPSix),
			string(servers.AzureManagedDiskPerformanceTierPOneZero),
			string(servers.AzureManagedDiskPerformanceTierPOneFive),
			string(servers.AzureManagedDiskPerformanceTierPTwoZero),
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{6, 10, 15, 20, 30, 40, 50}},
		int(math.Exp2(17)): {servers.AzureManagedDiskPerformanceTierPOneZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPOneZero),
			string(servers.AzureManagedDiskPerformanceTierPOneFive),
			string(servers.AzureManagedDiskPerformanceTierPTwoZero),
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{10, 15, 20, 30, 40, 50}},
		int(math.Exp2(18)): {servers.AzureManagedDiskPerformanceTierPOneFive, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPOneFive),
			string(servers.AzureManagedDiskPerformanceTierPTwoZero),
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{15, 20, 30, 40, 50}},
		int(math.Exp2(19)): {servers.AzureManagedDiskPerformanceTierPTwoZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPTwoZero),
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{20, 30, 40, 50}},
		int(math.Exp2(20)): {servers.AzureManagedDiskPerformanceTierPThreeZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPThreeZero),
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{30, 40, 50}},
		int(math.Exp2(21)): {servers.AzureManagedDiskPerformanceTierPFourZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPFourZero),
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{40, 50}},
		4193280: {servers.AzureManagedDiskPerformanceTierPFiveZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{50}},
		int(math.Exp2(22)): {servers.AzureManagedDiskPerformanceTierPFiveZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPFiveZero),
		}, &[]int{50}},
		int(math.Exp2(23)): {servers.AzureManagedDiskPerformanceTierPSixZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPSixZero),
			string(servers.AzureManagedDiskPerformanceTierPSevenZero),
			string(servers.AzureManagedDiskPerformanceTierPEightZero),
		}, &[]int{60, 70, 80}},
		int(math.Exp2(24)): {servers.AzureManagedDiskPerformanceTierPSevenZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPSevenZero),
			string(servers.AzureManagedDiskPerformanceTierPEightZero),
		}, &[]int{70, 80}},
		33553408: {servers.AzureManagedDiskPerformanceTierPEightZero, &[]string{
			string(servers.AzureManagedDiskPerformanceTierPEightZero),
		}, &[]int{80}},
	}
}
