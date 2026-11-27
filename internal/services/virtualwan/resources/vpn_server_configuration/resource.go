// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_server_configuration

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVPNServerConfigurationCreate,
		Read:   resourceVPNServerConfigurationRead,
		Update: resourceVPNServerConfigurationUpdate,
		Delete: resourceVPNServerConfigurationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualwans.ParseVpnServerConfigurationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: vpnServerConfigurationSchema(),
	}
}
