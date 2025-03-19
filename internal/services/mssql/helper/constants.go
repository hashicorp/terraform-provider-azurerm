package helper

// sku.name is not an Enum in the REST API specs.
// These are the accepted values based on https://learn.microsoft.com/en-us/azure/azure-sql/database/elastic-jobs-overview?view=azuresql#concurrent-capacity-tiers
const (
	SqlJobAgentSkuNameJA100 string = "JA100"
	SqlJobAgentSkuNameJA200 string = "JA200"
	SqlJobAgentSkuNameJA400 string = "JA400"
	SqlJobAgentSkuNameJA800 string = "JA800"
)

func PossibleValuesForJobAgentSkuName() []string {
	return []string{
		SqlJobAgentSkuNameJA100,
		SqlJobAgentSkuNameJA200,
		SqlJobAgentSkuNameJA400,
		SqlJobAgentSkuNameJA800,
	}
}
