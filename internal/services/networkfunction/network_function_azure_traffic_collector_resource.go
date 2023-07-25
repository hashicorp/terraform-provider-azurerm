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
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
	CollectorPolicies []string          `tfschema:"collector_policy_ids"`
	VirtualHub        []string          `tfschema:"virtual_hub_id"`
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
		"collector_policy_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"virtual_hub_id": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
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

			state := NetworkFunctionAzureTrafficCollectorModel{
				Name:              id.AzureTrafficCollectorName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if properties := model.Properties; properties != nil {
					state.CollectorPolicies = flattenCollectorPolicyModelArray(properties.CollectorPolicies)
					state.VirtualHub = flattenVirtualHubModel(properties.VirtualHub)
				}

				if model.Tags != nil {
					state.Tags = *model.Tags
				}
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

func flattenCollectorPolicyModelArray(inputList *[]azuretrafficcollectors.ResourceReference) []string {
	var outputList []string
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		if input.Id != nil {
			outputList = append(outputList, *input.Id)
		}
	}

	return outputList
}

func flattenVirtualHubModel(input *azuretrafficcollectors.ResourceReference) []string {
	var outputList []string
	if input != nil && input.Id != nil {
		outputList = append(outputList, *input.Id)
	}

	return outputList
}
