// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/communicationservices"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ServiceV0ToV1{}

type ServiceV0ToV1 struct{}

func (ServiceV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"data_location": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "United States",
		},

		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"primary_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (ServiceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		parsedId, err := communicationservices.ParseCommunicationServiceIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}
		newId := parsedId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q..", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
