package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2020-01-01/postgresql"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	resourcesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePostgreSQLServerKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgreSQLServerKeyCreateUpdate,
		Read:   resourcePostgreSQLServerKeyRead,
		Update: resourcePostgreSQLServerKeyCreateUpdate,
		Delete: resourcePostgreSQLServerKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ServerKeyID(id)
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

func getPostgreSQLServerKeyName(ctx context.Context, keyVaultsClient *client.Client, resourcesClient *resourcesClient.Client, keyVaultKeyURI string) (*string, error) {
	keyVaultKeyID, err := keyVaultParse.ParseNestedItemID(keyVaultKeyURI)
	if err != nil {
		return nil, err
	}
	keyVaultIDRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyVaultKeyID.KeyVaultBaseUrl)
	if err != nil {
		return nil, err
	}
	// function azure.GetKeyVaultIDFromBaseUrl returns nil with nil error when it does not find the keyvault by the keyvault URL
	if keyVaultIDRaw == nil {
		return nil, fmt.Errorf("cannot get the keyvault ID from keyvault URL %q", keyVaultKeyID.KeyVaultBaseUrl)
	}
	keyVaultID, err := keyVaultParse.VaultID(*keyVaultIDRaw)
	if err != nil {
		return nil, err
	}
	return utils.String(fmt.Sprintf("%s_%s_%s", keyVaultID.Name, keyVaultKeyID.Name, keyVaultKeyID.Version)), nil
}

func resourcePostgreSQLServerKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keysClient := meta.(*clients.Client).Postgres.ServerKeysClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverId, err := parse.ServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}
	keyVaultKeyURI := d.Get("key_vault_key_id").(string)
	name, err := getPostgreSQLServerKeyName(ctx, keyVaultsClient, resourcesClient, keyVaultKeyURI)
	if err != nil {
		return fmt.Errorf("cannot compose name for PostgreSQL Server Key (Resource Group %q / Server %q): %+v", serverId.ResourceGroup, serverId.Name, err)
	}

	locks.ByName(serverId.Name, postgreSQLServerResourceName)
	defer locks.UnlockByName(serverId.Name, postgreSQLServerResourceName)

	if d.IsNewResource() {
		// This resource is a singleton, but its name can be anything.
		// If you create a new key with different name with the old key, the service will not give you any warning but directly replace the old key with the new key.
		// Therefore sometimes you cannot get the old key using the GET API since you may not know the name of the old key
		resp, err := keysClient.List(ctx, serverId.ResourceGroup, serverId.Name)
		if err != nil {
			return fmt.Errorf("listing existing PostgreSQL Server Keys in Resource Group %q / Server %q: %+v", serverId.ResourceGroup, serverId.Name, err)
		}
		keys := resp.Values()
		if len(keys) >= 1 {
			if rawId := keys[0].ID; rawId != nil && *rawId != "" {
				id, err := parse.ServerKeyID(*rawId)
				if err != nil {
					return fmt.Errorf("parsing existing Server Key ID %q: %+v", *rawId, err)
				}

				return tf.ImportAsExistsError("azurerm_postgresql_server_key", id.ID())
			}
		}
	}

	param := postgresql.ServerKey{
		ServerKeyProperties: &postgresql.ServerKeyProperties{
			ServerKeyType: utils.String("AzureKeyVault"),
			URI:           utils.String(d.Get("key_vault_key_id").(string)),
		},
	}

	id := parse.NewServerKeyID(serverId.SubscriptionId, serverId.ResourceGroup, serverId.Name, *name)
	future, err := keysClient.CreateOrUpdate(ctx, id.ServerName, id.KeyName, param, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, keysClient.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePostgreSQLServerKeyRead(d, meta)
}

func resourcePostgreSQLServerKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keysClient := meta.(*clients.Client).Postgres.ServerKeysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerKeyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := keysClient.Get(ctx, id.ResourceGroup, id.ServerName, id.KeyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("server_id", parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName).ID())
	if props := resp.ServerKeyProperties; props != nil {
		d.Set("key_vault_key_id", props.URI)
	}

	return nil
}

func resourcePostgreSQLServerKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServerKeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerKeyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ServerName, postgreSQLServerResourceName)
	defer locks.UnlockByName(id.ServerName, postgreSQLServerResourceName)

	future, err := client.Delete(ctx, id.ServerName, id.KeyName, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
