package mssql

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	keyVaultParser "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMsSqlTransparentDataEncryption() *schema.Resource {
	return &schema.Resource{
		Create: resourceMsSqlTransparentDataEncryptionCreateUpdate,
		Read:   resourceMsSqlTransparentDataEncryptionRead,
		Update: resourceMsSqlTransparentDataEncryptionCreateUpdate,
		Delete: resourceMsSqlTransparentDataEncryptionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.EncryptionProtectorID(id)

			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),
			"key_vault_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.NestedItemId,
			},
		},
	}
}

func resourceMsSqlTransparentDataEncryptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	encryptionProtectorClient := meta.(*clients.Client).MSSQL.EncryptionProtectorClient
	serverKeysClient := meta.(*clients.Client).MSSQL.ServerKeysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MSSQL Encrypted Protector creation.")

	// encryptedProtectorName := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	// Normally we would check if this is a new resource, but the way encryption protector works, it always overwrites
	// whatever is there anyways. Compounding the issue is that SQL Server creates an instance of encryption protector
	// which causes the isNewResource check to fail because we are trying to create the encryption as a separate resource
	// and encryption protector is already present. The reason we create encryption protector as a separate resource is
	// because after the SQL server is created, we need to grant it permissions to AKV, so encryption protector can use those
	// keys are part of setting up TDE

	var serverKey sql.ServerKey
	var serverKeyName string
	var encryptionProtectorProperties sql.EncryptionProtectorProperties

	if v, ok := d.GetOk("key_vault_uri"); ok {
		keyUri := v.(string)

		// Set the SQL Server Key properties
		serverKeyProperties := sql.ServerKeyProperties{
			ServerKeyType: sql.AzureKeyVault,
			URI:           &keyUri,
		}
		serverKey.ServerKeyProperties = &serverKeyProperties

		// Set the encryption protector properties
		keyDetails, err := keyVaultParser.ParseNestedItemID(keyUri)

		if err != nil {
			return fmt.Errorf("Unable to parse key uri: %q: %+v", keyUri, err)
		}

		// Make sure it's a key, if not, throw an error
		if keyDetails.NestedItemType == "keys" {
			keyName := keyDetails.Name
			keyVersion := keyDetails.Version

			// Extract the vault name from the keyvault base url
			idURL, err := url.ParseRequestURI(keyDetails.KeyVaultBaseUrl)

			if err != nil {
				return fmt.Errorf("Unable to parse key vault hostname: %s", keyDetails.KeyVaultBaseUrl)
			}

			hostParts := strings.Split(idURL.Host, ".")
			vaultName := hostParts[0]

			// Create the key path as for the Encryption Protector. Format is: {vaultname}_{key}_{key_version}
			serverKeyName = fmt.Sprintf("%s_%s_%s", vaultName, keyName, keyVersion)

			// Set the encryption protector properties
			encryptionProtectorProperties = sql.EncryptionProtectorProperties{
				ServerKeyType: sql.AzureKeyVault,
				ServerKeyName: &serverKeyName,
			}
		} else {
			return fmt.Errorf("Key vault uri must be a reference to a key, but got: %s", keyDetails.NestedItemType)
		}
	} else {
		serverKeyName = ""

		// Service managed doesn't require a key name
		encryptionProtectorProperties = sql.EncryptionProtectorProperties{
			ServerKeyType: sql.ServiceManaged,
			ServerKeyName: &serverKeyName,
		}
	}

	// Only create  aserver key
	if serverKey.ServerKeyProperties != nil {
		// Create a key on the server
		futureServers, err := serverKeysClient.CreateOrUpdate(ctx, resGroup, serverName, serverKeyName, serverKey)
		if err != nil {
			return err
		}

		if err = futureServers.WaitForCompletionRef(ctx, serverKeysClient.Client); err != nil {
			return err
		}
	}

	encryptionProtectorObject := sql.EncryptionProtector{
		EncryptionProtectorProperties: &encryptionProtectorProperties,
	}

	futureEncryptionProtector, err := encryptionProtectorClient.CreateOrUpdate(ctx, resGroup, serverName, encryptionProtectorObject)
	if err != nil {
		return err
	}

	if err = futureEncryptionProtector.WaitForCompletionRef(ctx, encryptionProtectorClient.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for Encryption Protector on server %q (Resource Group %q): %+v", serverName, resGroup, err)
	}

	resp, err := encryptionProtectorClient.Get(ctx, resGroup, serverName)
	if err != nil {
		return fmt.Errorf("issuing get request for Encryption protector on server %q (Resource Group %q): %+v", serverName, resGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceMsSqlTransparentDataEncryptionRead(d, meta)
}

func resourceMsSqlTransparentDataEncryptionRead(d *schema.ResourceData, meta interface{}) error {
	encryptionProtectorClient := meta.(*clients.Client).MSSQL.EncryptionProtectorClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EncryptionProtectorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := encryptionProtectorClient.Get(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request for Encryption Protector on server %s: %v", id.ServerName, err)
	}

	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroup)

	log.Printf("[INFO] Encryption protector key type is %s", resp.EncryptionProtectorProperties.ServerKeyType)

	// Only set the key type if it's an AKV key. For service managed, we can omit the setting the key_vault_uri
	if resp.EncryptionProtectorProperties.ServerKeyType == sql.AzureKeyVault {
		log.Printf("[INFO] Setting Key Vault URI to %s", *resp.EncryptionProtectorProperties.URI)

		if err := d.Set("key_vault_uri", resp.EncryptionProtectorProperties.URI); err != nil {
			return fmt.Errorf("setting key_vault_uri`: %+v", err)
		}
	}

	return nil
}

func resourceMsSqlTransparentDataEncryptionDelete(d *schema.ResourceData, meta interface{}) error {
	// Note that encryption protector cannot be deleted. It can only be updated between AzureKeyVault
	// and SystemManaged.
	return nil
}
