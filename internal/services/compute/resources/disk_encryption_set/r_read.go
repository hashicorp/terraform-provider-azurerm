// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package disk_encryption_set

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDiskEncryptionSetRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseDiskEncryptionSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Disk Encryption Set %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.DiskEncryptionSetName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			rotationToLatestKeyVersionEnabled := pointer.From(props.RotationToLatestKeyVersionEnabled)
			d.Set("auto_key_rotation_enabled", rotationToLatestKeyVersionEnabled)

			encryptionType := string(diskencryptionsets.DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey)
			if props.EncryptionType != nil {
				encryptionType = string(*props.EncryptionType)
			}
			d.Set("encryption_type", encryptionType)
			d.Set("federated_client_id", pointer.From(props.FederatedClientId))

			if props.ActiveKey != nil && props.ActiveKey.KeyURL != "" {
				key, err := keyvault.ParseNestedItemID(props.ActiveKey.KeyURL, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
				if err != nil {
					return err
				}
				d.Set("key_vault_key_url", key.ID())

				if rotationToLatestKeyVersionEnabled {
					key.Version = ""
				}

				d.Set("key_vault_key_id", key.ID())
			}
		}

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func getKeyURL(ctx context.Context, id *keyvault.NestedItemID, rotationToLatestKeyVersionEnabled bool, meta any) (string, error) {
	keyVaultClient := meta.(*clients.Client).KeyVault.ManagementClient
	managedHSMClient := meta.(*clients.Client).ManagedHSMs.DataPlaneKeysClient

	keyURL := ""

	if rotationToLatestKeyVersionEnabled {
		if id.Version != "" {
			return keyURL, fmt.Errorf("`auto_key_rotation_enabled` field is set to `true` expected a key vault key with a versionless ID but version information was found: %s", id)
		}

		if id.IsManagedHSM() {
			keyBundle, err := managedHSMClient.GetKey(ctx, id.KeyVaultBaseURL, id.Name, "")
			if err != nil {
				return keyURL, err
			}

			if keyBundle.Key != nil {
				keyURL = pointer.From(keyBundle.Key.Kid)
			}
		} else {
			keyBundle, err := keyVaultClient.GetKey(ctx, id.KeyVaultBaseURL, id.Name, "")
			if err != nil {
				return keyURL, err
			}

			if keyBundle.Key != nil {
				keyURL = pointer.From(keyBundle.Key.Kid)
			}
		}
	} else {
		if id.Version == "" {
			return keyURL, fmt.Errorf("`auto_key_rotation_enabled` field is set to `false` expected a key vault key with a versioned ID but no version information was found: %s", id)
		}
		keyURL = id.ID()
	}

	if keyURL == "" {
		return keyURL, errors.New("internal-error: received an unexpected empty key URL")
	}

	return keyURL, nil
}
