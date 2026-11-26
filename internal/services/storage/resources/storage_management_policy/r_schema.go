// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_management_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func storageManagementPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},
		"rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"filters": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"blob_types": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Set:      pluginsdk.HashString,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"blockBlob",
											"appendBlob",
										}, false),
									},
								},
								"prefix_match": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Set:      pluginsdk.HashString,
									Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
								},
								"match_blob_index_tag": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 128),
											},

											"operation": {
												Type:     pluginsdk.TypeString,
												Optional: true,
												Default:  "==",
												ValidateFunc: validation.StringInSlice([]string{
													"==",
												}, false),
											},

											"value": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(0, 256),
											},
										},
									},
								},
							},
						},
					},
					// lintignore:XS003
					"actions": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								// lintignore:XS003
								"base_blob": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"tier_to_cool_after_days_since_modification_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_cool_after_days_since_last_access_time_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"auto_tier_to_hot_from_cool_enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},
											"tier_to_cool_after_days_since_creation_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_archive_after_days_since_modification_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_archive_after_days_since_last_access_time_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_archive_after_days_since_last_tier_change_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_archive_after_days_since_creation_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_cold_after_days_since_modification_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_cold_after_days_since_last_access_time_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_cold_after_days_since_creation_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"delete_after_days_since_modification_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"delete_after_days_since_last_access_time_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"delete_after_days_since_creation_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
										},
									},
								},
								// lintignore:XS003
								"snapshot": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"change_tier_to_archive_after_days_since_creation": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_archive_after_days_since_last_tier_change_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"change_tier_to_cool_after_days_since_creation": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_cold_after_days_since_creation_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"delete_after_days_since_creation_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
										},
									},
								},
								"version": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"change_tier_to_archive_after_days_since_creation": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_archive_after_days_since_last_tier_change_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"change_tier_to_cool_after_days_since_creation": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"tier_to_cold_after_days_since_creation_greater_than": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
											"delete_after_days_since_creation": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      -1,
												ValidateFunc: validation.IntBetween(0, 99999),
											},
										},
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
