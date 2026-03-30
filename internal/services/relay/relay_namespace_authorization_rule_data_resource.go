// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package relay

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = RelayNamespaceAuthorizationRuleDataResource{}

type RelayNamespaceAuthorizationRuleDataResource struct{}

type RelayNamespaceAuthorizationRuleDataResourceModel struct {
	Name                      string `tfschema:"name"`
	ResourceGroupName         string `tfschema:"resource_group_name"`
	RelayNamespaceName        string `tfschema:"namespace_name"`
	PrimaryConnectionString   string `tfschema:"primary_connection_string"`
	SecondaryConnectionString string `tfschema:"secondary_connection_string"`
	PrimaryKey                string `tfschema:"primary_key"`
	SecondaryKey              string `tfschema:"secondary_key"`
	Listen                    bool   `tfschema:"listen"`
	Send                      bool   `tfschema:"send"`
	Manage                    bool   `tfschema:"manage"`
}

func (RelayNamespaceAuthorizationRuleDataResource) Arguments() map[string]*pluginsdk.Schema {
	return authorizationRuleArgumentsFrom(map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"namespace_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	})
}

func (RelayNamespaceAuthorizationRuleDataResource) Attributes() map[string]*pluginsdk.Schema {
	return authorizationRuleAttributesFrom(map[string]*pluginsdk.Schema{})
}

func (RelayNamespaceAuthorizationRuleDataResource) ModelObject() interface{} {
	return &RelayNamespaceAuthorizationRuleDataResourceModel{}
}

func (RelayNamespaceAuthorizationRuleDataResource) ResourceType() string {
	return "azurerm_relay_namespace_authorization_rule"
}

func (r RelayNamespaceAuthorizationRuleDataResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(5 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config RelayNamespaceAuthorizationRuleDataResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := namespaces.NewAuthorizationRuleID(subscriptionId, config.ResourceGroupName, config.RelayNamespaceName, config.Name)

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

			state := RelayNamespaceAuthorizationRuleDataResourceModel{}

			state.Name = id.AuthorizationRuleName
			state.RelayNamespaceName = id.NamespaceName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				listen, send, manage := flattenAuthorizationRuleRights(model.Properties.Rights)
				state.Manage = manage
				state.Listen = listen
				state.Send = send
			}

			state.PrimaryConnectionString = pointer.From(keysResp.Model.PrimaryConnectionString)
			state.PrimaryKey = pointer.From(keysResp.Model.PrimaryKey)
			state.SecondaryConnectionString = pointer.From(keysResp.Model.SecondaryConnectionString)
			state.SecondaryKey = pointer.From(keysResp.Model.SecondaryKey)

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
