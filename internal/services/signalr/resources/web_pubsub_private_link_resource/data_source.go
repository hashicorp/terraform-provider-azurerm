// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_private_link_resource

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about the Private Link Resource supported by the Web Pubsub Resource.",
		Read:        dataSourceWebPubsubPrivateLinkResourceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"web_pubsub_id": {
				Description:  "The ID of an existing Web Pubsub Resource which Private Link Resource should be retrieved for.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: webpubsub.ValidateWebPubSubID,
			},

			"shared_private_link_resource_types": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"subresource_name": {
							Description: "The name for the resource that has been onboarded to Private Link Service.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"description": {
							Description: "The description of the resource type that has been onboarded to Private Link Service.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceWebPubsubPrivateLinkResourceRead(d *pluginsdk.ResourceData, meta any) error {
	privateLinkResourceClient := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	webPubSubIdRaw := d.Get("web_pubsub_id").(string)
	webPubSubId, err := webpubsub.ParseWebPubSubID(webPubSubIdRaw)
	if err != nil {
		return err
	}

	resourceList, err := privateLinkResourceClient.PrivateLinkResourcesListComplete(ctx, *webPubSubId)
	if err != nil {
		return fmt.Errorf("retrieving Private Link Resources for %s: %+v", *webPubSubId, err)
	}

	if resourceList.Items == nil {
		return fmt.Errorf("retrieving Private Link Resource for %s: `items` was nil", webPubSubId)
	}

	d.SetId(webPubSubId.ID())

	linkTypeList := make([]any, 0)
	for _, v := range resourceList.Items {
		if v.Properties != nil {
			if v.Properties.ShareablePrivateLinkResourceTypes == nil {
				continue
			}

			for _, resource := range *v.Properties.ShareablePrivateLinkResourceTypes {
				description := ""
				subResourceName := ""
				if props := resource.Properties; props != nil {
					if props.GroupId != nil {
						subResourceName = *props.GroupId
					}
					if props.Description != nil {
						description = *props.Description
					}
				}
				linkTypeList = append(linkTypeList, map[string]any{
					"description":      description,
					"subresource_name": subResourceName,
				})
			}
		}
	}

	if err := d.Set("shared_private_link_resource_types", linkTypeList); err != nil {
		return fmt.Errorf("setting `shared_private_link_resource_types` error: %+v", err)
	}
	return nil
}
