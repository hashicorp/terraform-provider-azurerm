// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceCdnFrontDoorProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCdnFrontDoorProfileRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"response_timeout_seconds": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"log_scrubbing": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"scrubbing_rule": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"match_variable": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),

			"resource_guid": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCdnFrontDoorProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := profiles.NewProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("name", id.ProfileName)
		d.Set("resource_group_name", id.ResourceGroupName)

		if skuName := model.Sku.Name; skuName != nil {
			d.Set("sku_name", string(pointer.From(skuName)))
		}

		if identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity); err == nil {
			d.Set("identity", identity)
		}

		if props := model.Properties; props != nil {
			d.Set("response_timeout_seconds", int(pointer.From(props.OriginResponseTimeoutSeconds)))

			// whilst this is returned in the API as FrontDoorID other resources refer to
			// this as the Resource GUID, so we will for consistency
			d.Set("resource_guid", pointer.From(props.FrontDoorId))

			if err := d.Set("log_scrubbing", flattenCdnFrontDoorProfileLogScrubbing(props.LogScrubbing)); err != nil {
				return fmt.Errorf("setting `log_scrubbing`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
