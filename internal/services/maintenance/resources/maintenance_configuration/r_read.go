// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_configuration

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMaintenanceConfigurationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] maintenance %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.MaintenanceConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("scope", pointer.FromEnum(props.MaintenanceScope))
			d.Set("visibility", pointer.FromEnum(props.Visibility))

			properties := flattenExtensionProperties(props.ExtensionProperties)
			if properties["InGuestPatchMode"] != nil {
				d.Set("in_guest_user_patch_mode", properties["InGuestPatchMode"])
				delete(properties, "InGuestPatchMode")
			}
			d.Set("properties", properties)

			if err := d.Set("window", flattenMaintenanceConfigurationWindow(props.MaintenanceWindow)); err != nil {
				return fmt.Errorf("setting `window`: %+v", err)
			}

			if err := d.Set("install_patches", flattenMaintenanceConfigurationInstallPatches(props.InstallPatches)); err != nil {
				return fmt.Errorf("setting `install_patches`: %+v", err)
			}
		}
		d.Set("location", location.NormalizeNilable(model.Location))
		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}
