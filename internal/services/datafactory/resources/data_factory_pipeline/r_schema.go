// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_pipeline

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryPipelineSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Factory Pipeline. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/naming-rules) for all restrictions.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
		},

		"data_factory_id": {
			Description:  "The Data Factory ID in which to associate the Linked Service with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"parameters": {
			Description: "A map of parameters to associate with the Data Factory Pipeline.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"variables": {
			Description: "A map of variables to associate with the Data Factory Pipeline.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Description: "The description for the Data Factory Pipeline.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"activities_json": {
			Description:      "A JSON object that contains the activities that will be associated with the Data Factory Pipeline.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			StateFunc:        pluginsdk.NormalizeJson,
			DiffSuppressFunc: helpers.SuppressJsonOrderingDifference,
		},

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Pipeline.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"concurrency": {
			Description:  "The max number of concurrent runs for the Data Factory Pipeline.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 50),
		},

		"folder": {
			Description:  "The folder that this Pipeline is in. If not specified, the Pipeline will appear at the root level.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"monitor_metrics_after_duration": {
			Description:  "The TimeSpan value after which an Azure Monitoring Metric is fired.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
