// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dashboard_grafana_managed_private_endpoint

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedgrafanas"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Dashboard Grafana Managed Private Endpoint.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`\A([a-zA-Z]{1}[a-zA-Z0-9\-]{1,19}[a-zA-Z0-9]{1})\z`),
				`Name length can only consist of alphanumeric characters or dashes, and must be between 2 and 20 characters long. It must begin with a letter and end with a letter or digit.`,
			),
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),

		"grafana_id": {
			Description:  "The ID of the associated managed Grafana.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedgrafanas.ValidateGrafanaID,
		},

		"private_link_resource_id": {
			Description:  "The ID of the resource to which this Dashboard Grafana Managed Private Endpoint will connect.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"group_ids": {
			Description: "A list of Private Link group IDs. The value of this will depend on the Private Link resource to which you are connecting.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			// no validation here because an empty group id is valid for a generic private link resource
		},

		"private_link_resource_region": {
			Description:  "The region in which to create the Private Link.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"request_message": {
			Description:  "A message to provide in the request which will be seen by approvers.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"private_link_service_url": {
			Description: "A domain name for this endpoint to be used within Grafana.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^([0-9A-Za-z\-]+\.){2,}([0-9A-Za-z\-]+)\.?$`),
				"The URL must contain at least 3 parts separated by dots, containing only alphanumeric characters and hyphens",
			),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
