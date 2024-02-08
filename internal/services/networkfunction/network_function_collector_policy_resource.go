// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkfunction

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/collectorpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetworkFunctionCollectorPolicyModel struct {
	Name                                   string                 `tfschema:"name"`
	NetworkFunctionAzureTrafficCollectorId string                 `tfschema:"traffic_collector_id"`
	IpfxEmission                           []IpfxEmissionModel    `tfschema:"ipfx_emission"`
	IpfxIngestion                          []IpfxIngestionModel   `tfschema:"ipfx_ingestion"`
	Location                               string                 `tfschema:"location"`
	Tags                                   map[string]interface{} `tfschema:"tags"`
}

type IpfxEmissionModel struct {
	DestinationTypes []string `tfschema:"destination_types"`
}

type IpfxIngestionModel struct {
	SourceResourceIds []string `tfschema:"source_resource_ids"`
}

type NetworkFunctionCollectorPolicyResource struct{}

var _ sdk.ResourceWithUpdate = NetworkFunctionCollectorPolicyResource{}

func (r NetworkFunctionCollectorPolicyResource) ResourceType() string {
	return "azurerm_network_function_collector_policy"
}

func (r NetworkFunctionCollectorPolicyResource) ModelObject() interface{} {
	return &NetworkFunctionCollectorPolicyModel{}
}

func (r NetworkFunctionCollectorPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return collectorpolicies.ValidateCollectorPolicyID
}

func (r NetworkFunctionCollectorPolicyResource) Arguments() map[string]*pluginsdk.Schema {
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

		"location": commonschema.Location(),

		"traffic_collector_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azuretrafficcollectors.ValidateAzureTrafficCollectorID,
		},

		"ipfx_emission": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"destination_types": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(collectorpolicies.PossibleValuesForDestinationType(), false),
						},
					},
				},
			},
		},

		"ipfx_ingestion": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source_resource_ids": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r NetworkFunctionCollectorPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetworkFunctionCollectorPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetworkFunctionCollectorPolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetworkFunction.CollectorPoliciesClient
			azureTrafficCollectorId, err := azuretrafficcollectors.ParseAzureTrafficCollectorID(model.NetworkFunctionAzureTrafficCollectorId)
			if err != nil {
				return err
			}

			id := collectorpolicies.NewCollectorPolicyID(azureTrafficCollectorId.SubscriptionId, azureTrafficCollectorId.ResourceGroupName, azureTrafficCollectorId.AzureTrafficCollectorName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &collectorpolicies.CollectorPolicy{
				Location: location.Normalize(model.Location),
				Properties: &collectorpolicies.CollectorPolicyPropertiesFormat{
					EmissionPolicies: expandEmissionPoliciesPropertiesFormatModelArray(model.IpfxEmission),
					IngestionPolicy: &collectorpolicies.IngestionPolicyPropertiesFormat{
						IngestionSources: expandIngestionSourcesPropertiesFormatModelArray(model.IpfxIngestion),
						IngestionType:    pointer.To(collectorpolicies.IngestionTypeIPFIX),
					},
				},

				Tags: tags.Expand(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkFunctionCollectorPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkFunction.CollectorPoliciesClient

			id, err := collectorpolicies.ParseCollectorPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkFunctionCollectorPolicyModel
			if err = metadata.Decode(&model); err != nil {
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
				properties.Tags = tags.Expand(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NetworkFunctionCollectorPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkFunction.CollectorPoliciesClient

			id, err := collectorpolicies.ParseCollectorPolicyID(metadata.ResourceData.Id())
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

			state := NetworkFunctionCollectorPolicyModel{
				Name:                                   id.CollectorPolicyName,
				NetworkFunctionAzureTrafficCollectorId: azuretrafficcollectors.NewAzureTrafficCollectorID(id.SubscriptionId, id.ResourceGroupName, id.AzureTrafficCollectorName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if properties := model.Properties; properties != nil {
					state.IpfxEmission = flattenEmissionPoliciesPropertiesFormatModelArray(properties.EmissionPolicies)
					if properties.IngestionPolicy != nil {
						state.IpfxIngestion = flattenIngestionSourcesPropertiesFormatModelArray(properties.IngestionPolicy.IngestionSources)
					}

				}

				state.Tags = tags.Flatten(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NetworkFunctionCollectorPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkFunction.CollectorPoliciesClient

			id, err := collectorpolicies.ParseCollectorPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			// API has bug, which appears to be eventually consistent. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/25152
			log.Printf("[DEBUG] Waiting for %s to be fully deleted..", *id)
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			stateConf := &pluginsdk.StateChangeConf{
				Pending:                   []string{"Exists"},
				Target:                    []string{"NotFound"},
				Refresh:                   collectorPolicyDeletedRefreshFunc(ctx, client, *id),
				MinTimeout:                10 * time.Second,
				ContinuousTargetOccurence: 20,
				Timeout:                   time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be fully deleted: %+v", *id, err)
			}

			return nil
		},
	}
}

func collectorPolicyDeletedRefreshFunc(ctx context.Context, client *collectorpolicies.CollectorPoliciesClient, id collectorpolicies.CollectorPolicyId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("checking if %s has been deleted: %+v", id, err)
		}

		return res, "Exists", nil
	}
}

func expandEmissionPoliciesPropertiesFormatModelArray(inputList []IpfxEmissionModel) *[]collectorpolicies.EmissionPoliciesPropertiesFormat {
	var outputList []collectorpolicies.EmissionPoliciesPropertiesFormat
	for _, v := range inputList {
		input := v
		output := collectorpolicies.EmissionPoliciesPropertiesFormat{
			EmissionType:         pointer.To(collectorpolicies.EmissionTypeIPFIX),
			EmissionDestinations: expandEmissionPolicyDestinationModelArray(input.DestinationTypes),
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func expandEmissionPolicyDestinationModelArray(inputList []string) *[]collectorpolicies.EmissionPolicyDestination {
	var outputList []collectorpolicies.EmissionPolicyDestination
	for _, v := range inputList {
		output := collectorpolicies.EmissionPolicyDestination{
			DestinationType: pointer.To(collectorpolicies.DestinationType(v)),
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func expandIngestionSourcesPropertiesFormatModelArray(inputList []IpfxIngestionModel) *[]collectorpolicies.IngestionSourcesPropertiesFormat {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []collectorpolicies.IngestionSourcesPropertiesFormat
	for _, v := range inputList[0].SourceResourceIds {
		output := collectorpolicies.IngestionSourcesPropertiesFormat{
			SourceType: pointer.To(collectorpolicies.SourceTypeResource),
			ResourceId: &v,
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenEmissionPoliciesPropertiesFormatModelArray(inputList *[]collectorpolicies.EmissionPoliciesPropertiesFormat) []IpfxEmissionModel {
	outputList := make([]IpfxEmissionModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := IpfxEmissionModel{
			DestinationTypes: flattenEmissionPolicyDestinationModelArray(input.EmissionDestinations),
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenEmissionPolicyDestinationModelArray(inputList *[]collectorpolicies.EmissionPolicyDestination) []string {
	outputList := make([]string, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		if input.DestinationType != nil {
			outputList = append(outputList, string(*input.DestinationType))
		}
	}

	return outputList
}

func flattenIngestionSourcesPropertiesFormatModelArray(inputList *[]collectorpolicies.IngestionSourcesPropertiesFormat) []IpfxIngestionModel {
	outputList := make([]IpfxIngestionModel, 0)
	if inputList == nil {
		return outputList
	}

	output := IpfxIngestionModel{
		SourceResourceIds: make([]string, 0),
	}

	for _, input := range *inputList {
		if input.ResourceId != nil && input.SourceType != nil && *input.SourceType == collectorpolicies.SourceTypeResource {
			output.SourceResourceIds = append(output.SourceResourceIds, *input.ResourceId)
		}
	}

	if len(output.SourceResourceIds) > 0 {
		outputList = append(outputList, output)
	}

	return outputList
}
