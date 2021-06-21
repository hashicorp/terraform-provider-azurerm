package eventhub

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/sdk/namespaces"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceEventHubNamespaceCustomerManagedKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubNamespaceCustomerManagedKeyCreateUpdate,
		Read:   resourceEventHubNamespaceCustomerManagedKeyRead,
		Update: resourceEventHubNamespaceCustomerManagedKeyCreateUpdate,
		Delete: resourceEventHubNamespaceCustomerManagedKeyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.DefaultImporter(),

		Schema: map[string]*pluginsdk.Schema{
			"eventhub_namespace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NamespaceID,
			},

			"key_vault_key_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
				},
			},
		},
	}
}

func resourceEventHubNamespaceCustomerManagedKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Get("eventhub_namespace_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.Name, "azurerm_eventhub_namespace")
	defer locks.UnlockByName(id.Name, "azurerm_eventhub_namespace")

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if d.IsNewResource() {
		if resp.Model.Properties != nil && resp.Model.Properties.Encryption != nil {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace_customer_managed_key", id.ID())
		}
	}

	namespace := resp.Model

	keySource := namespaces.KeySourceMicrosoftKeyVault
	namespace.Properties.Encryption = &namespaces.Encryption{
		KeySource: &keySource,
	}

	keyVaultProps, err := expandEventHubNamespaceKeyVaultKeyIds(d.Get("key_vault_key_ids").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}
	namespace.Properties.Encryption.KeyVaultProperties = keyVaultProps

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *namespace); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceEventHubNamespaceCustomerManagedKeyRead(d, meta)
}

func resourceEventHubNamespaceCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if resp.Model.Properties == nil && resp.Model.Properties.Encryption == nil {
		d.SetId("")
		return nil
	}

	d.Set("eventhub_namespace_id", id.ID())

	if props := resp.Model.Properties; props != nil {
		keyVaultKeyIds, err := flattenEventHubNamespaceKeyVaultKeyIds(props.Encryption)
		if err != nil {
			return err
		}

		d.Set("key_vault_key_ids", keyVaultKeyIds)
	}

	return nil
}

func resourceEventHubNamespaceCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, "azurerm_eventhub_namespace")
	defer locks.UnlockByName(id.Name, "azurerm_eventhub_namespace")

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// Since this isn't a real object and it cannot be disabled once Customer Managed Key at rest has been enabled
	// And it must keep at least one key once Customer Managed Key is enabled
	// So for the delete operation, it has to recreate the EventHub Namespace with disabled Customer Managed Key
	future, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(future.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := waitForEventHubNamespaceToBeDeleted(ctx, client, *id); err != nil {
		return err
	}

	namespace := resp.Model
	namespace.Properties.Encryption = nil

	if err = client.CreateOrUpdateThenPoll(ctx, *id, *namespace); err != nil {
		return fmt.Errorf("removing %s: %+v", *id, err)
	}

	return nil
}

func expandEventHubNamespaceKeyVaultKeyIds(input []interface{}) (*[]namespaces.KeyVaultProperties, error) {
	if len(input) == 0 {
		return nil, nil
	}

	results := make([]namespaces.KeyVaultProperties, 0)

	for _, item := range input {
		keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(item.(string))
		if err != nil {
			return nil, err
		}

		results = append(results, namespaces.KeyVaultProperties{
			KeyName:     utils.String(keyId.Name),
			KeyVaultUri: utils.String(keyId.KeyVaultBaseUrl),
			KeyVersion:  utils.String(keyId.Version),
		})
	}

	return &results, nil
}

func flattenEventHubNamespaceKeyVaultKeyIds(input *namespaces.Encryption) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil || input.KeyVaultProperties == nil {
		return results, nil
	}

	for _, item := range *input.KeyVaultProperties {
		var keyName string
		if item.KeyName != nil {
			keyName = *item.KeyName
		}

		var keyVaultUri string
		if item.KeyVaultUri != nil {
			keyVaultUri = *item.KeyVaultUri
		}

		var keyVersion string
		if item.KeyVersion != nil {
			keyVersion = *item.KeyVersion
		}

		keyVaultKeyId, err := keyVaultParse.NewNestedItemID(keyVaultUri, "keys", keyName, keyVersion)
		if err != nil {
			return nil, err
		}

		results = append(results, keyVaultKeyId.ID())
	}

	return results, nil
}
