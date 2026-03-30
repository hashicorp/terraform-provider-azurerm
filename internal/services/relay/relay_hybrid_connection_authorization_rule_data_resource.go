// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package relay

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = RelayHybridConnectionAuthorizationRuleDataResource{}

type RelayHybridConnectionAuthorizationRuleDataResource struct{}

type RelayHybridConnectionAuthorizationRuleDataResourceModel struct {
	Name                      string `tfschema:"name"`
	ResourceGroupName         string `tfschema:"resource_group_name"`
	RelayNamespaceName        string `tfschema:"relay_namespace_name"`
	HybridConnectionName      string `tfschema:"hybrid_connection_name"`
	PrimaryConnectionString   string `tfschema:"primary_connection_string"`
	SecondaryConnectionString string `tfschema:"secondary_connection_string"`
	PrimaryKey                string `tfschema:"primary_key"`
	SecondaryKey              string `tfschema:"secondary_key"`
	Listen                    bool   `tfschema:"listen"`
	Send                      bool   `tfschema:"send"`
	Manage                    bool   `tfschema:"manage"`
}

func (r RelayHybridConnectionAuthorizationRuleDataResource) Arguments() map[string]*pluginsdk.Schema {
	return authorizationRuleArgumentsFrom(map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"relay_namespace_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"hybrid_connection_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	})
}

func (r RelayHybridConnectionAuthorizationRuleDataResource) Attributes() map[string]*pluginsdk.Schema {
	return authorizationRuleAttributesFrom(map[string]*pluginsdk.Schema{})
}

func (RelayHybridConnectionAuthorizationRuleDataResource) ModelObject() interface{} {
	return &RelayHybridConnectionAuthorizationRuleDataResourceModel{}
}

func (r RelayHybridConnectionAuthorizationRuleDataResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(5 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			log.Printf("[INFO] preparing arguments for Relay HybridConnection Authorization Rule creation.")

			var config RelayHybridConnectionAuthorizationRuleDataResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := hybridconnections.NewHybridConnectionAuthorizationRuleID(subscriptionId, config.ResourceGroupName, config.RelayNamespaceName, config.HybridConnectionName, config.Name)

			resp, err := client.GetAuthorizationRule(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving: %s: %+v", id, err)
			}

			keysResp, err := client.ListKeys(ctx, id)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", id, err)
			}

			state := RelayHybridConnectionAuthorizationRuleDataResourceModel{}

			state.ResourceGroupName = id.ResourceGroupName
			state.RelayNamespaceName = id.NamespaceName
			state.HybridConnectionName = id.HybridConnectionName

			if model := resp.Model; model != nil {
				state.Name = pointer.From(model.Name)
				listen, send, manage := flattenHybridConnectionAuthorizationRuleRights(model.Properties.Rights)
				state.Listen = listen
				state.Send = send
				state.Manage = manage
			}

			state.PrimaryConnectionString = pointer.From(keysResp.Model.PrimaryConnectionString)
			state.PrimaryKey = pointer.From(keysResp.Model.PrimaryKey)
			state.SecondaryConnectionString = pointer.From(keysResp.Model.SecondaryConnectionString)
			state.SecondaryKey = pointer.From(keysResp.Model.SecondaryKey)

			return metadata.Encode(&state)
		},
	}
}

func (r RelayHybridConnectionAuthorizationRuleDataResource) ResourceType() string {
	return "azurerm_relay_hybrid_connection_authorization_rule"
}
