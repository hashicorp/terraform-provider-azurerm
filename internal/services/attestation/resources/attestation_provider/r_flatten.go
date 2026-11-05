// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package attestation_provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/attestation/2022-08-01/attestation"
)

func resourceAttestationProviderFlatten(ctx context.Context, attestationClients *client.Client, d *pluginsdk.ResourceData, id *attestationproviders.AttestationProvidersId, model *attestationproviders.AttestationProviders, includeResource bool) error {
	d.Set("name", id.AttestationProviderName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("attestation_uri", props.AttestUri)
			d.Set("trust_model", props.TrustModel)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	if includeResource {
		dataPlaneUri, err := attestationClients.DataPlaneEndpointForProvider(ctx, *id)
		if err != nil {
			return fmt.Errorf("determining Data Plane URI for %s: %+v", *id, err)
		}
		dataPlaneClient, err := attestationClients.DataPlaneClientWithEndpoint(*dataPlaneUri)
		if err != nil {
			return fmt.Errorf("building Data Plane Client for %s: %+v", *id, err)
		}

		// Status=400 Code="Bad request" Message="Tpm attestation is not supported in the 'UKSouth' region"
		openEnclavePolicy, err := dataPlaneClient.Get(ctx, *dataPlaneUri, attestation.TypeOpenEnclave)
		if err != nil && !response.WasBadRequest(openEnclavePolicy.Response.Response) {
			return fmt.Errorf("retrieving OpenEnclave Policy for %s: %+v", *id, err)
		}
		sgxEnclavePolicy, err := dataPlaneClient.Get(ctx, *dataPlaneUri, attestation.TypeSgxEnclave)
		if err != nil && !response.WasBadRequest(sgxEnclavePolicy.Response.Response) {
			return fmt.Errorf("retrieving SgxEnclave Policy for %s: %+v", *id, err)
		}
		tpmPolicy, err := dataPlaneClient.Get(ctx, *dataPlaneUri, attestation.TypeTpm)
		if err != nil && !response.WasBadRequest(tpmPolicy.Response.Response) {
			return fmt.Errorf("retrieving Tpm Policy for %s: %+v", *id, err)
		}
		sevSnpPolicy, err := dataPlaneClient.Get(ctx, *dataPlaneUri, attestation.TypeSevSnpVM)
		if err != nil && !response.WasBadRequest(sevSnpPolicy.Response.Response) {
			return fmt.Errorf("retrieving SEV-SNP Policy for %s: %+v", *id, err)
		}

		openEnclavePolicyData, err := base64DataFromAttestationJWT(openEnclavePolicy.Token)
		if err != nil {
			return fmt.Errorf("parsing OpenEnclave Policy for %s: %+v", *id, err)
		}
		d.Set("open_enclave_policy_base64", pointer.From(openEnclavePolicyData))

		sgxEnclavePolicyData, err := base64DataFromAttestationJWT(sgxEnclavePolicy.Token)
		if err != nil {
			return fmt.Errorf("parsing SgxEnclave Policy for %s: %+v", *id, err)
		}
		d.Set("sgx_enclave_policy_base64", pointer.From(sgxEnclavePolicyData))

		tpmPolicyData, err := base64DataFromAttestationJWT(tpmPolicy.Token)
		if err != nil {
			return fmt.Errorf("parsing Tpm Policy for %s: %+v", *id, err)
		}
		d.Set("tpm_policy_base64", pointer.From(tpmPolicyData))

		sevSnpPolicyData, err := base64DataFromAttestationJWT(sevSnpPolicy.Token)
		if err != nil {
			return fmt.Errorf("parsing SEV-SNP policy for %s: %+v", *id, err)
		}
		d.Set("sev_snp_policy_base64", pointer.From(sevSnpPolicyData))
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
