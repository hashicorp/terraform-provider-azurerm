// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
)

func OrchestratedVirtualMachineScaleSetPublicIPSku(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	publicIpSkus := []string{
		fmt.Sprintf("%s_%s", string(virtualmachinescalesets.PublicIPAddressSkuNameBasic), string(virtualmachinescalesets.PublicIPAddressSkuTierRegional)),
		fmt.Sprintf("%s_%s", string(virtualmachinescalesets.PublicIPAddressSkuNameStandard), string(virtualmachinescalesets.PublicIPAddressSkuTierRegional)),
		fmt.Sprintf("%s_%s", string(virtualmachinescalesets.PublicIPAddressSkuNameBasic), string(virtualmachinescalesets.PublicIPAddressSkuTierGlobal)),
		fmt.Sprintf("%s_%s", string(virtualmachinescalesets.PublicIPAddressSkuNameStandard), string(virtualmachinescalesets.PublicIPAddressSkuTierGlobal)),
	}

	for _, sku := range publicIpSkus {
		if v == sku {
			return
		}
	}

	errors = append(errors, fmt.Errorf("%q is not valid, expected to be one of %+v, got %q", key, publicIpSkus, v))
	return
}
