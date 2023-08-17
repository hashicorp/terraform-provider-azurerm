// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	identity2 "github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/automationaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceAutomationAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAutomationAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"primary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"private_endpoint_connection": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"hybrid_service_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAutomationAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	iclient := meta.(*clients.Client).Automation.AgentRegistrationInfoClient
	client := meta.(*clients.Client).Automation.AutomationAccount
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := automationaccount.NewAutomationAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retreiving %s: %+v", id, err)
	}
	d.SetId(id.ID())

	iresp, err := iclient.Get(ctx, id.ResourceGroupName, id.AutomationAccountName)
	if err != nil {
		if utils.ResponseWasNotFound(iresp.Response) {
			return fmt.Errorf("%q Account Registration Information was not found", id)
		}
		return fmt.Errorf("retreiving Automation Account Registration Information %s: %+v", id, err)
	}
	if iresp.Keys != nil {
		d.Set("primary_key", iresp.Keys.Primary)
		d.Set("secondary_key", iresp.Keys.Secondary)
	}

	identity, err := identity2.FlattenSystemAndUserAssignedMap(resp.Model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	d.Set("endpoint", iresp.Endpoint)
	if resp.Model != nil && resp.Model.Properties != nil {
		d.Set("private_endpoint_connection", flattenPrivateEndpointConnections(resp.Model.Properties.PrivateEndpointConnections))
		d.Set("hybrid_service_url", resp.Model.Properties.AutomationHybridServiceUrl)
	}
	return nil
}

func flattenPrivateEndpointConnections(conns *[]automationaccount.PrivateEndpointConnection) (res []interface{}) {
	if conns == nil || len(*conns) == 0 {
		return
	}
	for _, con := range *conns {
		res = append(res, map[string]interface{}{
			"id":   con.Id,
			"name": con.Name,
		})
	}
	return res
}
