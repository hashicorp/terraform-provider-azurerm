// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudAppCreate(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).AppsClient
	servicesClient := springcloudClient.From(meta).ServicesClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := appplatform_rm.NewAppID(subscriptionId, d.Get("resource_group_name").(string), d.Get("service_name").(string), d.Get("name").(string))
	serviceResp, err := servicesClient.Get(ctx, id.ResourceGroupName, id.SpringName)
	if err != nil {
		return fmt.Errorf("unable to retrieve %q: %+v", id, err)
	}

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, "")
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_app", id.ID())
		}
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
		Location: serviceResp.Location,
		Identity: identity,
		Properties: &appplatform.AppResourceProperties{
			AddonConfigs:          addonConfig,
			EnableEndToEndTLS:     pointer.To(d.Get("tls_enabled").(bool)),
			Public:                pointer.To(d.Get("is_public").(bool)),
			CustomPersistentDisks: expandAppCustomPersistentDiskResourceArray(d.Get("custom_persistent_disk").([]any), id),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.AppName, app)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %q: %+v", id, err)
	}

	// HTTPSOnly and PersistentDisk could only be set by update
	app.Properties.HTTPSOnly = pointer.To(d.Get("https_only").(bool))
	app.Properties.PersistentDisk = expandSpringCloudAppPersistentDisk(d.Get("persistent_disk").([]any))
	// VNetAddons.PublicEndpoint could only be set by update
	if enabled := d.Get("public_endpoint_enabled").(bool); enabled {
		app.Properties.VnetAddons = &appplatform.AppVNetAddons{
			PublicEndpoint: pointer.To(enabled),
		}
	}
	// IngressSettings could only be set by update
	// Issue: https://github.com/Azure/azure-rest-api-specs/issues/21536
	app.Properties.IngressSettings = expandSpringCloudAppIngressSetting(d.Get("ingress_settings").([]any))
	future, err = client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.AppName, app)
	if err != nil {
		return fmt.Errorf("update %q: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %q: %+v", id, err)
	}

	return resourceSpringCloudAppRead(d, meta)
}
