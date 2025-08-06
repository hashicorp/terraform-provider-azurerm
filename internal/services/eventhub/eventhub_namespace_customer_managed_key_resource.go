// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := namespaces.ParseNamespaceID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).Eventhub.NamespacesClient

			var cancel context.CancelFunc
			ctx, cancel = timeouts.ForRead(ctx, d)
			defer cancel()

			id, err := namespaces.ParseNamespaceID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Encryption == nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving %s: no customer managed key present", *id)
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: map[string]*pluginsdk.Schema{
			"eventhub_namespace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: namespaces.ValidateNamespaceID,
			},

			"key_vault_key_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
				},
			},

			"infrastructure_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"user_assigned_identity_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateUserAssignedIdentityID,
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

	locks.ByName(id.NamespaceName, "azurerm_eventhub_namespace")
	defer locks.UnlockByName(id.NamespaceName, "azurerm_eventhub_namespace")

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
	keySource := namespaces.KeySourceMicrosoftPointKeyVault
	namespace.Properties.Encryption = &namespaces.Encryption{
		KeySource: &keySource,
	}

	keyVaultProps, err := expandEventHubNamespaceKeyVaultKeyIds(d.Get("key_vault_key_ids").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}

	userAssignedIdentity := d.Get("user_assigned_identity_id").(string)
	if userAssignedIdentity != "" && keyVaultProps != nil {
		userAssignedIdentityId, err := commonids.ParseUserAssignedIdentityID(userAssignedIdentity)
		if err != nil {
			return err
		}

		// this provides a more helpful error message than the API response
		if namespace.Identity == nil {
			return fmt.Errorf("user assigned identity '%s' must also be assigned to the parent event hub - currently no user assigned identities are assigned to the parent event hub", userAssignedIdentity)
		}

		isIdentityAssignedToParent := false
		for item := range namespace.Identity.IdentityIds {
			parentEhnUaiId, err := commonids.ParseUserAssignedIdentityIDInsensitively(item)
			if err != nil {
				return fmt.Errorf("parsing %q as a User Assigned Identity ID: %+v", item, err)
			}
			if resourceids.Match(parentEhnUaiId, userAssignedIdentityId) {
				isIdentityAssignedToParent = true
			}
		}

		// this provides a more helpful error message than the API response
		if !isIdentityAssignedToParent {
			return fmt.Errorf("user assigned identity '%s' must also be assigned to the parent event hub", userAssignedIdentity)
		}

		for i := 0; i < len(*keyVaultProps); i++ {
			(*keyVaultProps)[i].Identity = &namespaces.UserAssignedIdentityProperties{
				UserAssignedIdentity: &userAssignedIdentity,
			}
		}
	}

	namespace.Properties.Encryption.KeyVaultProperties = keyVaultProps
	namespace.Properties.Encryption.RequireInfrastructureEncryption = utils.Bool(d.Get("infrastructure_encryption_enabled").(bool))

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
		d.Set("infrastructure_encryption_enabled", props.Encryption.RequireInfrastructureEncryption)

		if kvprops := props.Encryption.KeyVaultProperties; kvprops != nil {
			// we can only have a single user managed id for N number of keys, azure portal only allows setting a single one and then applies it to each key
			for _, item := range *kvprops {
				if item.Identity != nil && item.Identity.UserAssignedIdentity != nil {
					userAssignedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*item.Identity.UserAssignedIdentity)
					if err != nil {
						return fmt.Errorf("parsing `user_assigned_identity_id`: %+v", err)
					}
					if err := d.Set("user_assigned_identity_id", userAssignedId.ID()); err != nil {
						return fmt.Errorf("setting `user_assigned_identity_id`: %+v", err)
					}

					break
				}
			}
		}
	}

	return nil
}

func resourceEventHubNamespaceCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	log.Printf(`[INFO] Customer Managed Keys cannot be removed from EventHub Namespaces once added. To remove the Customer Managed Key delete and recreate the parent EventHub Namespace`)
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

func flattenEventHubNamespaceKeyVaultKeyIds(input *namespaces.Encryption) ([]string, error) {
	results := make([]string, 0)
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

		keyVaultKeyId, err := keyVaultParse.NewNestedItemID(keyVaultUri, keyVaultParse.NestedItemTypeKey, keyName, keyVersion)
		if err != nil {
			return nil, err
		}

		results = append(results, keyVaultKeyId.ID())
	}

	return results, nil
}
