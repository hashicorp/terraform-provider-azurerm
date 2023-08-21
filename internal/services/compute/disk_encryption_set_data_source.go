// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := diskencryptionsets.NewDiskEncryptionSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

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

		if props.ActiveKey != nil && props.ActiveKey.KeyUrl != "" {
			d.Set("key_vault_key_url", props.ActiveKey.KeyUrl)
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
