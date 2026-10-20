// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package attestation_provider

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/attestation/2022-08-01/attestation"
)

func resourceAttestationProviderCreate(d *pluginsdk.ResourceData, meta any) error {
	attestationClients := meta.(*clients.Client).Attestation
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := attestationproviders.NewAttestationProvidersID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := attestationClients.ProviderClient.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError(attestationProviderResourceName, id.ID())
		}
	}

	props := attestationproviders.AttestationServiceCreationParams{
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: attestationproviders.AttestationServiceCreationSpecificParams{},
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
	}

	// NOTE: This maybe an slice in a future release or even a slice of slices
	//    The service team does not currently have any user data for this
	//    pluginsdk.
	policySigningCertificate := d.Get("policy_signing_certificate_data").(string)

	if policySigningCertificate != "" {
		block, _ := pem.Decode([]byte(policySigningCertificate))
		if block == nil {
			return fmt.Errorf("invalid X.509 certificate, unable to decode")
		}

		v := base64.StdEncoding.EncodeToString(block.Bytes)
		props.Properties.PolicySigningCertificates = expandArmAttestationProviderJSONWebKeySet(v)
	}

	if _, err := attestationClients.ProviderClient.Create(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	dataPlaneUri, err := attestationClients.DataPlaneEndpointForProvider(ctx, id)
	if err != nil {
		return fmt.Errorf("determining Data Plane URI for %s: %+v", id, err)
	}
	dataPlaneClient, err := attestationClients.DataPlaneClientWithEndpoint(*dataPlaneUri)
	if err != nil {
		return fmt.Errorf("building Data Plane Client for %s: %+v", id, err)
	}

	if v := d.Get("open_enclave_policy_base64"); v != "" {
		if _, err = dataPlaneClient.Set(ctx, *dataPlaneUri, attestation.TypeOpenEnclave, d.Get("open_enclave_policy_base64").(string)); err != nil {
			return fmt.Errorf("updating value for `open_enclave_policy_base64`: %+v", err)
		}
	}
	if v := d.Get("sgx_enclave_policy_base64"); v != "" {
		if _, err = dataPlaneClient.Set(ctx, *dataPlaneUri, attestation.TypeSgxEnclave, d.Get("sgx_enclave_policy_base64").(string)); err != nil {
			return fmt.Errorf("updating value for `sgx_enclave_policy_base64`: %+v", err)
		}
	}
	if v := d.Get("tpm_policy_base64"); v != "" {
		if _, err = dataPlaneClient.Set(ctx, *dataPlaneUri, attestation.TypeTpm, d.Get("tpm_policy_base64").(string)); err != nil {
			return fmt.Errorf("updating value for `tpm_policy_base64`: %+v", err)
		}
	}
	if v := d.Get("sev_snp_policy_base64"); v != "" {
		if _, err = dataPlaneClient.Set(ctx, *dataPlaneUri, attestation.TypeSevSnpVM, d.Get("sev_snp_policy_base64").(string)); err != nil {
			return fmt.Errorf("updating value for `sev_snp_policy_base64`: %+v", err)
		}
	}

	return resourceAttestationProviderRead(d, meta)
}

func expandArmAttestationProviderJSONWebKeySet(pem string) *attestationproviders.JSONWebKeySet {
	if len(pem) == 0 {
		return nil
	}

	result := attestationproviders.JSONWebKeySet{
		Keys: expandArmAttestationProviderJSONWebKeyArray(pem),
	}

	return &result
}

func expandArmAttestationProviderJSONWebKeyArray(pem string) *[]attestationproviders.JSONWebKey {
	results := make([]attestationproviders.JSONWebKey, 0)

	result := attestationproviders.JSONWebKey{
		Kty: "RSA",
		X5c: pointer.To([]string{pem}),
	}

	results = append(results, result)

	return &results
}
