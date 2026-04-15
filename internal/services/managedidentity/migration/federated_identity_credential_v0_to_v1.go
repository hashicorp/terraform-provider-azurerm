// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = FederatedIdentityCredentialV0ToV1{}

type FederatedIdentityCredentialV0ToV1 struct{}

func (FederatedIdentityCredentialV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"audience": {
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Required: true,
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
		},
		"issuer": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"parent_id": {
			Type:     pluginsdk.TypeString,
			ForceNew: true,
			Required: true,
		},
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"subject": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (FederatedIdentityCredentialV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		if v, ok := rawState["user_assigned_identity_id"].(string); !ok || v == "" {
			if parentId, ok := rawState["parent_id"].(string); ok && parentId != "" {
				log.Printf("Copying `parent_id` to `user_assigned_identity_id`: %q", parentId)
				rawState["user_assigned_identity_id"] = parentId
			}
		}
		return rawState, nil
	}
}
