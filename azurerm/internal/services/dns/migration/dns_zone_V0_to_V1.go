package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DnsZoneV0ToV1{}

type DnsZoneV0ToV1 struct{}

func (DnsZoneV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"number_of_record_sets": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"max_number_of_record_sets": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"name_servers": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:      pluginsdk.HashString,
		},

		"soa_record": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"host_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"expire_time": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  2419200,
					},

					"minimum_ttl": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  300,
					},

					"refresh_time": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  3600,
					},

					"retry_time": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  300,
					},

					"serial_number": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  1,
					},

					"ttl": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  3600,
					},

					"tags": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"fqdn": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (DnsZoneV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		groupsClient := meta.(*clients.Client).Resource.GroupsClient
		oldId := rawState["id"].(string)
		id, err := parse.DnsZoneID(oldId)
		if err != nil {
			return rawState, err
		}
		resGroup, err := groupsClient.Get(ctx, id.ResourceGroup)
		if err != nil {
			return rawState, err
		}
		if resGroup.Name == nil {
			return rawState, fmt.Errorf("`name` was nil for Resource Group %q", id.ResourceGroup)
		}
		resourceGroup := *resGroup.Name
		name := rawState["name"].(string)
		newId := parse.NewDnsZoneID(id.SubscriptionId, resourceGroup, name).ID()
		log.Printf("Updating `id` from %q to %q", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
