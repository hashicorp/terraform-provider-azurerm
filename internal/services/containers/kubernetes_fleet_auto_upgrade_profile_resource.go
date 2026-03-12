package containers

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetupdatestrategies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = KubernetesFleetAutoUpgradeProfileResource{}
	_ sdk.ResourceWithUpdate = KubernetesFleetAutoUpgradeProfileResource{}
)

type KubernetesFleetAutoUpgradeProfileResource struct{}

func (r KubernetesFleetAutoUpgradeProfileResource) ModelObject() interface{} {
	return &KubernetesFleetAutoUpgradeProfileResourceModel{}
}

type KubernetesFleetAutoUpgradeProfileResourceModel struct {
	Name                     string  `tfschema:"name"`
	KubernetesFleetManagerId string  `tfschema:"kubernetes_fleet_manager_id"`
	Channel                  string  `tfschema:"channel"`
	NodeImageSelectionType   *string `tfschema:"node_image_selection_type"`
	UpdateStrategyId         *string `tfschema:"update_strategy_id"`
	Disabled                 *bool   `tfschema:"disabled"`
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
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 50),
				validation.StringMatch(
					regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`),
					"must start and end with a lowercase letter or number, and can only contain lowercase letters, numbers, and hyphens",
				),
			),
		},

		"kubernetes_fleet_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KubernetesFleetId{}),

		"channel": {
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(autoupgradeprofiles.PossibleValuesForUpgradeChannel(), false),
		},

		"node_image_selection_type": {
			Optional:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(autoupgradeprofiles.PossibleValuesForAutoUpgradeNodeImageSelectionType(), false),
		},

		"update_strategy_id": {
			Optional:     true,
			ForceNew:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: fleetupdatestrategies.ValidateUpdateStrategyID,
		},

		"disabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
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
			client := metadata.Client.Containers.FleetAutoUpgradeProfilesClient

			var config KubernetesFleetAutoUpgradeProfileResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			fleetId, err := commonids.ParseKubernetesFleetID(config.KubernetesFleetManagerId)
			if err != nil {
				return err
			}

			id := autoupgradeprofiles.NewAutoUpgradeProfileID(fleetId.SubscriptionId, fleetId.ResourceGroupName, fleetId.FleetName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := autoupgradeprofiles.AutoUpgradeProfile{
				Properties: &autoupgradeprofiles.AutoUpgradeProfileProperties{
					Channel: autoupgradeprofiles.UpgradeChannel(config.Channel),
				},
			}

			if config.NodeImageSelectionType != nil {
				payload.Properties.NodeImageSelection = &autoupgradeprofiles.AutoUpgradeNodeImageSelection{
					Type: autoupgradeprofiles.AutoUpgradeNodeImageSelectionType(pointer.From(config.NodeImageSelectionType)),
				}
			}

			payload.Properties.UpdateStrategyId = config.UpdateStrategyId
			payload.Properties.Disabled = config.Disabled

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
			client := metadata.Client.Containers.FleetAutoUpgradeProfilesClient

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

			state := KubernetesFleetAutoUpgradeProfileResourceModel{
				Name:                     id.AutoUpgradeProfileName,
				KubernetesFleetManagerId: commonids.NewKubernetesFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Channel = string(props.Channel)
					state.UpdateStrategyId = props.UpdateStrategyId
					state.Disabled = props.Disabled

					if props.NodeImageSelection != nil {
						nodeImageSelectionType := string(props.NodeImageSelection.Type)
						state.NodeImageSelectionType = &nodeImageSelectionType
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r KubernetesFleetAutoUpgradeProfileResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetAutoUpgradeProfilesClient

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
			client := metadata.Client.Containers.FleetAutoUpgradeProfilesClient

			id, err := autoupgradeprofiles.ParseAutoUpgradeProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config KubernetesFleetAutoUpgradeProfileResourceModel
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

			if payload.Properties == nil {
				payload.Properties = &autoupgradeprofiles.AutoUpgradeProfileProperties{}
			}

			if metadata.ResourceData.HasChange("channel") {
				payload.Properties.Channel = autoupgradeprofiles.UpgradeChannel(config.Channel)
			}

			if metadata.ResourceData.HasChange("node_image_selection_type") {
				if config.NodeImageSelectionType != nil {
					payload.Properties.NodeImageSelection = &autoupgradeprofiles.AutoUpgradeNodeImageSelection{
						Type: autoupgradeprofiles.AutoUpgradeNodeImageSelectionType(pointer.From(config.NodeImageSelectionType)),
					}
				} else {
					payload.Properties.NodeImageSelection = nil
				}
			}

			if metadata.ResourceData.HasChange("disabled") {
				payload.Properties.Disabled = config.Disabled
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload, autoupgradeprofiles.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
