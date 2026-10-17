// Copyright IBM Corp. 2014, 2025
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
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_id": commonschema.ResourceIDReferenceRequiredForceNew(&apimanagementservice.ServiceId{}),

		"format": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(policyfragment.PossibleValuesForPolicyFragmentContentFormat(), false),
			Default:      policyfragment.PolicyFragmentContentFormatXml,
		},

		"value": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: helpers.XmlWhitespaceDiffSuppress,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}
