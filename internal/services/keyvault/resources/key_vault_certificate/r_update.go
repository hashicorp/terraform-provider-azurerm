// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/certificates"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultCertificateUpdate(d *schema.ResourceData, meta any) error {
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keyvault.ParseNestedItemID(d.Id(), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeCertificate)
	if err != nil {
		return err
	}

	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	meta.(*clients.Client).KeyVault.AddToCache(*keyVaultId, id.KeyVaultBaseURL)

	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Certificates.Clone(id.KeyVaultBaseURL)

	if d.HasChange("certificate") {
		if v, ok := d.GetOk("certificate"); ok {
			// Import new version of certificate
			certificate := expandKeyVaultCertificate(v)
			importParameters := certificates.CertificateImportParameters{
				Value: certificate.CertificateData,
				Pwd:   pointer.To(certificate.CertificatePassword),
			}
			resp, err := client.ImportCertificate(ctx, certificates.NewCertificateID(id.KeyVaultBaseURL, id.Name), importParameters)
			if err != nil {
				return err
			}

			if resp.Model == nil || resp.Model.Id == nil {
				return fmt.Errorf("error: Certificate %q in Vault %q get nil ID from server", id.Name, id.KeyVaultBaseURL)
			}

			certificateId, err := keyvault.ParseNestedItemID(*resp.Model.Id, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeCertificate)
			if err != nil {
				return err
			}
			d.SetId(certificateId.ID())
		}
	}

	// update lifetime_action only should not recreate a certificate
	var lifeTimeOld, lifeTimeNew any
	var policyOld, policyNew map[string]any

	policyOldRaw, policyNewRaw := d.GetChange("certificate_policy")
	policyOldList, policyNewList := policyOldRaw.([]any), policyNewRaw.([]any)

	if len(policyOldList) > 0 {
		policyOld = policyOldList[0].(map[string]any)
		lifeTimeOld = policyOld["lifetime_action"]
		delete(policyOld, "lifetime_action")
	}
	if len(policyNewList) > 0 {
		policyNew = policyNewList[0].(map[string]any)
		lifeTimeNew = policyNew["lifetime_action"]
		delete(policyNew, "lifetime_action")
	}

	// do not recreate cerfiticate when only lifetime_action changes
	if !cmp.Equal(policyNewList, policyOldList) {
		policyNew["lifetime_action"] = lifeTimeNew
		newCert, _, err := createCertificate(d, meta)
		if err != nil {
			return err
		}
		if newCert == nil || newCert.Id == nil {
			return fmt.Errorf("retrieving Certificate %q in Vault %q: unable to determine certificate ID", id.Name, id.KeyVaultBaseURL)
		}
		certificateId, err := keyvault.ParseNestedItemID(*newCert.Id, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeCertificate)
		if err != nil {
			return err
		}
		d.SetId(certificateId.ID())
	}

	if updateLifetime := !cmp.Equal(lifeTimeOld, lifeTimeNew); d.HasChange("tags") || updateLifetime {
		patch := certificates.CertificateUpdateParameters{}
		if d.HasChange("tags") {
			if t, ok := d.GetOk("tags"); ok {
				patch.Tags = tags.Expand(t.(map[string]any))
			}
		}

		if updateLifetime {
			patch.Policy = &certificates.CertificatePolicy{
				LifetimeActions: expandKeyVaultCertificatePolicyLifetimeAction(lifeTimeNew),
			}
		}

		if _, err = client.UpdateCertificate(ctx, certificates.NewCertificateversionID(id.KeyVaultBaseURL, id.Name, ""), patch); err != nil {
			return err
		}
	}
	return resourceKeyVaultCertificateRead(d, meta)
}
