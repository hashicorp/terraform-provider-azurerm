package apimanagement

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementStandaloneGatewayModel struct {
	Name               string            `tfschema:"name"`
	ResourceGroupName  string            `tfschema:"resource_group_name"`
	BackendSubnetId    string            `tfschema:"backend_subnet_id"`
	Location           string            `tfschema:"location"`
	Sku                []GatewaySkuModel `tfschema:"sku"`
	Tags               map[string]string `tfschema:"tags"`
	VirtualNetworkType string            `tfschema:"virtual_network_type"`
}

type GatewaySkuModel struct {
	Capacity int64                        `tfschema:"capacity"`
	Name     apigateway.ApiGatewaySkuType `tfschema:"name"`
}

type ApiManagementStandaloneGatewayResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementStandaloneGatewayResource{}

func (r ApiManagementStandaloneGatewayResource) ResourceType() string {
	return "azurerm_api_management_standalone_gateway"
}

func (r ApiManagementStandaloneGatewayResource) ModelObject() interface{} {
	return &ApiManagementStandaloneGatewayModel{}
}

func (r ApiManagementStandaloneGatewayResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return apigateway.ValidateGatewayID
}

func (r ApiManagementStandaloneGatewayResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z](?:[a-zA-Z0-9-]{0,43}[a-zA-Z0-9])?$"),
				"The `name` must be between 1 and 45 characters long and can only include letters, numbers, and hyphens. The first character must be a letter and last character must be a letter or a number.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

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
							// `Standard`: Deprecated - will be removed in a future version.
							// `WorkspaceGatewayStandard`: Private Preview – not currently for general use.
							string(apigateway.ApiGatewaySkuTypeWorkspaceGatewayPremium),
						}, false),
					},
					"capacity": {
						Type:         pluginsdk.TypeInt,
						Default:      1,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"backend_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
			RequiredWith: []string{"virtual_network_type"},
		},

		"tags": commonschema.TagsForceNew(),

		"virtual_network_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				// Note: Whilst the `None` value exists it's handled in the Create/Update and Read functions.
				string(apigateway.VirtualNetworkTypeExternal),
				string(apigateway.VirtualNetworkTypeInternal),
			}, false),
			RequiredWith: []string{"backend_subnet_id"},
		},
	}
}

func (r ApiManagementStandaloneGatewayResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementStandaloneGatewayResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiGatewayClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ApiManagementStandaloneGatewayModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := apigateway.NewGatewayID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			virtualNetworkType := string(apigateway.VirtualNetworkTypeNone)
			if v := model.VirtualNetworkType; v != "" {
				virtualNetworkType = v
			}

			properties := &apigateway.ApiManagementGatewayResource{
				Location: location.Normalize(model.Location),
				Properties: apigateway.ApiManagementGatewayBaseProperties{
					VirtualNetworkType: pointer.To(apigateway.VirtualNetworkType(virtualNetworkType)),
				},
				Sku:  pointer.From(expandGatewaySkuModel(model.Sku)),
				Tags: pointer.To(model.Tags),
			}

			if model.BackendSubnetId != "" {
				properties.Properties.Backend = &apigateway.BackendConfiguration{
					Subnet: &apigateway.BackendSubnetConfiguration{
						Id: pointer.To(model.BackendSubnetId),
					},
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApiManagementStandaloneGatewayResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiGatewayClient

			id, err := apigateway.ParseGatewayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApiManagementStandaloneGatewayModel
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
			if metadata.ResourceData.HasChange("sku") {
				payload.Sku = pointer.From(expandGatewaySkuModel(model.Sku))
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementStandaloneGatewayResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiGatewayClient

			id, err := apigateway.ParseGatewayID(metadata.ResourceData.Id())
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

			state := ApiManagementStandaloneGatewayModel{
				Name:              id.GatewayName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if backend := model.Properties.Backend; backend != nil {
					if subnet := backend.Subnet; subnet != nil {
						state.BackendSubnetId = pointer.From(subnet.Id)
					}
				}

				virtualNetworkType := ""
				if v := model.Properties.VirtualNetworkType; v != nil && *v != apigateway.VirtualNetworkTypeNone {
					virtualNetworkType = string(*v)
				}
				state.VirtualNetworkType = virtualNetworkType

				state.Sku = flattenGatewaySkuModel(&model.Sku)
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementStandaloneGatewayResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiGatewayClient

			id, err := apigateway.ParseGatewayID(metadata.ResourceData.Id())
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

func expandGatewaySkuModel(inputList []GatewaySkuModel) *apigateway.ApiManagementGatewaySkuProperties {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	return &apigateway.ApiManagementGatewaySkuProperties{
		Capacity: pointer.To(input.Capacity),
		Name:     input.Name,
	}
}

func flattenGatewaySkuModel(input *apigateway.ApiManagementGatewaySkuProperties) []GatewaySkuModel {
	outputList := make([]GatewaySkuModel, 0)
	if input == nil {
		return outputList
	}
	output := GatewaySkuModel{
		Name:     input.Name,
		Capacity: pointer.From(input.Capacity),
	}

	return append(outputList, output)
}
