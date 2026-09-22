// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_application

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/application"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBatchApplicationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.ApplicationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := application.ParseApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Batch Application %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceBatchApplicationFlatten(d, id, resp.Model)
}

func resourceBatchApplicationFlatten(d *pluginsdk.ResourceData, id *application.ApplicationId, model *application.Application) error {
	d.Set("name", id.ApplicationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.BatchAccountName)

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("allow_updates", props.AllowUpdates)
			d.Set("default_version", props.DefaultVersion)
			d.Set("display_name", props.DisplayName)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
