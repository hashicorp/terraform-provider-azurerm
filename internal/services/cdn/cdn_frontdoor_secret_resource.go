// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/secrets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
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
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := secrets.ParseSecretID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CdnFrontDoorSecretName,
			},

			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: profiles.ValidateProfileID,
			},

			"secret": {
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
									"expiration_date": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"key_vault_certificate_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeCertificate),
									},

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

	profile, err := profiles.ParseProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := secrets.NewSecretID(profile.SubscriptionId, profile.ResourceGroupName, profile.ProfileName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_secret", id.ID())
		}
	}

	secretParams, err := expandCdnFrontDoorSecretParameters(ctx, d.Get("secret").([]interface{}), meta.(*clients.Client))
	if err != nil {
		return fmt.Errorf("expanding `secret`: %+v", err)
	}

	props := secrets.Secret{
		Properties: &secrets.SecretProperties{
			Parameters: secretParams,
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorSecretRead(d, meta)
}

func resourceCdnFrontDoorSecretRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecretsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := secrets.ParseSecretID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.SecretName)
	d.Set("cdn_frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())
	d.Set("cdn_frontdoor_profile_name", id.ProfileName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			customerCertificate, err := flattenCdnFrontDoorSecretParameters(ctx, props.Parameters, meta)
			if err != nil {
				return fmt.Errorf("flattening `secret`: %+v", err)
			}

			if err := d.Set("secret", customerCertificate); err != nil {
				return fmt.Errorf("setting `secret`: %+v", err)
			}
		}
	}

	return nil
}

func resourceCdnFrontDoorSecretDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecretsClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := secrets.ParseSecretID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandCdnFrontDoorSecretParameters(ctx context.Context, input []interface{}, clients *clients.Client) (secrets.SecretParameters, error) {
	v := input[0].(map[string]interface{})

	cc := v["customer_certificate"].(map[string]interface{})

	certificateId, err := keyvault.ParseNestedItemID(cc["key_vault_certificate_id"].(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeCertificate)
	if err != nil {
		return nil, err
	}

	keyVaultBaseId, err := clients.KeyVault.KeyVaultIDFromBaseUrl(ctx, commonids.NewSubscriptionID(clients.Account.SubscriptionId), certificateId.KeyVaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Key Vault Resource ID from the Key Vault Base URL (`%s`): %w", certificateId.KeyVaultBaseURL, err)
	}

	if keyVaultBaseId == nil {
		return nil, fmt.Errorf("retrieving the Key Vault Resource ID from the Key Vault Base URL (`%s`): id was nil", certificateId.KeyVaultBaseURL)
	}

	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultBaseId)
	if err != nil {
		return nil, err
	}

	useLatest := certificateId.Version == ""
	customerCertificate := &secrets.CustomerCertificateParameters{
		Type: secrets.SecretTypeCustomerCertificate,
		SecretSource: secrets.ResourceReference{
			Id: pointer.To(keyVaultParse.NewSecretVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, certificateId.Name).ID()),
		},
		UseLatestVersion: pointer.To(useLatest),
	}

	if !useLatest {
		customerCertificate.SecretVersion = pointer.To(certificateId.Version)
	}

	return customerCertificate, nil
}

func flattenCdnFrontDoorSecretParameters(ctx context.Context, input secrets.SecretParameters, meta interface{}) ([]interface{}, error) {
	client := meta.(*clients.Client).KeyVault

	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	result := make(map[string]interface{})

	customerCertificate, ok := input.(secrets.CustomerCertificateParameters)
	if !ok {
		return nil, fmt.Errorf("received an unexpected type (`%T`)", input.(secrets.CustomerCertificateParameters))
	}

	secretSourceId, err := keyVaultParse.SecretVersionlessID(pointer.From(customerCertificate.SecretSource.Id))
	if err != nil {
		return nil, err
	}

	keyVaultId := commonids.NewKeyVaultID(secretSourceId.SubscriptionId, secretSourceId.ResourceGroup, secretSourceId.VaultName)
	keyVaultBaseUri, err := client.BaseUriForKeyVault(ctx, keyVaultId)
	if err != nil {
		return nil, fmt.Errorf("looking up Base URI for Certificate %q in %s: %+v", secretSourceId.SecretName, keyVaultId, err)
	}

	keyVaultCertificateId, err := keyvault.NewNestedItemID(*keyVaultBaseUri, keyvault.NestedItemTypeCertificate, secretSourceId.SecretName, pointer.From(customerCertificate.SecretVersion))
	if err != nil {
		return nil, err
	}

	certificateID := keyVaultCertificateId.ID()
	if pointer.From(customerCertificate.UseLatestVersion) {
		certificateID = keyVaultCertificateId.VersionlessID()
	}

	result["customer_certificate"] = []interface{}{
		map[string]interface{}{
			"expiration_date":           pointer.From(customerCertificate.ExpirationDate),
			"key_vault_certificate_id":  certificateID,
			"subject_alternative_names": utils.FlattenStringSlice(customerCertificate.SubjectAlternativeNames),
		},
	}
	results = append(results, result)

	return results, nil
}
