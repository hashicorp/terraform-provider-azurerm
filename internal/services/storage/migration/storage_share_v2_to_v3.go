// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileshares"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = StorageShareV2ToV3{}

type StorageShareV2ToV3 struct{}

func (StorageShareV2ToV3) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"storage_account_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"quota": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"metadata": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"enabled_protocol": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"acl": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"access_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"start": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"expiry": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"permissions": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"access_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"rbac_scope_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (StorageShareV2ToV3) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldID, ok := rawState["id"].(string)
		if !ok {
			return rawState, fmt.Errorf("expected `id` to be of type string, got %T", rawState["id"])
		}

		// Some instances of this resource may already be using the resource manager ID, in which case this is a no-op
		if _, err := fileshares.ParseShareID(oldID); err == nil {
			return rawState, nil
		}

		storageAccountID, err := resolveStorageAccountIDForStateUpgrade(ctx, meta, rawState, func(idStr string) (*commonids.StorageAccountId, error) {
			id, err := fileshares.ParseShareIDInsensitively(idStr)
			if err != nil {
				return nil, err
			}
			return new(commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName)), nil
		})
		if err != nil {
			return rawState, err
		}

		nameRaw, ok := rawState["name"]
		if !ok {
			return rawState, errors.New("expected a `name` attribute to be present in state")
		}

		name, ok := nameRaw.(string)
		if !ok {
			return rawState, fmt.Errorf("expected `name` to be of type string, got %T", nameRaw)
		}

		rawState["id"] = fileshares.NewShareID(storageAccountID.SubscriptionId, storageAccountID.ResourceGroupName, storageAccountID.StorageAccountName, name).ID()

		return rawState, nil
	}
}
