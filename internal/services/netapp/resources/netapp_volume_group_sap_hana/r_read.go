// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_group_sap_hana

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
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeGroupClient

			id, err := volumegroups.ParseVolumeGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state models.NetAppVolumeGroupSAPHanaModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := models.NetAppVolumeGroupSAPHanaModel{
				Name:              id.VolumeGroupName,
				AccountName:       id.NetAppAccountName,
				Location:          location.NormalizeNilable(existing.Model.Location),
				ResourceGroupName: id.ResourceGroupName,
			}

			if props := existing.Model.Properties; props != nil {
				model.GroupDescription = pointer.From(props.GroupMetaData.GroupDescription)
				model.ApplicationIdentifier = pointer.From(props.GroupMetaData.ApplicationIdentifier)

				volumes, err := helpers.FlattenNetAppVolumeGroupSAPHanaVolumes(ctx, props.Volumes, metadata)
				if err != nil {
					return fmt.Errorf("setting `volume`: %+v", err)
				}

				model.Volumes = volumes
			}

			metadata.SetID(id)

			return metadata.Encode(&model)
		},
	}
}
