// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_spark_pool

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_spark_pool/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseSparkPoolCreate,
		Read:   resourceSynapseSparkPoolRead,
		Update: resourceSynapseSparkPoolUpdate,
		Delete: resourceSynapseSparkPoolDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := bigdatapools.ParseBigDataPoolID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SynapseSparkPoolV0ToV1{},
		}),

		SchemaVersion: 1,
		Schema:        synapseSparkPoolSchema(),
	}
}
