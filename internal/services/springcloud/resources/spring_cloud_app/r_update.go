// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudAppUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform_rm.ParseAppIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	identity, err := expandSpringCloudAppIdentity(d.Get("identity").([]any))
	if err != nil {
		return err
	}

	addonConfig, err := helpers.ExpandSpringCloudAppAddon(d.Get("addon_json").(string))
	if err != nil {
		return err
	}

	app := appplatform.AppResource{
		Identity: identity,
		Properties: &appplatform.AppResourceProperties{
			AddonConfigs:          addonConfig,
			EnableEndToEndTLS:     pointer.To(d.Get("tls_enabled").(bool)),
			Public:                pointer.To(d.Get("is_public").(bool)),
			HTTPSOnly:             pointer.To(d.Get("https_only").(bool)),
			IngressSettings:       expandSpringCloudAppIngressSetting(d.Get("ingress_settings").([]any)),
			PersistentDisk:        expandSpringCloudAppPersistentDisk(d.Get("persistent_disk").([]any)),
			CustomPersistentDisks: expandAppCustomPersistentDiskResourceArray(d.Get("custom_persistent_disk").([]any), *id),
		},
	}
	if enabled := d.Get("public_endpoint_enabled").(bool); enabled {
		app.Properties.VnetAddons = &appplatform.AppVNetAddons{
			PublicEndpoint: pointer.To(enabled),
		}
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.AppName, app)
	if err != nil {
		return fmt.Errorf("update %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	return resourceSpringCloudAppRead(d, meta)
}
