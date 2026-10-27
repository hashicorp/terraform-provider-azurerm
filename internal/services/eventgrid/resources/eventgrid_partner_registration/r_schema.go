// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_registration

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
					"EventGrid partner registration name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			),
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"tags":                commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"partner_registration_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
