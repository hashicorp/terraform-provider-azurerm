// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_private_link_scope

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-07-01-preview/privatelinkscopesapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorPrivateLinkScopeCreateUpdate,
		Read:   resourceMonitorPrivateLinkScopeRead,
		Update: resourceMonitorPrivateLinkScopeCreateUpdate,
		Delete: resourceMonitorPrivateLinkScopeDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := privatelinkscopesapis.ParsePrivateLinkScopeID(id)
			return err
		}),

		Schema: monitorPrivateLinkScopeSchema(),
	}
}
