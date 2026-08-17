// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/applicationgateways"
)

var DeprecatedV1SkuNames = []string{
	string(applicationgateways.ApplicationGatewaySkuNameStandardSmall),
	string(applicationgateways.ApplicationGatewaySkuNameStandardMedium),
	string(applicationgateways.ApplicationGatewaySkuNameStandardLarge),
	string(applicationgateways.ApplicationGatewaySkuNameWAFMedium),
	string(applicationgateways.ApplicationGatewaySkuNameWAFMedium),
}

var DeprecatedV1SkuTiers = []string{
	string(applicationgateways.ApplicationGatewayTierStandard),
	string(applicationgateways.ApplicationGatewayTierWAF),
}
