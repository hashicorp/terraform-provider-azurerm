// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cosmosdbCassandraClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Cassandra Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"delegated_management_subnet_id": {
			Description:  "The ID of the delegated management subnet for this Cassandra Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"default_admin_password": {
			Description:  "The initial admin password for this Cassandra Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"authentication_method": {
			Description: "The authentication method that is used to authenticate clients.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(managedcassandras.AuthenticationMethodCassandra),
			ValidateFunc: validation.StringInSlice([]string{
				string(managedcassandras.AuthenticationMethodNone),
				string(managedcassandras.AuthenticationMethodCassandra),
			}, false),
		},

		"client_certificate_pems": {
			Description: "A list of TLS certificates that is used to authorize client connecting to the Cassandra Cluster.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.IsCert,
			},
		},

		"external_gossip_certificate_pems": {
			Description: "A list of TLS certificates that is used to authorize gossip from unmanaged Cassandra Data Center.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.IsCert,
			},
		},

		"external_seed_node_ip_addresses": {
			Description: "A list of IP Addresses of the seed nodes in the unmanaged Cassandra Data Center which will be added to the seed node lists of all managed nodes.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPv4Address,
			},
		},

		"hours_between_backups": {
			Description: "The number of hours to wait between taking a backup of the Cassandra Cluster.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     24,
		},

		"identity": commonschema.SystemAssignedIdentityOptional(),

		"repair_enabled": {
			Description: "Whether the automatic repair is enabled on the Cassandra Cluster.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"version": {
			Description: "The version of Cassandra that the Cluster converges to run.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "3.11",
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"3.11",
				"4.0",
				"4.1",
				"5.0",
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}
