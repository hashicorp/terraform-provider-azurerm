// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_policy_fragment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policyfragment"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementPolicyFragmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName("The name which should be used for this API Management Policy Fragment."),

		"api_management_id": commonschema.ResourceIDReferenceRequiredForceNew(&apimanagementservice.ServiceId{}, "API Management Service"),

		"format": {
			Description:  "The format of the Policy Fragment.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(policyfragment.PossibleValuesForPolicyFragmentContentFormat(), false),
			Default:      policyfragment.PolicyFragmentContentFormatXml,
		},

		"value": {
			Description:      "The value of the Policy Fragment.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: helpers.XmlWhitespaceDiffSuppress,
		},

		"description": {
			Description: "The description for the Policy Fragment.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},
	}
}
