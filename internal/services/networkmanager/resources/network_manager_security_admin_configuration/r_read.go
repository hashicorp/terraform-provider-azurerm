// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_security_admin_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/securityadminconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.SecurityAdminConfigurations

			id, err := securityadminconfigurations.ParseSecurityAdminConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties

			state := ManagerSecurityAdminConfigurationModel{
				Name:                                    id.SecurityAdminConfigurationName,
				NetworkManagerId:                        securityadminconfigurations.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID(),
				ApplyOnNetworkIntentPolicyBasedServices: flattenNetworkIntentPolicyBasedServiceModel(properties.ApplyOnNetworkIntentPolicyBasedServices),
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenNetworkIntentPolicyBasedServiceModel(inputList *[]securityadminconfigurations.NetworkIntentPolicyBasedService) []string {
	if inputList == nil {
		return []string{}
	}

	outputList := make([]string, 0, len(*inputList))
	for _, input := range *inputList {
		outputList = append(outputList, string(input))
	}

	return outputList
}
