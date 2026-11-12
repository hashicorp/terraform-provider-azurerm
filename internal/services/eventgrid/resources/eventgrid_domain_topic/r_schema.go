// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_domain_topic

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func eventgridDomainTopicSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,128}$"),
					"EventGrid domain topic name must be 3 - 128 characters long, contain only letters, numbers and hyphens.",
				),
			),
		},

		"domain_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
					"EventGrid domain name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}
