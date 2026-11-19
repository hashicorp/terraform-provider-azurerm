// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultCertificateCreate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
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

	client := keyVaultsClient.DataPlaneKeyVaultClient.Certificates.Clone(*keyVaultBaseUrl)
	certificateId := certificates.NewCertificateID(*keyVaultBaseUrl, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetCertificate(ctx, certificates.NewCertificateversionID(*keyVaultBaseUrl, name, ""))
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Certificate %q in %s: %s", name, *keyVaultBaseUrl, err)
			}
		}

		if existing.Model != nil && pointer.From(existing.Model.Id) != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_certificate", *existing.Model.Id)
		}
	}

	t := d.Get("tags").(map[string]any)
	policy, err := expandKeyVaultCertificatePolicy(d)
	if err != nil {
		return fmt.Errorf("expanding certificate policy: %s", err)
	}

	var newCert *certificates.CertificateBundle
	if v, ok := d.GetOk("certificate"); ok {
		// Import
		certificate := expandKeyVaultCertificate(v)
		importParameters := certificates.CertificateImportParameters{
			Value:  certificate.CertificateData,
			Pwd:    pointer.To(certificate.CertificatePassword),
			Policy: policy,
			Tags:   tags.Expand(t),
		}
		resp, err := client.ImportCertificate(ctx, certificateId, importParameters)
		if err != nil {
			if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedCerts && response.WasConflict(resp.HttpResponse) {
				if err = recoverDeletedCertificate(ctx, d, meta, *keyVaultBaseUrl, name); err != nil {
					return fmt.Errorf("recover deleted certificate: %+v", err)
				}
				resp, err = client.ImportCertificate(ctx, certificateId, importParameters)
				if err != nil {
					return fmt.Errorf("update recovered certificate: %+v", err)
				}
			} else {
				return err
			}
		}
		newCert = resp.Model
	} else {
		// Generate new
		var httpResp *http.Response
		newCert, httpResp, err = createCertificate(d, meta)
		if err != nil {
			if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedCerts && response.WasConflict(httpResp) {
				if err = recoverDeletedCertificate(ctx, d, meta, *keyVaultBaseUrl, name); err != nil {
					return fmt.Errorf("recover deleted certificate: %+v", err)
				}
				// after we recovered the existing certificate we still have to apply our changes
				newCert, _, err = createCertificate(d, meta)
				if err != nil {
					return fmt.Errorf("update recovered certificate: %+v", err)
				}
			} else {
				return err
			}
		}
	}

	if newCert == nil || newCert.Id == nil {
		return fmt.Errorf("retrieving Certificate %q in %s: unable to determine certificate ID", name, *keyVaultBaseUrl)
	}

	parsedCertificateId, err := keyvault.ParseNestedItemID(*newCert.Id, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeCertificate)
	if err != nil {
		return err
	}
	d.SetId(parsedCertificateId.ID())

	return resourceKeyVaultCertificateRead(d, meta)
}
