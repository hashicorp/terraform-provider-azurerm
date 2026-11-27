// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_cluster

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
)

func flattenCassandraClusterCertificate(input *[]managedcassandras.Certificate) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, pointer.From(item.Pem))
	}

	return results
}

func flattenCassandraClusterExternalSeedNode(input *[]managedcassandras.SeedNode) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, pointer.From(item.IPAddress))
	}

	return results
}

func flattenCassandraClusterIdentity(input *identity.SystemAssigned) []any {
	var transform *identity.SystemAssigned

	if input != nil {
		transform = &identity.SystemAssigned{
			Type: input.Type,
		}
		if input.PrincipalId != "" {
			transform.PrincipalId = input.PrincipalId
		}
		if input.TenantId != "" {
			transform.TenantId = input.TenantId
		}
	}

	return identity.FlattenSystemAssigned(transform)
}
