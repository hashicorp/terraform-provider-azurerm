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
			Description:  "The name of the Data Factory Custom Event Trigger.",
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

		"eventgrid_topic_id": {
			Description:  "The ID of Event Grid Topic in which event will be listened.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: topics.ValidateTopicID,
		},

		"events": {
			Description: "List of events that will fire this trigger. At least one event must be specified.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MinItems:    1,
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
			Description: "Whether the Data Factory Custom Event Trigger is activated.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"additional_properties": {
			Description: "A map of additional properties to associate with the Data Factory Custom Event Trigger.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Custom Event Trigger.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"description": {
			Description:  "The description for the Data Factory Custom Event Trigger.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"subject_begins_with": {
			Description:  "The pattern that event subject starts with for trigger to fire.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"subject_ends_with": {
			Description:  "The pattern that event subject ends with for trigger to fire.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
