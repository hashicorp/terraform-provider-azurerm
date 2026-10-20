// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_group_oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeGroupClient
			volumeClient := metadata.Client.NetApp.VolumeClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model models.NetAppVolumeGroupOracleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := volumegroups.NewVolumeGroupID(subscriptionId, model.ResourceGroupName, model.AccountName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			volumeList, err := helpers.ExpandNetAppVolumeGroupOracleVolumes(model.Volumes)
			if err != nil {
				return err
			}

			// Performing some basic validations that are not possible in the schema
			if errorList := validate.ValidateNetAppVolumeGroupOracleVolumes(volumeList); len(errorList) > 0 {
				return fmt.Errorf("one or more issues found while performing deeper validations for %s:\n%+v", id, errorList)
			}

			parameters := volumegroups.VolumeGroupDetails{
				Location: pointer.To(location.Normalize(model.Location)),
				Properties: &volumegroups.VolumeGroupProperties{
					GroupMetaData: &volumegroups.VolumeGroupMetaData{
						GroupDescription:      pointer.To(model.GroupDescription),
						ApplicationType:       pointer.To(volumegroups.ApplicationTypeORACLE),
						ApplicationIdentifier: pointer.To(model.ApplicationIdentifier),
					},
					Volumes: volumeList,
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			// Waiting for volume group be completely provisioned
			if err := helpers.WaitForVolumeGroupCreateOrUpdate(ctx, client, id); err != nil {
				return fmt.Errorf("waiting creation %s: %+v", id, err)
			}

			// CRR - Authorizing secondaries from primary volumes
			if err := helpers.AuthorizeVolumeReplication(ctx, volumeList, volumeClient, subscriptionId, model.ResourceGroupName, model.AccountName); err != nil {
				return err
			}

			return nil
		},
	}
}
