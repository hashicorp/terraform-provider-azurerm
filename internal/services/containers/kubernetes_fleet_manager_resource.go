package containers

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type KubernetesFleetManagerModel struct {
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Location          string                 `tfschema:"location"`
	HubProfile        []FleetHubProfileModel `tfschema:"hub_profile"`
	Tags              map[string]string      `tfschema:"tags"`
}

type FleetHubProfileModel struct {
	DnsPrefix         string `tfschema:"dns_prefix"`
	Fqdn              string `tfschema:"fqdn"`
	KubernetesVersion string `tfschema:"kubernetes_version"`
	PortalFqdn        string `tfschema:"portal_fqdn"`
}

type KubernetesFleetManagerResource struct{}

var _ sdk.ResourceWithUpdate = KubernetesFleetManagerResource{}

func (r KubernetesFleetManagerResource) ResourceType() string {
	return "azurerm_kubernetes_fleet_manager"
}

func (r KubernetesFleetManagerResource) ModelObject() interface{} {
	return &KubernetesFleetManagerModel{}
}

func (r KubernetesFleetManagerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fleets.ValidateFleetID
}

func (r KubernetesFleetManagerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"hub_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_prefix": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.All(
							validation.StringLenBetween(1, 54),
							validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]$|^[a-zA-Z0-9][a-zA-Z0-9-]{0,52}[a-zA-Z0-9]$`), "must match the pattern ^[a-zA-Z0-9]$|^[a-zA-Z0-9][a-zA-Z0-9-]{0,52}[a-zA-Z0-9]$"),
						),
					},

					"fqdn": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"kubernetes_version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"portal_fqdn": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r KubernetesFleetManagerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r KubernetesFleetManagerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20240401.Fleets
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model KubernetesFleetManagerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := fleets.NewFleetID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			params := fleets.Fleet{
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Properties: &fleets.FleetProperties{
					HubProfile: expandFleetHubProfileModel(model.HubProfile),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, params, fleets.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KubernetesFleetManagerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20240401.Fleets

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KubernetesFleetManagerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			payload := resp.Model
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}
			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload, fleets.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesFleetManagerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20240401.Fleets

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := KubernetesFleetManagerModel{
				Name:              id.FleetName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if properties := model.Properties; properties != nil {
					state.HubProfile = flattenFleetHubProfileModel(properties.HubProfile)
				}

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r KubernetesFleetManagerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20240401.Fleets

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, fleets.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandFleetHubProfileModel(inputList []FleetHubProfileModel) *fleets.FleetHubProfile {
	if len(inputList) == 0 {
		return nil
	}

	input := inputList[0]
	output := &fleets.FleetHubProfile{
		DnsPrefix: pointer.To(input.DnsPrefix),
	}

	return output
}

func flattenFleetHubProfileModel(input *fleets.FleetHubProfile) []FleetHubProfileModel {
	if input == nil {
		return []FleetHubProfileModel{}
	}

	output := FleetHubProfileModel{
		DnsPrefix:         pointer.From(input.DnsPrefix),
		Fqdn:              pointer.From(input.Fqdn),
		KubernetesVersion: pointer.From(input.KubernetesVersion),
		PortalFqdn:        pointer.From(input.PortalFqdn),
	}

	return []FleetHubProfileModel{output}
}
