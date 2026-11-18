// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package advisor_suppression

import (
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/advisor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The Name which should be used for this Advisor suppression.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"recommendation_id": {
			Description: "The ID of the Advisor recommendation to suppress.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"resource_id": {
			Description:  "The ID of the Resource to suppress the Advisor recommendation for.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
		"ttl": {
			Description:  "A optional time to live value. If omitted, the suppression will not expire.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.Duration,
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"suppression_id": {
			Description: "The GUID of the suppression.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
