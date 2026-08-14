// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package attestation

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/attestation/2022-08-01/attestation"
)

func dataSourceAttestationProvider() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmAttestationProviderRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"attestation_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"trust_model": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sev_snp_policy_base64": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"open_enclave_policy_base64": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"sgx_enclave_policy_base64": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tpm_policy_base64": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceArmAttestationProviderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	attestationClients := meta.(*clients.Client).Attestation
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := attestationproviders.NewAttestationProvidersID(subscriptionId, resourceGroup, name)

	resp, err := attestationClients.ProviderClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	dataPlaneUri, err := attestationClients.DataPlaneEndpointForProvider(ctx, id)
	if err != nil {
		return fmt.Errorf("determining Data Plane URI for %s: %+v", id, err)
	}

	dataPlaneClient, err := attestationClients.DataPlaneClientWithEndpoint(*dataPlaneUri)
	if err != nil {
		return fmt.Errorf("building Data Plane Client for %s: %+v", id, err)
	}

	// Status=400 Code="Bad request" Message="Tpm attestation is not supported in the 'UKSouth' region"
	openEnclavePolicy, err := dataPlaneClient.Get(ctx, *dataPlaneUri, attestation.TypeOpenEnclave)
	if err != nil && !utils.ResponseWasBadRequest(openEnclavePolicy.Response) {
		return fmt.Errorf("retrieving OpenEnclave Policy for %s: %+v", id, err)
	}
	sgxEnclavePolicy, err := dataPlaneClient.Get(ctx, *dataPlaneUri, attestation.TypeSgxEnclave)
	if err != nil && !utils.ResponseWasBadRequest(sgxEnclavePolicy.Response) {
		return fmt.Errorf("retrieving SgxEnclave Policy for %s: %+v", id, err)
	}
	tpmPolicy, err := dataPlaneClient.Get(ctx, *dataPlaneUri, attestation.TypeTpm)
	if err != nil && !utils.ResponseWasBadRequest(tpmPolicy.Response) {
		return fmt.Errorf("retrieving Tpm Policy for %s: %+v", id, err)
	}
	sevSnpPolicy, err := dataPlaneClient.Get(ctx, *dataPlaneUri, attestation.TypeSevSnpVM)
	if err != nil && !utils.ResponseWasBadRequest(sevSnpPolicy.Response) {
		return fmt.Errorf("retrieving SEV-SNP Policy for %s: %+v", id, err)
	}

	openEnclavePolicyData, err := base64DataFromAttestationJWT(openEnclavePolicy.Token)
	if err != nil {
		return fmt.Errorf("parsing OpenEnclave Policy for %s: %+v", id, err)
	}
	d.Set("open_enclave_policy_base64", pointer.From(openEnclavePolicyData))

	sgxEnclavePolicyData, err := base64DataFromAttestationJWT(sgxEnclavePolicy.Token)
	if err != nil {
		return fmt.Errorf("parsing SgxEnclave Policy for %s: %+v", id, err)
	}
	d.Set("sgx_enclave_policy_base64", pointer.From(sgxEnclavePolicyData))

	tpmPolicyData, err := base64DataFromAttestationJWT(tpmPolicy.Token)
	if err != nil {
		return fmt.Errorf("parsing Tpm Policy for %s: %+v", id, err)
	}
	d.Set("tpm_policy_base64", pointer.From(tpmPolicyData))

	sevSnpPolicyData, err := base64DataFromAttestationJWT(sevSnpPolicy.Token)
	if err != nil {
		return fmt.Errorf("parsing SEV-SNP policy for %s: %+v", id, err)
	}
	d.Set("sev_snp_policy_base64", pointer.From(sevSnpPolicyData))

	if resp.Model != nil {
		d.Set("location", location.Normalize(resp.Model.Location))

		if props := resp.Model.Properties; props != nil {
			d.Set("attestation_uri", props.AttestUri)
			d.Set("trust_model", props.TrustModel)
		}
		return tags.FlattenAndSet(d, resp.Model.Tags)
	}

	return nil
}
