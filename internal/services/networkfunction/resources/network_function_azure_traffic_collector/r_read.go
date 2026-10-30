// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_function_azure_traffic_collector

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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

func flattenCollectorPolicyModelArray(inputList *[]azuretrafficcollectors.ResourceReference) []string {
	var outputList []string
	if inputList == nil {
		return []string{}
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
