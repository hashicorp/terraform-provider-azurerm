// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managedredis

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
)

func TestClusteringPolicyRequiresDatabaseRecreation(t *testing.T) {
	testCases := []struct {
		name           string
		oldPolicy      redisenterprise.ClusteringPolicy
		newPolicy      redisenterprise.ClusteringPolicy
		expectedResult bool
	}{
		{
			name:      "NoCluster to OSSCluster",
			oldPolicy: redisenterprise.ClusteringPolicyNoCluster,
			newPolicy: redisenterprise.ClusteringPolicyOSSCluster,
		},
		{
			name:      "NoCluster to EnterpriseCluster",
			oldPolicy: redisenterprise.ClusteringPolicyNoCluster,
			newPolicy: redisenterprise.ClusteringPolicyEnterpriseCluster,
		},
		{
			name:           "OSSCluster to NoCluster",
			oldPolicy:      redisenterprise.ClusteringPolicyOSSCluster,
			newPolicy:      redisenterprise.ClusteringPolicyNoCluster,
			expectedResult: true,
		},
		{
			name:           "EnterpriseCluster to NoCluster",
			oldPolicy:      redisenterprise.ClusteringPolicyEnterpriseCluster,
			newPolicy:      redisenterprise.ClusteringPolicyNoCluster,
			expectedResult: true,
		},
		{
			name:           "OSSCluster to EnterpriseCluster",
			oldPolicy:      redisenterprise.ClusteringPolicyOSSCluster,
			newPolicy:      redisenterprise.ClusteringPolicyEnterpriseCluster,
			expectedResult: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualResult := clusteringPolicyRequiresDatabaseRecreation(string(testCase.oldPolicy), string(testCase.newPolicy))
			if actualResult != testCase.expectedResult {
				t.Fatalf("expected %t, got %t", testCase.expectedResult, actualResult)
			}
		})
	}
}
