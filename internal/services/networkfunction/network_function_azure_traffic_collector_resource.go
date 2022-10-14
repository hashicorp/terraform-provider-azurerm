package networkfunction

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetworkFunctionAzureTrafficCollectorModel struct {
	Name              string                   `tfschema:"name"`
	ResourceGroupName string                   `tfschema:"resource_group_name"`
	Location          string                   `tfschema:"location"`
	Tags              map[string]string        `tfschema:"tags"`
	VirtualHub        []ResourceReferenceModel `tfschema:"virtual_hub"`
	CollectorPolicies []ResourceReferenceModel `tfschema:"collector_policies"`
}

type ResourceReferenceModel struct {
	Id string `tfschema:"id"`
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9]([-._a-zA-Z0-9]{0,78}[a-zA-Z0-9_])?$"),
				"The name can contain only letters, numbers, periods (.), hyphens (-),and underscores (_), up to 80 characters, and it must begin with a letter or number and end with a letter, number or underscore.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r NetworkFunctionAzureTrafficCollectorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"collector_policies": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"virtual_hub": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
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

			properties.SystemData = nil

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
				collectorPoliciesValue, err := flattenResourceReferenceModelArray(properties.CollectorPolicies)
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

func flattenResourceReferenceModelArray(inputList *[]azuretrafficcollectors.ResourceReference) ([]ResourceReferenceModel, error) {
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
