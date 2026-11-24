// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_blob_event

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryTriggerBlobEventSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Factory Blob Event Trigger.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
		},

		"data_factory_id": {
			Description:  "The ID of Data Factory in which to associate the Trigger with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"storage_account_id": {
			Description:  "The ID of Storage Account in which blob event will be listened.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"events": {
			Description: "List of events that will fire this trigger.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"Microsoft.Storage.BlobCreated",
					"Microsoft.Storage.BlobDeleted",
				}, false),
			},
		},

		"pipeline": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The Data Factory Pipeline name that the trigger will act on.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
					},

					"parameters": {
						Description: "The Data Factory Pipeline parameters that the trigger will act on.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"activated": {
			Description: "Whether the Data Factory Blob Event Trigger is activated.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"additional_properties": {
			Description: "A map of additional properties to associate with the Data Factory Blob Event Trigger.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Blob Event Trigger.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"blob_path_begins_with": {
			Description:  "The pattern that blob path starts with for trigger to fire.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			AtLeastOneOf: []string{"blob_path_begins_with", "blob_path_ends_with"},
		},

		"blob_path_ends_with": {
			Description:  "The pattern that blob path ends with for trigger to fire.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			AtLeastOneOf: []string{"blob_path_begins_with", "blob_path_ends_with"},
		},

		"description": {
			Description:  "The description for the Data Factory Blob Event Trigger.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"ignore_empty_blobs": {
			Description: "Are blobs with zero bytes ignored?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},
	}
}
