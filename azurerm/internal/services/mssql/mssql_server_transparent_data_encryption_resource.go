package mssql

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	keyVaultParser "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	mssqlValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMsSqlTransparentDataEncryption() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlTransparentDataEncryptionCreateUpdate,
		Read:   resourceMsSqlTransparentDataEncryptionRead,
		Update: resourceMsSqlTransparentDataEncryptionCreateUpdate,
		Delete: resourceMsSqlTransparentDataEncryptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EncryptionProtectorID(id)

			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: mssqlValidate.ServerID,
			},
			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.NestedItemId,
			},
		},
	}
}

func resourceMsSqlTransparentDataEncryptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	encryptionProtectorClient := meta.(*clients.Client).MSSQL.EncryptionProtectorClient
	serverKeysClient := meta.(*clients.Client).MSSQL.ServerKeysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverId, err := parse.ServerID(d.Get("server_id").(string))

	if err != nil {
		return err
	}

	// Normally we would check if this is a new resource, but the way encryption protector works, it always overwrites
	// whatever is there anyways. Compounding the issue is that SQL Server creates an instance of encryption protector
	// which causes the isNewResource check to fail because we are trying to create the encryption as a separate resource
	// and encryption protector is already present. The reason we create encryption protector as a separate resource is
	// because after the SQL server is created, we need to grant it permissions to AKV, so encryption protector can use those
	// keys are part of setting up TDE

	var serverKey sql.ServerKey

	// Default values for Service Managed keys. Will update to AKV values if key_vault_key_id references a key.
	serverKeyName := ""
	serverKeyType := sql.ServiceManaged

	keyVaultKeyId := strings.TrimSpace(d.Get("key_vault_key_id").(string))

	// If it has content, then we assume it's a key vault key id
	if keyVaultKeyId != "" {
		// Update the server key type to AKV
		serverKeyType = sql.AzureKeyVault

		// Set the SQL Server Key properties
		serverKeyProperties := sql.ServerKeyProperties{
			ServerKeyType: serverKeyType,
			URI:           &keyVaultKeyId,
		}
		serverKey.ServerKeyProperties = &serverKeyProperties

		// Set the encryption protector properties
		keyId, err := keyVaultParser.ParseNestedItemID(keyVaultKeyId)

		if err != nil {
			return fmt.Errorf("Unable to parse key: %q: %+v", keyVaultKeyId, err)
		}

		// Make sure it's a key, if not, throw an error
		if keyId.NestedItemType == "keys" {
			keyName := keyId.Name
			keyVersion := keyId.Version

			// Extract the vault name from the keyvault base url
			idURL, err := url.ParseRequestURI(keyId.KeyVaultBaseUrl)

			if err != nil {
				return fmt.Errorf("Unable to parse key vault hostname: %s", keyId.KeyVaultBaseUrl)
			}

			hostParts := strings.Split(idURL.Host, ".")
			vaultName := hostParts[0]

			// Create the key path for the Encryption Protector. Format is: {vaultname}_{key}_{key_version}
			serverKeyName = fmt.Sprintf("%s_%s_%s", vaultName, keyName, keyVersion)
		} else {
			return fmt.Errorf("Key vault key id must be a reference to a key, but got: %s", keyId.NestedItemType)
		}
	}

	// Service managed doesn't require a key name
	encryptionProtectorProperties := sql.EncryptionProtectorProperties{
		ServerKeyType: serverKeyType,
		ServerKeyName: &serverKeyName,
	}

	// Only create a server key if the properties have been set
	if serverKey.ServerKeyProperties != nil {
		// Create a key on the server
		futureServers, err := serverKeysClient.CreateOrUpdate(ctx, serverId.ResourceGroup, serverId.Name, serverKeyName, serverKey)
		if err != nil {
			return fmt.Errorf("creating/updating server key for %s: %+v", serverId, err)
		}

		if err = futureServers.WaitForCompletionRef(ctx, serverKeysClient.Client); err != nil {
			return fmt.Errorf("waiting on update of %s: %+v", serverId, err)
		}
	}

	encryptionProtectorObject := sql.EncryptionProtector{
		EncryptionProtectorProperties: &encryptionProtectorProperties,
	}

	// Encryption protector always uses "current" for the name
	id := parse.NewEncryptionProtectorID(serverId.SubscriptionId, serverId.ResourceGroup, serverId.Name, "current")

	futureEncryptionProtector, err := encryptionProtectorClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, encryptionProtectorObject)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = futureEncryptionProtector.WaitForCompletionRef(ctx, encryptionProtectorClient.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMsSqlTransparentDataEncryptionRead(d, meta)
}

func resourceMsSqlTransparentDataEncryptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("making Read request for %s: %v", id, err)
	}

	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)
	d.Set("server_id", serverId.ID())

	log.Printf("[INFO] Encryption protector key type is %s", resp.EncryptionProtectorProperties.ServerKeyType)

	keyVaultKeyId := ""

	// Only set the key type if it's an AKV key. For service managed, we can omit the setting the key_vault_key_id
	if resp.EncryptionProtectorProperties != nil && resp.EncryptionProtectorProperties.ServerKeyType == sql.AzureKeyVault {
		log.Printf("[INFO] Setting Key Vault URI to %s", *resp.EncryptionProtectorProperties.URI)

		keyVaultKeyId = *resp.EncryptionProtectorProperties.URI
	}

	if err := d.Set("key_vault_key_id", keyVaultKeyId); err != nil {
		return fmt.Errorf("setting key_vault_key_id`: %+v", err)
	}

	return nil
}

func resourceMsSqlTransparentDataEncryptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// Note that encryption protector cannot be deleted. It can only be updated between AzureKeyVault
	// and SystemManaged. For safety, when this resource is deleted, we're resetting the key type
	// to service managed to prevent accidental lockout if someone were to delete the keys from key vault

	encryptionProtectorClient := meta.(*clients.Client).MSSQL.EncryptionProtectorClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EncryptionProtectorID(d.Id())

	if err != nil {
		return err
	}

	serverKeyName := ""

	// Service managed doesn't require a key name
	encryptionProtector := sql.EncryptionProtector{
		EncryptionProtectorProperties: &sql.EncryptionProtectorProperties{
			ServerKeyType: sql.ServiceManaged,
			ServerKeyName: &serverKeyName,
		},
	}

	futureEncryptionProtector, err := encryptionProtectorClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, encryptionProtector)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = futureEncryptionProtector.WaitForCompletionRef(ctx, encryptionProtectorClient.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for %s: %+v", id, err)
	}

	return nil
}
