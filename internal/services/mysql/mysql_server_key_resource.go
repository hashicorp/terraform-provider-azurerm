// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2020-01-01/serverkeys"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
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
			_, err := serverkeys.ParseKeyID(id)
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
				ValidateFunc: servers.ValidateServerID,
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
	keyVaultID, err := commonids.ParseKeyVaultID(*keyVaultIDRaw)
	if err != nil {
		return nil, err
	}
	return utils.String(fmt.Sprintf("%s_%s_%s", keyVaultID.VaultName, keyVaultKeyID.Name, keyVaultKeyID.Version)), nil
}

func resourceMySQLServerKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keysClient := meta.(*clients.Client).MySQL.ServerKeysClient.ServerKeys
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverID, err := serverkeys.ParseServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(serverID.ServerName, mySQLServerResourceName)
	defer locks.UnlockByName(serverID.ServerName, mySQLServerResourceName)

	if d.IsNewResource() {
		// This resource is a singleton, but its name can be anything.
		// If you create a new key with different name with the old key, the service will not give you any warning but directly replace the old key with the new key.
		// Therefore sometimes you cannot get the old key using the GET API since you may not know the name of the old key
		resp, err := keysClient.List(ctx, *serverID)
		if err != nil {
			return fmt.Errorf("listing existing MySQL Server Keys in %s: %s", serverID, err)
		}
		if resp.Model == nil {
			return fmt.Errorf("model was nil for %s", serverID)
		}
		keys := *resp.Model
		if len(keys) > 0 {
			if len(keys) > 1 {
				return fmt.Errorf("expecting at most one MySQL Server Key, but got %q", len(keys))
			}
			if keys[0].Id == nil || *keys[0].Id == "" {
				return fmt.Errorf("missing ID for existing MySQL Server Key")
			}

			id, err := serverkeys.ParseKeyID(*keys[0].Id)
			if err != nil {
				return err
			}

			return tf.ImportAsExistsError("azurerm_mysql_server_key", id.ID())
		}
	}

	keyVaultKeyURI := d.Get("key_vault_key_id").(string)
	name, err := getMySQLServerKeyName(ctx, keyVaultsClient, resourcesClient, keyVaultKeyURI)
	if err != nil {
		return fmt.Errorf("cannot compose name for MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroupName, serverID.ServerName, err)
	}

	id := serverkeys.NewKeyID(serverID.SubscriptionId, serverID.ResourceGroupName, serverID.ServerName, *name)
	param := serverkeys.ServerKey{
		Properties: &serverkeys.ServerKeyProperties{
			ServerKeyType: serverkeys.ServerKeyTypeAzureKeyVault,
			Uri:           &keyVaultKeyURI,
		},
	}

	if err = keysClient.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMySQLServerKeyRead(d, meta)
}

func resourceMySQLServerKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keysClient := meta.(*clients.Client).MySQL.ServerKeysClient.ServerKeys
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serverkeys.ParseKeyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := keysClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	serverId := servers.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
	d.Set("server_id", serverId.ID())
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("key_vault_key_id", props.Uri)
		}
	}

	return nil
}

func resourceMySQLServerKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServerKeysClient.ServerKeys
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serverkeys.ParseKeyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ServerName, mySQLServerResourceName)
	defer locks.UnlockByName(id.ServerName, mySQLServerResourceName)

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
