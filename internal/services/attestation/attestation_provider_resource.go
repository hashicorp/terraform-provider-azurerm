package attestation

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/attestation/2020-10-01/attestation"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

		Schema: map[string]*pluginsdk.Schema{
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

			"policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// one type MUST have only one policy as most, add this validation in Create/Update
						"environment_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(attestation.OpenEnclave),
								string(attestation.SgxEnclave),
								string(attestation.Tpm),
							}, false),
						},

						"data": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`[A-Za-z0-9_-]+\.[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*`), ""),
						},
					},
				},
			},
		},
	}
}

type policyDef struct {
	Type attestation.Type
	Data string
}

func resourceAttestationProviderCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := attestationproviders.NewAttestationProvidersID(subscriptionId, resourceGroup, name)
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of exisiting %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_attestation_provider", id.ID())
	}

	// each type of policy should have 1 item at most.
	policies, err := expandPolicies(d.Get("policy").([]interface{}))
	if err != nil {
		return fmt.Errorf("configuration in policy invalid: %+v", err)
	}

	props := attestationproviders.AttestationServiceCreationParams{
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: attestationproviders.AttestationServiceCreationSpecificParams{
			// AttestationPolicy was deprecated in October of 2019
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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

	resp, err := client.Create(ctx, id, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// set policies
	if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.AttestUri != nil {
		url := *resp.Model.Properties.AttestUri
		cli := meta.(*clients.Client).Attestation.PolicyClient
		for _, policy := range policies {
			if _, err = cli.Set(ctx, url, policy.Type, policy.Data); err != nil {
				return fmt.Errorf("set policy: %+v", err)
			}
		}
	}

	d.SetId(id.ID())
	return resourceAttestationProviderRead(d, meta)
}

func resourceAttestationProviderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := attestationproviders.ParseAttestationProvidersID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AttestationProviderName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(resp.Model.Location))

		if props := resp.Model.Properties; props != nil {
			d.Set("attestation_uri", props.AttestUri)
			d.Set("trust_model", props.TrustModel)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceAttestationProviderUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := attestationproviders.ParseAttestationProvidersID(d.Id())
	if err != nil {
		return err
	}

	var hasChange bool
	updateParams := attestationproviders.AttestationServicePatchParams{}
	if d.HasChange("tags") {
		hasChange = true
		updateParams.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if hasChange {
		if _, err := client.Update(ctx, *id, updateParams); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	if d.HasChange("policy") {
		policies, err := expandPolicies(d.Get("policy").([]interface{}))
		if err != nil {
			return fmt.Errorf("expand policies: %+v", err)
		}
		policyClient := meta.(*clients.Client).Attestation.PolicyClient

		url := d.Get("attestation_uri").(string)
		if url == "" {
			log.Printf("[Warn] got empty attestation instance url")
		} else {
			for _, policy := range policies {
				if _, err = policyClient.Set(ctx, url, policy.Type, policy.Data); err != nil {
					return fmt.Errorf("set policy in %s: %+v", url, err)
				}
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

func expandArmAttestationProviderJSONWebKeyArray(pem string) *[]attestationproviders.JsonWebKey {
	results := make([]attestationproviders.JsonWebKey, 0)
	certs := []string{pem}

	result := attestationproviders.JsonWebKey{
		Kty: "RSA",
		X5c: &certs,
	}

	results = append(results, result)

	return &results
}

func expandPolicies(input []interface{}) (res []policyDef, err error) {
	for _, ins := range input {
		if ins == nil {
			continue
		}
		policy := ins.(map[string]interface{})
		typ := attestation.Type(policy["environment_type"].(string))

		for _, def := range res {
			if def.Type == typ {
				return nil, fmt.Errorf("repeated policy environment_type: %s", typ)
			}
		}
		res = append(res, policyDef{
			Type: typ,
			Data: policy["data"].(string),
		})
	}
	return
}
