// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_version

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_id": {
			Type:         pluginsdk.TypeString,
			Description:  "The ID of the API to which this API version belongs.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apis.ValidateApiID,
		},
		"name": {
			Type:         pluginsdk.TypeString,
			Description:  "The name of the API version.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},
		"lifecycle_stage": {
			Type:         pluginsdk.TypeString,
			Description:  "The lifecycle stage of the API version.",
			Required:     true,
			ValidateFunc: validation.StringInSlice(apiversions.PossibleValuesForLifecycleStage(), false),
		},
		"display_name": {
			Type:        pluginsdk.TypeString,
			Description: "The display name of the API version.",
			Optional:    true,
			// Note: O+C if display name is not set, it will default to the name of the API version
			Computed:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
