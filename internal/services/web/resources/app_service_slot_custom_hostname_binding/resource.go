// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_slot_custom_hostname_binding

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const appServiceSlotCustomHostnameBindingResourceName = "azurerm_app_service_slot_custom_hostname_binding"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceSlotCustomHostnameBindingCreate,
		Read:   resourceAppServiceSlotCustomHostnameBindingRead,
		Delete: resourceAppServiceSlotCustomHostnameBindingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webapps.ParseSlotHostNameBindingID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: appServiceSlotCustomHostnameBindingSchema(),
	}
}
