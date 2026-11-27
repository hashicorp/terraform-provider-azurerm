// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client_config

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
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
		},
	}
}

func dataSourceArmClientConfigRead(d *pluginsdk.ResourceData, meta any) error {
	_, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := fmt.Sprintf("clientConfigs/clientId=%s;objectId=%s;subscriptionId=%s;tenantId=%s", account.From(meta).ClientId, account.From(meta).ObjectId, account.From(meta).SubscriptionId, account.From(meta).TenantId)
	d.SetId(base64.StdEncoding.EncodeToString([]byte(id)))
	d.Set("client_id", account.From(meta).ClientId)
	d.Set("object_id", account.From(meta).ObjectId)
	d.Set("subscription_id", account.From(meta).SubscriptionId)
	d.Set("tenant_id", account.From(meta).TenantId)

	return nil
}
