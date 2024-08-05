// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupvaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDataProtectionBackupVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDataProtectionBackupVaultRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{2,50}$"),
					"DataProtection BackupVault name must be 2 - 50 characters long, contain only letters, numbers and hyphens.).",
				),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"datastore_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"redundancy": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedIdentityComputed(),

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceDataProtectionBackupVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := backupvaults.NewBackupVaultID(subscriptionId, resourceGroup, name)

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving DataProtection BackupVault (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.BackupVaultName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		props := model.Properties
		if props.StorageSettings != nil && len(props.StorageSettings) > 0 {
			d.Set("datastore_type", string(pointer.From((props.StorageSettings)[0].DatastoreType)))
			d.Set("redundancy", string(pointer.From((props.StorageSettings)[0].Type)))
		}

		if err = d.Set("identity", dataSourceFlattenBackupVaultDppIdentityDetails(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func dataSourceFlattenBackupVaultDppIdentityDetails(input *backupvaults.DppIdentityDetails) []interface{} {
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
