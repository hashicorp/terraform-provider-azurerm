// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package attestation_provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-name basicForResourceIdentity

const attestationProviderResourceName = "azurerm_attestation_provider"

func RegisterResource() *pluginsdk.Resource {
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

		Importer: pluginsdk.ImporterValidatingIdentity(&attestationproviders.AttestationProvidersId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&attestationproviders.AttestationProvidersId{}),
		},

		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, i any) error {
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
			return map[string]*pluginsdk.Schema{
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
		}(),
	}
}
