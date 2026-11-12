// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/ctyhelpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func NodePoolSecurityEnabledCustomizeDiff(securityPath string) pluginsdk.CustomizeDiffFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
		rawConfig := d.GetRawConfig()
		if !rawConfig.IsKnown() || rawConfig.IsNull() {
			return nil
		}

		security, err := ctyhelpers.ConstructCtyPath(securityPath).Apply(rawConfig)
		if err != nil {
			return fmt.Errorf("retrieving `%s` from configuration: %+v", securityPath, err)
		}
		if !security.IsKnown() || security.IsNull() || security.LengthInt() == 0 {
			return nil
		}
		security = security.AsValueSlice()[0]
		if !security.IsKnown() || security.IsNull() {
			return nil
		}

		settings := security.AsValueMap()
		secureBoot := settings["secure_boot_enabled"]
		vtpm := settings["vtpm_enabled"]
		if !secureBoot.IsKnown() || !vtpm.IsKnown() {
			return nil
		}
		if (secureBoot.IsNull() || !secureBoot.True()) && (vtpm.IsNull() || !vtpm.True()) {
			return fmt.Errorf("at least one of `%s.0.secure_boot_enabled` or `%s.0.vtpm_enabled` must be set to `true`; omit `%s` to disable both", securityPath, securityPath, securityPath)
		}

		return nil
	}
}
