// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationruleseventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func EventHubAuthorizationRuleDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: EventHubAuthorizationRuleDataSourceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: eventHubAuthorizationRuleSchemaFrom(map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateEventHubAuthorizationRuleName(),
			},

			"namespace_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateEventHubNamespaceName(),
			},

			"eventhub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateEventHubName(),
			},

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
		}),
	}
}

func EventHubAuthorizationRuleDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	eventHubsClient := meta.(*clients.Client).Eventhub.EventHubsClient
	rulesClient := meta.(*clients.Client).Eventhub.EventHubAuthorizationRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	eventHubName := d.Get("eventhub_name").(string)
	namespaceName := d.Get("namespace_name").(string)

	id := eventhubs.NewEventhubAuthorizationRuleID(subscriptionId, resourceGroup, namespaceName, eventHubName, name)
	resp, err := eventHubsClient.GetAuthorizationRule(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.AuthorizationRuleName)
	d.Set("eventhub_name", id.EventhubName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	localId := authorizationruleseventhubs.NewEventhubAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.EventhubName, id.AuthorizationRuleName)
	keysResp, err := rulesClient.EventHubsListKeys(ctx, localId)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	if model := keysResp.Model; model != nil {
		d.Set("primary_key", model.PrimaryKey)
		d.Set("secondary_key", model.SecondaryKey)
		d.Set("primary_connection_string", model.PrimaryConnectionString)
		d.Set("secondary_connection_string", model.SecondaryConnectionString)
		d.Set("primary_connection_string_alias", model.AliasPrimaryConnectionString)
		d.Set("secondary_connection_string_alias", model.AliasSecondaryConnectionString)
	}

	return nil
}
