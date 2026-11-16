// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool_registration_info

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualDesktopHostPoolRegistrationInfoSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hostpool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: hostpool.ValidateHostPoolID,
		},

		"expiration_date": {
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.IsRFC3339Time,
			Required:     true,
		},

		"token": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},
	}
}
