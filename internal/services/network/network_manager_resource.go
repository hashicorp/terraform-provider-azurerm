package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	managementGroupValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type ManagerModel struct {
	CrossTenantScopes []ManagerCrossTenantScopeModel `tfschema:"cross_tenant_scopes"`
	Scope             []ManagerScopeModel            `tfschema:"scope"`
	ScopeAccesses     []string                       `tfschema:"scope_accesses"`
	Description       string                         `tfschema:"description"`
	Name              string                         `tfschema:"name"`
	Location          string                         `tfschema:"location"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Tags              map[string]interface{}         `tfschema:"tags"`
}

type ManagerScopeModel struct {
	ManagementGroups []string `tfschema:"management_group_ids"`
	Subscriptions    []string `tfschema:"subscription_ids"`
}

type ManagerCrossTenantScopeModel struct {
	TenantId         string   `tfschema:"tenant_id"`
	ManagementGroups []string `tfschema:"management_groups"`
	Subscriptions    []string `tfschema:"subscriptions"`
}

type ManagerResource struct{}

func (r ManagerResource) ResourceType() string {
	return "azurerm_network_manager"
}

func (r ManagerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NetworkManagerID
}

func (r ManagerResource) ModelObject() interface{} {
	return &ManagerModel{}
}

func (r ManagerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"scope": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"management_group_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: managementGroupValidate.ManagementGroupID,
						},
						AtLeastOneOf: []string{"scope.0.management_group_ids", "scope.0.subscription_ids"},
					},
					"subscription_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: commonids.ValidateSubscriptionID,
						},
						AtLeastOneOf: []string{"scope.0.management_group_ids", "scope.0.subscription_ids"},
					},
				},
			},
		},

		"scope_accesses": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.ConfigurationTypeConnectivity),
					string(network.ConfigurationTypeSecurityAdmin),
				}, false),
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r ManagerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cross_tenant_scopes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"subscriptions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"management_groups": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
	}
}

func (r ManagerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state ManagerModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Network.ManagersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewNetworkManagerID(subscriptionId, state.ResourceGroupName, state.Name)
			metadata.Logger.Infof("creating %s", id)

			existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := network.Manager{
				Location: utils.String(azure.NormalizeLocation(state.Location)),
				Name:     utils.String(state.Name),
				ManagerProperties: &network.ManagerProperties{
					Description:                 utils.String(state.Description),
					NetworkManagerScopes:        expandNetworkManagerScope(state.Scope),
					NetworkManagerScopeAccesses: expandNetworkManagerScopeAccesses(state.ScopeAccesses),
				},
				Tags: tags.Expand(state.Tags),
			}

			if _, err := client.CreateOrUpdate(ctx, input, id.ResourceGroup, id.Name); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ManagerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagersClient
			id, err := parse.NetworkManagerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving %s", *id)
			resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var description string
			var scope []ManagerScopeModel
			var ScopeAccesses []string
			if prop := resp.ManagerProperties; prop != nil {
				if prop.Description != nil {
					description = *resp.Description
				}
				scope = flattenNetworkManagerScope(resp.NetworkManagerScopes)
				ScopeAccesses = flattenNetworkManagerScopeAccesses(resp.NetworkManagerScopeAccesses)
			}

			return metadata.Encode(&ManagerModel{
				Description:       description,
				Location:          location.NormalizeNilable(resp.Location),
				Name:              id.Name,
				ResourceGroupName: id.ResourceGroup,
				ScopeAccesses:     ScopeAccesses,
				Scope:             scope,
				Tags:              tags.Flatten(resp.Tags),
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ManagerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.NetworkManagerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Network.ManagersClient
			existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.ManagerProperties == nil {
				return fmt.Errorf("unexpected null properties of %s", *id)
			}
			var state ManagerModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("description") {
				existing.ManagerProperties.Description = utils.String(state.Description)
			}

			if metadata.ResourceData.HasChange("scope") {
				existing.ManagerProperties.NetworkManagerScopes = expandNetworkManagerScope(state.Scope)
			}

			if metadata.ResourceData.HasChange("scope_accesses") {
				existing.ManagerProperties.NetworkManagerScopeAccesses = expandNetworkManagerScopeAccesses(state.ScopeAccesses)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.Expand(state.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, existing, id.ResourceGroup, id.Name); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ManagerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagersClient
			id, err := parse.NetworkManagerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)
			future, err := client.Delete(ctx, id.ResourceGroup, id.Name, utils.Bool(true))
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func stringSlice(input []string) *[]string {
	return &input
}

func expandNetworkManagerScope(input []ManagerScopeModel) *network.ManagerPropertiesNetworkManagerScopes {
	if len(input) == 0 {
		return nil
	}

	return &network.ManagerPropertiesNetworkManagerScopes{
		ManagementGroups: stringSlice(input[0].ManagementGroups),
		Subscriptions:    stringSlice(input[0].Subscriptions),
	}
}

func expandNetworkManagerScopeAccesses(input []string) *[]network.ConfigurationType {
	result := make([]network.ConfigurationType, 0)
	for _, v := range input {
		result = append(result, network.ConfigurationType(v))
	}
	return &result
}

func flattenStringSlicePtr(input *[]string) []string {
	if input == nil {
		return make([]string, 0)
	}
	return *input
}

func flattenNetworkManagerScope(input *network.ManagerPropertiesNetworkManagerScopes) []ManagerScopeModel {
	if input == nil {
		return make([]ManagerScopeModel, 0)
	}

	return []ManagerScopeModel{{
		ManagementGroups: flattenStringSlicePtr(input.ManagementGroups),
		Subscriptions:    flattenStringSlicePtr(input.Subscriptions),
	}}
}

func flattenNetworkManagerScopeAccesses(input *[]network.ConfigurationType) []string {
	var result []string
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}
