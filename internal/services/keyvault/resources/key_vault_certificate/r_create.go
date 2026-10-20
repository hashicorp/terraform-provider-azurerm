// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	kv "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultCertificateCreate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Base URI for Certificate %q in %s: %+v", name, *keyVaultId, err)
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetCertificate(ctx, *keyVaultBaseUrl, name, "")
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing Certificate %q in %s: %s", name, *keyVaultBaseUrl, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_certificate", *existing.ID)
		}
	}

	t := d.Get("tags").(map[string]any)
	policy, err := expandKeyVaultCertificatePolicy(d)
	if err != nil {
		return fmt.Errorf("expanding certificate policy: %s", err)
	}

	var newCert kv.CertificateBundle
	if v, ok := d.GetOk("certificate"); ok {
		// Import
		certificate := expandKeyVaultCertificate(v)
		importParameters := kv.CertificateImportParameters{
			Base64EncodedCertificate: pointer.To(certificate.CertificateData),
			Password:                 pointer.To(certificate.CertificatePassword),
			CertificatePolicy:        policy,
			Tags:                     tags.Expand(t),
		}
		newCert, err = client.ImportCertificate(ctx, *keyVaultBaseUrl, name, importParameters)
		if err != nil {
			if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedCerts && response.WasConflict(newCert.Response.Response) {
				if err = recoverDeletedCertificate(ctx, d, meta, *keyVaultBaseUrl, name); err != nil {
					return fmt.Errorf("recover deleted certificate: %+v", err)
				}
				newCert, err = client.ImportCertificate(ctx, *keyVaultBaseUrl, name, importParameters)
				if err != nil {
					return fmt.Errorf("update recovered certificate: %+v", err)
				}
			} else {
				return err
			}
		}
	} else {
		// Generate new
		newCert, err = createCertificate(d, meta)
		if err != nil {
			if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedCerts && response.WasConflict(newCert.Response.Response) {
				if err = recoverDeletedCertificate(ctx, d, meta, *keyVaultBaseUrl, name); err != nil {
					return fmt.Errorf("recover deleted certificate: %+v", err)
				}
				// after we recovered the existing certificate we still have to apply our changes
				newCert, err = createCertificate(d, meta)
				if err != nil {
					return fmt.Errorf("update recovered certificate: %+v", err)
				}
			} else {
				return err
			}
		}
	}

	certificateId, err := keyvault.ParseNestedItemID(*newCert.ID, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeCertificate)
	if err != nil {
		return err
	}
	d.SetId(certificateId.ID())

	return resourceKeyVaultCertificateRead(d, meta)
}
