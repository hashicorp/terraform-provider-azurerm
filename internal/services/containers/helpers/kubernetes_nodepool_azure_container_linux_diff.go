// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func NodePoolAzureContainerLinuxCustomizeDiff(nodePoolPath string) pluginsdk.CustomizeDiffFunc {
	prefix := ""
	if nodePoolPath != "" {
		prefix = nodePoolPath + ".0."
	}

	return func(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
		rawConfig := d.GetRawConfig()
		if !rawConfig.IsKnown() || rawConfig.IsNull() {
			return nil
		}
		if nodePoolPath != "" {
			nodePools := rawConfig.GetAttr(nodePoolPath)
			if !nodePools.IsKnown() || nodePools.IsNull() || len(nodePools.AsValueSlice()) == 0 {
				return nil
			}
			rawConfig = nodePools.AsValueSlice()[0]
			if !rawConfig.IsKnown() || rawConfig.IsNull() {
				return nil
			}
		}

		osSKU := rawConfig.GetAttr("os_sku")
		if !osSKU.IsKnown() || d.Get(prefix+"os_sku").(string) != string(agentpools.OSSKUAzureContainerLinux) {
			return nil
		}

		security := rawConfig.GetAttr("security")
		if !security.IsKnown() {
			return nil
		}
		if security.IsNull() || len(security.AsValueSlice()) == 0 {
			return fmt.Errorf("`%ssecurity.0.secure_boot_enabled` must be set to `true` when `%sos_sku` is `AzureContainerLinux`", prefix, prefix)
		}
		security = security.AsValueSlice()[0]
		if !security.IsKnown() {
			return nil
		}

		for _, flag := range []string{"secure_boot_enabled", "vtpm_enabled"} {
			path := prefix + "security.0." + flag
			if security.IsNull() {
				return fmt.Errorf("`%s` must be set to `true` when `%sos_sku` is `AzureContainerLinux`", path, prefix)
			}
			value := security.GetAttr(flag)
			if value.IsKnown() && (value.IsNull() || !value.True()) {
				return fmt.Errorf("`%s` must be set to `true` when `%sos_sku` is `AzureContainerLinux`", path, prefix)
			}
		}
		return nil
	}
}
