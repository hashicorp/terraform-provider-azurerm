// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_datacenter

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a Cassandra Datacenter.",
		Create:      resourceCassandraDatacenterCreate,
		Read:        resourceCassandraDatacenterRead,
		Update:      resourceCassandraDatacenterUpdate,
		Delete:      resourceCassandraDatacenterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := managedcassandras.ParseDataCenterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: cosmosdbCassandraDatacenterSchema(),
	}
}
