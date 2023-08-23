// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = SRVRecordV0ToV1{}

type SRVRecordV0ToV1 struct{}

func (SRVRecordV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"record": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"priority": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"weight": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"port": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"target": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
			Set: resourceDnsSrvRecordHash,
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

func (SRVRecordV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
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

func resourceDnsSrvRecordHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%d-", m["priority"].(int)))
		buf.WriteString(fmt.Sprintf("%d-", m["weight"].(int)))
		buf.WriteString(fmt.Sprintf("%d-", m["port"].(int)))
		buf.WriteString(fmt.Sprintf("%s-", m["target"].(string)))
	}

	return pluginsdk.HashString(buf.String())
}
