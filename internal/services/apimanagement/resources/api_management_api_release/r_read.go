// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_release

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apirelease"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	apimanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiReleaseRead(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := account.From(meta).SubscriptionId
	client := apimanagementClient.From(meta).ApiReleasesClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := apirelease.ParseReleaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] apimanagement %s does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.ReleaseId)
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("api_id", api.NewApiID(subscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId).ID())
			d.Set("notes", pointer.From(props.Notes))
		}
	}
	return nil
}
