// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attestation

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/attestation/2022-08-01/attestation"
)

func resourceAttestationProvider() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAttestationProviderCreate,
		Read:   resourceAttestationProviderRead,
		Update: resourceAttestationProviderUpdate,
		Delete: resourceAttestationProviderDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := attestationproviders.ParseAttestationProvidersID(id)
			return err
		}),

		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
			if o, n := diff.GetChange("open_enclave_policy_base64"); o.(string) != "" && n.(string) == "" {
				return fmt.Errorf("`open_enclave_policy_base64` can not be removed, add it to `ignore_changes` block to keep the default values")
			}

			if o, n := diff.GetChange("sgx_enclave_policy_base64"); o.(string) != "" && n.(string) == "" {
				return fmt.Errorf("`sgx_enclave_policy_base64` can not be removed, add it to `ignore_changes` block to keep the default values")
			}

			if o, n := diff.GetChange("tpm_policy_base64"); o.(string) != "" && n.(string) == "" {
				return fmt.Errorf("`tpm_policy_base64` can not be removed, add it to `ignore_changes` block to keep the default values")
			}

			if o, n := diff.GetChange("sev_snp_policy_base64"); o.(string) != "" && n.(string) == "" {
				return fmt.Errorf("`sev_snp_policy_base64` can not be removed, add it to `ignore_changes` block to keep the default values")
			}

			return nil
		},

		Schema: func() map[string]*pluginsdk.Schema {
			s := map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validate.AttestationProviderName,
				},

				"resource_group_name": commonschema.ResourceGroupName(),

				"location": commonschema.Location(),

				"policy_signing_certificate_data": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validate.IsCert,
				},

				"tags": commonschema.Tags(),

				"attestation_uri": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"trust_model": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"open_enclave_policy_base64": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.ContainsABase64UriEncodedJWTOfAStoredAttestationPolicy,
				},

				"sgx_enclave_policy_base64": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.ContainsABase64UriEncodedJWTOfAStoredAttestationPolicy,
				},

				"tpm_policy_base64": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.ContainsABase64UriEncodedJWTOfAStoredAttestationPolicy,
				},

				"sev_snp_policy_base64": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.ContainsABase64UriEncodedJWTOfAStoredAttestationPolicy,
				},
			}

			if !features.FourPointOhBeta() {
				s["policy"] = &pluginsdk.Schema{
					Type:       pluginsdk.TypeList,
					Optional:   true,
					Deprecated: "This field is no longer used and will be removed in v4.0 of the Azure Provider - use `open_enclave_policy_base64`, `sgx_enclave_policy_base64`, `tpm_policy_base64` and `sev_snp_policy_base64` instead.",
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"environment_type": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},

							"data": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},
						},
					},
				}
			}

			return s
		}(),
	}
}

func resourceAttestationProviderCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	attestationClients := meta.(*clients.Client).Attestation
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := attestationproviders.NewAttestationProvidersID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := attestationClients.ProviderClient.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of exisiting %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_attestation_provider", id.ID())
	}

	props := attestationproviders.AttestationServiceCreationParams{
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: attestationproviders.AttestationServiceCreationSpecificParams{},
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	// NOTE: This maybe an slice in a future release or even a slice of slices
	//       The service team does not currently have any user data for this
	//       pluginsdk.
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

func resourceAttestationProviderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	attestationClients := meta.(*clients.Client).Attestation
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := attestationproviders.ParseAttestationProvidersID(d.Id())
	if err != nil {
		return err
	}

	resp, err := attestationClients.ProviderClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

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
	if err != nil && !utils.ResponseWasBadRequest(openEnclavePolicy.Response) {
		return fmt.Errorf("retrieving OpenEnclave Policy for %s: %+v", *id, err)
	}
	sgxEnclavePolicy, err := dataPlaneClient.Get(ctx, *dataPlaneUri, attestation.TypeSgxEnclave)
	if err != nil && !utils.ResponseWasBadRequest(sgxEnclavePolicy.Response) {
		return fmt.Errorf("retrieving SgxEnclave Policy for %s: %+v", *id, err)
	}
	tpmPolicy, err := dataPlaneClient.Get(ctx, *dataPlaneUri, attestation.TypeTpm)
	if err != nil && !utils.ResponseWasBadRequest(tpmPolicy.Response) {
		return fmt.Errorf("retrieving Tpm Policy for %s: %+v", *id, err)
	}
	sevSnpPolicy, err := dataPlaneClient.Get(ctx, *dataPlaneUri, attestation.TypeSevSnpVM)
	if err != nil && !utils.ResponseWasBadRequest(sevSnpPolicy.Response) {
		return fmt.Errorf("retrieving SEV-SNP Policy for %s: %+v", *id, err)
	}

	d.Set("name", id.AttestationProviderName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("attestation_uri", props.AttestUri)
			d.Set("trust_model", props.TrustModel)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	openEnclavePolicyData, err := base64DataFromAttestationJWT(openEnclavePolicy.Token)
	if err != nil {
		return fmt.Errorf("parsing OpenEnclave Policy for %s: %+v", *id, err)
	}
	d.Set("open_enclave_policy_base64", utils.NormalizeNilableString(openEnclavePolicyData))

	sgxEnclavePolicyData, err := base64DataFromAttestationJWT(sgxEnclavePolicy.Token)
	if err != nil {
		return fmt.Errorf("parsing SgxEnclave Policy for %s: %+v", *id, err)
	}
	d.Set("sgx_enclave_policy_base64", utils.NormalizeNilableString(sgxEnclavePolicyData))

	tpmPolicyData, err := base64DataFromAttestationJWT(tpmPolicy.Token)
	if err != nil {
		return fmt.Errorf("parsing Tpm Policy for %s: %+v", *id, err)
	}
	d.Set("tpm_policy_base64", utils.NormalizeNilableString(tpmPolicyData))

	sevSnpPolicyData, err := base64DataFromAttestationJWT(sevSnpPolicy.Token)
	if err != nil {
		return fmt.Errorf("parsing SEV-SNP policy for %s: %+v", *id, err)
	}
	d.Set("sev_snp_policy_base64", utils.NormalizeNilableString(sevSnpPolicyData))

	if !features.FourPointOhBeta() {
		if err := d.Set("policy", []interface{}{}); err != nil {
			return fmt.Errorf("setting `policy`: %+v", err)
		}
	}

	return nil
}

func resourceAttestationProviderUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	attestationClients := meta.(*clients.Client).Attestation
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := attestationproviders.ParseAttestationProvidersID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("tags") {
		payload := attestationproviders.AttestationServicePatchParams{
			Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		}
		if _, err := attestationClients.ProviderClient.Update(ctx, *id, payload); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	if d.HasChanges("open_enclave_policy_base64", "sgx_enclave_policy_base64", "tpm_policy_base64", "sev_snp_policy_base64") {
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

func resourceAttestationProviderDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := attestationproviders.ParseAttestationProvidersID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}

func expandArmAttestationProviderJSONWebKeySet(pem string) *attestationproviders.JsonWebKeySet {
	if len(pem) == 0 {
		return nil
	}

	result := attestationproviders.JsonWebKeySet{
		Keys: expandArmAttestationProviderJSONWebKeyArray(pem),
	}

	return &result
}

func expandArmAttestationProviderJSONWebKeyArray(pem string) *[]attestationproviders.JSONWebKey {
	results := make([]attestationproviders.JSONWebKey, 0)
	certs := []string{pem}

	result := attestationproviders.JSONWebKey{
		Kty: "RSA",
		X5c: &certs,
	}

	results = append(results, result)

	return &results
}

func base64DataFromAttestationJWT(input *string) (*string, error) {
	if input == nil {
		return nil, nil
	}

	split := strings.Split(*input, ".")
	if len(split) != 3 {
		return nil, fmt.Errorf("expected the first token to have 3 segments but got %d", len(split))
	}
	// decode the JWT into a PolicyResult object
	decodedJwtSegment, err := base64.RawURLEncoding.DecodeString(split[1])
	if err != nil {
		return nil, fmt.Errorf("base64-decoding the first JWT Segment %q: %+v", split[1], err)
	}
	var firstResult attestation.PolicyResult
	if err := json.Unmarshal(decodedJwtSegment, &firstResult); err != nil {
		return nil, fmt.Errorf("unmarshaling into PolicyResult: %+v", err)
	}
	if firstResult.Policy == nil {
		return nil, nil
	}

	out := *firstResult.Policy
	return &out, nil
}
