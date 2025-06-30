// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

func ManagedRedisClusterSkuName(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	skuParts := strings.Split(value, "-")

	if len(skuParts) > 2 {
		errors = append(errors, fmt.Errorf("expected %q to have at most 2 segments separated by '-', got %q", k, value))
		return warnings, errors
	}

	skuName := skuParts[0]
	validSkuName := false
	validSkuNames := redisenterprise.PossibleValuesForSkuName()
	for _, str := range validSkuNames {
		if skuName == str {
			validSkuName = true
			break
		}
	}
	if !validSkuName {
		errors = append(errors, fmt.Errorf("expected %q sku name segment to be one of [%s], got %q", k, azure.QuotedStringSlice(validSkuNames), value))
		return warnings, errors
	}

	enterpriseOrEnterpriseFlash := strings.HasPrefix(skuName, "Enterprise_") || strings.HasPrefix(skuName, "EnterpriseFlash_")
	if len(skuParts) == 1 && enterpriseOrEnterpriseFlash {
		errors = append(errors, fmt.Errorf("expected %q to have a capacity segment for Enterprise_ and EnterpriseFlash_ SKUs, got %q", k, value))
		return warnings, errors
	}

	if len(skuParts) == 2 {
		if !enterpriseOrEnterpriseFlash {
			errors = append(errors, fmt.Errorf("capacity segment is only valid for Enterprise_ and EnterpriseFlash_ SKUs: %q", value))
			return warnings, errors
		}

		skuCapacity := skuParts[1]
		validCapacityHints := "2, 4, 6, ..."
		isEnterpriseFlash := strings.HasPrefix(skuName, "EnterpriseFlash_")
		if isEnterpriseFlash {
			validCapacityHints = "3, 9, 15, ..."
		}
		i, err := strconv.ParseInt(skuCapacity, 10, 32)
		if err != nil {
			errors = append(errors, fmt.Errorf("expected %q %q segment to fit into an int32 value, got %q", k, "capacity", skuCapacity))
			return warnings, errors
		}
		skuCapacityInt := int32(i)
		validCapacity := false
		if isEnterpriseFlash {
			if skuCapacityInt%6 == 3 {
				validCapacity = true
			}
		} else if skuCapacityInt%2 == 0 {
			validCapacity = true
		}
		if !validCapacity {
			errors = append(errors, fmt.Errorf("expected %q %q segment to be one of [%s], got %q", k, "capacity", validCapacityHints, value))
			return warnings, errors
		}
	}

	return warnings, errors
}
