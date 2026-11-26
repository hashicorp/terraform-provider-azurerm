// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_function_collector_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/collectorpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
