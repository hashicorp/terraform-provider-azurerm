package cdn

import (
	"context"
	"fmt"
	"time"

	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	CdnFrontDoorsecretparams "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/frontdoorsecretparams"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	keyValutValidation "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontDoorSecret() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorSecretCreate,
		Read:   resourceCdnFrontDoorSecretRead,
		Delete: resourceCdnFrontDoorSecretDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontDoorSecretID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				// TODO: validation
				// WS: Fixed
				ValidateFunc: validate.CdnFrontDoorSecretName,
			},

			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorProfileID,
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
							Required: true,
							ForceNew: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"key_vault_certificate_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: keyValutValidation.KeyVaultChildIDWithOptionalVersion,
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

func resourceCdnFrontDoorSecretCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecretsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := parse.FrontDoorProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorSecretID(profileId.SubscriptionId, profileId.ResourceGroup, profileId.ProfileName, d.Get("name").(string))

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

	secretParams, err := expandCdnFrontDoorBasicSecretParameters(ctx, d.Get("secret_parameters").([]interface{}), meta.(*clients.Client))
	if err != nil {
		return fmt.Errorf("expanding %q: %+v", "secret_parameters", err)
	}

	props := cdn.Secret{
		SecretProperties: &cdn.SecretProperties{
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
	return resourceCdnFrontDoorSecretRead(d, meta)
}

func resourceCdnFrontDoorSecretRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecretsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorSecretID(d.Id())
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
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontDoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.SecretProperties; props != nil {
		var customerCertificate []interface{}
		if customerCertificate, err = flattenSecretSecretParameters(ctx, props.Parameters, meta); err != nil {
			return fmt.Errorf("flattening `secret_parameters`: %+v", err)
		}

		if err := d.Set("secret_parameters", customerCertificate); err != nil {
			return fmt.Errorf("setting `secret_parameters`: %+v", err)
		}

		d.Set("cdn_frontdoor_profile_name", props.ProfileName)
	}

	return nil
}

func resourceCdnFrontDoorSecretDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecretsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorSecretID(d.Id())
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

func expandCdnFrontDoorBasicSecretParameters(ctx context.Context, input []interface{}, clients *clients.Client) (cdn.BasicSecretParameters, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("%[1]q is invalid, expected to receive a %q, got %d", "secret_parameter", "Customer Certificate Parameter", len(input))
	}

	secretParameters := input[0].(map[string]interface{})
	m := *CdnFrontDoorsecretparams.InitializeCdnFrontDoorSecretMappings()
	config := secretParameters[m.CustomerCertificate.ConfigName]

	customerCertificate, err := CdnFrontDoorsecretparams.ExpandCdnFrontDoorCustomerCertificateParameters(ctx, config.([]interface{}), clients)
	if err != nil {
		return nil, err
	}

	return customerCertificate, nil
}

func flattenSecretSecretParameters(ctx context.Context, input cdn.BasicSecretParameters, meta interface{}) ([]interface{}, error) {
	client := meta.(*clients.Client).KeyVault

	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	result := make(map[string]interface{})
	fields := make(map[string]interface{})

	customerCertificate, ok := input.AsCustomerCertificateParameters()
	if !ok {
		return nil, fmt.Errorf("expected a Customer Certificate Parameter")
	}

	// Secret Source ID is what comes back from Frontdoor, now I need to build the URL from that...
	// secretSourceId: /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/keyVaultName1/secrets/certificateName1
	// id            : "https://[vaultName].vault.azure.net/certificates/[certificateName]/[certificateVersion]"
	// versionless id: "https://[vaultName].vault.azure.net/certificates/[certificateName]"

	secretSourceId, err := keyVaultParse.SecretVersionlessID(*customerCertificate.SecretSource.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse the %q field of the %q, got %q", "Secret Source", "Customer Certificate", *customerCertificate.SecretSource.ID)
	}

	var keyVaultCertificateId string

	if customerCertificate.UseLatestVersion != nil {
		// The API always sends back the version...
		var certificateVersion string
		var useLatest bool

		if customerCertificate.SecretVersion != nil {
			certificateVersion = *customerCertificate.SecretVersion
		}

		if customerCertificate.UseLatestVersion != nil {
			useLatest = *customerCertificate.UseLatestVersion
		}

		keyVaultId := keyVaultParse.NewVaultID(secretSourceId.SubscriptionId, secretSourceId.ResourceGroup, secretSourceId.VaultName)
		keyVaultBaseUri, err := client.BaseUriForKeyVault(ctx, keyVaultId)
		if err != nil {
			return nil, fmt.Errorf("looking up Base URI for Certificate %q in %s: %+v", secretSourceId.SecretName, keyVaultId, err)
		}

		if useLatest {
			// Build the versionless certificate id
			keyVaultCertificateId = fmt.Sprintf("%scertificates/%s", *keyVaultBaseUri, secretSourceId.SecretName)
		} else {
			// Build the certificate id with the version information
			keyVaultCertificateId = fmt.Sprintf("%scertificates/%s/%s", *keyVaultBaseUri, secretSourceId.SecretName, certificateVersion)
		}
	}

	fields["key_vault_certificate_id"] = keyVaultCertificateId

	if customerCertificate.SubjectAlternativeNames != nil {
		fields["subject_alternative_names"] = utils.FlattenStringSlice(customerCertificate.SubjectAlternativeNames)
	} else {
		fields["subject_alternative_names"] = make([]string, 0)
	}

	result["customer_certificate"] = []interface{}{fields}
	results = append(results, result)

	return results, nil
}
