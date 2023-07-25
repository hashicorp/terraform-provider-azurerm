// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2017-04-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = NotificationHubNamespaceResourceV0ToV1{}

type NotificationHubNamespaceResourceV0ToV1 struct{}

func (NotificationHubNamespaceResourceV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"namespace_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"servicebus_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (NotificationHubNamespaceResourceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := namespaces.ParseNamespaceIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, fmt.Errorf("parsing ID %q to upgrade: %+v", oldIdRaw, err)
		}

		rawState["id"] = oldId.ID()
		return rawState, nil
	}
}
