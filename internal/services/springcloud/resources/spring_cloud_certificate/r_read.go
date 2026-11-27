// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_certificate

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

func resourceSpringCloudCertificateRead(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).CertificatesClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform.ParseCertificateIDInsensitively(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.CertificateName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Spring Cloud Certificate %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Spring Cloud Certificate %q (Spring Cloud Service %q / Resource Group %q): %+v", id.CertificateName, id.SpringName, id.ResourceGroupName, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("service_name", id.SpringName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if props, ok := resp.Properties.AsKeyVaultCertificateProperties(); ok && props != nil {
		d.Set("thumbprint", props.Thumbprint)
		d.Set("exclude_private_key", props.ExcludePrivateKey)
	}
	if props, ok := resp.Properties.AsContentCertificateProperties(); ok && props != nil {
		d.Set("thumbprint", props.Thumbprint)
	}

	return nil
}
