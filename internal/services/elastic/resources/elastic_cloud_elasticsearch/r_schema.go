// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name which should be used for this Elasticsearch resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticsearchName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Description: "The name of the SKU for this Elasticsearch.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"elastic_cloud_email_address": {
			Description:  "The Email Address which should be associated with this Elasticsearch account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsEmailAddress,
		},

		"monitoring_enabled": {
			Description: "Whether the Elasticsearch should have monitoring configured.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     true,
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
									Description:  "The name (key) of the Tag which should be filtered.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"value": {
									Description:  "The value of the Tag which should be filtered.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"action": {
									Description:  "The type of action which should be taken when the Tag matches the `name` and `value`.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(rules.PossibleValuesForTagAction(), false),
								},
							},
						},
					},

					"send_activity_logs": {
						Description: "Whether the Azure Activity Logs should be sent to the Elasticsearch cluster.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"send_azuread_logs": {
						Description: "Whether the AzureAD Logs should be sent to the Elasticsearch cluster.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"send_subscription_logs": {
						Description: "Whether the Azure Subscription Logs should be sent to the Elasticsearch cluster.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"elastic_cloud_deployment_id": {
			Description: "The ID of the Deployment within Elastic Cloud.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"elastic_cloud_sso_default_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"elastic_cloud_user_id": {
			Description: "The ID of the User Account within Elastic Cloud.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"elasticsearch_service_url": {
			Description: "The URL to the Elasticsearch Service associated with this Elasticsearch.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"kibana_service_url": {
			Description: "The URL to the Kibana Dashboard associated with this Elasticsearch.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"kibana_sso_uri": {
			Description: "The URI used for SSO to the Kibana Dashboard associated with this Elasticsearch.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
