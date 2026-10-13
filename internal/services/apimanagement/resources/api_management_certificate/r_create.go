// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/certificate"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementCertificateCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.CertificatesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := certificate.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	data := d.Get("data").(string)
	password := d.Get("password").(string)
	keyVaultSecretId := d.Get("key_vault_secret_id").(string)
	keyVaultIdentity := d.Get("key_vault_identity_client_id").(string)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_certificate", id.ID())
			}
		}
	}

	parameters := certificate.CertificateCreateOrUpdateParameters{
		Properties: &certificate.CertificateCreateOrUpdateProperties{},
	}

	if keyVaultSecretId != "" {
		parsedSecretId, err := keyvault.ParseNestedItemID(keyVaultSecretId, keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret)
		if err != nil {
			return err
		}

		parameters.Properties.KeyVault = &certificate.KeyVaultContractCreateProperties{
			SecretIdentifier: pointer.To(parsedSecretId.ID()),
		}

		if keyVaultIdentity != "" {
			parameters.Properties.KeyVault.IdentityClientId = pointer.To(keyVaultIdentity)
		}
	}

	if data != "" {
		parameters.Properties.Data = pointer.To(data)
		parameters.Properties.Password = pointer.To(password)
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, certificate.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementCertificateRead(d, meta)
}
