// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceArmClientConfig() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmClientConfigRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"client_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"subscription_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"object_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmClientConfigRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client)
	_, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := fmt.Sprintf("clientConfigs/clientId=%s;objectId=%s;subscriptionId=%s;tenantId=%s", client.Account.ClientId, client.Account.ObjectId, client.Account.SubscriptionId, client.Account.TenantId)
	d.SetId(base64.StdEncoding.EncodeToString([]byte(id)))
	d.Set("client_id", client.Account.ClientId)
	d.Set("object_id", client.Account.ObjectId)
	d.Set("subscription_id", client.Account.SubscriptionId)
	d.Set("tenant_id", client.Account.TenantId)

	return nil
}
