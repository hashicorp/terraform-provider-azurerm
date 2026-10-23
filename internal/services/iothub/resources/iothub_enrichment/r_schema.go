// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_enrichment

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func iothubEnrichmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"key": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
				"Enrichment Key name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"iothub_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IoTHubName,
		},

		"value": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"endpoint_names": {
			Type:     pluginsdk.TypeList,
			MaxItems: 100,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			Required: true,
		},
	}
}
