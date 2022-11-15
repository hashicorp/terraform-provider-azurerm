package mysql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	resourcesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMySQLServerKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySQLServerKeyCreateUpdate,
		Read:   resourceMySQLServerKeyRead,
		Update: resourceMySQLServerKeyCreateUpdate,
		Delete: resourceMySQLServerKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.KeyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerID,
			},

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.NestedItemId,
			},
		},
	}
}

func getMySQLServerKeyName(ctx context.Context, keyVaultsClient *client.Client, resourcesClient *resourcesClient.Client, keyVaultKeyURI string) (*string, error) {
	keyVaultKeyID, err := keyVaultParse.ParseNestedItemID(keyVaultKeyURI)
	if err != nil {
		return nil, err
	}
	keyVaultIDRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyVaultKeyID.KeyVaultBaseUrl)
	if err != nil {
		return nil, err
	}
	keyVaultID, err := keyVaultParse.VaultID(*keyVaultIDRaw)
	if err != nil {
		return nil, err
	}
	return utils.String(fmt.Sprintf("%s_%s_%s", keyVaultID.Name, keyVaultKeyID.Name, keyVaultKeyID.Version)), nil
}

func resourceMySQLServerKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keysClient := meta.(*clients.Client).MySQL.ServerKeysClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverID, err := parse.ServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(serverID.Name, mySQLServerResourceName)
	defer locks.UnlockByName(serverID.Name, mySQLServerResourceName)

	if d.IsNewResource() {
		// This resource is a singleton, but its name can be anything.
		// If you create a new key with different name with the old key, the service will not give you any warning but directly replace the old key with the new key.
		// Therefore sometimes you cannot get the old key using the GET API since you may not know the name of the old key
		resp, err := keysClient.List(ctx, serverID.ResourceGroup, serverID.Name)
		if err != nil {
			return fmt.Errorf("listing existing MySQL Server Keys in Resource Group %q / Server %q: %+v", serverID.ResourceGroup, serverID.Name, err)
		}
		keys := resp.Values()
		if len(keys) > 0 {
			if len(keys) > 1 {
				return fmt.Errorf("expecting at most one MySQL Server Key, but got %q", len(keys))
			}
			if keys[0].ID == nil || *keys[0].ID == "" {
				return fmt.Errorf("missing ID for existing MySQL Server Key")
			}

			id, err := parse.KeyID(*keys[0].ID)
			if err != nil {
				return err
			}

			return tf.ImportAsExistsError("azurerm_mysql_server_key", id.ID())
		}
	}

	keyVaultKeyURI := d.Get("key_vault_key_id").(string)
	name, err := getMySQLServerKeyName(ctx, keyVaultsClient, resourcesClient, keyVaultKeyURI)
	if err != nil {
		return fmt.Errorf("cannot compose name for MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
	}

	id := parse.NewKeyID(serverID.SubscriptionId, serverID.ResourceGroup, serverID.Name, *name)
	param := mysql.ServerKey{
		ServerKeyProperties: &mysql.ServerKeyProperties{
			ServerKeyType: utils.String("AzureKeyVault"),
			URI:           &keyVaultKeyURI,
		},
	}

	future, err := keysClient.CreateOrUpdate(ctx, id.ServerName, id.Name, param, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, keysClient.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMySQLServerKeyRead(d, meta)
}

func resourceMySQLServerKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keysClient := meta.(*clients.Client).MySQL.ServerKeysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KeyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := keysClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)
	d.Set("server_id", serverId.ID())
	if props := resp.ServerKeyProperties; props != nil {
		d.Set("key_vault_key_id", props.URI)
	}

	return nil
}

func resourceMySQLServerKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServerKeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KeyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ServerName, mySQLServerResourceName)
	defer locks.UnlockByName(id.ServerName, mySQLServerResourceName)

	future, err := client.Delete(ctx, id.ServerName, id.Name, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
