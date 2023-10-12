// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultManagedStorageAccountSasTokenDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultManagedStorageAccountSasTokenDefinitionCreateUpdate,
		Read:   resourceKeyVaultManagedStorageAccountSasTokenDefinitionRead,
		Update: resourceKeyVaultManagedStorageAccountSasTokenDefinitionCreateUpdate,
		Delete: resourceKeyVaultManagedStorageAccountSasTokenDefinitionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SasDefinitionID(id)
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
				ValidateFunc: keyVaultValidate.NestedItemName,
			},

			"managed_storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.VersionlessNestedItemId,
			},

			"sas_template_uri": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"sas_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"service",
					"account",
				}, false),
			},

			"secret_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"validity_period": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceKeyVaultManagedStorageAccountSasTokenDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	resourcesClient := meta.(*clients.Client).Resource
	defer cancel()

	name := d.Get("name").(string)
	storageAccount, err := parse.ParseOptionallyVersionedNestedItemID(d.Get("managed_storage_account_id").(string))
	if err != nil {
		return err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, storageAccount.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID of the Key Vault at URL %q: %s", storageAccount.KeyVaultBaseUrl, err)
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Base URI for Managed Storage Account Key Vault %s: %+v", *keyVaultId, err)
	}

	if d.IsNewResource() {
		existing, err := client.GetSasDefinition(ctx, *keyVaultBaseUri, storageAccount.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Managed Storage Account Sas Defition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_managed_storage_account_sas_token_definition", *existing.ID)
		}
	}

	t := d.Get("tags").(map[string]interface{})
	parameters := keyvault.SasDefinitionCreateParameters{
		TemplateURI:    utils.String(d.Get("sas_template_uri").(string)),
		SasType:        keyvault.SasTokenType(d.Get("sas_type").(string)),
		ValidityPeriod: utils.String(d.Get("validity_period").(string)),
		SasDefinitionAttributes: &keyvault.SasDefinitionAttributes{
			Enabled: utils.Bool(true),
		},
		Tags: tags.Expand(t),
	}

	if resp, err := client.SetSasDefinition(ctx, *keyVaultBaseUri, storageAccount.Name, name, parameters); err != nil {
		// In the case that the Storage Account already exists in a Soft Deleted / Recoverable state we check if `recover_soft_deleted_key_vaults` is set
		// and attempt recovery where appropriate
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeyVaults && utils.ResponseWasConflict(resp.Response) {
			recoveredStorageAccount, err := client.RecoverDeletedSasDefinition(ctx, *keyVaultBaseUri, storageAccount.Name, name)
			if err != nil {
				return fmt.Errorf("recovery of Managed Storage Account SAS Definition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
			}
			log.Printf("[DEBUG] Recovering Managed Storage Account Sas Definition %q (Storage Account %q, Key Vault %q)", name, storageAccount.Name, *keyVaultId)
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
					return fmt.Errorf("waiting for Key Vault Managed Storage Account Sas Definition %q (Storage Account %q, Key Vault %q) to become available: %s", name, storageAccount.Name, *keyVaultId, err)
				}
				log.Printf("[DEBUG] Managed Storage Account Sas Definition %q (Storage Account %q, Key Vault %q) recovered", name, storageAccount.Name, *keyVaultId)

				if _, err := client.SetSasDefinition(ctx, *keyVaultBaseUri, storageAccount.Name, name, parameters); err != nil {
					return fmt.Errorf("creation of Managed Storage Account SAS Definition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
				}
			}
		} else {
			return fmt.Errorf("creation of Managed Storage Account SAS Definition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
		}
	}

	read, err := client.GetSasDefinition(ctx, *keyVaultBaseUri, storageAccount.Name, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read Managed Storage Account Sas Definition  %q (Storage Account %q, Key Vault %q)", name, storageAccount.Name, *keyVaultId)
	}

	sasId, err := parse.SasDefinitionID(*read.ID)
	if err != nil {
		return err
	}

	d.SetId(sasId.ID())

	return resourceKeyVaultManagedStorageAccountSasTokenDefinitionRead(d, meta)
}

func resourceKeyVaultManagedStorageAccountSasTokenDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SasDefinitionID(d.Id())
	if err != nil {
		return err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID of the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}

	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking for presence of existing Managed Storage Account Sas Defition %q (Key Vault %q): %+v", id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Managed Storage Account Sas Definition %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	resp, err := client.GetSasDefinition(ctx, id.KeyVaultBaseUrl, id.StorageAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Managed Storage Account Sas Definition %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("cannot read Managed Storage Account Sas Definition %s: %+v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("sas_template_uri", resp.TemplateURI)
	d.Set("sas_type", resp.SasType)
	d.Set("secret_id", resp.SecretID)
	d.Set("validity_period", resp.ValidityPeriod)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceKeyVaultManagedStorageAccountSasTokenDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SasDefinitionID(d.Id())
	if err != nil {
		return err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID of the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		return fmt.Errorf("unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Managed Storage Account Sas Definition %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Managed Storage Account Sas Definition %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	_, err = client.DeleteSasDefinition(ctx, id.KeyVaultBaseUrl, id.StorageAccountName, id.Name)
	return err
}
