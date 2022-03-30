package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	cdnfrontdoorsecretparams "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/frontdoorsecretparams"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontdoorSecret() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorSecretCreate,
		Read:   resourceCdnFrontdoorSecretRead,
		Delete: resourceCdnFrontdoorSecretDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontdoorSecretID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorProfileID,
			},

			"secret_parameters": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						// NOTE: Not supported at GA
						// "url_signing_key": {
						// 	Type:     pluginsdk.TypeList,
						// 	Optional: true,
						// 	MaxItems: 1,

						// 	Elem: &pluginsdk.Resource{
						// 		Schema: map[string]*pluginsdk.Schema{

						// 			"key_id": {
						// 				Type:         pluginsdk.TypeString,
						// 				Required:     true,
						// 				ValidateFunc: validation.StringIsNotEmpty,
						// 			},

						// 			"secret_source_id": {
						// 				Type:         pluginsdk.TypeString,
						// 				Required:     true,
						// 				ValidateFunc: validation.StringIsNotEmpty,
						// 			},

						// 			"secret_version": {
						// 				Type:         pluginsdk.TypeString,
						// 				Required:     true,
						// 				ValidateFunc: validation.StringIsNotEmpty,
						// 			},
						// 		},
						// 	},
						// },

						// NOTE: Not supported at GA
						// "managed_certificate": {
						// 	Type:     pluginsdk.TypeList,
						// 	Optional: true,

						// 	Elem: &pluginsdk.Resource{
						// 		Schema: map[string]*pluginsdk.Schema{

						// 			"certificate_type": {
						// 				Type:     pluginsdk.TypeString,
						// 				Optional: true,
						// 				Default:  string(track1.TypeBasicSecretParametersTypeManagedCertificate),
						// 				ValidateFunc: validation.StringInSlice([]string{
						// 					string(track1.TypeBasicSecretParametersTypeManagedCertificate),
						// 				}, false),
						// 			},
						// 		},
						// 	},
						// },

						"customer_certificate": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									// TODO: Add validation: No Azure Key Vault name found in provided SecretSource id
									"secret_source_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"secret_version": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"use_latest": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
									},

									"subject_alternative_names": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 100,

										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},

						// NOTE: Not supported at GA
						// "azure_first_party_managed_certificate": {
						// 	Type:     pluginsdk.TypeList,
						// 	Optional: true,

						// 	Elem: &pluginsdk.Resource{
						// 		Schema: map[string]*pluginsdk.Schema{

						// 			"certificate_type": {
						// 				Type:     pluginsdk.TypeString,
						// 				Optional: true,
						// 				Default:  string(track1.TypeBasicSecretParametersTypeAzureFirstPartyManagedCertificate),
						// 				ValidateFunc: validation.StringInSlice([]string{
						// 					string(track1.TypeBasicSecretParametersTypeAzureFirstPartyManagedCertificate),
						// 				}, false),
						// 			},
						// 		},
						// 	},
						// },
					},
				},
			},

			"cdn_frontdoor_profile_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCdnFrontdoorSecretCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorSecretsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := parse.FrontdoorProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontdoorSecretID(profileId.SubscriptionId, profileId.ResourceGroup, profileId.ProfileName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.SecretName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_secret", id.ID())
		}
	}

	secretParams, err := expandCdnFrontdoorBasicSecretParameters(d.Get("secret_parameters").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding %q: %+v", "secret_parameters", err)
	}

	props := track1.Secret{
		SecretProperties: &track1.SecretProperties{
			Parameters: secretParams,
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.SecretName, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontdoorSecretRead(d, meta)
}

func resourceCdnFrontdoorSecretRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorSecretsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorSecretID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.SecretName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.SecretName)

	d.Set("cdn_frontdoor_profile_id", parse.NewFrontdoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.SecretProperties; props != nil {
		if err := d.Set("parameters", flattenSecretSecretParameters(&props.Parameters)); err != nil {
			return fmt.Errorf("setting `parameters`: %+v", err)
		}
		d.Set("cdn_frontdoor_profile_name", props.ProfileName)
	}

	return nil
}

func resourceCdnFrontdoorSecretDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorSecretsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorSecretID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.SecretName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandCdnFrontdoorBasicSecretParameters(input []interface{}) (track1.BasicSecretParameters, error) {
	results := make([]track1.BasicSecretParameters, 0)

	type expandfunc func(input []interface{}) (*track1.BasicSecretParameters, error)

	m := *cdnfrontdoorsecretparams.InitializeCdnFrontdoorSecretMappings()

	secrets := map[string]expandfunc{
		m.UrlSigningKey.ConfigName:                     cdnfrontdoorsecretparams.ExpandCdnFrontdoorUrlSigningKeyParameters,
		m.ManagedCertificate.ConfigName:                cdnfrontdoorsecretparams.ExpandCdnFrontdoorManagedCertificateParameters,
		m.CustomerCertificate.ConfigName:               cdnfrontdoorsecretparams.ExpandCdnFrontdoorCustomerCertificateParameters,
		m.AzureFirstPartyManagedCertificate.ConfigName: cdnfrontdoorsecretparams.ExpandCdnFrontdoorAzureFirstPartyManagedCertificateyParameters,
	}

	basicSecretParameters := input[0].(map[string]interface{})

	for secretName, expand := range secrets {
		if config := basicSecretParameters[secretName].([]interface{}); len(config) > 0 {
			expanded, err := expand(config)
			if err != nil {
				return nil, err
			}

			if expanded != nil {
				results = append(results, *expanded)
			}
		}
	}

	if len(results) > 1 {
		return nil, fmt.Errorf("%q is invalid, you may only have one %q defined in the configuration file, got %d", "secret_parameters", "secret_parameter", len(results))
	}

	for _, basicSecretParameter := range results {
		if basicSecretParameter != nil {
			return basicSecretParameter, nil
		}
	}

	return nil, fmt.Errorf("unknown secret parameter type encountered")
}

func flattenSecretSecretParameters(input *track1.BasicSecretParameters) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	// TODO: Not returned pull from state if exists
	// result := make(map[string]interface{})
	// result["type"] = input.Type
	// return append(results, result)
	return nil
}
