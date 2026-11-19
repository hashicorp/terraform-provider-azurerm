// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_setting

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-05-01/settings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func securityCenterSettingSchema() map[string]*pluginsdk.Schema {
	validSettingName := settings.PossibleValuesForSettingName()

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
