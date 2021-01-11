package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAppServiceCertificateOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppServiceCertificateOrderCreateUpdate,
		Read:   resourceAppServiceCertificateOrderRead,
		Update: resourceAppServiceCertificateOrderCreateUpdate,
		Delete: resourceAppServiceCertificateOrderDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CertificateOrderID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
				Default:          "Standard",
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

			"signed_certificate_thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"root_thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"intermediate_thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceAppServiceCertificateOrderCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesOrderClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Certificate creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
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

	switch d.Get("product_type").(string) {
	case "Standard":
		properties.ProductType = web.StandardDomainValidatedSsl
	case "WildCard":
		properties.ProductType = web.StandardDomainValidatedWildCardSsl
	default:
		return fmt.Errorf("Error setting `product_type` for App Service Certificate Order %q (Resource Group %q), either `Standard` or `WildCard`.", name, resourceGroup)
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

	return resourceAppServiceCertificateOrderRead(d, meta)
}

func resourceAppServiceCertificateOrderRead(d *schema.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("Error making Read request on AzureRM App Service Certificate Order %q: %+v", id.Name, err)
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

		if productType := props.ProductType; productType == web.StandardDomainValidatedSsl {
			d.Set("product_type", "Standard")
		} else if productType == web.StandardDomainValidatedWildCardSsl {
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

func resourceAppServiceCertificateOrderDelete(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Error deleting App Service Certificate Order %q (Resource Group %q): %s)", id.Name, id.ResourceGroup, err)
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
