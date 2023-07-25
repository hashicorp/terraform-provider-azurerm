// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backupvaults"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataProtectionBackupVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataProtectionBackupVaultCreateUpdate,
		Read:   resourceDataProtectionBackupVaultRead,
		Update: resourceDataProtectionBackupVaultCreateUpdate,
		Delete: resourceDataProtectionBackupVaultDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := backupvaults.ParseBackupVaultIDInsensitively(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{2,50}$"),
					"DataProtection BackupVault name must be 2 - 50 characters long, contain only letters, numbers and hyphens.).",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"datastore_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(backupvaults.StorageSettingStoreTypesArchiveStore),
					string(backupvaults.StorageSettingStoreTypesSnapshotStore),
					string(backupvaults.StorageSettingStoreTypesVaultStore),
				}, false),
			},

			"redundancy": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(backupvaults.StorageSettingTypesGeoRedundant),
					string(backupvaults.StorageSettingTypesLocallyRedundant),
				}, false),
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"tags": tags.Schema(),
		},
	}
}

func resourceDataProtectionBackupVaultCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := backupvaults.NewBackupVaultID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing DataProtection BackupVault (%q): %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_protection_backup_vault", id.ID())
		}
	}

	expandedIdentity, err := expandBackupVaultDppIdentityDetails(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	datastoreType := backupvaults.StorageSettingStoreTypes(d.Get("datastore_type").(string))
	storageSettingType := backupvaults.StorageSettingTypes(d.Get("redundancy").(string))

	parameters := backupvaults.BackupVaultResource{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: backupvaults.BackupVault{
			StorageSettings: []backupvaults.StorageSetting{
				{
					DatastoreType: &datastoreType,
					Type:          &storageSettingType,
				},
			},
		},
		Identity: expandedIdentity,
		Tags:     expandTags(d.Get("tags").(map[string]interface{})),
	}
	err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating DataProtection BackupVault (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupVaultRead(d, meta)
}

func resourceDataProtectionBackupVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupvaults.ParseBackupVaultID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] DataProtection BackupVault %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupVault (%q): %+v", id, err)
	}
	d.Set("name", id.BackupVaultName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))
		props := model.Properties
		if props.StorageSettings != nil && len(props.StorageSettings) > 0 {
			d.Set("datastore_type", string(pointer.From((props.StorageSettings)[0].DatastoreType)))
			d.Set("redundancy", string(pointer.From((props.StorageSettings)[0].Type)))
		}

		if err = d.Set("identity", flattenBackupVaultDppIdentityDetails(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}

func resourceDataProtectionBackupVaultDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupvaults.ParseBackupVaultID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting DataProtection BackupVault (%q): %+v", id, err)
	}
	return nil
}

func expandBackupVaultDppIdentityDetails(input []interface{}) (*backupvaults.DppIdentityDetails, error) {
	config, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &backupvaults.DppIdentityDetails{
		Type: utils.String(string(config.Type)),
	}, nil
}

func flattenBackupVaultDppIdentityDetails(input *backupvaults.DppIdentityDetails) []interface{} {
	var config *identity.SystemAssigned
	if input != nil {
		principalId := ""
		if input.PrincipalId != nil {
			principalId = *input.PrincipalId
		}

		tenantId := ""
		if input.TenantId != nil {
			tenantId = *input.TenantId
		}
		config = &identity.SystemAssigned{
			Type:        identity.Type(*input.Type),
			PrincipalId: principalId,
			TenantId:    tenantId,
		}
	}
	return identity.FlattenSystemAssigned(config)
}
