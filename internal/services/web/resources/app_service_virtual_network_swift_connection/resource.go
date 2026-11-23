// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_virtual_network_swift_connection

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceVirtualNetworkSwiftConnectionCreate,
		Read:   resourceAppServiceVirtualNetworkSwiftConnectionRead,
		Update: resourceAppServiceVirtualNetworkSwiftConnectionUpdate,
		Delete: resourceAppServiceVirtualNetworkSwiftConnectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualNetworkSwiftConnectionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: appServiceVirtualNetworkSwiftConnectionSchema(),
	}
}
