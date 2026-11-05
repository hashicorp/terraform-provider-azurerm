// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

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
