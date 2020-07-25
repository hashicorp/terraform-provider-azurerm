package attestation

import (
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

func resourceArmAttestation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAttestationCreate,
		Read:   resourceArmAttestationRead,
		Update: resourceArmAttestationUpdate,
		Delete: resourceArmAttestationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AttestationId(id)
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

			"attestation_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"policy_signing_certificate": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alg": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"kid": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"kty": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"use": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"crv": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"d": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"dp": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"dq": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"e": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"k": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"n": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"p": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"q": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"qi": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"x": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"x5cs": {
										Type:     schema.TypeSet,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"y": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": tags.Schema(),

			"attest_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"trust_model": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceArmAttestationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Attestation %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_attestation", *existing.ID)
	}

	props := attestation.ServiceCreationParams{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &attestation.ServiceCreationSpecificParams{
			AttestationPolicy:         utils.String(d.Get("attestation_policy").(string)),
			PolicySigningCertificates: expandArmAttestationProviderJSONWebKeySet(d.Get("policy_signing_certificate").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if _, err := client.Create(ctx, resourceGroup, name, props); err != nil {
		return fmt.Errorf("creating Attestation %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Attestation %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Attestation %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	return resourceArmAttestationRead(d, meta)
}

func resourceArmAttestationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AttestationId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] attestation %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Attestation %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.StatusResult; props != nil {
		d.Set("attest_uri", props.AttestURI)
		d.Set("trust_model", props.TrustModel)
	}
	d.Set("type", resp.Type)
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAttestationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AttestationId(d.Id())
	if err != nil {
		return err
	}

	updateParams := attestation.ServicePatchParams{}
	if d.HasChange("tags") {
		updateParams.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, updateParams); err != nil {
		return fmt.Errorf("updating Attestation %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return resourceArmAttestationRead(d, meta)
}

func resourceArmAttestationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AttestationId(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Attestation %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}

func expandArmAttestationProviderJSONWebKeySet(input []interface{}) *attestation.JSONWebKeySet {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	result := attestation.JSONWebKeySet{
		Keys: expandArmAttestationProviderJSONWebKeyArray(v["key"].(*schema.Set).List()),
	}
	return &result
}

func expandArmAttestationProviderJSONWebKeyArray(input []interface{}) *[]attestation.JSONWebKey {
	results := make([]attestation.JSONWebKey, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		result := attestation.JSONWebKey{
			Alg: utils.String(v["alg"].(string)),
			Crv: utils.String(v["crv"].(string)),
			D:   utils.String(v["d"].(string)),
			Dp:  utils.String(v["dp"].(string)),
			Dq:  utils.String(v["dq"].(string)),
			E:   utils.String(v["e"].(string)),
			K:   utils.String(v["k"].(string)),
			Kid: utils.String(v["kid"].(string)),
			Kty: utils.String(v["kty"].(string)),
			N:   utils.String(v["n"].(string)),
			P:   utils.String(v["p"].(string)),
			Q:   utils.String(v["q"].(string)),
			Qi:  utils.String(v["qi"].(string)),
			Use: utils.String(v["use"].(string)),
			X:   utils.String(v["x"].(string)),
			X5c: utils.ExpandStringSlice(v["x5cs"].(*schema.Set).List()),
			Y:   utils.String(v["y"].(string)),
		}
		results = append(results, result)
	}
	return &results
}
