// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/resources/frontdoor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontDoorCreate,
		Read:   resourceFrontDoorRead,
		Update: resourceFrontDoorUpdate,
		Delete: resourceFrontDoorDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := frontdoors.ParseFrontDoorID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.FrontDoorUpgradeV0ToV1{},
			1: migration.FrontDoorUpgradeV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(6 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(6 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Schema: resourceFrontDoorSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(frontDoorCustomizeDiff),
	}
}
