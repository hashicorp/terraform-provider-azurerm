// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func OrchestratedVirtualMachineScaleSetPublicIPSku(input interface{}, key string) ([]string, []error) {
	publicIpSkus := []string{
		fmt.Sprintf("%s_%s", string(virtualmachinescalesets.PublicIPAddressSkuNameBasic), string(virtualmachinescalesets.PublicIPAddressSkuTierRegional)),
		fmt.Sprintf("%s_%s", string(virtualmachinescalesets.PublicIPAddressSkuNameStandard), string(virtualmachinescalesets.PublicIPAddressSkuTierRegional)),
		fmt.Sprintf("%s_%s", string(virtualmachinescalesets.PublicIPAddressSkuNameBasic), string(virtualmachinescalesets.PublicIPAddressSkuTierGlobal)),
		fmt.Sprintf("%s_%s", string(virtualmachinescalesets.PublicIPAddressSkuNameStandard), string(virtualmachinescalesets.PublicIPAddressSkuTierGlobal)),
	}

	return validation.StringInSlice(publicIpSkus, false)(input, key)
}
