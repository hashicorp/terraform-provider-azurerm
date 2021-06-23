package resource

import (
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func dataSourceResourceId() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceResourceIdRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Second),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_id": {
				Type:        pluginsdk.TypeString,
				Required:    true,
				Description: "The Azure resource id to parse.",
			},
			"subscription_id": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "The parsed Azure subscription.",
			},
			"resource_group_name": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "The parsed Azure resource group name.",
			},
			"resource_type": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "The type of the Resource. (e.g. `Microsoft.Network/virtualNetworks`).",
			},
			"secondary_resource_type": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "The type of the child resource",
			},
			"parts": {
				Type:        pluginsdk.TypeMap,
				Computed:    true,
				Description: "A map of any additional key-value pairs in the path, this includes the resource name, accessed using an index of the key name.",
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceResourceIdRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id := d.Get("resource_id").(string)
	resourceId, err := azure.ParseAzureResourceID(id)
	if err != nil {
		return err
	}
	d.Set("subscription_id", resourceId.SubscriptionID)
	d.Set("resource_group_name", resourceId.ResourceGroup)
	d.Set("resource_type", resourceId.Provider)
	d.Set("secondary_resource_type", resourceId.SecondaryProvider)
	d.Set("parts", resourceId.Path)
	d.SetId(id)
	return nil
}
