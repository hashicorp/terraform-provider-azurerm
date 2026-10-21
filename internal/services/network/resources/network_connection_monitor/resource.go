// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_connection_monitor

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectionmonitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkConnectionMonitorCreate,
		Read:   resourceNetworkConnectionMonitorRead,
		Update: resourceNetworkConnectionMonitorUpdate,
		Delete: resourceNetworkConnectionMonitorDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := connectionmonitors.ParseConnectionMonitorID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceNetworkConnectionMonitorSchema(),
	}
}
