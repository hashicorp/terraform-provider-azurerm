// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image_gallery

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/gallerysharingupdate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageGalleryDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.GalleriesClient
	gallerySharingUpdateClient := meta.(*clients.Client).Compute.GallerySharingUpdateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSharedImageGalleryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, galleries.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		if prop := model.Properties; prop != nil && prop.SharingProfile != nil && prop.SharingProfile.Permissions != nil {
			if pointer.From(prop.SharingProfile.Permissions) == galleries.GallerySharingPermissionTypesCommunity {
				updatePayload := gallerysharingupdate.SharingUpdate{
					OperationType: gallerysharingupdate.SharingUpdateOperationTypesReset,
				}
				if err = gallerySharingUpdateClient.GallerySharingProfileUpdateThenPoll(ctx, *id, updatePayload); err != nil {
					return fmt.Errorf("resetting community sharing of %s: %+v", id, err)
				}
			}
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
