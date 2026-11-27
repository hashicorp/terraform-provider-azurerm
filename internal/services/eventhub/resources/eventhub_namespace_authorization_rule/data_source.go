// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_authorization_rule

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationrulesnamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	eventhubClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an Authorization Rule for an Event Hub Namespace.",
		Read:        EventHubNamespaceDataSourceAuthorizationRuleRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the EventHub Authorization Rule resource.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateEventHubAuthorizationRuleName(),
			},

			"namespace_name": {
				Description:  "The name of the EventHub Namespace.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateEventHubNamespaceName(),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"listen": {
				Description: "Whether this Authorization Rule has permissions to Listen to the Event Hub.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"manage": {
				Description: "Whether this Authorization Rule has permissions to Manage to the Event Hub.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"primary_connection_string": {
				Description: "The Primary Connection String for the Event Hubs authorization Rule.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"primary_connection_string_alias": {
				Description: "The alias of the Primary Connection String for the Event Hubs authorization Rule.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"primary_key": {
				Description: "The Primary Key for the Event Hubs authorization Rule.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"secondary_connection_string": {
				Description: "The Secondary Connection String for the Event Hubs authorization Rule.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"secondary_connection_string_alias": {
				Description: "The alias of the Secondary Connection String for the Event Hubs authorization Rule.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"secondary_key": {
				Description: "The Secondary Key for the Event Hubs authorization Rule.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"send": {
				Description: "Whether this Authorization Rule has permissions to Send to the Event Hub.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},
		},
	}
}

func EventHubNamespaceDataSourceAuthorizationRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := eventhubClient.From(meta).NamespaceAuthorizationRulesClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := authorizationrulesnamespaces.NewAuthorizationRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("name").(string))
	resp, err := client.NamespacesGetAuthorizationRule(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.AuthorizationRuleName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			listen, send, manage := flattenEventHubNamespaceAuthorizationRuleRights(props.Rights)
			d.Set("manage", manage)
			d.Set("listen", listen)
			d.Set("send", send)
		}
	}

	keysResp, err := client.NamespacesListKeys(ctx, id)
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
