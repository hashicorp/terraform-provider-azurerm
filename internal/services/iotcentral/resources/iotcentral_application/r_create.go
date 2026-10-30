// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_application

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotCentralAppCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := apps.NewIotAppID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_iotcentral_application", id.ID())
		}
	}

	inputs := apps.OperationInputs{
		Name: id.IotAppName,
	}

	resp, err := client.CheckNameAvailability(ctx, commonids.NewSubscriptionID(id.SubscriptionId), inputs)
	if err != nil {
		return fmt.Errorf("checking if the name %q was globally available: %+v", id.IotAppName, err)
	}
	if model := resp.Model; model != nil {
		if !*model.NameAvailable {
			return fmt.Errorf("the name %q cannot be used. Reason: %q Message: %q", id.IotAppName, *model.Reason, *model.Message)
		}
	}

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = id.IotAppName
	}

	identity, err := identity.ExpandSystemAssigned(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	subdomain := d.Get("sub_domain").(string)
	template := d.Get("template").(string)
	app := apps.App{
		Properties: &apps.AppProperties{
			DisplayName:         &displayName,
			PublicNetworkAccess: pointer.To(apps.PublicNetworkAccessEnabled),
			Subdomain:           &subdomain,
			Template:            &template,
		},
		Sku: apps.AppSkuInfo{
			Name: apps.AppSku(d.Get("sku").(string)),
		},
		Identity: identity,
		Location: d.Get("location").(string),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, app, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	// Public Network Access can only be disabled after creation
	if !d.Get("public_network_access_enabled").(bool) {
		app.Properties.PublicNetworkAccess = pointer.To(apps.PublicNetworkAccessDisabled)
		if err := client.CreateOrUpdateThenPoll(ctx, id, app); err != nil {
			return fmt.Errorf("updating `public_network_access_enabled` to false for %s: %+v", id, err)
		}
	}

	return resourceIotCentralAppRead(d, meta)
}
