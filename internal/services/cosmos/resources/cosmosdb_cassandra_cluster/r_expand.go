// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_cluster

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
)

func expandCassandraClusterIdentity(input []any) (*identity.SystemAssigned, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &identity.SystemAssigned{
		Type: expanded.Type,
	}, nil
}

func expandCassandraClusterCertificate(input []any) *[]managedcassandras.Certificate {
	results := make([]managedcassandras.Certificate, 0)

	for _, pem := range input {
		result := managedcassandras.Certificate{
			Pem: pointer.To(pem.(string)),
		}
		results = append(results, result)
	}

	return &results
}

func expandCassandraClusterExternalSeedNode(input []any) *[]managedcassandras.SeedNode {
	results := make([]managedcassandras.SeedNode, 0)

	for _, ipAddress := range input {
		result := managedcassandras.SeedNode{
			IPAddress: pointer.To(ipAddress.(string)),
		}
		results = append(results, result)
	}

	return &results
}
