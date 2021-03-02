package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

func DnsZoneV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Version: 0,
		Type:    DnsZoneV0Schema().CoreConfigSchema().ImpliedType(),
		Upgrade: DnsZoneUpgradeV0ToV1,
	}
}

func DnsZoneV0Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"name_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"soa_record": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.DnsZoneSOARecordEmail,
						},

						"host_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"expire_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      2419200,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"minimum_ttl": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"refresh_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"retry_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"serial_number": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"ttl": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntBetween(0, 2147483647),
						},

						"tags": tags.Schema(),

						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func DnsZoneUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	ctx := context.TODO()
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
