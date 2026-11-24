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
			Description:  "A unique UUID/GUID which identifies this lighthouse assignment- one will be generated if not specified.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IsUUID,
		},

		"lighthouse_definition_id": {
			Description:  "A Fully qualified path of the lighthouse definition, such as `/subscriptions/0afefe50-734e-4610-8c82-a144aff49dea/providers/Microsoft.ManagedServices/registrationDefinitions/26c128c2-fefa-4340-9bb1-8e081c90ada2`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: registrationdefinitions.ValidateScopedRegistrationDefinitionID,
		},

		"scope": {
			Description:  "The scope at which the Lighthouse Assignment applies to, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333` or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.Any(commonids.ValidateSubscriptionID, commonids.ValidateResourceGroupID),
		},
	}
}
