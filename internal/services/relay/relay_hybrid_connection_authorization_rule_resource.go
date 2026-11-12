// Copyright IBM Corp. 2014, 2026
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

type RelayHybridConnectionAuthorizationRuleResource struct{}

var (
	_ sdk.ResourceWithUpdate        = RelayHybridConnectionAuthorizationRuleResource{}
	_ sdk.ResourceWithCustomizeDiff = RelayHybridConnectionAuthorizationRuleResource{}
)

type RelayHybridConnectionAuthorizationRuleResourceModel struct {
	Name                      string `tfschema:"name"`
	ResourceGroupName         string `tfschema:"resource_group_name"`
	NamespaceName             string `tfschema:"namespace_name"`
	HybridConnectionName      string `tfschema:"hybrid_connection_name"`
	Listen                    bool   `tfschema:"listen"`
	Manage                    bool   `tfschema:"manage"`
	Send                      bool   `tfschema:"send"`
	PrimaryConnectionString   string `tfschema:"primary_connection_string"`
	PrimaryKey                string `tfschema:"primary_key"`
	SecondaryConnectionString string `tfschema:"secondary_connection_string"`
	SecondaryKey              string `tfschema:"secondary_key"`
}

func (r RelayHybridConnectionAuthorizationRuleResource) ModelObject() any {
	return &RelayHybridConnectionAuthorizationRuleResourceModel{}
}

func (r RelayHybridConnectionAuthorizationRuleResource) ResourceType() string {
	return "azurerm_relay_hybrid_connection_authorization_rule"
}

func (r RelayHybridConnectionAuthorizationRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return hybridconnections.ValidateHybridConnectionAuthorizationRuleID
}

func (r RelayHybridConnectionAuthorizationRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"namespace_name": {
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

		"resource_group_name": commonschema.ResourceGroupName(),

		"listen": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"manage": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"send": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (r RelayHybridConnectionAuthorizationRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (r RelayHybridConnectionAuthorizationRuleResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return authorizationRuleCustomizeDiff(ctx, metadata.ResourceDiff, nil)
		},
	}
}

func (r RelayHybridConnectionAuthorizationRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model RelayHybridConnectionAuthorizationRuleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := hybridconnections.NewHybridConnectionAuthorizationRuleID(subscriptionId, model.ResourceGroupName, model.NamespaceName, model.HybridConnectionName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetAuthorizationRule(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := hybridconnections.AuthorizationRule{
				Name: pointer.To(model.Name),
				Properties: &hybridconnections.AuthorizationRuleProperties{
					Rights: expandHybridConnectionAuthorizationRuleRights(metadata.ResourceData),
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

func (r RelayHybridConnectionAuthorizationRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient

			id, err := hybridconnections.ParseHybridConnectionAuthorizationRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetAuthorizationRule(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			keysResp, err := client.ListKeys(ctx, *id)
			if err != nil {
				log.Printf("[WARN] listing keys for %s: %+v\n", *id, err)
			}

			state := RelayHybridConnectionAuthorizationRuleResourceModel{
				Name:                 id.AuthorizationRuleName,
				HybridConnectionName: id.HybridConnectionName,
				ResourceGroupName:    id.ResourceGroupName,
				NamespaceName:        id.NamespaceName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Listen, state.Send, state.Manage = flattenHybridConnectionAuthorizationRuleRights(props.Rights)
				}
			}

			if keysModel := keysResp.Model; keysModel != nil {
				state.PrimaryConnectionString = pointer.From(keysModel.PrimaryConnectionString)
				state.SecondaryConnectionString = pointer.From(keysModel.SecondaryConnectionString)
				state.PrimaryKey = pointer.From(keysModel.PrimaryKey)
				state.SecondaryKey = pointer.From(keysModel.SecondaryKey)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r RelayHybridConnectionAuthorizationRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient

			id, err := hybridconnections.ParseHybridConnectionAuthorizationRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.DeleteAuthorizationRule(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r RelayHybridConnectionAuthorizationRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient

			id, err := hybridconnections.ParseHybridConnectionAuthorizationRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetAuthorizationRule(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			existingModel := *resp.Model

			if metadata.ResourceData.HasChange("listen") {
				existingModel.Properties.Rights = expandHybridConnectionAuthorizationRuleRights(metadata.ResourceData)
			}

			if metadata.ResourceData.HasChange("manage") {
				existingModel.Properties.Rights = expandHybridConnectionAuthorizationRuleRights(metadata.ResourceData)
			}

			if metadata.ResourceData.HasChange("send") {
				existingModel.Properties.Rights = expandHybridConnectionAuthorizationRuleRights(metadata.ResourceData)
			}

			if _, err := client.CreateOrUpdateAuthorizationRule(ctx, *id, existingModel); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}
