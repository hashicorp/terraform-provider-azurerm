// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_cluster

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func eventhubClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9]([-._a-zA-Z0-9]{0,48}[a-zA-Z0-9])?$"),
				"The event hub name can contain only letters, numbers, periods (.), hyphens (-),and underscores (_), up to 50 characters, and it must begin and end with a letter or number.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^Dedicated_[1-9][0-9]*$`),
				"SKU name must match /^Dedicated_[1-9][0-9]*$/.",
			),
		},

		"tags": commonschema.Tags(),
	}
}
