// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package communication_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Communication Service resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CommunicationServiceName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"data_location": {
			Description: "The location where the Communication Service stores its data at rest.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"Africa",
				"Asia Pacific",
				"Australia",
				"Brazil",
				"Canada",
				"Europe",
				"France",
				"Germany",
				"India",
				"Japan",
				"Korea",
				"Norway",
				"Switzerland",
				"UAE",
				"UK",
				"United States",
				"usgov",
			}, false),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"primary_connection_string": {
			Description: "The primary connection string of the Communication Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_connection_string": {
			Description: "The secondary connection string of the Communication Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_key": {
			Description: "The primary key of the Communication Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_key": {
			Description: "The secondary key of the Communication Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"hostname": {
			Description: "The hostname of the Communication Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
