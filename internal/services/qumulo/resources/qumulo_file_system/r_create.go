// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package qumulo_file_system

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/qumulostorage/2026-04-16/filesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Qumulo.FileSystemsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config FileSystemResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := filesystems.NewFileSystemID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			// check the subnet is valid: https://learn.microsoft.com/azure/partner-solutions/qumulo/qumulo-troubleshoot#you-cant-create-a-resource
			if err := checkSubnet(ctx, config.SubnetId, metadata); err != nil {
				return err
			}

			payload := filesystems.LiftrBaseStorageFileSystemResource{
				Location: location.Normalize(config.Location),
				Tags:     pointer.To(config.Tags),
				Properties: &filesystems.LiftrBaseStorageFileSystemResourceProperties{
					AdminPassword:     config.AdminPassword,
					AvailabilityZone:  pointer.To(config.Zone),
					DelegatedSubnetId: config.SubnetId,
					StorageSku:        config.StorageSku,
					UserDetails: filesystems.LiftrBaseUserDetails{
						Email: config.Email,
					},
					MarketplaceDetails: filesystems.LiftrBaseMarketplaceDetails{
						OfferId:     config.OfferId,
						PlanId:      config.PlanId,
						PublisherId: pointer.To(config.PublisherId),
					},
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
