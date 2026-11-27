// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_builder

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSpringCloudBuildServiceBuilderRead(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).BuildServiceBuilderClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform.ParseBuilderIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.BuildServiceName, id.BuilderName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] appplatform %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.BuilderName)
	d.Set("spring_cloud_service_id", commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName).ID())
	if props := resp.Properties; props != nil {
		if err := d.Set("build_pack_group", flattenBuildServiceBuilderBuildPacksGroupPropertiesArray(props.BuildpackGroups)); err != nil {
			return fmt.Errorf("setting `build_pack_group`: %+v", err)
		}
		if err := d.Set("stack", flattenBuildServiceBuilderStackProperties(props.Stack)); err != nil {
			return fmt.Errorf("setting `stack`: %+v", err)
		}
	}
	return nil
}
