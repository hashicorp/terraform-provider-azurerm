// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_setting

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func securityCenterSettingSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"setting_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(validSettingName, false),
		},
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},
	}
}
