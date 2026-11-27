// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_slot_custom_hostname_binding

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func appServiceSlotCustomHostnameBindingSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"app_service_slot_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: webapps.ValidateSlotID,
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"ssl_state": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapps.SslStateIPBasedEnabled),
				string(webapps.SslStateSniEnabled),
			}, false),
		},

		"thumbprint": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"virtual_ip": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
