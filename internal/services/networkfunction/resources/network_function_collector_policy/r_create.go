// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_function_collector_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/collectorpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func expandEmissionPoliciesPropertiesFormatModelArray(inputList []IpfxEmissionModel) *[]collectorpolicies.EmissionPoliciesPropertiesFormat {
	outputList := make([]collectorpolicies.EmissionPoliciesPropertiesFormat, 0, len(inputList))
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
	outputList := make([]collectorpolicies.EmissionPolicyDestination, 0, len(inputList))
	for _, v := range inputList {
		output := collectorpolicies.EmissionPolicyDestination{
			DestinationType: pointer.ToEnum[collectorpolicies.DestinationType](v),
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func expandIngestionSourcesPropertiesFormatModelArray(inputList []IpfxIngestionModel) *[]collectorpolicies.IngestionSourcesPropertiesFormat {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]collectorpolicies.IngestionSourcesPropertiesFormat, 0, len(inputList))
	for _, v := range inputList[0].SourceResourceIds {
		output := collectorpolicies.IngestionSourcesPropertiesFormat{
			SourceType: pointer.To(collectorpolicies.SourceTypeResource),
			ResourceId: &v,
		}

		outputList = append(outputList, output)
	}

	return &outputList
}
