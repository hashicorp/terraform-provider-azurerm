// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_manager

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"hub_profile": {
			Deprecated: "The service team has indicated this field is now deprecated and not to be used, as such we are marking it as such and no longer sending it to the API, please see url: https://learn.microsoft.com/azure/kubernetes-fleet/architectural-overview",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_prefix": {
						Required: true,
						Type:     pluginsdk.TypeString,
					},
					"fqdn": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"kubernetes_version": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
				},
			},
			ForceNew: true,
			MaxItems: 1,
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
