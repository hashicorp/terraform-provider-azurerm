// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/consumergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ConsumerGroupsV0ToV1{}

type ConsumerGroupsV0ToV1 struct{}

func (ConsumerGroupsV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumergroups/consumergroup1
		// new:
		// 	/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumerGroups/consumergroup1
		oldId := rawState["id"].(string)
		parsed, err := consumergroups.ParseConsumerGroupIDInsensitively(oldId)
		if err != nil {
			return rawState, fmt.Errorf("parsing existing Consumer Group ID %q: %+v", oldId, err)
		}
		rawState["id"] = parsed.ID()

		return rawState, nil
	}
}

func (ConsumerGroupsV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return consumerGroupsSchemaForV0AndV1()
}

func consumerGroupsSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"namespace_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"eventhub_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"user_metadata": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}
