// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/communicationservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/resources/communication_service/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
)

type Resource struct{}

func (Resource) Identity() resourceids.ResourceId {
	return &communicationservices.CommunicationServiceId{}
}

func (Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.ServiceV0ToV1{},
		},
	}
}

func (Resource) ResourceType() string {
	return "azurerm_communication_service"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return communicationservices.ValidateCommunicationServiceID
}
