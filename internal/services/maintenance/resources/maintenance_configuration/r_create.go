// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_configuration

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMaintenanceConfigurationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := maintenanceconfigurations.NewMaintenanceConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_maintenance_configuration", id.ID())
		}
	}

	scope := maintenanceconfigurations.MaintenanceScope(d.Get("scope").(string))
	visibility := maintenanceconfigurations.Visibility(d.Get("visibility").(string))
	windowRaw := d.Get("window").([]any)
	window := expandMaintenanceConfigurationWindow(windowRaw)
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

	configuration := maintenanceconfigurations.MaintenanceConfiguration{
		Name:     pointer.To(id.MaintenanceConfigurationName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &maintenanceconfigurations.MaintenanceConfigurationProperties{
			MaintenanceScope:    &scope,
			Visibility:          &visibility,
			Namespace:           pointer.To("Microsoft.Maintenance"),
			MaintenanceWindow:   window,
			ExtensionProperties: extensionProperties,
			InstallPatches:      installPatches,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if _, err := client.CreateOrUpdate(ctx, id, configuration); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMaintenanceConfigurationRead(d, meta)
}
