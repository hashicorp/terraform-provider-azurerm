package attestation

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/attestation/mgmt/2018-09-01-preview/attestation"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/attestation/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/attestation/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAttestationProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAttestationProviderCreate,
		Read:   resourceAttestationProviderRead,
		Update: resourceAttestationProviderUpdate,
		Delete: resourceAttestationProviderDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ProviderID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AttestationName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"policy_signing_certificate_data": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.IsCert,
			},

			"tags": tags.Schema(),

			"attestation_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"trust_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAttestationProviderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resourceId := parse.NewProviderID(subscriptionId, resourceGroup, name).ID("")
	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Attestation Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_attestation_provider", resourceId)
	}

	props := attestation.ServiceCreationParams{
		Location:   utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &attestation.ServiceCreationSpecificParams{
			// AttestationPolicy was deprecated in October of 2019
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	// NOTE: This maybe an slice in a future release or even a slice of slices
	//       The service team does not currently have any user data for this
	//       resource.
	policySigningCertificate := d.Get("policy_signing_certificate_data").(string)

	if policySigningCertificate != "" {
		block, _ := pem.Decode([]byte(policySigningCertificate))
		if block == nil {
			return fmt.Errorf("invalid X.509 certificate, unable to decode")
		}

		v := base64.StdEncoding.EncodeToString(block.Bytes)
		props.Properties.PolicySigningCertificates = expandArmAttestationProviderJSONWebKeySet(v)
	}

	if _, err := client.Create(ctx, resourceGroup, name, props); err != nil {
		return fmt.Errorf("creating Attestation Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(resourceId)
	return resourceAttestationProviderRead(d, meta)
}

func resourceAttestationProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProviderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AttestationProviderName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Attestation Provider %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Attestation Provider %q (Resource Group %q): %+v", id.AttestationProviderName, id.ResourceGroup, err)
	}

	d.Set("name", id.AttestationProviderName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.StatusResult; props != nil {
		d.Set("attestation_uri", props.AttestURI)
		d.Set("trust_model", props.TrustModel)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAttestationProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProviderID(d.Id())
	if err != nil {
		return err
	}

	updateParams := attestation.ServicePatchParams{}
	if d.HasChange("tags") {
		updateParams.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.AttestationProviderName, updateParams); err != nil {
		return fmt.Errorf("updating Attestation Provider %q (Resource Group %q): %+v", id.AttestationProviderName, id.ResourceGroup, err)
	}
	return resourceAttestationProviderRead(d, meta)
}

func resourceAttestationProviderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProviderID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.AttestationProviderName); err != nil {
		return fmt.Errorf("deleting Attestation Provider %q (Resource Group %q): %+v", id.AttestationProviderName, id.ResourceGroup, err)
	}
	return nil
}

func expandArmAttestationProviderJSONWebKeySet(pem string) *attestation.JSONWebKeySet {
	if len(pem) == 0 {
		return nil
	}

	result := attestation.JSONWebKeySet{
		Keys: expandArmAttestationProviderJSONWebKeyArray(pem),
	}

	return &result
}

func expandArmAttestationProviderJSONWebKeyArray(pem string) *[]attestation.JSONWebKey {
	results := make([]attestation.JSONWebKey, 0)
	certs := []string{pem}

	result := attestation.JSONWebKey{
		Kty: utils.String("RSA"),
		X5c: &certs,
	}

	results = append(results, result)

	return &results
}
