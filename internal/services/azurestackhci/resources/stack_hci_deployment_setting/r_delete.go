// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_deployment_setting

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/storagecontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 1 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.DeploymentSettings

			id, err := deploymentsettings.ParseDeploymentSettingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			applianceName := fmt.Sprintf("%s-arcbridge", id.ClusterName)
			applianceId := appliances.NewApplianceID(id.SubscriptionId, id.ResourceGroupName, applianceName)

			applianceClient := metadata.Client.ArcResourceBridge.AppliancesClient
			if err := applianceClient.DeleteThenPoll(ctx, applianceId); err != nil {
				return fmt.Errorf("deleting Arc Resource Bridge Appliance generated during deployment: deleting %s: %+v", applianceId, err)
			}

			var customLocationName string
			if resp.Model != nil && resp.Model.Properties != nil &&
				len(resp.Model.Properties.DeploymentConfiguration.ScaleUnits) > 0 &&
				resp.Model.Properties.DeploymentConfiguration.ScaleUnits[0].DeploymentData.OptionalServices != nil {
				customLocationName = pointer.From(resp.Model.Properties.DeploymentConfiguration.ScaleUnits[0].DeploymentData.OptionalServices.CustomLocation)
			}

			if customLocationName != "" {
				customLocationId := customlocations.NewCustomLocationID(id.SubscriptionId, id.ResourceGroupName, customLocationName)

				// we need to delete the Storage Paths before the Custom Location, otherwise the Custom Location cannot be deleted if there is any Resource in it

				storageContainerClient := metadata.Client.AzureStackHCI.StorageContainers
				resourceGroupId := commonids.NewResourceGroupID(id.SubscriptionId, id.ResourceGroupName)
				storageContainers, err := storageContainerClient.ListComplete(ctx, resourceGroupId)
				if err != nil {
					return fmt.Errorf("deleting Stack HCI Storage Paths generated during deployment: retrieving Storage Path under %s: %+v", resourceGroupId, err)
				}

				// match Storage Paths under the Custom Location, the generated Storage Path name should match below pattern
				storageContainerNamePattern := regexp.MustCompile(`UserStorage[0-9]+-[a-z0-9]{32}`)
				for _, v := range storageContainers.Items {
					if v.Id != nil && v.ExtendedLocation != nil && v.ExtendedLocation.Name != nil && strings.EqualFold(*v.ExtendedLocation.Name, customLocationId.ID()) && v.Name != nil && storageContainerNamePattern.MatchString(*v.Name) {
						storageContainerId, err := storagecontainers.ParseStorageContainerIDInsensitively(*v.Id)
						if err != nil {
							return fmt.Errorf("parsing the Stack HCI Storage Path ID generated during deployment: %+v", err)
						}

						if err := storageContainerClient.DeleteThenPoll(ctx, *storageContainerId); err != nil {
							return fmt.Errorf("deleting the Stack HCI Storage Path generated during deployment: deleting %s: %+v", storageContainerId, err)
						}
					}
				}

				customLocationsClient := metadata.Client.ExtendedLocation.CustomLocationsClient
				if err := customLocationsClient.DeleteThenPoll(ctx, customLocationId); err != nil {
					return fmt.Errorf("deleting the Custom Location generated during deployment: deleting %s: %+v", customLocationId, err)
				}
			}

			return nil
		},
	}
}
