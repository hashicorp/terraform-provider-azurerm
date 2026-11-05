// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/ctyhelpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func NodePoolSecurityCustomizeDiff(nodePoolPath string) pluginsdk.CustomizeDiffFunc {
	prefix := ""
	if nodePoolPath != "" {
		prefix = nodePoolPath + ".0."
	}

	return func(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
		rawConfig := d.GetRawConfig()
		if !rawConfig.IsKnown() || rawConfig.IsNull() {
			return nil
		}
		osSKU, err := ctyhelpers.ConstructCtyPath(prefix + "os_sku").Apply(rawConfig)
		if err != nil {
			return fmt.Errorf("retrieving `%sos_sku` from configuration: %+v", prefix, err)
		}
		if !osSKU.IsKnown() || d.Get(prefix+"os_sku").(string) != string(agentpools.OSSKUAzureContainerLinux) {
			return nil
		}

		for _, flag := range []string{"secure_boot_enabled", "vtpm_enabled"} {
			path := prefix + "security.0." + flag
			value, err := ctyhelpers.ConstructCtyPath(path).Apply(rawConfig)
			if err != nil || (value.IsKnown() && (value.IsNull() || !value.True())) {
				return fmt.Errorf("`%s` must be set to `true` when `%sos_sku` is `AzureContainerLinux`", path, prefix)
			}
		}
		return nil
	}
}

func NodePoolSecurityEnabledCustomizeDiff(securityPath string) pluginsdk.CustomizeDiffFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
		if !d.NewValueKnown(securityPath) || len(d.Get(securityPath).([]any)) == 0 {
			return nil
		}

		secureBootPath := securityPath + ".0.secure_boot_enabled"
		vtpmPath := securityPath + ".0.vtpm_enabled"
		if d.NewValueKnown(secureBootPath) && d.NewValueKnown(vtpmPath) && !d.Get(secureBootPath).(bool) && !d.Get(vtpmPath).(bool) {
			return fmt.Errorf("at least one of `%s` or `%s` must be set to `true`; omit `%s` to disable both", secureBootPath, vtpmPath, securityPath)
		}

		return nil
	}
}
