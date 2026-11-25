// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func NodePoolAzureContainerLinuxCustomizeDiff(pathPrefix string) pluginsdk.CustomizeDiffFunc {
	osSKUPath := pathPrefix + "os_sku"
	securityPath := pathPrefix + "security"

	return func(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
		osSKU, ok := d.GetOkExists(osSKUPath)
		if !ok || osSKU.(string) != string(agentpools.OSSKUAzureContainerLinux) {
			return nil
		}

		if !d.NewValueKnown(securityPath) {
			return nil
		}

		for _, flag := range []string{"secure_boot_enabled", "vtpm_enabled"} {
			path := securityPath + ".0." + flag
			value, ok := d.GetOkExists(path)
			if ok && !value.(bool) {
				return fmt.Errorf("`%s` must be set to `true` when `%s` is `AzureContainerLinux`", path, osSKUPath)
			}
		}
		return nil
	}
}
