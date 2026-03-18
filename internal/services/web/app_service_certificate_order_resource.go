// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/certificateregistration/2023-12-01/appservicecertificateorders"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppServiceCertificateOrder() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceCertificateOrderCreate,
		Read:   resourceAppServiceCertificateOrderRead,
		Update: resourceAppServiceCertificateOrderUpdate,
		Delete: resourceAppServiceCertificateOrderDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := appservicecertificateorders.ParseCertificateOrderID(id)
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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceAppServiceCertificateOrderCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceCertificateOrdersClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := appservicecertificateorders.NewCertificateOrderID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
		return tf.ImportAsExistsError("azurerm_app_service_certificate_order", id.ID())
	}

	certificateOrder := appservicecertificateorders.AppServiceCertificateOrder{
		Properties: &appservicecertificateorders.AppServiceCertificateOrderProperties{
			AutoRenew:         pointer.To(d.Get("auto_renew").(bool)),
			Csr:               pointer.To(d.Get("csr").(string)),
			DistinguishedName: pointer.To(d.Get("distinguished_name").(string)),
			KeySize:           pointer.To(int64(d.Get("key_size").(int))),
			ProductType:       expandProductType(d.Get("product_type").(string)),
			ValidityInYears:   pointer.To(int64(d.Get("validity_in_years").(int))),
		},
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     utils.ExpandPtrMapStringString(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, certificateOrder); err != nil {
		return fmt.Errorf("creating %s: %w", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceCertificateOrderRead(d, meta)
}

func resourceAppServiceCertificateOrderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceCertificateOrdersClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := appservicecertificateorders.ParseCertificateOrderID(d.Id())
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

	d.Set("name", id.CertificateOrderName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("tags", model.Tags)

		if props := model.Properties; props != nil {
			d.Set("auto_renew", props.AutoRenew)
			d.Set("csr", props.Csr)
			d.Set("distinguished_name", props.DistinguishedName)
			d.Set("key_size", props.KeySize)
			d.Set("validity_in_years", props.ValidityInYears)
			d.Set("domain_verification_token", props.DomainVerificationToken)
			d.Set("status", pointer.FromEnum(props.Status))
			d.Set("is_private_key_external", props.IsPrivateKeyExternal)
			d.Set("certificates", flattenArmCertificateOrderCertificate(props.Certificates))
			d.Set("app_service_certificate_not_renewable_reasons", flattenAppServiceCertificateNotRenewableReasons(props.AppServiceCertificateNotRenewableReasons))
			d.Set("expiration_time", props.ExpirationTime)
			d.Set("product_type", flattenProductType(props.ProductType))

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
	}

	return nil
}

func resourceAppServiceCertificateOrderUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceCertificateOrdersClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := appservicecertificateorders.ParseCertificateOrderID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	props := existing.Model.Properties

	if d.HasChange("auto_renew") {
		props.AutoRenew = pointer.To(d.Get("auto_renew").(bool))
	}

	if d.HasChange("csr") {
		props.Csr = pointer.To(d.Get("csr").(string))
	}

	if d.HasChange("distinguished_name") {
		props.DistinguishedName = pointer.To(d.Get("distinguished_name").(string))
	}

	if d.HasChange("key_size") {
		props.KeySize = pointer.To(int64(d.Get("key_size").(int)))
	}

	if d.HasChange("product_type") {
		props.ProductType = expandProductType(d.Get("product_type").(string))
	}

	if d.HasChange("validity_in_years") {
		props.ValidityInYears = pointer.To(int64(d.Get("validity_in_years").(int)))
	}

	if d.HasChange("tags") {
		existing.Model.Tags = utils.ExpandPtrMapStringString(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %w", id, err)
	}

	return resourceAppServiceCertificateOrderRead(d, meta)
}

func resourceAppServiceCertificateOrderDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceCertificateOrdersClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := appservicecertificateorders.ParseCertificateOrderID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func flattenArmCertificateOrderCertificate(input *map[string]appservicecertificateorders.AppServiceCertificate) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for k, v := range *input {
		result := make(map[string]interface{})

		result["certificate_name"] = k

		if keyVaultID := v.KeyVaultId; keyVaultID != nil {
			result["key_vault_id"] = *keyVaultID
		}
		if keyVaultSecretName := v.KeyVaultSecretName; keyVaultSecretName != nil {
			result["key_vault_secret_name"] = *keyVaultSecretName
		}
		result["provisioning_state"] = pointer.FromEnum(v.ProvisioningState)

		results = append(results, result)
	}

	return results
}

func flattenAppServiceCertificateNotRenewableReasons(input *[]appservicecertificateorders.ResourceNotRenewableReason) []string {
	results := make([]string, 0)

	if input == nil {
		return results
	}

	for _, v := range *input {
		results = append(results, string(v))
	}

	return results
}

func expandProductType(input string) appservicecertificateorders.CertificateProductType {
	if input == "WildCard" {
		return appservicecertificateorders.CertificateProductTypeStandardDomainValidatedWildCardSsl
	}
	return appservicecertificateorders.CertificateProductTypeStandardDomainValidatedSsl
}

func flattenProductType(input appservicecertificateorders.CertificateProductType) string {
	if input == appservicecertificateorders.CertificateProductTypeStandardDomainValidatedWildCardSsl {
		return "WildCard"
	}
	return "Standard"
}
