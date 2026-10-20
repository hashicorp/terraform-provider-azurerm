// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_managed_private_endpoint

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/managedprivateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_cluster_managed_private_endpoint/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoClusterManagedPrivateEndpointCreateUpdate,
		Read:   resourceKustoClusterManagedPrivateEndpointRead,
		Update: resourceKustoClusterManagedPrivateEndpointCreateUpdate,
		Delete: resourceKustoClusterManagedPrivateEndpointDelete,

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoManagedPrivateEndpointV0ToV1{},
			1: migration.KustoManagedPrivateEndpointV1ToV2{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := managedprivateendpoints.ParseManagedPrivateEndpointID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: kustoClusterManagedPrivateEndpointSchema(),
	}
}
