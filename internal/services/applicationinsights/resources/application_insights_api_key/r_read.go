// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_api_key

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentapikeysapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsAPIKeyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.APIKeysClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := componentapikeysapis.ParseApiKeyID(d.Id())
	if err != nil {
		return err
	}

	appInsightsId := componentsapis.NewComponentID(subscriptionId, id.ResourceGroupName, id.ComponentName)

	result, err := client.APIKeysGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			log.Printf("[DEBUG] %s not found, removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("application_insights_id", appInsightsId.ID())

	if model := result.Model; model != nil {
		d.Set("name", model.Name)
		if props := model.LinkedReadProperties; props != nil {
			if err := d.Set("read_permissions", flattenApplicationInsightsAPIKeyLinkedProperties(props)); err != nil {
				return fmt.Errorf("flattening `read_permissions `: %s", err)
			}
		}
		if props := model.LinkedWriteProperties; props != nil {
			if err := d.Set("write_permissions", flattenApplicationInsightsAPIKeyLinkedProperties(props)); err != nil {
				return fmt.Errorf("flattening `write_permissions `: %s", err)
			}
		}
	}

	return nil
}
