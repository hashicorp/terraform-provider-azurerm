// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BlobInventoryPolicyV0ToV1 struct {
}

func (BlobInventoryPolicyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_account_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"rules": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"storage_container_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"format": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"schedule": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"scope": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"schema_fields": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"filter": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"blob_types": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"include_blob_versions": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"include_deleted": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"include_snapshots": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"prefix_match": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									MaxItems: 10,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"exclude_prefixes": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									MaxItems: 10,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (BlobInventoryPolicyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)

		// This now uses the Storage Account ID since it's a 1:1 resource
		storageAccountIdRaw := rawState["storage_account_id"].(string)
		newId, err := commonids.ParseStorageAccountID(storageAccountIdRaw)
		if err != nil {
			return nil, err
		}
		newIdRaw := newId.ID()

		log.Printf("[DEBUG] Updating the Resource ID from %q to %q", oldIdRaw, newIdRaw)
		rawState["id"] = newIdRaw
		return rawState, nil
	}
}
