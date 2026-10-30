// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_application

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotCentralAppRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apps.ParseIotAppID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.IotAppName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		d.Set("sku", model.Sku.Name)

		if props := model.Properties; props != nil {
			d.Set("sub_domain", props.Subdomain)
			d.Set("display_name", props.DisplayName)
			d.Set("template", props.Template)

			publicNetworkAccess := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == apps.PublicNetworkAccessDisabled {
				publicNetworkAccess = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccess)
		}

		if err := d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
