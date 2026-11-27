// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_api_portal_custom_domain

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSpringCloudAPIPortalCustomDomainRead(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).APIPortalCustomDomainClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform.ParseApiPortalDomainIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.ApiPortalName, id.DomainName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] appplatform %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.DomainName)
	d.Set("spring_cloud_api_portal_id", appplatform.NewApiPortalID(id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ApiPortalName).ID())
	if props := resp.Properties; props != nil {
		d.Set("thumbprint", props.Thumbprint)
	}
	return nil
}
