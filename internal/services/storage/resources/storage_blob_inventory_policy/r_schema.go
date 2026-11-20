// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_blob_inventory_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobinventorypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func storageBlobInventoryPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.StorageAccountId{}),

		"rules": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"storage_container_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.StorageContainerName,
					},

					"format": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(blobinventorypolicies.PossibleValuesForFormat(), false),
					},

					"schedule": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(blobinventorypolicies.PossibleValuesForSchedule(), false),
					},

					"scope": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(blobinventorypolicies.PossibleValuesForObjectType(), false),
					},

					"schema_fields": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
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
										ValidateFunc: validation.StringInSlice([]string{
											"blockBlob",
											"appendBlob",
											"pageBlob",
										}, false),
									},
								},

								"created_within_days": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(1, 36500),
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
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},

								"exclude_prefixes": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									MaxItems: 10,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
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
