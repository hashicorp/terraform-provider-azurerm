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
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyValutValidation "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
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
						"customer_certificate": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"key_vault_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: keyValutValidation.VaultID,
									},

									"key_vault_certificate_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: keyValutValidation.NestedItemName,
									},

									"key_vault_certificate_version": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"use_latest": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
									},

									// This is a READ-ONLY field as the secret resource is reading these from the certificate in the key vault...
									"subject_alternative_names": {
										Type:     pluginsdk.TypeList,
										Computed: true,

										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
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
		var customerCertificate []interface{}
		if customerCertificate, err = flattenSecretSecretParameters(props.Parameters); err != nil {
			return fmt.Errorf("flattening `secret_parameters`: %+v", err)
		}

		if err := d.Set("secret_parameters", customerCertificate); err != nil {
			return fmt.Errorf("setting `secret_parameters`: %+v", err)
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
		m.CustomerCertificate.ConfigName: cdnfrontdoorsecretparams.ExpandCdnFrontdoorCustomerCertificateParameters,
	}

	secretParameters := input[0].(map[string]interface{})

	// I will leave this in a loop as there will be more of these coming...
	for secretName, expand := range secrets {
		if config := secretParameters[secretName].([]interface{}); config != nil {
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
		return nil, fmt.Errorf("%[1]q is invalid, you may only have one %[1]q defined in the configuration file, got %d", "secret_parameter", len(results))
	}

	// There can be only one, so the first one that is not nil is the right one...
	for _, basicSecretParameter := range results {
		if basicSecretParameter != nil {
			return basicSecretParameter, nil
		}
	}

	return nil, fmt.Errorf("unknown secret parameter type encountered")
}

func flattenSecretSecretParameters(input track1.BasicSecretParameters) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	result := make(map[string]interface{})
	wrapper := make([]interface{}, 0)
	fields := make(map[string]interface{})

	customerCertificate, ok := input.AsCustomerCertificateParameters()
	if !ok {
		return nil, fmt.Errorf("expected a Customer Certificate Parameter")
	}

	secretSourceId, err := parse.FrontdoorKeyVaultSecretID(*customerCertificate.SecretSource.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse the %q field of the %q, got %q", "Secret Source", "Customer Certificate", *customerCertificate.SecretSource.ID)
	}

	keyVaultId := keyVaultParse.NewVaultID(secretSourceId.SubscriptionId, secretSourceId.ResourceGroup, secretSourceId.VaultName)

	fields["key_vault_id"] = keyVaultId.ID()
	fields["key_vault_certificate_name"] = secretSourceId.SecretName

	if customerCertificate.SecretVersion != nil {
		fields["key_vault_certificate_version"] = *customerCertificate.SecretVersion
	}

	if customerCertificate.UseLatestVersion != nil {
		fields["use_latest"] = *customerCertificate.UseLatestVersion
	}

	if customerCertificate.SubjectAlternativeNames != nil {
		fields["subject_alternative_names"] = utils.FlattenStringSlice(customerCertificate.SubjectAlternativeNames)
	}

	wrapper = append(wrapper, fields)
	result["customer_certificate"] = wrapper
	results = append(results, result)

	return results, nil
}
