// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_auto_upgrade_profile

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetupdatestrategies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 50),
				validation.StringMatch(
					regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`),
					"`name` must start and end with a lowercase letter or number, and can only contain lowercase letters, numbers, and hyphens",
				),
			),
		},

		"kubernetes_fleet_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KubernetesFleetId{}),

		"channel": {
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(autoupgradeprofiles.PossibleValuesForUpgradeChannel(), false),
		},

		"enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
			Default:  true,
		},

		"node_image_selection_type": {
			Optional:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(autoupgradeprofiles.PossibleValuesForAutoUpgradeNodeImageSelectionType(), false),
		},

		"update_strategy_id": commonschema.ResourceIDReferenceOptionalForceNew(&fleetupdatestrategies.UpdateStrategyId{}),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
