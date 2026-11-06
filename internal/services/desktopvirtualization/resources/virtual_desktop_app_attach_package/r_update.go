// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_app_attach_package

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/appattachpackage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.AppAttachPackagesClient

			id, err := appattachpackage.ParseAppAttachPackageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model VirtualDesktopAppAttachPackageModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			param := *existing.Model

			// Update `param.Properties.Image` if any of the child properties has changed
			if metadata.ResourceData.HasChanges("display_name", "msix_package_name", "storage_share_file_id", "register_at_log_on_enabled", "state_enabled") {
				msixImageProperties, err := helpers.GetMsixImageProperties(ctx, metadata, model.HostPoolIds, model.StorageShareFileId, model.MsixPackageName)
				if err != nil {
					return fmt.Errorf("retrieving MSIX image properties: %+v", err)
				}

				param.Properties.Image = r.expandVirtualDesktopAppAttachPackageImage(model, msixImageProperties)
			}

			if metadata.ResourceData.HasChange("host_pool_ids") {
				param.Properties.HostPoolReferences = pointer.To(model.HostPoolIds)
			}

			if metadata.ResourceData.HasChange("health_check_status_on_failure") {
				param.Properties.FailHealthCheckOnStagingFailure = pointer.ToEnum[appattachpackage.FailHealthCheckOnStagingFailure](model.HealthCheckStatusOnFailure)
			}

			if metadata.ResourceData.HasChange("tags") {
				param.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
