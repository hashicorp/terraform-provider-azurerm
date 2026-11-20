// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_serverless

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"configuration_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(elasticmonitorresources.PossibleValuesForConfigurationType(), false),
		},

		"elastic_cloud_email_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsEmailAddress,
		},

		"kind": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"offer_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z0-9_-]{1,50}$`),
				"must be between 1 and 50 characters in length and contain only lowercase letters, numbers, hyphens and underscores"),
		},

		"project_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(elasticmonitorresources.ProjectTypeElasticsearch),
				string(elasticmonitorresources.ProjectTypeObservability),
				string(elasticmonitorresources.ProjectTypeSecurity),
			}, false),
		},

		"sku": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"term_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"generate_api_key": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"monitoring_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},

		"plan_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "ess-consumption-2024",
			ForceNew: true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z0-9_-]{1,50}$`),
				"must be between 1 and 50 characters in length and contain only lowercase letters, numbers, hyphens and underscores"),
		},

		"publisher_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "elastic",
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"tags": commonschema.TagsForceNew(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"elastic_cloud_deployment_id": {
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
