// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/recordsets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = PTRRecordV0ToV1{}

type PTRRecordV0ToV1 struct{}

func (PTRRecordV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"zone_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"records": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Set: pluginsdk.HashString,
		},

		"ttl": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func (PTRRecordV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		parsedId, err := recordsets.ParseRecordTypeIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}
		newId := parsedId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
