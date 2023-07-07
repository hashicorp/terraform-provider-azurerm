// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2022-01-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

// RedisEnterpriseClusterSkuName - validates if passed input string contains a valid Redis Enterprise Cluster Sku
func RedisEnterpriseClusterSkuName(v interface{}, k string) (warnings []string, errors []error) {
	validSku := false
	validCapacity := false

	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	skuParts := strings.Split(value, "-")
	validSkus := getValidRedisEnterpriseClusterSkus()
	validValues := "2, 4, 6, ..."
	// Validate the SKU Name section
	for _, str := range validSkus {
		if skuParts[0] == str {
			validSku = true
			break
		}
	}

	isFlash := strings.Contains(skuParts[0], "Flash")

	if isFlash {
		validValues = "3, 9, 15, ..."
	}

	if len(skuParts) > 1 {
		i, err := strconv.ParseInt(skuParts[1], 10, 32)
		if err != nil {
			errors = append(errors, fmt.Errorf("expected %q %q segment to fit into an int32 value, got %q", k, "capacity", skuParts[1]))
			return warnings, errors
		}

		skuCapacity := int32(i)

		if isFlash {
			if skuCapacity%6 == 3 {
				validCapacity = true
			}
		} else if skuCapacity%2 == 0 {
			validCapacity = true
		}
	}

	if !validSku {
		errors = append(errors, fmt.Errorf("expected %q %q segment to be one of [%s], got %q", k, "name", azure.QuotedStringSlice(validSkus), value))
	}

	if !validCapacity {
		errors = append(errors, fmt.Errorf("expected %q %q segment to be one of [%s], got %q", k, "capacity", validValues, value))
	}

	return warnings, errors
}

func getValidRedisEnterpriseClusterSkus() []string {
	return []string{
		string(redisenterprise.SkuNameEnterpriseEOneZero),
		string(redisenterprise.SkuNameEnterpriseETwoZero),
		string(redisenterprise.SkuNameEnterpriseEFiveZero),
		string(redisenterprise.SkuNameEnterpriseEOneHundred),
		string(redisenterprise.SkuNameEnterpriseFlashFThreeHundred),
		string(redisenterprise.SkuNameEnterpriseFlashFSevenHundred),
		string(redisenterprise.SkuNameEnterpriseFlashFOneFiveHundred),
	}
}
