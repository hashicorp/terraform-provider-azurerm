// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBatchAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := batchaccount.NewBatchAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := location.Normalize(d.Get("location").(string))
	storageAccountId := d.Get("storage_account_id").(string)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_batch_account", id.ID())
		}
	}

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}

	encryptionRaw := d.Get("encryption").([]interface{})
	encryption := expandEncryption(encryptionRaw)

	poolAllocationMode := batchaccount.PoolAllocationMode(d.Get("pool_allocation_mode").(string))
	parameters := batchaccount.BatchAccountCreateParameters{
		Location: location,
		Properties: &batchaccount.BatchAccountCreateProperties{
			PoolAllocationMode:         &poolAllocationMode,
			PublicNetworkAccess:        pointer.To(batchaccount.PublicNetworkAccessTypeEnabled),
			Encryption:                 encryption,
			AllowedAuthenticationModes: expandAllowedAuthenticationModes(d.Get("allowed_authentication_modes").(*pluginsdk.Set).List()),
		},
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		parameters.Properties.PublicNetworkAccess = pointer.To(batchaccount.PublicNetworkAccessTypeDisabled)
	}

	if v, ok := d.GetOk("network_profile"); ok {
		parameters.Properties.NetworkProfile = expandBatchAccountNetworkProfile(v.([]interface{}))
	}

	// if pool allocation mode is UserSubscription, a key vault reference needs to be set
	if poolAllocationMode == batchaccount.PoolAllocationModeUserSubscription {
		keyVaultReferenceSet := d.Get("key_vault_reference").([]interface{})
		keyVaultReference, err := expandBatchAccountKeyVaultReference(keyVaultReferenceSet)
		if err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		if keyVaultReference == nil {
			return fmt.Errorf("creating %s: When setting pool allocation mode to UserSubscription, a Key Vault reference needs to be set", id)
		}

		parameters.Properties.KeyVaultReference = keyVaultReference

		if v, ok := d.GetOk("allowed_authentication_modes"); ok {
			authModes := v.(*pluginsdk.Set).List()
			for _, mode := range authModes {
				if batchaccount.AuthenticationMode(mode.(string)) == batchaccount.AuthenticationModeSharedKey {
					return fmt.Errorf("creating %s: When setting pool allocation mode to UserSubscription, `allowed_authentication_modes=[StorageKeys]` is not allowed. ", id)
				}
			}
		}
	}

	authMode := d.Get("storage_account_authentication_mode").(string)
	if batchaccount.AutoStorageAuthenticationMode(authMode) == batchaccount.AutoStorageAuthenticationModeBatchAccountManagedIdentity && identity.Type == "None" {
		return fmt.Errorf(" storage_account_authentication_mode=`BatchAccountManagedIdentity` can only be set when identity.type is `SystemAssigned` or `UserAssigned`")
	}

	if storageAccountId != "" {
		if authMode == "" {
			return fmt.Errorf("`storage_account_authentication_mode` is required when `storage_account_id` ")
		}
		parameters.Properties.AutoStorage = &batchaccount.AutoStorageBaseProperties{
			StorageAccountId:   storageAccountId,
			AuthenticationMode: pointer.ToEnum[batchaccount.AutoStorageAuthenticationMode](authMode),
		}
	}

	nodeIdentity := d.Get("storage_account_node_identity").(string)
	if nodeIdentity != "" {
		parameters.Properties.AutoStorage.NodeIdentityReference = &batchaccount.ComputeNodeIdentityReference{
			ResourceId: pointer.To(nodeIdentity),
		}
	}

	if err := client.CreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceBatchAccountRead(d, meta)
}
