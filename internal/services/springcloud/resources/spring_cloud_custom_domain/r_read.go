// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_custom_domain

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSpringCloudCustomDomainRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.CustomDomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform_rm.ParseDomainIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.DomainName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Spring Cloud Custom Domain %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Spring Cloud Custom Domain %q (Spring Cloud service %q / App %q / rcsource group %q): %+v", id.DomainName, id.SpringName, id.AppName, id.ResourceGroupName, err)
	}

	d.Set("name", id.DomainName)
	d.Set("spring_cloud_app_id", appplatform_rm.NewAppID(id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.AppName).ID())
	if props := resp.Properties; props != nil {
		d.Set("certificate_name", props.CertName)
		d.Set("thumbprint", props.Thumbprint)
	}

	return nil
}
