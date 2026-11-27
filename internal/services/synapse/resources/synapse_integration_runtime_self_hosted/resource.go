// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_integration_runtime_self_hosted

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_integration_runtime_self_hosted/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseIntegrationRuntimeSelfHostedCreateUpdate,
		Read:   resourceSynapseIntegrationRuntimeSelfHostedRead,
		Update: resourceSynapseIntegrationRuntimeSelfHostedCreateUpdate,
		Delete: resourceSynapseIntegrationRuntimeSelfHostedDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := integrationruntimes.ParseIntegrationRuntimeID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SynapseIntegrationRuntimeSelfHostedV0ToV1{},
			1: migration.SynapseIntegrationRuntimeSelfHostedV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: synapseIntegrationRuntimeSelfHostedSchema(),
	}
}
