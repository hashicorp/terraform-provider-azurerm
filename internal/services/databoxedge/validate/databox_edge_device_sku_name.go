// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/devices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

func DataboxEdgeDeviceSkuName(v interface{}, k string) (warnings []string, errors []error) {
	validSku := false
	validTier := false

	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	skuParts := strings.Split(value, "-")
	validSkus := devices.PossibleValuesForSkuName()
	validTiers := devices.PossibleValuesForSkuTier()

	// Validate the SKU Name section
	if slices.Contains(validSkus, skuParts[0]) {
		validSku = true
	}

	if len(skuParts) > 1 {
		// Validate the SKU Tier section
		if slices.Contains(validTiers, skuParts[1]) {
			validTier = true
		}
	}

	if !validSku {
		errors = append(errors, fmt.Errorf("expected %q %q segment to be one of [%s], got %q", k, "name", azure.QuotedStringSlice(validSkus), value))
	}
	if !validTier {
		errors = append(errors, fmt.Errorf("expected %q %q segment to be one of [%s], got %q", k, "tier", azure.QuotedStringSlice(validTiers), value))
	}

	return warnings, errors
}
