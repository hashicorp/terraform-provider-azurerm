package mysql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMySQLServerKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMySQLServerKeyCreateUpdate,
		Read:   resourceArmMySQLServerKeyRead,
		Update: resourceArmMySQLServerKeyCreateUpdate,
		Delete: resourceArmMySQLServerKeyDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.MySQLServerKeyID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MySQLServerID,
			},

			"key_vault_key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateKeyVaultChildId,
			},
		},
	}
}

func getMySQLServerKeyName(ctx context.Context, vaultsClient *keyvault.VaultsClient, keyVaultKeyURI string) (*string, error) {
	keyVaultKeyID, err := azure.ParseKeyVaultChildID(keyVaultKeyURI)
	if err != nil {
		return nil, err
	}
	keyVaultIDRaw, err := azure.GetKeyVaultIDFromBaseUrl(ctx, vaultsClient, keyVaultKeyID.KeyVaultBaseUrl)
	if err != nil {
		return nil, err
	}
	keyVaultID, err := keyVaultParse.KeyVaultID(*keyVaultIDRaw)
	if err != nil {
		return nil, err
	}
	return utils.String(fmt.Sprintf("%s_%s_%s", keyVaultID.Name, keyVaultKeyID.Name, keyVaultKeyID.Version)), nil
}

func resourceArmMySQLServerKeyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	keysClient := meta.(*clients.Client).MySQL.ServerKeysClient
	vaultsClient := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverID, err := parse.MySQLServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}
	keyVaultKeyURI := d.Get("key_vault_key_id").(string)
	name, err := getMySQLServerKeyName(ctx, vaultsClient, keyVaultKeyURI)
	if err != nil {
		return fmt.Errorf("cannot compose name for MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
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
		if len(keys) > 1 {
			return fmt.Errorf("expecting at most one MySQL Server Key, but got %q", len(keys))
		}
		if len(keys) == 1 && keys[0].ID != nil && *keys[0].ID != "" {
			return tf.ImportAsExistsError("azurerm_mysql_server_key", *keys[0].ID)
		}
	}

	param := mysql.ServerKey{
		ServerKeyProperties: &mysql.ServerKeyProperties{
			ServerKeyType: utils.String("AzureKeyVault"),
			URI:           &keyVaultKeyURI,
		},
	}

	future, err := keysClient.CreateOrUpdate(ctx, serverID.Name, *name, param, serverID.ResourceGroup)
	if err != nil {
		return fmt.Errorf("creating/updating MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
	}
	if err := future.WaitForCompletionRef(ctx, keysClient.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
	}

	resp, err := keysClient.Get(ctx, serverID.ResourceGroup, serverID.Name, *name)
	if err != nil {
		return fmt.Errorf("retrieving MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
	}

	d.SetId(*resp.ID)

	return resourceArmMySQLServerKeyRead(d, meta)
}

func resourceArmMySQLServerKeyRead(d *schema.ResourceData, meta interface{}) error {
	serversClient := meta.(*clients.Client).MySQL.ServersClient
	keysClient := meta.(*clients.Client).MySQL.ServerKeysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MySQLServerKeyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := keysClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] MySQL Server Key %q was not found (Resource Group %q / Server %q)", id.Name, id.ResourceGroup, id.ServerName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving MySQL Server Key %q (Resource Group %q / Server %q): %+v", id.Name, id.ResourceGroup, id.ServerName, err)
	}

	respServer, err := serversClient.Get(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		return fmt.Errorf("cannot get MySQL Server ID: %+v", err)
	}

	d.Set("server_id", respServer.ID)
	if props := resp.ServerKeyProperties; props != nil {
		d.Set("key_vault_key_id", props.URI)
	}

	return nil
}

func resourceArmMySQLServerKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServerKeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MySQLServerKeyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ServerName, mySQLServerResourceName)
	defer locks.UnlockByName(id.ServerName, mySQLServerResourceName)

	future, err := client.Delete(ctx, id.ServerName, id.Name, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("deleting MySQL Server Key %q (Resource Group %q / Server %q): %+v", id.Name, id.ResourceGroup, id.ServerName, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of MySQL Server Key %q (Resource Group %q / Server %q): %+v", id.Name, id.ResourceGroup, id.ServerName, err)
	}

	return nil
}
