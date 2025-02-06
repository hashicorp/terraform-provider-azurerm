// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	managedHsmHelpers "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceDiskEncryptionSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDiskEncryptionSetRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"auto_key_rotation_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"key_vault_key_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceDiskEncryptionSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	env := meta.(*clients.Client).Account.Environment
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewDiskEncryptionSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	model := resp.Model
	if err != nil || model == nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.DiskEncryptionSetName)
	d.Set("resource_group_name", id.ResourceGroupName)

	d.Set("location", location.NormalizeNilable(utils.String(model.Location)))

	if props := model.Properties; props != nil {
		d.Set("auto_key_rotation_enabled", props.RotationToLatestKeyVersionEnabled)

		if props.ActiveKey != nil && props.ActiveKey.KeyURL != "" {
			keyVaultURI := props.ActiveKey.KeyURL
			isHSMURI, err, _, _ := managedHsmHelpers.IsManagedHSMURI(env, keyVaultURI)
			if err != nil {
				return fmt.Errorf("parsing key vault URI: %+v", err)
			}
			if isHSMURI {
				d.Set("managed_hsm_key_id", keyVaultURI)
			} else {
				d.Set("key_vault_key_url", keyVaultURI)
			}
		}
	}

	flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}

	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, model.Tags)
}
