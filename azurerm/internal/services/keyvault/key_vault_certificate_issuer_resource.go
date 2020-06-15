package keyvault

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultCertificateIssuer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultCertificateIssuerCreate,
		Update: resourceArmKeyVaultCertificateIssuerUpdate,
		Read:   resourceArmKeyVaultCertificateIssuerRead,
		Delete: resourceArmKeyVaultCertificateIssuerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateKeyVaultCertificateIssuerName,
			},

			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DigiCert",
					"GlobalSign",
				}, true),
			},

			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"org_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"admins": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"phone": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmKeyVaultCertificateIssuerCreate(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	keyVaultBaseUri, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Certificate Issuer %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	if features.ShouldResourcesBeImported() {
		existing, err := client.GetCertificateIssuer(ctx, keyVaultBaseUri, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Certificate Issuer %q (Key Vault %q): %s", name, keyVaultBaseUri, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_certificate", *existing.ID)
		}
	}

	parameter := keyvault.CertificateIssuerSetParameters{}
	parameter.Provider = utils.String(d.Get("provider_name").(string))
	parameter.OrganizationDetails = &keyvault.OrganizationDetails{ID: utils.String(d.Get("org_id").(string))}
	parameter.Credentials = &keyvault.IssuerCredentials{AccountID: utils.String(d.Get("account_id").(string)), Password: utils.String(d.Get("password").(string))}
	resp, err := client.SetCertificateIssuer(ctx, keyVaultBaseUri, name, parameter)
	if err != nil {
		return fmt.Errorf("Error setting Certificate Issuer %q (Key Vault %q): %s", name, keyVaultId, err)
	}
	d.SetId(*resp.ID)

	return resourceArmKeyVaultCertificateIssuerRead(d, meta)
}

func resourceArmKeyVaultCertificateIssuerUpdate(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	keyVaultBaseUri, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Certificate Issuer %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	existing, err := client.GetCertificateIssuer(ctx, keyVaultBaseUri, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing Certificate Issuer %q (Key Vault %q): %s", name, keyVaultBaseUri, err)
		} else {
			return fmt.Errorf("KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist", name, keyVaultBaseUri)
		}
	}

	parameter := keyvault.CertificateIssuerSetParameters{}
	parameter.Provider = utils.String(d.Get("provider_name").(string))
	parameter.OrganizationDetails = &keyvault.OrganizationDetails{ID: utils.String(d.Get("org_id").(string))}
	parameter.Credentials = &keyvault.IssuerCredentials{AccountID: utils.String(d.Get("account_id").(string)), Password: utils.String(d.Get("password").(string))}
	resp, err := client.SetCertificateIssuer(ctx, keyVaultBaseUri, name, parameter)
	if err != nil {
		return fmt.Errorf("Error setting Certificate Issuer %q (Key Vault %q): %s", name, keyVaultId, err)
	}
	d.SetId(*resp.ID)

	return resourceArmKeyVaultCertificateIssuerRead(d, meta)
}

func resourceArmKeyVaultCertificateIssuerRead(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	keyVaultBaseUri, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Certificate Issuer %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	resp, err := client.GetCertificateIssuer(ctx, keyVaultBaseUri, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist", name, keyVaultBaseUri)
		}
		return fmt.Errorf("Error making Read request on Azure KeyVault Certificate Issuer %s: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("provider_name", resp.Provider)
	d.Set("org_id", resp.OrganizationDetails.ID)
	d.Set("account_id", resp.Credentials.AccountID)
	d.Set("password", resp.Credentials.Password)
	if resp.OrganizationDetails.AdminDetails != nil {
		d.Set("admins", flattenKeyVaultCertificateIssuerAdmins(resp.OrganizationDetails.AdminDetails))
	}

	return nil
}

func flattenKeyVaultCertificateIssuerAdmins(input *[]keyvault.AdministratorDetails) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, admin := range *input {
		result := make(map[string]interface{})
		if admin.FirstName != nil {
			result["first_name"] = admin.FirstName
		}
		if admin.LastName != nil {
			result["last_name"] = admin.LastName
		}
		if admin.EmailAddress != nil {
			result["email_address"] = admin.EmailAddress
		}
		if admin.Phone != nil {
			result["phone"] = admin.Phone
		}
		results = append(results, result)
	}

	return results
}

func resourceArmKeyVaultCertificateIssuerDelete(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	keyVaultBaseUri, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Certificate Issuer %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	// we verify it exists
	resp, err := client.GetCertificateIssuer(ctx, keyVaultBaseUri, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist", name, keyVaultBaseUri)
		}
		return fmt.Errorf("Error making Read request on Azure KeyVault Certificate Issuer %s: %+v", name, err)
	}

	resp, err = client.DeleteCertificateIssuer(ctx, keyVaultBaseUri, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("Error deleting Certificate Issuer %q (KeyVault URI %q) from Key Vault: %+v", name, keyVaultBaseUri, err)
	}

	return nil
}
