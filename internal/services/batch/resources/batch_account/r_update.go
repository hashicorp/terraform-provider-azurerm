// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBatchAccountUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Batch.AccountClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := batchaccount.ParseBatchAccountID(d.Id())
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]any)

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}

	encryptionRaw := d.Get("encryption").([]any)
	encryption := expandEncryption(encryptionRaw)

	parameters := batchaccount.BatchAccountUpdateParameters{
		Properties: &batchaccount.BatchAccountUpdateProperties{
			Encryption: encryption,
		},
		Identity: identity,
		Tags:     tags.Expand(t),
	}

	if d.HasChange("allowed_authentication_modes") {
		allowedAuthModes := d.Get("allowed_authentication_modes").(*pluginsdk.Set).List()
		if len(allowedAuthModes) == 0 {
			parameters.Properties.AllowedAuthenticationModes = &[]batchaccount.AuthenticationMode{} // remove all modes need explicit set it to empty array not nil
		} else {
			parameters.Properties.AllowedAuthenticationModes = expandAllowedAuthenticationModes(d.Get("allowed_authentication_modes").(*pluginsdk.Set).List())
		}
	}

	if d.HasChange("public_network_access_enabled") {
		if d.Get("public_network_access_enabled").(bool) {
			parameters.Properties.PublicNetworkAccess = pointer.To(batchaccount.PublicNetworkAccessTypeEnabled)
		} else {
			parameters.Properties.PublicNetworkAccess = pointer.To(batchaccount.PublicNetworkAccessTypeDisabled)
		}
	}

	if d.HasChange("network_profile") {
		parameters.Properties.NetworkProfile = expandBatchAccountNetworkProfile(d.Get("network_profile").([]any))
	}

	if d.HasChange("storage_account_id") {
		if v, ok := d.GetOk("storage_account_id"); ok {
			parameters.Properties.AutoStorage = &batchaccount.AutoStorageBaseProperties{
				StorageAccountId: v.(string),
			}
		} else {
			parameters.Properties.AutoStorage = &batchaccount.AutoStorageBaseProperties{}
		}
	}

	authMode := d.Get("storage_account_authentication_mode").(string)
	if batchaccount.AutoStorageAuthenticationMode(authMode) == batchaccount.AutoStorageAuthenticationModeBatchAccountManagedIdentity && identity.Type == "None" {
		return fmt.Errorf(" storage_account_authentication_mode=`BatchAccountManagedIdentity` can only be set when identity.type is `SystemAssigned` or `UserAssigned`")
	}

	storageAccountId := d.Get("storage_account_id").(string)
	if storageAccountId != "" {
		parameters.Properties.AutoStorage = &batchaccount.AutoStorageBaseProperties{
			StorageAccountId:   storageAccountId,
			AuthenticationMode: pointer.ToEnum[batchaccount.AutoStorageAuthenticationMode](authMode),
		}
	}

	nodeIdentity := d.Get("storage_account_node_identity").(string)
	if nodeIdentity != "" {
		parameters.Properties.AutoStorage.NodeIdentityReference = &batchaccount.ComputeNodeIdentityReference{
			ResourceId: new(nodeIdentity),
		}
	}

	if _, err = client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceBatchAccountRead(d, meta)
}
