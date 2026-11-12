// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_consumer_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/consumergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub_consumer_group/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var (
	_ sdk.ResourceWithIdentity       = Resource{}
	_ sdk.Resource                   = Resource{}
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &consumergroups.ConsumerGroupId{}
}

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_eventhub_consumer_group"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return consumergroups.ValidateConsumerGroupID
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.ConsumerGroupsV0ToV1{},
		},
	}
}
