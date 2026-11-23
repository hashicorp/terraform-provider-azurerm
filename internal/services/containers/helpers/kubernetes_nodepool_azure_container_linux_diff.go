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
	osSKUPath := prefix + "os_sku"
	securityPath := prefix + "security.0."

	return func(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
		osSKU, ok := d.GetOkExists(osSKUPath)
		if !ok || osSKU.(string) != string(agentpools.OSSKUAzureContainerLinux) {
			return nil
		}

		rawConfig := d.GetRawConfig()
		if !rawConfig.IsKnown() || rawConfig.IsNull() {
			return nil
		}
		if nodePoolPath != "" {
			nodePools := rawConfig.GetAttr(nodePoolPath)
			if !nodePools.IsKnown() || nodePools.IsNull() || nodePools.LengthInt() == 0 {
				return nil
			}
			rawConfig = nodePools.AsValueSlice()[0]
			if !rawConfig.IsKnown() || rawConfig.IsNull() {
				return nil
			}
		}

		security := rawConfig.GetAttr("security")
		if !security.IsKnown() {
			return nil
		}
		if security.IsNull() || security.LengthInt() == 0 {
			return fmt.Errorf("`%ssecure_boot_enabled` must be set to `true` when `%s` is `AzureContainerLinux`", securityPath, osSKUPath)
		}
		security = security.AsValueSlice()[0]
		if !security.IsKnown() {
			return nil
		}
		if security.IsNull() {
			return fmt.Errorf("`%ssecure_boot_enabled` must be set to `true` when `%s` is `AzureContainerLinux`", securityPath, osSKUPath)
		}

		for _, flag := range []string{"secure_boot_enabled", "vtpm_enabled"} {
			path := securityPath + flag
			value, ok := d.GetOkExists(path)
			if ok && value.(bool) {
				continue
			}
			if !ok && !security.GetAttr(flag).IsNull() {
				continue
			}
			return fmt.Errorf("`%s` must be set to `true` when `%s` is `AzureContainerLinux`", path, osSKUPath)
		}
		return nil
	}
}
