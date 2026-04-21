// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestAccCassandraSequential(t *testing.T) {
	if os.Getenv("CASSANDRA_VERY_EXPENSIVE_TEST_ENABLED") != "true" {
		t.Skip("Skipping due to massive cost of this resource test.")
	}
	// NOTE: this is a combined test rather than separate split out tests due to all below TCs sharing same sp
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"cluster": {
			"basic":          testAccCassandraCluster_basic,
			"complete":       testAccCassandraCluster_complete,
			"update":         testAccCassandraCluster_update,
			"requiresImport": testAccCassandraCluster_requiresImport,
		},
		"dataCenter": {
			"basic":     testAccCassandraDatacenter_basic,
			"update":    testAccCassandraDatacenter_update,
			"updateSku": testAccCassandraDatacenter_updateSku,
		},
	})
}
