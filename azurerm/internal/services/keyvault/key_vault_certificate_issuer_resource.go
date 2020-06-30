package keyvault

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
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
				}, false),
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
		return fmt.Errorf("failed to look up Certificate Issuer %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	existing, err := client.GetCertificateIssuer(ctx, keyVaultBaseUri, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("failed to check for presence of existing Certificate Issuer %q (Key Vault %q): %s", name, keyVaultBaseUri, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_key_vault_certificate_issuer", *existing.ID)
	}

	parameter := keyvault.CertificateIssuerSetParameters{}
	parameter.Provider = utils.String(d.Get("provider_name").(string))
	if adminsRaw, ok := d.GetOk("admins"); ok {
		parameter.OrganizationDetails = &keyvault.OrganizationDetails{
			ID:           utils.String(d.Get("org_id").(string)),
			AdminDetails: expandKeyVaultCertificateIssuerOrganizationDetailsAdminDetails(adminsRaw.([]interface{})),
		}
	}
	accountId, gotAccountId := d.GetOk("account_id")
	password, gotPassword := d.GetOk("password")

	if gotAccountId && gotPassword {
		parameter.Credentials = &keyvault.IssuerCredentials{
			AccountID: utils.String(accountId.(string)),
			Password:  utils.String(password.(string)),
		}
	}
	resp, err := client.SetCertificateIssuer(ctx, keyVaultBaseUri, name, parameter)
	if err != nil {
		return fmt.Errorf("failed to set Certificate Issuer %q (Key Vault %q): %s", name, keyVaultId, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("failure reading Key Vault Certificate Issuer ID for %q", name)
	}
	d.SetId(*resp.ID)

	return resourceArmKeyVaultCertificateIssuerRead(d, meta)
}

func expandKeyVaultCertificateIssuerOrganizationDetailsAdminDetails(vs []interface{}) *[]keyvault.AdministratorDetails {
	results := make([]keyvault.AdministratorDetails, 0, len(vs))

	for _, v := range vs {
		administratorDetails := keyvault.AdministratorDetails{}
		args := v.(map[string]interface{})
		if firstName, ok := args["first_name"]; ok {
			administratorDetails.FirstName = utils.String(firstName.(string))
		}
		if lastName, ok := args["last_name"]; ok {
			administratorDetails.LastName = utils.String(lastName.(string))
		}
		if emailAddress, ok := args["email_address"]; ok {
			administratorDetails.EmailAddress = utils.String(emailAddress.(string))
		}
		if phone, ok := args["phone"]; ok {
			administratorDetails.Phone = utils.String(phone.(string))
		}
		results = append(results, administratorDetails)
	}

	return &results
}

func resourceArmKeyVaultCertificateIssuerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseKeyVaultCertificateIssuerID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.GetCertificateIssuer(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("failed to check for presence of existing Certificate Issuer %q (Key Vault %q): %s", id.Name, id.KeyVaultBaseUrl, err)
		}
		return fmt.Errorf("KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist", id.Name, id.KeyVaultBaseUrl)
	}

	parameter := keyvault.CertificateIssuerSetParameters{
		Provider:            utils.String(d.Get("provider_name").(string)),
		OrganizationDetails: &keyvault.OrganizationDetails{},
		Credentials:         &keyvault.IssuerCredentials{},
	}

	if orgID, ok := d.GetOk("org_id"); ok {
		parameter.OrganizationDetails.ID = utils.String(orgID.(string))
	}

	if adminDetails, ok := d.GetOk("admins"); ok {
		parameter.OrganizationDetails.AdminDetails = expandKeyVaultCertificateIssuerOrganizationDetailsAdminDetails(adminDetails.([]interface{}))
	}

	if accountID, ok := d.GetOk("account_id"); ok {
		parameter.Credentials.AccountID = utils.String(accountID.(string))
	}

	if password, ok := d.GetOk("password"); ok {
		parameter.Credentials.Password = utils.String(password.(string))
	}
	resp, err := client.SetCertificateIssuer(ctx, id.KeyVaultBaseUrl, id.Name, parameter)
	if err != nil {
		return fmt.Errorf("failed to set Certificate Issuer %q (Key Vault %q): %s", id.Name, id.KeyVaultBaseUrl, err)
	}
	if resp.ID != nil || *resp.ID == "" {
		d.SetId(*resp.ID)
	}

	return resourceArmKeyVaultCertificateIssuerRead(d, meta)
}

func resourceArmKeyVaultCertificateIssuerRead(d *schema.ResourceData, meta interface{}) error {
	keyVaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseKeyVaultCertificateIssuerID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, keyVaultClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultId == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	ok, err := azure.KeyVaultExists(ctx, keyVaultClient, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error checking if key vault %q for Certificate %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Certificate %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	resp, err := client.GetCertificateIssuer(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("failed to make Read request on Azure KeyVault Certificate Issuer %s: %+v", id.Name, err)
	}

	if resp.Provider != nil {
		d.Set("provider_name", resp.Provider)
	}
	if resp.OrganizationDetails != nil {
		if resp.OrganizationDetails.ID != nil {
			d.Set("org_id", resp.OrganizationDetails.ID)
		}
		if resp.OrganizationDetails.AdminDetails != nil {
			adminDetails, err := flattenKeyVaultCertificateIssuerAdmins(resp.OrganizationDetails.AdminDetails)
			if err != nil {
				return fmt.Errorf("failed to flatten Azure KeyVault Certificate Issuer Admin Details: %v", err)
			}
			d.Set("admins", adminDetails)
		}
	}
	if resp.Credentials != nil {
		if resp.Credentials.AccountID != nil {
			d.Set("account_id", resp.Credentials.AccountID)
		}
		if resp.Credentials.Password != nil {
			d.Set("password", resp.Credentials.Password)
		}
	}

	return nil
}

func flattenKeyVaultCertificateIssuerAdmins(input *[]keyvault.AdministratorDetails) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
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
		} else {
			return nil, fmt.Errorf("email_address is required, but %q was nil in API response with the full struct being %#v", "EmailAddress", admin)
		}
		if admin.Phone != nil {
			result["phone"] = admin.Phone
		}
		results = append(results, result)
	}

	return results, nil
}

func resourceArmKeyVaultCertificateIssuerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseKeyVaultCertificateIssuerID(d.Id())
	if err != nil {
		return err
	}

	// we verify it exists
	resp, err := client.GetCertificateIssuer(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist", id.Name, id.KeyVaultBaseUrl)
		}
		return fmt.Errorf("failed to make Read request on Azure KeyVault Certificate Issuer %s: %+v", id.Name, err)
	}

	resp, err = client.DeleteCertificateIssuer(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("failed to delete Certificate Issuer %q (KeyVault URI %q) from Key Vault: %+v", id.Name, id.KeyVaultBaseUrl, err)
	}

	return nil
}
