package keyvault

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultCertificateContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultCertificateContactCreateOrUpdate,
		Update: resourceArmKeyVaultCertificateContactCreateOrUpdate,
		Read:   resourceArmKeyVaultCertificateContactRead,
		Delete: resourceArmKeyVaultCertificateContactDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.KeyVaultCertificateContactID(id)
			return err
		}),

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
				ValidateFunc: validate.KeyVaultID,
			},

			"contact": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
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

func resourceArmKeyVaultCertificateContactCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	keyVaultId := d.Get("key_vault_id").(string)

	keyVaultBaseUri, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("failed to look up Certificate Contacts vault url from id %q: %+v", keyVaultId, err)
	}

	if d.IsNewResource() {
		existing, err := client.GetCertificateContacts(ctx, keyVaultBaseUri)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failed to check for presence of existing Certificate Contacts (Key Vault %q): %s", keyVaultBaseUri, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_certificate_contact", *existing.ID)
		}
	}

	contacts := keyvault.Contacts{
		ContactList: expandKeyVaultCertificateContactList(d.Get("contact").(*schema.Set).List()),
	}

	_, err = client.SetCertificateContacts(ctx, keyVaultBaseUri, contacts)
	if err != nil {
		return fmt.Errorf("failed to set Certificate Contacts (Key Vault %q): %s", keyVaultId, err)
	}

	resp, err := client.GetCertificateContacts(ctx, keyVaultBaseUri)
	if err != nil {
		return err
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("failure reading Certificate Contacts ID for Key Vault %q", keyVaultBaseUri)
	}
	d.SetId(*resp.ID)

	return resourceArmKeyVaultCertificateContactRead(d, meta)
}

func resourceArmKeyVaultCertificateContactRead(d *schema.ResourceData, meta interface{}) error {
	keyVaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KeyVaultCertificateContactID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, keyVaultClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID of Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultId == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	ok, err := azure.KeyVaultExists(ctx, keyVaultClient, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Certificate Contacts in Vault at url %q exists: %v", *keyVaultId, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Certificate contacts of Key Vault %q was not found in Key Vault at URI %q - removing from state", *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	resp, err := client.GetCertificateContacts(ctx, id.KeyVaultBaseUrl)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] KeyVault Certificate Contacts (KeyVault URI %q) does not exist - removing from state", id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("failed to make Read request on Azure KeyVault Certificate Contacts: %+v", err)
	}

	d.Set("key_vault_id", keyVaultId)
	if err := d.Set("contact", flattenKeyVaultCertificateContactList(resp.ContactList)); err != nil {
		return fmt.Errorf("setting `contact` for KeyVault Certificate Contacts: %+v", err)
	}

	return nil
}

func resourceArmKeyVaultCertificateContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagementClient
	keyVaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KeyVaultCertificateContactID(d.Id())
	if err != nil {
		return err
	}

	// we verify it exists
	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, keyVaultClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}

	ok, err := azure.KeyVaultExists(ctx, keyVaultClient, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q at url %q exists: %v", *keyVaultId, id.KeyVaultBaseUrl, err)
	}

	if !ok {
		log.Printf("[DEBUG] Certificate Contacts (Key Vault %q) was not found in Key Vault at URI %q - removing from state", *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	if _, err := client.DeleteCertificateContacts(ctx, id.KeyVaultBaseUrl); err != nil {
		return fmt.Errorf("deleting key vault %q Certificate Contacts: %+v", keyVaultId, err)
	}

	return err
}

func expandKeyVaultCertificateContactList(input []interface{}) *[]keyvault.Contact {
	results := make([]keyvault.Contact, 0)
	if len(input) == 0 || input[0] == nil {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, keyvault.Contact{
			Name:         utils.String(v["name"].(string)),
			EmailAddress: utils.String(v["email"].(string)),
			Phone:        utils.String(v["phone"].(string)),
		})
	}

	return &results
}

func flattenKeyVaultCertificateContactList(input *[]keyvault.Contact) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, contact := range *input {
		emailAddress := ""
		if contact.EmailAddress != nil {
			emailAddress = *contact.EmailAddress
		}

		name := ""
		if contact.Name != nil {
			name = *contact.Name
		}

		phone := ""
		if contact.Phone != nil {
			phone = *contact.Phone
		}

		results = append(results, map[string]interface{}{
			"email": emailAddress,
			"name":  name,
			"phone": phone,
		})
	}

	return results
}
