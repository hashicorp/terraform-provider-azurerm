// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func OrchestratedVirtualMachineScaleSetPublicIPSku(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	publicIpSkus := []string{
		fmt.Sprintf("%s_%s", string(compute.PublicIPAddressSkuNameBasic), string(compute.PublicIPAddressSkuTierRegional)),
		fmt.Sprintf("%s_%s", string(compute.PublicIPAddressSkuNameStandard), string(compute.PublicIPAddressSkuTierRegional)),
		fmt.Sprintf("%s_%s", string(compute.PublicIPAddressSkuNameBasic), string(compute.PublicIPAddressSkuTierGlobal)),
		fmt.Sprintf("%s_%s", string(compute.PublicIPAddressSkuNameStandard), string(compute.PublicIPAddressSkuTierGlobal)),
	}

	for _, sku := range publicIpSkus {
		if v == sku {
			return
		}
	}

	errors = append(errors, fmt.Errorf("%q is not valid, expected to be one of %+v, got %q", key, publicIpSkus, v))
	return
}
