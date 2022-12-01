package packetcorecontrolplane

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingSku string

const (
	BillingSkuEdgeSiteFourGBPS       BillingSku = "EdgeSite4GBPS"
	BillingSkuEdgeSiteThreeGBPS      BillingSku = "EdgeSite3GBPS"
	BillingSkuEdgeSiteTwoGBPS        BillingSku = "EdgeSite2GBPS"
	BillingSkuEvaluationPackage      BillingSku = "EvaluationPackage"
	BillingSkuFlagshipStarterPackage BillingSku = "FlagshipStarterPackage"
	BillingSkuLargePackage           BillingSku = "LargePackage"
	BillingSkuMediumPackage          BillingSku = "MediumPackage"
)

func PossibleValuesForBillingSku() []string {
	return []string{
		string(BillingSkuEdgeSiteFourGBPS),
		string(BillingSkuEdgeSiteThreeGBPS),
		string(BillingSkuEdgeSiteTwoGBPS),
		string(BillingSkuEvaluationPackage),
		string(BillingSkuFlagshipStarterPackage),
		string(BillingSkuLargePackage),
		string(BillingSkuMediumPackage),
	}
}

func parseBillingSku(input string) (*BillingSku, error) {
	vals := map[string]BillingSku{
		"edgesite4gbps":          BillingSkuEdgeSiteFourGBPS,
		"edgesite3gbps":          BillingSkuEdgeSiteThreeGBPS,
		"edgesite2gbps":          BillingSkuEdgeSiteTwoGBPS,
		"evaluationpackage":      BillingSkuEvaluationPackage,
		"flagshipstarterpackage": BillingSkuFlagshipStarterPackage,
		"largepackage":           BillingSkuLargePackage,
		"mediumpackage":          BillingSkuMediumPackage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BillingSku(input)
	return &out, nil
}

type CoreNetworkType string

const (
	CoreNetworkTypeEPC    CoreNetworkType = "EPC"
	CoreNetworkTypeFiveGC CoreNetworkType = "5GC"
)

func PossibleValuesForCoreNetworkType() []string {
	return []string{
		string(CoreNetworkTypeEPC),
		string(CoreNetworkTypeFiveGC),
	}
}

func parseCoreNetworkType(input string) (*CoreNetworkType, error) {
	vals := map[string]CoreNetworkType{
		"epc": CoreNetworkTypeEPC,
		"5gc": CoreNetworkTypeFiveGC,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CoreNetworkType(input)
	return &out, nil
}

type PlatformType string

const (
	PlatformTypeAKSNegativeHCI PlatformType = "AKS-HCI"
	PlatformTypeBaseVM         PlatformType = "BaseVM"
)

func PossibleValuesForPlatformType() []string {
	return []string{
		string(PlatformTypeAKSNegativeHCI),
		string(PlatformTypeBaseVM),
	}
}

func parsePlatformType(input string) (*PlatformType, error) {
	vals := map[string]PlatformType{
		"aks-hci": PlatformTypeAKSNegativeHCI,
		"basevm":  PlatformTypeBaseVM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PlatformType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateDeleted   ProvisioningState = "Deleted"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUnknown   ProvisioningState = "Unknown"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnknown),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":  ProvisioningStateAccepted,
		"canceled":  ProvisioningStateCanceled,
		"deleted":   ProvisioningStateDeleted,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"unknown":   ProvisioningStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
