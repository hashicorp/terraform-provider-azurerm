// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2020-01-01/serverkeys"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
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
				ValidateFunc: serverkeys.ValidateServerID,
			},

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.NestedItemId,
			},
		},
	}
}

func getPostgreSQLServerKeyName(ctx context.Context, keyVaultsClient *client.Client, subscriptionId string, keyVaultKeyURI string) (*string, error) {
	keyVaultKeyID, err := keyVaultParse.ParseNestedItemID(keyVaultKeyURI)
	if err != nil {
		return nil, err
	}
	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIDRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyVaultKeyID.KeyVaultBaseUrl)
	if err != nil {
		return nil, err
	}
	// function azure.GetKeyVaultIDFromBaseUrl returns nil with nil error when it does not find the keyvault by the keyvault URL
	if keyVaultIDRaw == nil {
		return nil, fmt.Errorf("cannot get the keyvault ID from keyvault URL %q", keyVaultKeyID.KeyVaultBaseUrl)
	}
	keyVaultID, err := commonids.ParseKeyVaultID(*keyVaultIDRaw)
	if err != nil {
		return nil, err
	}
	return utils.String(fmt.Sprintf("%s_%s_%s", keyVaultID.VaultName, keyVaultKeyID.Name, keyVaultKeyID.Version)), nil
}

func resourcePostgreSQLServerKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keysClient := meta.(*clients.Client).Postgres.ServerKeysClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverId, err := serverkeys.ParseServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}
	keyVaultKeyURI := d.Get("key_vault_key_id").(string)
	name, err := getPostgreSQLServerKeyName(ctx, keyVaultsClient, serverId.SubscriptionId, keyVaultKeyURI)
	if err != nil {
		return fmt.Errorf("cannot compose name for %s: %+v", serverId, err)
	}

	locks.ByName(serverId.ServerName, postgreSQLServerResourceName)
	defer locks.UnlockByName(serverId.ServerName, postgreSQLServerResourceName)

	id := serverkeys.NewKeyID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, *name)

	if d.IsNewResource() {
		// This resource is a singleton, but its name can be anything.
		// If you create a new key with different name with the old key, the service will not give you any warning but directly replace the old key with the new key.
		// Therefore sometimes you cannot get the old key using the GET API since you may not know the name of the old key
		resp, err := keysClient.List(ctx, *serverId)
		if err != nil {
			return fmt.Errorf("listing existing Keys in %s: %+v", serverId, err)
		}
		if resp.Model != nil && len(*resp.Model) >= 1 {
			keys := *resp.Model
			if rawId := keys[0].Id; rawId != nil && *rawId != "" {
				id, err := serverkeys.ParseKeyID(*rawId)
				if err != nil {
					return fmt.Errorf("parsing existing Server Key ID %q: %+v", *rawId, err)
				}

				// API allows adding same key again with Create action, which would trigger revalidation of the key on the server.
				// This is required to revalidate Replica server after creation.
				if *rawId != id.ID() {
					return tf.ImportAsExistsError("azurerm_postgresql_server_key", id.ID())
				}
			}
		}
	}

	param := serverkeys.ServerKey{
		Properties: &serverkeys.ServerKeyProperties{
			ServerKeyType: serverkeys.ServerKeyTypeAzureKeyVault,
			Uri:           utils.String(d.Get("key_vault_key_id").(string)),
		},
	}

	if err = keysClient.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePostgreSQLServerKeyRead(d, meta)
}

func resourcePostgreSQLServerKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keysClient := meta.(*clients.Client).Postgres.ServerKeysClient
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

	d.Set("server_id", serverkeys.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName).ID())
	if resp.Model != nil && resp.Model.Properties != nil {
		d.Set("key_vault_key_id", resp.Model.Properties.Uri)
	}

	return nil
}

func resourcePostgreSQLServerKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServerKeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serverkeys.ParseKeyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ServerName, postgreSQLServerResourceName)
	defer locks.UnlockByName(id.ServerName, postgreSQLServerResourceName)

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
