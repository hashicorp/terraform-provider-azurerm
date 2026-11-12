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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`\A([a-zA-Z]{1}[a-zA-Z0-9\-]{1,19}[a-zA-Z0-9]{1})\z`),
				`Name length can only consist of alphanumeric characters or dashes, and must be between 2 and 20 characters long. It must begin with a letter and end with a letter or digit.`,
			),
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),

		"grafana_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedgrafanas.ValidateGrafanaID,
		},

		"private_link_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"group_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			// no validation here because an empty group id is valid for a generic private link resource
		},

		"private_link_resource_region": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"request_message": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"private_link_service_url": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
