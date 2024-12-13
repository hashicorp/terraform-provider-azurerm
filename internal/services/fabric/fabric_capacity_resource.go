package fabric

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fabric/2023-11-01/fabriccapacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type FabricCapacityResource struct{}

var _ sdk.ResourceWithUpdate = FabricCapacityResource{}

type FabricCapacityResourceModel struct {
	Name                  string            `tfschema:"name"`
	ResourceGroupName     string            `tfschema:"resource_group_name"`
	AdministrationMembers []string          `tfschema:"administration_members"`
	Location              string            `tfschema:"location"`
	Sku                   []SkuModel        `tfschema:"sku"`
	Tags                  map[string]string `tfschema:"tags"`
}

type SkuModel struct {
	Name string `tfschema:"name"`
	Tier string `tfschema:"tier"`
}

func (r FabricCapacityResource) ModelObject() interface{} {
	return &FabricCapacityResource{}
}

func (r FabricCapacityResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fabriccapacities.ValidateCapacityID
}

func (r FabricCapacityResource) ResourceType() string {
	return "azurerm_fabric_capacity"
}

func (r FabricCapacityResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z]([a-z\d]{2,62})$`),
				"`name` must be between 3 and 63 characters. It can contain only lowercase letters and numbers. It must start with a lowercase letter.",
			),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"F2",
							"F4",
							"F8",
							"F16",
							"F32",
							"F64",
							"F128",
							"F256",
							"F512",
							"F1024",
							"F2048",
						}, false),
					},
					"tier": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(fabriccapacities.RpSkuTierFabric),
						}, false),
					},
				},
			},
		},

		"administration_members": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r FabricCapacityResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r FabricCapacityResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Fabric.FabricCapacitiesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model FabricCapacityResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := fabriccapacities.NewCapacityID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			if len(model.AdministrationMembers) == 0 {
				return fmt.Errorf("`administration_members` is required when creating a fabric capacity")
			}

			properties := fabriccapacities.FabricCapacity{
				Location: location.Normalize(model.Location),
				Properties: fabriccapacities.FabricCapacityProperties{
					Administration: fabriccapacities.CapacityAdministration{
						Members: model.AdministrationMembers,
					},
				},
				Sku: expandSkuModel(model.Sku),
			}

			if model.Tags != nil {
				properties.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r FabricCapacityResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Fabric.FabricCapacitiesClient

			id, err := fabriccapacities.ParseCapacityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model FabricCapacityResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			payload := existing.Model
			if metadata.ResourceData.HasChange("administration_members") {
				payload.Properties.Administration = fabriccapacities.CapacityAdministration{
					Members: model.AdministrationMembers,
				}
			}

			if metadata.ResourceData.HasChange("sku") {
				payload.Sku = expandSkuModel(model.Sku)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r FabricCapacityResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Fabric.FabricCapacitiesClient

			id, err := fabriccapacities.ParseCapacityID(metadata.ResourceData.Id())
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

			state := FabricCapacityResourceModel{
				Name:              id.CapacityName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.AdministrationMembers = model.Properties.Administration.Members
				state.Sku = flattenSkuModel(model.Sku)
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r FabricCapacityResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Fabric.FabricCapacitiesClient

			id, err := fabriccapacities.ParseCapacityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandSkuModel(inputList []SkuModel) fabriccapacities.RpSku {
	input := &inputList[0]
	return fabriccapacities.RpSku{
		Name: input.Name,
		Tier: fabriccapacities.RpSkuTier(input.Tier),
	}
}

func flattenSkuModel(input fabriccapacities.RpSku) []SkuModel {
	outputList := make([]SkuModel, 0)
	output := SkuModel{
		Name: input.Name,
		Tier: string(input.Tier),
	}

	return append(outputList, output)
}
