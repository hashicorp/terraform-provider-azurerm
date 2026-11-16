// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lighthouse_assignment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func lighthouseAssignmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IsUUID,
		},

		"lighthouse_definition_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: registrationdefinitions.ValidateScopedRegistrationDefinitionID,
		},

		"scope": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.Any(commonids.ValidateSubscriptionID, commonids.ValidateResourceGroupID),
		},
	}
}
