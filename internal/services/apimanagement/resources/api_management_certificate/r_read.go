// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_certificate

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementCertificateRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.CertificatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := certificate.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)

	if model := resp.Model; model != nil {
		d.Set("name", pointer.From(model.Name))
		if props := model.Properties; props != nil {
			d.Set("expiration", props.ExpirationDate)
			d.Set("subject", props.Subject)
			d.Set("thumbprint", props.Thumbprint)

			if keyvault := props.KeyVault; keyvault != nil {
				d.Set("key_vault_secret_id", pointer.From(keyvault.SecretIdentifier))
				d.Set("key_vault_identity_client_id", pointer.From(keyvault.IdentityClientId))
			}
		}
	}

	return nil
}
