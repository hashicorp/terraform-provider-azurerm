// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_datacenter

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cosmosdbCassandraDatacenterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Cassandra Datacenter.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cassandra_cluster_id": {
			Description:  "The ID of the Cassandra Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedcassandras.ValidateCassandraClusterID,
		},

		"location": commonschema.Location(),

		"delegated_management_subnet_id": {
			Description:  "The ID of the delegated management subnet for this Cassandra Datacenter.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"backup_storage_customer_key_uri": {
			Description:  "The key URI of the customer key to use for the encryption of the backup Storage Account.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"base64_encoded_yaml_fragment": {
			Description:  "The fragment of the cassandra.yaml configuration file to be included in the cassandra.yaml for all nodes in this Cassandra Datacenter. The fragment should be Base64 encoded and only a subset of keys is allowed.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"disk_sku": {
			Description:  "The Disk SKU that is used for this Cassandra Datacenter.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "P30",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"managed_disk_customer_key_uri": {
			Description:  "The key URI of the customer key to use for the encryption of the Managed Disk.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"node_count": {
			Description:  "The number of nodes the Cassandra Datacenter should have. The number should be equal or greater than `3`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      3,
			ValidateFunc: validation.IntAtLeast(3),
		},

		"disk_count": {
			Description:  "The number of p30 disks that are attached to each node.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 10),
		},

		"availability_zones_enabled": {
			Description: "Whether Availability Zones are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"sku_name": {
			Description:  "The selected SKU.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "Standard_E16s_v5",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"seed_node_ip_addresses": {
			Description: "A list of IP Addresses for the seed nodes in this Cassandra Datacenter.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
