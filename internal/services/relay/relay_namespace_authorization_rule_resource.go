// Copyright IBM Corp. 2014, 2026
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

type RelayNamespaceAuthorizationRuleResource struct{}

var (
	_ sdk.ResourceWithUpdate        = RelayNamespaceAuthorizationRuleResource{}
	_ sdk.ResourceWithCustomizeDiff = RelayNamespaceAuthorizationRuleResource{}
)

type RelayNamespaceAuthorizationRuleResourceModel struct {
	Name                      string `tfschema:"name"`
	NamespaceName             string `tfschema:"namespace_name"`
	ResourceGroupName         string `tfschema:"resource_group_name"`
	Listen                    bool   `tfschema:"listen"`
	Manage                    bool   `tfschema:"manage"`
	Send                      bool   `tfschema:"send"`
	PrimaryConnectionString   string `tfschema:"primary_connection_string"`
	PrimaryKey                string `tfschema:"primary_key"`
	SecondaryConnectionString string `tfschema:"secondary_connection_string"`
	SecondaryKey              string `tfschema:"secondary_key"`
}

func (r RelayNamespaceAuthorizationRuleResource) ModelObject() any {
	return &RelayNamespaceAuthorizationRuleResourceModel{}
}

func (r RelayNamespaceAuthorizationRuleResource) ResourceType() string {
	return "azurerm_relay_namespace_authorization_rule"
}

func (r RelayNamespaceAuthorizationRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return namespaces.ValidateAuthorizationRuleID
}

func (r RelayNamespaceAuthorizationRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(6, 50),
		},

		"namespace_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
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

func (r RelayNamespaceAuthorizationRuleResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r RelayNamespaceAuthorizationRuleResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return authorizationRuleCustomizeDiff(ctx, metadata.ResourceDiff, nil)
		},
	}
}

func (r RelayNamespaceAuthorizationRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model RelayNamespaceAuthorizationRuleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := namespaces.NewAuthorizationRuleID(subscriptionId, model.ResourceGroupName, model.NamespaceName, model.Name)

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

			parameters := namespaces.AuthorizationRule{
				Name: pointer.To(model.Name),
				Properties: &namespaces.AuthorizationRuleProperties{
					Rights: expandAuthorizationRuleRights(metadata.ResourceData),
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

func (r RelayNamespaceAuthorizationRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
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

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			keysResp, err := client.ListKeys(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", *id, err)
			}

			state := RelayNamespaceAuthorizationRuleResourceModel{
				Name:              id.AuthorizationRuleName,
				NamespaceName:     id.NamespaceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Listen, state.Send, state.Manage = flattenAuthorizationRuleRights(props.Rights)
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

func (r RelayNamespaceAuthorizationRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
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

func (r RelayNamespaceAuthorizationRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			id, err := namespaces.ParseAuthorizationRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model RelayNamespaceAuthorizationRuleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := namespaces.AuthorizationRule{
				Name: pointer.To(model.Name),
				Properties: &namespaces.AuthorizationRuleProperties{
					Rights: expandAuthorizationRuleRights(metadata.ResourceData),
				},
			}

			if _, err := client.CreateOrUpdateAuthorizationRule(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
