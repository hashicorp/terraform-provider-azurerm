// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package disk_encryption_set

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDiskEncryptionSetCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewDiskEncryptionSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_disk_encryption_set", id.ID())
		}
	}

	rotationToLatestKeyVersionEnabled := d.Get("auto_key_rotation_enabled").(bool)
	activeKey := &diskencryptionsets.KeyForDiskEncryptionSet{}

	key, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return err
	}

	keyURL, err := getKeyURL(ctx, key, rotationToLatestKeyVersionEnabled, meta)
	if err != nil {
		return err
	}
	activeKey.KeyURL = keyURL

	expandedIdentity, err := expandDiskEncryptionSetIdentity(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	params := diskencryptionsets.DiskEncryptionSet{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &diskencryptionsets.EncryptionSetProperties{
			ActiveKey:                         activeKey,
			RotationToLatestKeyVersionEnabled: pointer.To(rotationToLatestKeyVersionEnabled),
			EncryptionType:                    pointer.ToEnum[diskencryptionsets.DiskEncryptionSetType](d.Get("encryption_type").(string)),
		},
		Identity: expandedIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("federated_client_id"); ok {
		params.Properties.FederatedClientId = pointer.To(v.(string))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, params, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDiskEncryptionSetRead(d, meta)
}

func expandDiskEncryptionSetIdentity(input []any) (*identity.SystemAndUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	return expanded, nil
}
