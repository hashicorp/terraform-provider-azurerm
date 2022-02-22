package signalr

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceWebPubsubPrivateLinkResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceWebPubsubPrivateLinkResourceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"web_pubsub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.WebPubsubID,
			},

			"shared_private_link_resource_types": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"subresource_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceWebPubsubPrivateLinkResourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	privateLinkResourceClient := meta.(*clients.Client).SignalR.WebPubsubPrivateLinkedResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	webPubsubID, err := parse.WebPubsubID(d.Get("web_pubsub_id").(string))
	if err != nil {
		return fmt.Errorf("parsing ID of %q: %+v", webPubsubID, err)
	}

	resourceList, err := privateLinkResourceClient.List(ctx, webPubsubID.ResourceGroup, webPubsubID.WebPubSubName)
	if err != nil {
		return fmt.Errorf("retrieving Private Link Resource for Resource %s: %+v", webPubsubID, err)
	}

	if resourceList.Values() == nil {
		return fmt.Errorf("retrieving Private Link Resource for Resource %s: `resource value` was nil", webPubsubID)
	}

	d.SetId(webPubsubID.ID())
	val := resourceList.Values()

	linkTypeList := make([]interface{}, 0)

	for _, v := range val {
		if v.PrivateLinkResourceProperties != nil && v.PrivateLinkResourceProperties.ShareablePrivateLinkResourceTypes != nil {
			for _, resource := range *v.PrivateLinkResourceProperties.ShareablePrivateLinkResourceTypes {
				item := make(map[string]interface{})
				if props := resource.Properties; props != nil {
					if props.GroupID != nil {
						item["subresource_name"] = props.GroupID
					}
					if props.Description != nil {
						item["description"] = props.Description
					}
				}
				linkTypeList = append(linkTypeList, item)
			}
		}
	}

	if err := d.Set("shared_private_link_resource_types", linkTypeList); err != nil {
		return fmt.Errorf("setting `shared_private_link_resource_types` error: %+v", err)
	}
	return nil
}
