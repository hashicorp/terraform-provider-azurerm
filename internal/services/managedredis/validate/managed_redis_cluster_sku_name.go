package validate

import (
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
)

func PossibleValuesForSkuName() []string {
	validSkus := make([]string, 0, len(redisenterprise.PossibleValuesForSkuName()))
	for _, sku := range redisenterprise.PossibleValuesForSkuName() {
		// Enterprise_ and EnterpriseFlash_ SKUs are retained in the SDK enums, but no longer supported by this resource
		// https://learn.microsoft.com/azure/redis/migrate/migrate-overview
		if strings.HasPrefix(sku, "Enterprise_") || strings.HasPrefix(sku, "EnterpriseFlash_") {
			continue
		}
		validSkus = append(validSkus, sku)
	}
	return validSkus
}
