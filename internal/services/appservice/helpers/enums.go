// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

// publicNetworkAccess Enums are missing in api-specs
// https://github.com/Azure/azure-rest-api-specs/issues/24680

const (
	PublicNetworkAccessEnabled  string = "Enabled"
	PublicNetworkAccessDisabled string = "Disabled"
)

const (
	ValidationTypeTXT   = "dns-txt-token"
	ValidationTypeCName = "cname-delegation"
)
