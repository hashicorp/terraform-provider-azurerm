// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

// sku.name is not an Enum in the REST API specs.
// These are the accepted values based on https://learn.microsoft.com/en-us/azure/azure-sql/database/elastic-jobs-overview?view=azuresql#concurrent-capacity-tiers
const (
	SqlJobAgentSkuJA100 string = "JA100"
	SqlJobAgentSkuJA200 string = "JA200"
	SqlJobAgentSkuJA400 string = "JA400"
	SqlJobAgentSkuJA800 string = "JA800"
)

func PossibleValuesForJobAgentSku() []string {
	return []string{
		SqlJobAgentSkuJA100,
		SqlJobAgentSkuJA200,
		SqlJobAgentSkuJA400,
		SqlJobAgentSkuJA800,
	}
}
