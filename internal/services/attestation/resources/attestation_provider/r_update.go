// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package attestation_provider

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/attestation/2022-08-01/attestation"
)

func resourceAttestationProviderUpdate(d *pluginsdk.ResourceData, meta any) error {
	attestationClients := meta.(*clients.Client).Attestation
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := attestationproviders.ParseAttestationProvidersID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("tags") {
		payload := attestationproviders.AttestationServicePatchParams{
			Tags: tags.Expand(d.Get("tags").(map[string]any)),
		}
		if _, err := attestationClients.ProviderClient.Update(ctx, *id, payload); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	if d.HasChanges(
		"open_enclave_policy_base64",
		"sgx_enclave_policy_base64",
		"tpm_policy_base64",
		"sev_snp_policy_base64",
	) {
		dataPlaneUri, err := attestationClients.DataPlaneEndpointForProvider(ctx, *id)
		if err != nil {
			return fmt.Errorf("determining Data Plane URI for %s: %+v", *id, err)
		}
		dataPlaneClient, err := attestationClients.DataPlaneClientWithEndpoint(*dataPlaneUri)
		if err != nil {
			return fmt.Errorf("building Data Plane Client for %s: %+v", *id, err)
		}

		if d.HasChange("open_enclave_policy_base64") {
			if _, err = dataPlaneClient.Set(ctx, *dataPlaneUri, attestation.TypeOpenEnclave, d.Get("open_enclave_policy_base64").(string)); err != nil {
				return fmt.Errorf("updating value for `open_enclave_policy_base64`: %+v", err)
			}
		}
		if d.HasChange("sgx_enclave_policy_base64") {
			if _, err = dataPlaneClient.Set(ctx, *dataPlaneUri, attestation.TypeSgxEnclave, d.Get("sgx_enclave_policy_base64").(string)); err != nil {
				return fmt.Errorf("updating value for `sgx_enclave_policy_base64`: %+v", err)
			}
		}
		if d.HasChange("tpm_policy_base64") {
			if _, err = dataPlaneClient.Set(ctx, *dataPlaneUri, attestation.TypeTpm, d.Get("tpm_policy_base64").(string)); err != nil {
				return fmt.Errorf("updating value for `tpm_policy_base64`: %+v", err)
			}
		}
		if d.HasChange("sev_snp_policy_base64") {
			if _, err = dataPlaneClient.Set(ctx, *dataPlaneUri, attestation.TypeSevSnpVM, d.Get("sev_snp_policy_base64").(string)); err != nil {
				return fmt.Errorf("updating value for `sev_snp_policy_base64`: %+v", err)
			}
		}
	}

	return resourceAttestationProviderRead(d, meta)
}
