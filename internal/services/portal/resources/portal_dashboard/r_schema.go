// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package portal_dashboard

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func portalDashboardSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Shared Dashboard.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DashboardName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),

		"dashboard_properties": {
			Description:  "JSON data representing dashboard body. See above for details on how to obtain this from the Portal.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.DashboardProperties,
			StateFunc:    pluginsdk.NormalizeJson,
		},
	}
}
