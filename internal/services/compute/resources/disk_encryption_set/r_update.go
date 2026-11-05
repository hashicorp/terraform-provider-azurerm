// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package disk_encryption_set

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDiskEncryptionSetUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseDiskEncryptionSetID(d.Id())
	if err != nil {
		return err
	}

	update := diskencryptionsets.DiskEncryptionSetUpdate{}

	if d.HasChange("identity") {
		expandedIdentity, err := expandDiskEncryptionSetIdentity(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		update.Identity = expandedIdentity
	}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	rotationToLatestKeyVersionEnabled := d.Get("auto_key_rotation_enabled").(bool)

	if d.HasChange("key_vault_key_id") {
		if update.Properties == nil {
			update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{
				ActiveKey: &diskencryptionsets.KeyForDiskEncryptionSet{},
			}
		}

		key, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
		if err != nil {
			return err
		}

		keyURL, err := getKeyURL(ctx, key, rotationToLatestKeyVersionEnabled, meta)
		if err != nil {
			return err
		}
		update.Properties.ActiveKey.KeyURL = keyURL
	}

	if d.HasChange("auto_key_rotation_enabled") {
		if update.Properties == nil {
			update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{}
		}

		update.Properties.RotationToLatestKeyVersionEnabled = pointer.To(rotationToLatestKeyVersionEnabled)
	}

	if d.HasChange("federated_client_id") {
		if update.Properties == nil {
			update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{}
		}

		update.Properties.FederatedClientId = pointer.To("None")
		if v, ok := d.GetOk("federated_client_id"); ok {
			update.Properties.FederatedClientId = pointer.To(v.(string))
		}
	}

	if err := client.UpdateThenPoll(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDiskEncryptionSetRead(d, meta)
}
