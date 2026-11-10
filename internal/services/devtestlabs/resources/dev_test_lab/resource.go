// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_lab

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/labs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/resources/dev_test_lab/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDevTestLabCreateUpdate,
		Read:   resourceDevTestLabRead,
		Update: resourceDevTestLabCreateUpdate,
		Delete: resourceDevTestLabDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := labs.ParseLabID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DevTestLabUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: devTestLabSchema(),
	}
}
