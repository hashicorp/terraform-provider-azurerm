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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"delegated_management_subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"default_admin_password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"authentication_method": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(managedcassandras.AuthenticationMethodCassandra),
			ValidateFunc: validation.StringInSlice([]string{
				string(managedcassandras.AuthenticationMethodNone),
				string(managedcassandras.AuthenticationMethodCassandra),
			}, false),
		},

		"client_certificate_pems": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.IsCert,
			},
		},

		"external_gossip_certificate_pems": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.IsCert,
			},
		},

		"external_seed_node_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPv4Address,
			},
		},

		"hours_between_backups": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  24,
		},

		"identity": commonschema.SystemAssignedIdentityOptional(),

		"repair_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "3.11",
			ForceNew: true,
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
