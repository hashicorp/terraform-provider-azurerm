// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container_immutability_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_container_resource_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageContainerID,
		},

		"immutability_period_in_days": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 146000),
		},

		"locked": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"protected_append_writes_all_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"protected_append_writes_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff

			protectedAppendWritesAllEnabled := diff.Get("protected_append_writes_all_enabled").(bool)
			protectedAppendWritesEnabled := diff.Get("protected_append_writes_enabled").(bool)

			if protectedAppendWritesAllEnabled && protectedAppendWritesEnabled {
				return fmt.Errorf("`protected_append_writes_all_enabled` and `protected_append_writes_enabled` cannot be set at the same time")
			}

			lockedOld, lockedNew := diff.GetChange("locked")

			if lockedOld.(bool) && !lockedNew.(bool) {
				return fmt.Errorf("unable to set `locked = false` - once an immutability policy locked it cannot be unlocked")
			}

			if lockedOld.(bool) {
				if diff.HasChange("immutability_period_in_days") {
					if periodOld, periodNew := diff.GetChange("immutability_period_in_days"); periodOld.(int) < periodNew.(int) {
						return fmt.Errorf("`immutability_period_in_days` cannot be decreased once an immutability policy has been locked")
					}
				}

				if diff.HasChange("protected_append_writes_all_enabled") {
					return fmt.Errorf("`protected_append_writes_all_enabled` cannot be changed once an immutability policy has been locked")
				}

				if diff.HasChange("protected_append_writes_enabled") {
					return fmt.Errorf("`protected_append_writes_enabled` cannot be changed once an immutability policy has been locked")
				}
			}

			return nil
		},
	}
}
