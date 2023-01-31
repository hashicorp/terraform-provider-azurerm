package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type ManagerAdminRuleCollectionModel struct {
	Name                         string   `tfschema:"name"`
	SecurityAdminConfigurationId string   `tfschema:"security_admin_configuration_id"`
	NetworkGroupIds              []string `tfschema:"network_group_ids"`
	Description                  string   `tfschema:"description"`
}

type ManagerAdminRuleCollectionResource struct{}

var _ sdk.ResourceWithUpdate = ManagerAdminRuleCollectionResource{}

func (r ManagerAdminRuleCollectionResource) ResourceType() string {
	return "azurerm_network_manager_admin_rule_collection"
}

func (r ManagerAdminRuleCollectionResource) ModelObject() interface{} {
	return &ManagerAdminRuleCollectionModel{}
}

func (r ManagerAdminRuleCollectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NetworkManagerAdminRuleCollectionID
}

func (r ManagerAdminRuleCollectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"security_admin_configuration_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NetworkManagerSecurityAdminConfigurationID,
		},

		"network_group_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.NetworkManagerNetworkGroupID,
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r ManagerAdminRuleCollectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerAdminRuleCollectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerAdminRuleCollectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.ManagerAdminRuleCollectionsClient
			configurationId, err := parse.NetworkManagerSecurityAdminConfigurationID(model.SecurityAdminConfigurationId)
			if err != nil {
				return err
			}

			id := parse.NewNetworkManagerAdminRuleCollectionID(configurationId.SubscriptionId, configurationId.ResourceGroup,
				configurationId.NetworkManagerName, configurationId.SecurityAdminConfigurationName, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName)

			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			adminRuleCollection := &network.AdminRuleCollection{
				AdminRuleCollectionPropertiesFormat: &network.AdminRuleCollectionPropertiesFormat{
					AppliesToGroups: expandNetworkManagerNetworkGroupIds(model.NetworkGroupIds),
				},
			}

			if model.Description != "" {
				adminRuleCollection.AdminRuleCollectionPropertiesFormat.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, *adminRuleCollection, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerAdminRuleCollectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerAdminRuleCollectionsClient

			id, err := parse.NetworkManagerAdminRuleCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerAdminRuleCollectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.AdminRuleCollectionPropertiesFormat
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("network_group_ids") {
				properties.AppliesToGroups = expandNetworkManagerNetworkGroupIds(model.NetworkGroupIds)
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Description = utils.String(model.Description)
			}

			if _, err := client.CreateOrUpdate(ctx, existing, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagerAdminRuleCollectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerAdminRuleCollectionsClient

			id, err := parse.NetworkManagerAdminRuleCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.AdminRuleCollectionPropertiesFormat
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := ManagerAdminRuleCollectionModel{
				Name:                         id.RuleCollectionName,
				SecurityAdminConfigurationId: parse.NewNetworkManagerSecurityAdminConfigurationID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName).ID(),
				NetworkGroupIds:              flattenNetworkManagerNetworkGroupIds(properties.AppliesToGroups),
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerAdminRuleCollectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerAdminRuleCollectionsClient

			id, err := parse.NetworkManagerAdminRuleCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName, utils.Bool(true))
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func expandNetworkManagerNetworkGroupIds(inputList []string) *[]network.ManagerSecurityGroupItem {
	var outputList []network.ManagerSecurityGroupItem
	for _, v := range inputList {
		input := v
		output := network.ManagerSecurityGroupItem{
			NetworkGroupID: utils.String(input),
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenNetworkManagerNetworkGroupIds(inputList *[]network.ManagerSecurityGroupItem) []string {
	var outputList []string
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		if input.NetworkGroupID != nil {
			outputList = append(outputList, *input.NetworkGroupID)
		}
	}

	return outputList
}
