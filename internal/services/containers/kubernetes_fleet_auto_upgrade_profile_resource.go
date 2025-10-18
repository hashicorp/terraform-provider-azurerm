package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleetupdatestrategies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KubernetesFleetAutoUpgradeProfileResource struct{}

func (r KubernetesFleetAutoUpgradeProfileResource) ModelObject() interface{} {
	return &KubernetesFleetAutoUpgradeProfileResourceSchema{}
}

type KubernetesFleetAutoUpgradeProfileResourceSchema struct {
	Name                 string  `tfschema:"name"`
	ResourceGroupName    string  `tfschema:"resource_group_name"`
	FleetName            string  `tfschema:"fleet_name"`
	Channel              string  `tfschema:"channel"`
	NodeImageUpgradeType *string `tfschema:"node_image_upgrade_type"`
	UpdateStrategyId     *string `tfschema:"update_strategy_id"`
}

func (r KubernetesFleetAutoUpgradeProfileResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autoupgradeprofiles.ValidateAutoUpgradeProfileID
}

func (r KubernetesFleetAutoUpgradeProfileResource) ResourceType() string {
	return "azurerm_kubernetes_fleet_auto_upgrade_profile"
}

func (r KubernetesFleetAutoUpgradeProfileResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{"default"}, false),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"fleet_name": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"channel": {
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(autoupgradeprofiles.PossibleValuesForUpgradeChannel(), false),
		},

		"node_image_upgrade_type": {
			Optional:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(autoupgradeprofiles.PossibleValuesForAutoUpgradeNodeImageSelectionType(), false),
		},

		"update_strategy_id": {
			Optional:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: fleetupdatestrategies.ValidateUpdateStrategyID,
		},
	}
}

func (r KubernetesFleetAutoUpgradeProfileResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r KubernetesFleetAutoUpgradeProfileResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20250301.AutoUpgradeProfiles

			var config KubernetesFleetAutoUpgradeProfileResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := autoupgradeprofiles.NewAutoUpgradeProfileID(subscriptionId, config.ResourceGroupName, config.FleetName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload autoupgradeprofiles.AutoUpgradeProfile
			r.mapKubernetesFleetAutoUpgradeProfileResourceSchemaToAutoUpgradeProfile(config, &payload)

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload, autoupgradeprofiles.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KubernetesFleetAutoUpgradeProfileResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20250301.AutoUpgradeProfiles
			schema := KubernetesFleetAutoUpgradeProfileResourceSchema{}

			id, err := autoupgradeprofiles.ParseAutoUpgradeProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.AutoUpgradeProfileName
				schema.ResourceGroupName = id.ResourceGroupName
				schema.FleetName = id.FleetName
				r.mapAutoUpgradeProfileToKubernetesFleetAutoUpgradeProfileResourceSchema(*model, &schema)
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r KubernetesFleetAutoUpgradeProfileResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20250301.AutoUpgradeProfiles

			id, err := autoupgradeprofiles.ParseAutoUpgradeProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, autoupgradeprofiles.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesFleetAutoUpgradeProfileResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20250301.AutoUpgradeProfiles

			id, err := autoupgradeprofiles.ParseAutoUpgradeProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config KubernetesFleetAutoUpgradeProfileResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: properties was nil", *id)
			}
			payload := *existing.Model

			r.mapKubernetesFleetAutoUpgradeProfileResourceSchemaToAutoUpgradeProfile(config, &payload)

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload, autoupgradeprofiles.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesFleetAutoUpgradeProfileResource) mapKubernetesFleetAutoUpgradeProfileResourceSchemaToAutoUpgradeProfile(input KubernetesFleetAutoUpgradeProfileResourceSchema, output *autoupgradeprofiles.AutoUpgradeProfile) {
	if output.Properties == nil {
		output.Properties = &autoupgradeprofiles.AutoUpgradeProfileProperties{}
	}

	output.Properties.Channel = autoupgradeprofiles.UpgradeChannel(input.Channel)

	if input.NodeImageUpgradeType != nil {
		output.Properties.NodeImageSelection = &autoupgradeprofiles.AutoUpgradeNodeImageSelection{
			Type: autoupgradeprofiles.AutoUpgradeNodeImageSelectionType(pointer.From(input.NodeImageUpgradeType)),
		}
	}

	output.Properties.UpdateStrategyId = input.UpdateStrategyId
}

func (r KubernetesFleetAutoUpgradeProfileResource) mapAutoUpgradeProfileToKubernetesFleetAutoUpgradeProfileResourceSchema(input autoupgradeprofiles.AutoUpgradeProfile, output *KubernetesFleetAutoUpgradeProfileResourceSchema) {
	if input.Properties == nil {
		input.Properties = &autoupgradeprofiles.AutoUpgradeProfileProperties{}
	}

	output.Channel = string(input.Properties.Channel)

	if input.Properties.NodeImageSelection != nil {
		output.NodeImageUpgradeType = pointer.To(string(input.Properties.NodeImageSelection.Type))
	}

	output.UpdateStrategyId = input.Properties.UpdateStrategyId
}
