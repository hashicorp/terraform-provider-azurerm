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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerSecurityAdminConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetworkManager.SecurityAdminConfigurations
			networkManagerId, err := securityadminconfigurations.ParseNetworkManagerID(model.NetworkManagerId)
			if err != nil {
				return err
			}

			id := securityadminconfigurations.NewSecurityAdminConfigurationID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			conf := securityadminconfigurations.SecurityAdminConfiguration{
				Properties: &securityadminconfigurations.SecurityAdminConfigurationPropertiesFormat{
					ApplyOnNetworkIntentPolicyBasedServices: expandNetworkIntentPolicyBasedServiceModel(model.ApplyOnNetworkIntentPolicyBasedServices),
				},
			}

			if model.Description != "" {
				conf.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, conf); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandNetworkIntentPolicyBasedServiceModel(inputList []string) *[]securityadminconfigurations.NetworkIntentPolicyBasedService {
	outputList := make([]securityadminconfigurations.NetworkIntentPolicyBasedService, 0, len(inputList))
	for _, input := range inputList {
		output := securityadminconfigurations.NetworkIntentPolicyBasedService(input)

		outputList = append(outputList, output)
	}

	return &outputList
}
