// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSpringCloudAppRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform.ParseAppIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Spring Cloud App %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", id.AppName, id.SpringName, id.ResourceGroupName, err)
	}

	d.Set("name", id.AppName)
	d.Set("service_name", id.SpringName)
	d.Set("resource_group_name", id.ResourceGroupName)

	identity, err := flattenSpringCloudAppIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	if prop := resp.Properties; prop != nil {
		d.Set("is_public", prop.Public)
		d.Set("https_only", prop.HTTPSOnly)
		d.Set("fqdn", prop.Fqdn)
		d.Set("url", prop.URL)
		d.Set("tls_enabled", prop.EnableEndToEndTLS)
		if err := d.Set("addon_json", helpers.FlattenSpringCloudAppAddon(prop.AddonConfigs)); err != nil {
			return fmt.Errorf("setting `addon_json`: %s", err)
		}
		if err := d.Set("persistent_disk", flattenSpringCloudAppPersistentDisk(prop.PersistentDisk)); err != nil {
			return fmt.Errorf("setting `persistent_disk`: %s", err)
		}
		if err := d.Set("custom_persistent_disk", flattenAppCustomPersistentDiskResourceArray(prop.CustomPersistentDisks)); err != nil {
			return fmt.Errorf("setting `custom_persistent_disk`: %+v", err)
		}
		if err := d.Set("ingress_settings", flattenSpringCloudAppIngressSettings(prop.IngressSettings)); err != nil {
			return fmt.Errorf("setting `ingress_settings`: %+v", err)
		}
		if prop.VnetAddons != nil {
			d.Set("public_endpoint_enabled", prop.VnetAddons.PublicEndpoint)
		}
	}

	return nil
}
