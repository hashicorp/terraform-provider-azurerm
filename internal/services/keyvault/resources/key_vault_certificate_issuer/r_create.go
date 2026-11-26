// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_issuer

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultCertificateIssuerCreateOrUpdate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("retrieving base uri for %s: %+v", *keyVaultId, err)
	}

	id := parse.NewIssuerID(*keyVaultBaseUri, name)
	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Certificates.Clone(*keyVaultBaseUri)
	issuerId := certificates.NewIssuernameID(*keyVaultBaseUri, name)
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetCertificateIssuer(ctx, issuerId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("failed to check for presence of existing Certificate Issuer %q (Key Vault %q): %s", name, *keyVaultBaseUri, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_key_vault_certificate_issuer", id.ID())
			}
		}
	}

	parameter := certificates.CertificateIssuerSetParameters{
		Provider:   d.Get("provider_name").(string),
		OrgDetails: &certificates.OrganizationDetails{},
	}

	if orgIdRaw, ok := d.GetOk("org_id"); ok {
		parameter.OrgDetails.Id = pointer.To(orgIdRaw.(string))
	}

	if adminsRaw, ok := d.GetOk("admin"); ok {
		parameter.OrgDetails.AdminDetails = expandKeyVaultCertificateIssuerOrganizationDetailsAdminDetails(adminsRaw.([]any))
	}

	accountId, gotAccountId := d.GetOk("account_id")
	password, gotPassword := d.GetOk("password")

	if gotAccountId && gotPassword {
		parameter.Credentials = &certificates.IssuerCredentials{
			AccountId: pointer.To(accountId.(string)),
			Pwd:       pointer.To(password.(string)),
		}
	}

	if _, err = client.SetCertificateIssuer(ctx, issuerId, parameter); err != nil {
		return fmt.Errorf("failed to set Certificate Issuer %q (Key Vault %q): %s", name, keyVaultId, err)
	}

	d.SetId(id.ID())

	return resourceKeyVaultCertificateIssuerRead(d, meta)
}

func expandKeyVaultCertificateIssuerOrganizationDetailsAdminDetails(vs []any) *[]certificates.AdministratorDetails {
	results := make([]certificates.AdministratorDetails, 0, len(vs))

	for _, v := range vs {
		administratorDetails := certificates.AdministratorDetails{}
		args := v.(map[string]any)
		if firstName, ok := args["first_name"]; ok {
			administratorDetails.FirstName = pointer.To(firstName.(string))
		}
		if lastName, ok := args["last_name"]; ok {
			administratorDetails.LastName = pointer.To(lastName.(string))
		}
		if emailAddress, ok := args["email_address"]; ok {
			administratorDetails.Email = pointer.To(emailAddress.(string))
		}
		if phone, ok := args["phone"]; ok {
			administratorDetails.Phone = pointer.To(phone.(string))
		}
		results = append(results, administratorDetails)
	}

	return &results
}
