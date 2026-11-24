// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_application

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotCentralAppUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apps.ParseIotAppID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil || existing.Model == nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model.Properties == nil {
		existing.Model.Properties = &apps.AppProperties{}
	}

	if d.HasChange("sub_domain") {
		existing.Model.Properties.Subdomain = pointer.To(d.Get("sub_domain").(string))
	}

	if d.HasChange("display_name") {
		existing.Model.Properties.DisplayName = pointer.To(d.Get("display_name").(string))
	}

	if d.HasChange("sku") {
		existing.Model.Sku = apps.AppSkuInfo{
			Name: apps.AppSku(d.Get("sku").(string)),
		}
	}

	if d.HasChange("template") {
		existing.Model.Properties.Template = pointer.To(d.Get("template").(string))
	}

	if d.HasChange("tags") {
		existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("identity") {
		identity, err := identity.ExpandSystemAssigned(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		existing.Model.Identity = identity
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := apps.PublicNetworkAccessDisabled
		if d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = apps.PublicNetworkAccessEnabled
		}
		existing.Model.Properties.PublicNetworkAccess = &publicNetworkAccess
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceIotCentralAppRead(d, meta)
}
