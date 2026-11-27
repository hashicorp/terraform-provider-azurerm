// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_managed_private_endpoint

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_managed_private_endpoint/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithStateMigration = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_stream_analytics_managed_private_endpoint"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return privateendpoints.ValidatePrivateEndpointID
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsManagedPrivateEndpointV0ToV1{},
		},
	}
}
