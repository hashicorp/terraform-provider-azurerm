// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

<<<<<<<< HEAD:internal/services/mssql/helpers/constants.go
package helpers
========
package mssql_job_agent
>>>>>>>> bb31be9448 (re run mssql for conflicts resolution):internal/services/mssql/resources/mssql_job_agent/constants.go

// sku.name is not an Enum in the REST API specs.
// These are the accepted values based on https://learn.microsoft.com/azure/azure-sql/database/elastic-jobs-overview?view=azuresql#concurrent-capacity-tiers
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
