// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultManagedStorageAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultManagedStorageAccountCreateUpdate,
		Read:   resourceKeyVaultManagedStorageAccountRead,
		Update: resourceKeyVaultManagedStorageAccountCreateUpdate,
		Delete: resourceKeyVaultManagedStorageAccountDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.NestedItemName,
			},

			"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),

			"storage_account_key": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"key1",
					"key2",
				}, false),
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"regenerate_key_automatically": {
				Type:         pluginsdk.TypeBool,
				Optional:     true,
				Default:      false,
				RequiredWith: []string{"regeneration_period"},
			},

			"regeneration_period": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ISO8601Duration,
				RequiredWith: []string{"regenerate_key_automatically"},
			},

			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceKeyVaultManagedStorageAccountCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Base URI for Managed Storage Account Key Vault %q in %s: %+v", name, *keyVaultId, err)
	}

	if d.IsNewResource() {
		existing, err := client.GetStorageAccount(ctx, *keyVaultBaseUrl, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Managed Storage Account %q (Key Vault %q): %s", name, *keyVaultId, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_managed_storage_account", *existing.ID)
		}
	}

	t := d.Get("tags").(map[string]interface{})

	parameters := keyvault.StorageAccountCreateParameters{
		ResourceID:         utils.String(d.Get("storage_account_id").(string)),
		ActiveKeyName:      utils.String(d.Get("storage_account_key").(string)),
		AutoRegenerateKey:  utils.Bool(d.Get("regenerate_key_automatically").(bool)),
		RegenerationPeriod: utils.String(d.Get("regeneration_period").(string)),
		Tags:               tags.Expand(t),
	}

	if resp, err := client.SetStorageAccount(ctx, *keyVaultBaseUrl, name, parameters); err != nil {
		// In the case that the Storage Account already exists in a Soft Deleted / Recoverable state we check if `recover_soft_deleted_key_vaults` is set
		// and attempt recovery where appropriate
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeyVaults && utils.ResponseWasConflict(resp.Response) {
			recoveredStorageAccount, err := client.RecoverDeletedStorageAccount(ctx, *keyVaultBaseUrl, name)
			if err != nil {
				return fmt.Errorf("recovery of Managed Storage Account %q (Key Vault %q): %s", name, *keyVaultBaseUrl, err)
			}
			log.Printf("[DEBUG] Recovering Managed Storage Account %q (Key Vault %q)", name, *keyVaultBaseUrl)
			// We need to wait for consistency, recovered Key Vault Child items are not as readily available as newly created
			if secret := recoveredStorageAccount.ID; secret != nil {
				stateConf := &pluginsdk.StateChangeConf{
					Pending:                   []string{"pending"},
					Target:                    []string{"available"},
					Refresh:                   keyVaultChildItemRefreshFunc(*secret),
					Delay:                     30 * time.Second,
					PollInterval:              10 * time.Second,
					ContinuousTargetOccurence: 10,
					Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
				}

				if _, err := stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("waiting for Managed Storage Account %q (Key Vault %q) to become available after recovery: %s", name, *keyVaultId, err)
				}
				log.Printf("[DEBUG] Managed Storage Account %q recovered with ID: %q", name, *recoveredStorageAccount.ID)

				if _, err := client.SetStorageAccount(ctx, *keyVaultBaseUrl, name, parameters); err != nil {
					return fmt.Errorf("creating Managed Storage Account %s (Key Vault %q): %+v", name, *keyVaultId, err)
				}
			}
		} else {
			return fmt.Errorf("creating Managed Storage Account %s: %+v", name, err)
		}
	}

	read, err := client.GetStorageAccount(ctx, *keyVaultBaseUrl, name)
	if err != nil {
		return fmt.Errorf("cannot read Managed Storage Account %q (Key Vault %q): %+v", name, *keyVaultId, err)
	}
	if read.ID == nil {
		return fmt.Errorf("cannot read Managed Storage Account ID %q (Key Vault %q)", name, *keyVaultId)
	}

	storageId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(*read.ID)
	if err != nil {
		return err
	}

	d.SetId(storageId.ID())

	return resourceKeyVaultManagedStorageAccountRead(d, meta)
}

func resourceKeyVaultManagedStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(d.Id())
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID of the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if Key Vault %q for Managed Storage Account %q exists: %v", *keyVaultId, id.Name, err)
	}
	if !ok {
		log.Printf("[DEBUG] Managed Storage Account %q (Key Vault %q) was not found - removing from state", id.Name, *keyVaultId)
		d.SetId("")
		return nil
	}

	resp, err := client.GetStorageAccount(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Managed Storage Account %q (Key Vault %q) was not found - removing from state", id.Name, *keyVaultId)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Managed Storage Account %q (Key Vault %q): %+v", id.Name, *keyVaultId, err)
	}

	d.Set("name", id.Name)
	d.Set("key_vault_id", keyVaultId.ID())
	d.Set("storage_account_id", resp.ResourceID)
	d.Set("storage_account_key", resp.ActiveKeyName)
	d.Set("regenerate_key_automatically", resp.AutoRegenerateKey)
	d.Set("regeneration_period", resp.RegenerationPeriod)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceKeyVaultManagedStorageAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(d.Id())
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		return fmt.Errorf("determining the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking for existence of Key Vault %q for Managed Storage Account %q in Vault at url %q: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Managed Storage Account %q (Key Vault %q) was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	shouldPurge := meta.(*clients.Client).Features.KeyVault.PurgeSoftDeleteOnDestroy
	description := fmt.Sprintf("Secret %q (Key Vault %q)", id.Name, id.KeyVaultBaseUrl)
	deleter := deleteAndPurgeSecret{
		client:      client,
		keyVaultUri: id.KeyVaultBaseUrl,
		name:        id.Name,
	}
	if err := deleteAndOptionallyPurge(ctx, description, shouldPurge, deleter); err != nil {
		return err
	}

	return nil
}
