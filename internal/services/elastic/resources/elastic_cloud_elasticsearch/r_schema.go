// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_elasticsearch

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func elasticCloudElasticsearchSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticsearchName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"elastic_cloud_email_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsEmailAddress,
		},

		"monitoring_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},

		"logs": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"filtering_tag": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"action": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(rules.PossibleValuesForTagAction(), false),
								},
							},
						},
					},

					"send_activity_logs": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"send_azuread_logs": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"send_subscription_logs": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"elastic_cloud_deployment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"elastic_cloud_sso_default_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"elastic_cloud_user_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"elasticsearch_service_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"kibana_service_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"kibana_sso_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
