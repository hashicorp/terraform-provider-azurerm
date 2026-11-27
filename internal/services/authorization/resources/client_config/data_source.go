// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client_config

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access the configuration of the AzureRM provider.",
		Read:        dataSourceArmClientConfigRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"client_id": {
				Description: "The Azure Client ID (Application Object ID).",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tenant_id": {
				Description: "The Azure Tenant ID.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"subscription_id": {
				Description: "The Azure Subscription ID.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"object_id": {
				Description: "The Azure Object ID.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"principal_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmClientConfigRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client)
	_, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := fmt.Sprintf("clientConfigs/clientId=%s;objectId=%s;subscriptionId=%s;tenantId=%s;principalType=%s", client.Account.ClientId, client.Account.ObjectId, client.Account.SubscriptionId, client.Account.TenantId, client.Account.PrincipalType)
	d.SetId(base64.StdEncoding.EncodeToString([]byte(id)))
	d.Set("client_id", client.Account.ClientId)
	d.Set("object_id", client.Account.ObjectId)
	d.Set("subscription_id", client.Account.SubscriptionId)
	d.Set("tenant_id", client.Account.TenantId)
	d.Set("principal_type", client.Account.PrincipalType)

	return nil
}
