// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespacetopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Event Grid Namespace Topic.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-]{3,50}$"),
					"Event Grid Namespace Topic name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			),
		},

		"eventgrid_namespace_id": {
			Description:  "The ID of the Event Grid Namespace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespacetopics.ValidateNamespaceID,
		},

		"event_retention_in_days": {
			Description:  "Event retention for the namespace topic expressed in days.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      7,
			ValidateFunc: validation.IntBetween(1, 7),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
