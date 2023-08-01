package automanage

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
	"time"
)

type ConfigurationHCIAssignmentResource struct {
	Name                   string `tfschema:"name"`
	ResourceGroupName      string `tfschema:"resource_group_name"`
	ClusterName            string `tfschema:"cluster_name"`
	ConfigurationProfileID string `tfschema:"configuration_id"`
}

func (c ConfigurationHCIAssignmentResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"cluster_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"configuration_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomanageConfigurationID,
		},
	}
}

func (c ConfigurationHCIAssignmentResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (c ConfigurationHCIAssignmentResource) ModelObject() interface{} {
	return &ConfigurationHCIAssignmentResource{}
}

func (c ConfigurationHCIAssignmentResource) ResourceType() string {
	return "azurerm_automanage_configuration_hci_assignment"
}

func (c ConfigurationHCIAssignmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ConfigurationHCIAssignmentResource
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.Automanage.HCIAssignmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := parse.NewAutomanageConfigurationHCIAssignmentID(subscriptionId, model.ResourceGroupName, model.ClusterName, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.ConfigurationProfileAssignmentName)

			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(c.ResourceType(), id)
			}

			properties := automanage.ConfigurationProfileAssignment{
				Properties: &automanage.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: &model.ConfigurationProfileID,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, properties, id.ResourceGroup, id.ClusterName, id.ConfigurationProfileAssignmentName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (c ConfigurationHCIAssignmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.HCIAssignmentClient

			id, err := parse.AutomanageConfigurationHCIAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.ConfigurationProfileAssignmentName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ConfigurationHCIAssignmentResource{
				Name:              id.ConfigurationProfileAssignmentName,
				ResourceGroupName: id.ResourceGroup,
				ClusterName:       id.ClusterName,
			}

			if existing.Properties != nil && existing.Properties.ConfigurationProfile != nil {
				profileId, err := parse.AutomanageConfigurationHCIAssignmentID(*existing.Properties.ConfigurationProfile)
				if err != nil {
					return err
				}
				state.ConfigurationProfileID = profileId.String()
			}

			return metadata.Encode(state)
		},
	}
}

func (c ConfigurationHCIAssignmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.HCIAssignmentClient

			id, err := parse.AutomanageConfigurationHCIAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, id.ResourceGroup, id.ClusterName, id.ConfigurationProfileAssignmentName); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (c ConfigurationHCIAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AutomanageConfigurationHCIAssignmentID
}

var _ sdk.Resource = ConfigurationHCIAssignmentResource{}
