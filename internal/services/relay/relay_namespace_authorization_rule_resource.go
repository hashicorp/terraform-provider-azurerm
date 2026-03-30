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
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = RelayNamespaceAuthorizationRuleResource{}

type RelayNamespaceAuthorizationRuleResource struct{}

type RelayNamespaceAuthorizationRuleResourceModel struct {
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

func (RelayNamespaceAuthorizationRuleResource) Arguments() map[string]*pluginsdk.Schema {
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

func (RelayNamespaceAuthorizationRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return authorizationRuleAttributesFrom(map[string]*pluginsdk.Schema{})
}

func (RelayNamespaceAuthorizationRuleResource) ModelObject() interface{} {
	return &RelayNamespaceAuthorizationRuleResourceModel{}
}

func (RelayNamespaceAuthorizationRuleResource) ResourceType() string {
	return "azurerm_relay_namespace_authorization_rule"
}

func (r RelayNamespaceAuthorizationRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(30 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			log.Printf("[INFO] preparing arguments for Relay Namespace Authorization Rule creation.")

			var config RelayNamespaceAuthorizationRuleResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := namespaces.NewAuthorizationRuleID(subscriptionId, config.ResourceGroupName, config.RelayNamespaceName, config.Name)

			existing, err := client.GetAuthorizationRule(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := namespaces.AuthorizationRule{
				Name: pointer.To(config.Name),
				Properties: &namespaces.AuthorizationRuleProperties{
					Rights: expandAuthorizationRuleRights(config),
				},
			}

			if _, err := client.CreateOrUpdateAuthorizationRule(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r RelayNamespaceAuthorizationRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(30 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			log.Printf("[INFO] preparing arguments for Relay Namespace Authorization Rule creation.")

			var config RelayNamespaceAuthorizationRuleResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := namespaces.ParseAuthorizationRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.GetAuthorizationRule(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving: %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			parameters := namespaces.AuthorizationRule{
				Name: pointer.To(config.Name),
				Properties: &namespaces.AuthorizationRuleProperties{
					Rights: expandAuthorizationRuleRights(config),
				},
			}

			if _, err := client.CreateOrUpdateAuthorizationRule(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r RelayNamespaceAuthorizationRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(5 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			id, err := namespaces.ParseAuthorizationRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetAuthorizationRule(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving: %s: %+v", id, err)
			}

			keysResp, err := client.ListKeys(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", id, err)
			}

			state := RelayNamespaceAuthorizationRuleResourceModel{}

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

			return metadata.Encode(&state)
		},
	}
}

func (r RelayNamespaceAuthorizationRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(30 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			id, err := namespaces.ParseAuthorizationRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.DeleteAuthorizationRule(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r RelayNamespaceAuthorizationRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		if _, err := namespaces.ParseAuthorizationRuleID(v); err != nil {
			errors = append(errors, err)
		}

		return
	}
}

func (r RelayNamespaceAuthorizationRuleResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(30 * time.Minute),

		Func: authorizationRuleCustomizeDiff,
	}
}
