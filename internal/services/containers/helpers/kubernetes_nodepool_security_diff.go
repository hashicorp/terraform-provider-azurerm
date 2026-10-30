// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// NodePoolSecurityCustomizeDiff requires explicit Azure Container Linux security settings
// so a plan after refresh cannot try to disable the required flags.
func NodePoolSecurityCustomizeDiff(osSKUPath, securityPath string) pluginsdk.CustomizeDiffFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
		if !d.NewValueKnown(osSKUPath) || d.Get(osSKUPath).(string) != string(agentpools.OSSKUAzureContainerLinux) {
			return nil
		}
		if !d.NewValueKnown(securityPath) {
			return nil
		}

		for _, flag := range []string{"secure_boot_enabled", "vtpm_enabled"} {
			path := securityPath + ".0." + flag
			if d.NewValueKnown(path) && !d.Get(path).(bool) {
				return fmt.Errorf("`%s` must be set to `true` when `%s` is `AzureContainerLinux`", path, osSKUPath)
			}
		}
		return nil
	}
}
