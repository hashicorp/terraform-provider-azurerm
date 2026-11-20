// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
)

func storageShareSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageShareName,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"quota": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 102400),
		},

		"metadata": helpers.MetaDataComputedSchema(),

		"enabled_protocol": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(shares.SMB),
				string(shares.NFS),
			}, false),
			Default: string(shares.SMB),
		},

		"acl": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 64),
					},
					"access_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"start": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.IsRFC3339Time,
								},
								"expiry": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.IsRFC3339Time,
								},
								"permissions": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
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
			Computed: true, // azignore:AZS007 - pre-existing violation
			Optional: true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					string(shares.PremiumAccessTier),
					string(shares.HotAccessTier),
					string(shares.CoolAccessTier),
					string(shares.TransactionOptimizedAccessTier),
				}, false,
			),
		},

		"rbac_scope_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
