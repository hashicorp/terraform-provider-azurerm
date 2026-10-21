// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_configuration

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMaintenanceConfigurationUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := existing.Model

	if d.HasChanges("scope", "window", "install_patches", "properties") {
		scope := maintenanceconfigurations.MaintenanceScope(d.Get("scope").(string))
		window := expandMaintenanceConfigurationWindow(d.Get("window").([]any))
		installPatches := expandMaintenanceConfigurationInstallPatches(d.Get("install_patches").([]any))
		extensionProperties := expandExtensionProperties(d.Get("properties").(map[string]any))
		if scope == maintenanceconfigurations.MaintenanceScopeInGuestPatch {
			if window == nil {
				return errors.New("`window` must be specified when `scope` is `InGuestPatch`")
			}
			if installPatches == nil {
				return errors.New("`install_patches` must be specified when `scope` is `InGuestPatch`")
			}
			if _, ok := (*extensionProperties)["InGuestPatchMode"]; !ok {
				if _, ok := d.GetOk("in_guest_user_patch_mode"); !ok {
					return errors.New("`in_guest_user_patch_mode` must be specified when `scope` is `InGuestPatch`")
				}
				(*extensionProperties)["InGuestPatchMode"] = d.Get("in_guest_user_patch_mode").(string)
			}
		}

		payload.Properties.MaintenanceScope = &scope
		payload.Properties.MaintenanceWindow = window
		payload.Properties.ExtensionProperties = extensionProperties
		payload.Properties.InstallPatches = installPatches
	}

	if d.HasChange("visibility") {
		payload.Properties.Visibility = pointer.ToEnum[maintenanceconfigurations.Visibility](d.Get("visibility").(string))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceMaintenanceConfigurationRead(d, meta)
}
