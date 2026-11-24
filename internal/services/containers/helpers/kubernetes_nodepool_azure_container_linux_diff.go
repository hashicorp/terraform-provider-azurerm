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

		// The SDK reports no count for both an omitted block and an unknown block.
		if _, ok := d.GetOkExists(prefix + "security.#"); !ok {
			config := d.GetRawConfig()
			if !config.IsKnown() || config.IsNull() {
				return nil
			}
			if nodePoolPath != "" {
				nodePools := config.GetAttr(nodePoolPath)
				if !nodePools.IsKnown() || nodePools.IsNull() || nodePools.LengthInt() == 0 {
					return nil
				}
				config = nodePools.AsValueSlice()[0]
				if !config.IsKnown() || config.IsNull() {
					return nil
				}
			}

			security := config.GetAttr("security")
			if security.IsKnown() && (security.IsNull() || security.LengthInt() == 0) {
				return fmt.Errorf("`%ssecure_boot_enabled` must be set to `true` when `%s` is `AzureContainerLinux`", securityPath, osSKUPath)
			}
			return nil
		}

		for _, flag := range []string{"secure_boot_enabled", "vtpm_enabled"} {
			path := securityPath + flag
			value, ok := d.GetOkExists(path)
			if ok && !value.(bool) {
				return fmt.Errorf("`%s` must be set to `true` when `%s` is `AzureContainerLinux`", path, osSKUPath)
			}
		}
		return nil
	}
}
