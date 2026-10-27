// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_custom_event

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryTriggerCustomEventSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
		},

		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"eventgrid_topic_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: topics.ValidateTopicID,
		},

		"events": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"pipeline": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
					},

					"parameters": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"activated": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"additional_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"subject_begins_with": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"subject_ends_with": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
