package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceCertificateOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceCertificateOrderCreateUpdate,
		Read:   resourceArmAppServiceCertificateOrderRead,
		Update: resourceArmAppServiceCertificateOrderCreateUpdate,
		Delete: resourceArmAppServiceCertificateOrderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"auto_renew": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"key_vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"key_vault_secret_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"provisioning_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"csr": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"distinguished_name"},
			},

			"distinguished_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"csr"},
			},

			"key_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2048,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"product_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "standard",
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"WildCard",
				}, true),
			},

			"validity_in_years": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 3),
			},

			"domain_verification_token": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_certificate_issuance_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"next_auto_renewal_time_stamp": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"is_private_key_external": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"app_service_certificate_not_renewable_reasons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"signed_certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: appServiceCertificateDetailsSchema(),
				},
			},

			"root": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: appServiceCertificateDetailsSchema(),
				},
			},

			"intermediate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: appServiceCertificateDetailsSchema(),
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmAppServiceCertificateOrderCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Web.CertificatesOrderClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for App Service Certificate creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing App Service Certificate Order %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_certificate_order", *existing.ID)
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

	if v, ok := d.GetOk("product_type"); ok {
		productType := v.(string)
		if productType == "standard" {
			properties.ProductType = web.StandardDomainValidatedSsl
		} else if productType == "wildcard" {
			properties.ProductType = web.StandardDomainValidatedWildCardSsl
		}
	}

	certificateOrder := web.AppServiceCertificateOrder{
		AppServiceCertificateOrderProperties: &properties,
		Location:                             utils.String(location),
		Tags:                                 tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, certificateOrder)
	if err != nil {
		return fmt.Errorf("Error creating/updating App Service Certificate Order %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creating/updating of App Service Certificate Order %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service Certificate Order %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service Certificate Order %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceCertificateOrderRead(d, meta)
}

func resourceArmAppServiceCertificateOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Web.CertificatesOrderClient

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["certificateOrders"]

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Certificate Order %q (resource group %q) was not found - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM App Service Certificate Order %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

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
		d.Set("serial_number", props.SerialNumber)
		d.Set("is_private_key_external", props.IsPrivateKeyExternal)
		d.Set("certificates", flattenArmCertificateOrderCertificate(props.Certificates))
		d.Set("app_service_certificate_not_renewable_reasons", utils.FlattenStringSlice(props.AppServiceCertificateNotRenewableReasons))

		if productType := props.ProductType; productType == web.StandardDomainValidatedSsl {
			d.Set("product_type", "standard")
		} else if productType == web.StandardDomainValidatedWildCardSsl {
			d.Set("product_type", "wildcard")
		}

		if lastCertificateIssuanceTime := props.LastCertificateIssuanceTime; lastCertificateIssuanceTime != nil {
			d.Set("last_certificate_issuance_time", lastCertificateIssuanceTime.Format(time.RFC3339))
		}

		if nextAutoRenewalTimeStamp := props.NextAutoRenewalTimeStamp; nextAutoRenewalTimeStamp != nil {
			d.Set("next_auto_renewal_time_stamp", nextAutoRenewalTimeStamp.Format(time.RFC3339))
		}

		if expirationTime := props.ExpirationTime; expirationTime != nil {
			d.Set("expiration_time", expirationTime.Format(time.RFC3339))
		}

		d.Set("signed_certificate", flattenArmCertificateOrderCertificateDetails(props.SignedCertificate))
		d.Set("root", flattenArmCertificateOrderCertificateDetails(props.Root))
		d.Set("intermediate", flattenArmCertificateOrderCertificateDetails(props.Intermediate))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAppServiceCertificateOrderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Web.CertificatesOrderClient

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["certificateOrders"]

	log.Printf("[DEBUG] Deleting App Service Certificate Order %q (Resource Group %q)", name, resourceGroup)

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting App Service Certificate Order %q (Resource Group %q): %s)", name, resourceGroup, err)
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

func flattenArmCertificateOrderCertificateDetails(input *web.CertificateDetails) []interface{} {
	results := make([]interface{}, 0)

	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if version := input.Version; version != nil {
		result["version"] = version
	}

	if serialNumber := input.SerialNumber; serialNumber != nil {
		result["serial_number"] = serialNumber
	}

	if thumbprint := input.Thumbprint; thumbprint != nil {
		result["thumbprint"] = thumbprint
	}

	if subject := input.Subject; subject != nil {
		result["subject"] = subject
	}

	if notBefore := input.NotBefore; notBefore != nil {
		result["not_before"] = notBefore.Format(time.RFC3339)
	}

	if notAfter := input.NotAfter; notAfter != nil {
		result["not_after"] = notAfter.Format(time.RFC3339)
	}

	if signatureAlgorithm := input.SignatureAlgorithm; signatureAlgorithm != nil {
		result["signature_algorithm"] = signatureAlgorithm
	}

	if issuer := input.Issuer; issuer != nil {
		result["issuer"] = issuer
	}

	if rawData := input.RawData; rawData != nil {
		result["raw_data"] = rawData
	}

	results = append(results, result)

	return results
}

func appServiceCertificateDetailsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},

		"serial_number": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"thumbprint": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"subject": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"signature_algorithm": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"issuer": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"raw_data": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"not_before": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"not_after": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
