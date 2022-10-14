package networkfunction

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkFunctionAzureTrafficCollectorModel struct {
	Name              string                   `tfschema:"name"`
	ResourceGroupName string                   `tfschema:"resource_group_name"`
	Location          string                   `tfschema:"location"`
	Tags              map[string]string        `tfschema:"tags"`
	VirtualHub        []ResourceReferenceModel `tfschema:"virtual_hub"`
	CollectorPolicies []ResourceReferenceModel `tfschema:"collector_policies"`
	SystemData        []SystemDataModel        `tfschema:"system_data"`
}

type ResourceReferenceModel struct {
	Id string `tfschema:"id"`
}

type SystemDataModel struct {
	CreatedAt          string                               `tfschema:"created_at"`
	CreatedBy          string                               `tfschema:"created_by"`
	CreatedByType      azuretrafficcollectors.CreatedByType `tfschema:"created_by_type"`
	LastModifiedBy     string                               `tfschema:"last_modified_by"`
	LastModifiedByType azuretrafficcollectors.CreatedByType `tfschema:"last_modified_by_type"`
}

type NetworkFunctionAzureTrafficCollectorResource struct{}

var _ sdk.ResourceWithUpdate = NetworkFunctionAzureTrafficCollectorResource{}

func (r NetworkFunctionAzureTrafficCollectorResource) ResourceType() string {
	return "azurerm_network_function_azure_traffic_collector"
}

func (r NetworkFunctionAzureTrafficCollectorResource) ModelObject() interface{} {
	return &NetworkFunctionAzureTrafficCollectorModel{}
}

func (r NetworkFunctionAzureTrafficCollectorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return azuretrafficcollectors.ValidateAzureTrafficCollectorID
}

func (r NetworkFunctionAzureTrafficCollectorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.LocationWithoutForceNew(),

		"tags": commonschema.Tags(),

		"virtual_hub": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func (r NetworkFunctionAzureTrafficCollectorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"collector_policies": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"system_data": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"created_at": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"created_by": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"created_by_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"last_modified_by": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"last_modified_by_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r NetworkFunctionAzureTrafficCollectorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetworkFunctionAzureTrafficCollectorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetworkFunction.AzureTrafficCollectorsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := azuretrafficcollectors.NewAzureTrafficCollectorID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &azuretrafficcollectors.AzureTrafficCollector{
				Location:   location.Normalize(model.Location),
				Properties: &azuretrafficcollectors.AzureTrafficCollectorPropertiesFormat{},
				Tags:       &model.Tags,
			}

			virtualHubValue, err := expandResourceReferenceModel(model.VirtualHub)
			if err != nil {
				return err
			}

			properties.Properties.VirtualHub = virtualHubValue

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkFunctionAzureTrafficCollectorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkFunction.AzureTrafficCollectorsClient

			id, err := azuretrafficcollectors.ParseAzureTrafficCollectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkFunctionAzureTrafficCollectorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("location") {
				properties.Location = location.Normalize(model.Location)
			}

			if metadata.ResourceData.HasChange("virtual_hub") {
				virtualHubValue, err := expandResourceReferenceModel(model.VirtualHub)
				if err != nil {
					return err
				}

				properties.Properties.VirtualHub = virtualHubValue
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NetworkFunctionAzureTrafficCollectorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkFunction.AzureTrafficCollectorsClient

			id, err := azuretrafficcollectors.ParseAzureTrafficCollectorID(metadata.ResourceData.Id())
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

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := NetworkFunctionAzureTrafficCollectorModel{
				Name:              id.AzureTrafficCollectorName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			if properties := model.Properties; properties != nil {
				collectorPoliciesValue, err := flattenResourceReferenceModel(properties.CollectorPolicies)
				if err != nil {
					return err
				}

				state.CollectorPolicies = collectorPoliciesValue

				virtualHubValue, err := flattenResourceReferenceModel(properties.VirtualHub)
				if err != nil {
					return err
				}

				state.VirtualHub = virtualHubValue
			}
			systemDataValue, err := flattenSystemDataModel(model.SystemData)
			if err != nil {
				return err
			}

			state.SystemData = systemDataValue
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NetworkFunctionAzureTrafficCollectorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkFunction.AzureTrafficCollectorsClient

			id, err := azuretrafficcollectors.ParseAzureTrafficCollectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandResourceReferenceModel(inputList []ResourceReferenceModel) (*azuretrafficcollectors.ResourceReference, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := azuretrafficcollectors.ResourceReference{}

	return &output, nil
}

func flattenResourceReferenceModel(inputList *[]azuretrafficcollectors.ResourceReference) ([]ResourceReferenceModel, error) {
	var outputList []ResourceReferenceModel
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := ResourceReferenceModel{}

		if input.Id != nil {
			output.Id = *input.Id
		}

		outputList = append(outputList, output)
	}

	return outputList, nil
}

func flattenSystemDataModel(input *azuretrafficcollectors.SystemData) ([]SystemDataModel, error) {
	var outputList []SystemDataModel
	if input == nil {
		return outputList, nil
	}

	output := SystemDataModel{}

	if input.CreatedAt != nil {
		output.CreatedAt = *input.CreatedAt
	}

	if input.CreatedBy != nil {
		output.CreatedBy = *input.CreatedBy
	}

	if input.CreatedByType != nil {
		output.CreatedByType = *input.CreatedByType
	}

	if input.LastModifiedBy != nil {
		output.LastModifiedBy = *input.LastModifiedBy
	}

	if input.LastModifiedByType != nil {
		output.LastModifiedByType = *input.LastModifiedByType
	}

	return append(outputList, output), nil
}

func flattenResourceReferenceModel(input *azuretrafficcollectors.ResourceReference) ([]ResourceReferenceModel, error) {
	var outputList []ResourceReferenceModel
	if input == nil {
		return outputList, nil
	}

	output := ResourceReferenceModel{}

	if input.Id != nil {
		output.Id = *input.Id
	}

	return append(outputList, output), nil
}
