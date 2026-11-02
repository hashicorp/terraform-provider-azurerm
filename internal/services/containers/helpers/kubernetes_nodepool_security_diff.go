// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func NodePoolSecurityCustomizeDiff(nodePoolPath string) pluginsdk.CustomizeDiffFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
		rawConfig := d.GetRawConfig()
		if !rawConfig.IsKnown() || rawConfig.IsNull() {
			return nil
		}

		prefix := ""
		if nodePoolPath != "" {
			nodePools := rawConfig.AsValueMap()[nodePoolPath]
			if !nodePools.IsKnown() || nodePools.IsNull() || nodePools.LengthInt() == 0 {
				return nil
			}
			rawConfig = nodePools.AsValueSlice()[0]
			if !rawConfig.IsKnown() || rawConfig.IsNull() {
				return nil
			}
			prefix = nodePoolPath + ".0."
		}

		config := rawConfig.AsValueMap()
		osSKU := config["os_sku"]
		if !osSKU.IsKnown() || osSKU.IsNull() || osSKU.AsString() != string(agentpools.OSSKUAzureContainerLinux) {
			return nil
		}

		security := config["security"]
		if !security.IsKnown() {
			return nil
		}
		if security.IsNull() || security.LengthInt() == 0 {
			return fmt.Errorf("`%ssecurity.0.secure_boot_enabled` must be set to `true` when `%sos_sku` is `AzureContainerLinux`", prefix, prefix)
		}
		security = security.AsValueSlice()[0]
		if !security.IsKnown() {
			return nil
		}

		for _, flag := range []string{"secure_boot_enabled", "vtpm_enabled"} {
			if security.IsNull() {
				return fmt.Errorf("`%ssecurity.0.%s` must be set to `true` when `%sos_sku` is `AzureContainerLinux`", prefix, flag, prefix)
			}
			value := security.AsValueMap()[flag]
			if value.IsKnown() && (value.IsNull() || !value.True()) {
				return fmt.Errorf("`%ssecurity.0.%s` must be set to `true` when `%sos_sku` is `AzureContainerLinux`", prefix, flag, prefix)
			}
		}
		return nil
	}
}
