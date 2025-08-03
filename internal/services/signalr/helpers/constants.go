// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

// PublicNetworkAccess Enum is missing in the REST API spec
// https://github.com/Azure/azure-rest-api-specs/issues/32981
const (
	PublicNetworkAccessEnabled  string = "Enabled"
	PublicNetworkAccessDisabled string = "Disabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		PublicNetworkAccessEnabled,
		PublicNetworkAccessDisabled,
	}
}

// ServiceMode Enum is missing in the REST API spec
// https://github.com/Azure/azure-rest-api-specs/issues/32994
const (
	SocketIOServiceModeDefault    string = "Default"
	SocketIOServiceModeServerless string = "Serverless"
)

func PossibleValuesForSocketIOServiceMode() []string {
	return []string{
		SocketIOServiceModeDefault,
		SocketIOServiceModeServerless,
	}
}

// ResourceSku.name Enum is missing in the REST API specs
const (
	SkuNameFreeF1     string = "Free_F1"
	SkuNameStandardS1 string = "Standard_S1"
	SkuNamePremiumP1  string = "Premium_P1"
	SkuNamePremiumP2  string = "Premium_P2"
)

func PossibleValuesForSkuName() []string {
	return []string{
		SkuNameFreeF1,
		SkuNameStandardS1,
		SkuNamePremiumP1,
		SkuNamePremiumP2,
	}
}
