// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package user_assigned_identity

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedidentity/resources/user_assigned_identity/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var (
	_ sdk.Resource                   = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
)

type Resource struct{}

func (r Resource) Identity() resourceids.ResourceId {
	return &commonids.UserAssignedIdentityId{}
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.UserAssignedIdentityV0ToV1{},
		},
	}
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateUserAssignedIdentityID
}

func (r Resource) ResourceType() string {
	return "azurerm_user_assigned_identity"
}
