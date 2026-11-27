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
	secureBootPath := pathPrefix + "security.0.secure_boot_enabled"
	vtpmPath := pathPrefix + "security.0.vtpm_enabled"

	return func(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
		if d.Get(osSKUPath).(string) != string(agentpools.OSSKUAzureContainerLinux) {
			return nil
		}

		if !d.Get(secureBootPath).(bool) {
			return fmt.Errorf("`%s` must be set to `true` when `%s` is `AzureContainerLinux`", secureBootPath, osSKUPath)
		}
		if !d.Get(vtpmPath).(bool) {
			return fmt.Errorf("`%s` must be set to `true` when `%s` is `AzureContainerLinux`", vtpmPath, osSKUPath)
		}
		return nil
	}
}
