// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-01-01/bastionhosts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceBastionHost() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBastionHostRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.BastionHostName,
			},

			"copy_paste_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"file_copy_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"ip_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"public_ip_address_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ip_connect_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"scale_units": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"shareable_link_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tunneling_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"session_recording_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"dns_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceBastionHostRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.BastionHostsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := bastionhosts.NewBastionHostID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		skuName := ""
		if sku := model.Sku; sku != nil && sku.Name != nil {
			skuName = string(*sku.Name)
		}
		d.Set("sku", skuName)

		if props := model.Properties; props != nil {
			d.Set("dns_name", props.DnsName)
			d.Set("scale_units", props.ScaleUnits)
			d.Set("file_copy_enabled", props.EnableFileCopy)
			d.Set("ip_connect_enabled", props.EnableIPConnect)
			d.Set("shareable_link_enabled", props.EnableShareableLink)
			d.Set("tunneling_enabled", props.EnableTunneling)
			d.Set("session_recording_enabled", props.EnableSessionRecording)

			copyPasteEnabled := true
			if props.DisableCopyPaste != nil {
				copyPasteEnabled = !*props.DisableCopyPaste
			}
			d.Set("copy_paste_enabled", copyPasteEnabled)

			if ipConfigs := props.IPConfigurations; ipConfigs != nil {
				if err := d.Set("ip_configuration", flattenBastionHostIPConfiguration(ipConfigs)); err != nil {
					return fmt.Errorf("flattening `ip_configuration`: %+v", err)
				}
			}
		}

		d.SetId(id.ID())

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
