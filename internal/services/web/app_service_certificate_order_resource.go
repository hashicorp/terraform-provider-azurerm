// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppServiceCertificateOrder() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceCertificateOrderCreateUpdate,
		Read:   resourceAppServiceCertificateOrderRead,
		Update: resourceAppServiceCertificateOrderCreateUpdate,
		Delete: resourceAppServiceCertificateOrderDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CertificateOrderID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AppServiceCertificateOrderResourceV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"auto_renew": {
				Type:     pluginsdk.TypeBool,
				Default:  true,
				Optional: true,
			},

			"certificates": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"certificate_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"key_vault_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"key_vault_secret_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"provisioning_state": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"csr": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"distinguished_name"},
			},

			"distinguished_name": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"csr"},
			},

			"key_size": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      2048,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"product_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Standard",
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"WildCard",
				}, false),
			},

			"validity_in_years": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 3),
			},

			"domain_verification_token": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"expiration_time": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"is_private_key_external": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"app_service_certificate_not_renewable_reasons": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"signed_certificate_thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"root_thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"intermediate_thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceAppServiceCertificateOrderCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesOrderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Certificate creation.")

	id := parse.NewCertificateOrderID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_app_service_certificate_order", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	distinguishedName := d.Get("distinguished_name").(string)
	csr := d.Get("csr").(string)
	keySize := d.Get("key_size").(int)
	autoRenew := d.Get("auto_renew").(bool)
	validityInYears := d.Get("validity_in_years").(int)

	properties := web.AppServiceCertificateOrderProperties{
		DistinguishedName: utils.String(distinguishedName),
		Csr:               utils.String(csr),
		KeySize:           utils.Int32(int32(keySize)),
		AutoRenew:         utils.Bool(autoRenew),
		ValidityInYears:   utils.Int32(int32(validityInYears)),
	}

	switch d.Get("product_type").(string) {
	case "Standard":
		properties.ProductType = web.CertificateProductTypeStandardDomainValidatedSsl
	case "WildCard":
		properties.ProductType = web.CertificateProductTypeStandardDomainValidatedWildCardSsl
	default:
		return fmt.Errorf("`product_type` must be `Standard` or `WildCard`")
	}

	certificateOrder := web.AppServiceCertificateOrder{
		AppServiceCertificateOrderProperties: &properties,
		Location:                             utils.String(location),
		Tags:                                 tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, certificateOrder)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %s", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceCertificateOrderRead(d, meta)
}

func resourceAppServiceCertificateOrderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesOrderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateOrderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Certificate Order %q (resource group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on AzureRM App Service Certificate Order %q: %+v", id.Name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.AppServiceCertificateOrderProperties; props != nil {
		d.Set("auto_renew", props.AutoRenew)
		d.Set("csr", props.Csr)
		d.Set("distinguished_name", props.DistinguishedName)
		d.Set("key_size", props.KeySize)
		d.Set("validity_in_years", props.ValidityInYears)
		d.Set("domain_verification_token", props.DomainVerificationToken)
		d.Set("status", string(props.Status))
		d.Set("is_private_key_external", props.IsPrivateKeyExternal)
		d.Set("certificates", flattenArmCertificateOrderCertificate(props.Certificates))
		d.Set("app_service_certificate_not_renewable_reasons", utils.FlattenStringSlice(props.AppServiceCertificateNotRenewableReasons))

		if productType := props.ProductType; productType == web.CertificateProductTypeStandardDomainValidatedSsl {
			d.Set("product_type", "Standard")
		} else if productType == web.CertificateProductTypeStandardDomainValidatedWildCardSsl {
			d.Set("product_type", "WildCard")
		}

		if expirationTime := props.ExpirationTime; expirationTime != nil {
			d.Set("expiration_time", expirationTime.Format(time.RFC3339))
		}

		if signedCertificate := props.SignedCertificate; signedCertificate != nil {
			d.Set("signed_certificate_thumbprint", signedCertificate.Thumbprint)
		}

		if root := props.Root; root != nil {
			d.Set("root_thumbprint", root.Thumbprint)
		}

		if intermediate := props.Intermediate; intermediate != nil {
			d.Set("intermediate_thumbprint", intermediate.Thumbprint)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAppServiceCertificateOrderDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesOrderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateOrderID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service Certificate Order %q (Resource Group %q)", id.Name, id.ResourceGroup)

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting App Service Certificate Order %q (Resource Group %q): %s)", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func flattenArmCertificateOrderCertificate(input map[string]*web.AppServiceCertificate) []interface{} {
	results := make([]interface{}, 0)

	for k, v := range input {
		result := make(map[string]interface{})

		result["certificate_name"] = k

		if keyVaultID := v.KeyVaultID; keyVaultID != nil {
			result["key_vault_id"] = *keyVaultID
		}
		if keyVaultSecretName := v.KeyVaultSecretName; keyVaultSecretName != nil {
			result["key_vault_secret_name"] = *keyVaultSecretName
		}
		result["provisioning_state"] = string(v.ProvisioningState)

		results = append(results, result)
	}

	return results
}
