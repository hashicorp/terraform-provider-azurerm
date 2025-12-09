package validate

import (
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
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
	slices.Sort(validSkus)
	return validSkus
}

func SKUsSupportingGeoReplication() []string {
	skus := make([]string, 0, len(PossibleValuesForSkuName()))
	for _, sku := range PossibleValuesForSkuName() {
		if sku == string(redisenterprise.SkuNameBalancedBZero) ||
			sku == string(redisenterprise.SkuNameBalancedBOne) {
			continue
		}
		skus = append(skus, sku)
	}
	return skus
}
