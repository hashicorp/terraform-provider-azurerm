package cosmos_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestAccCassandraSequential(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to all below TCs sharing same sp
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"cluster": {
			"basic":          testAccCassandraCluster_basic,
			"complete":       testAccCassandraCluster_complete,
			"update":         testAccCassandraCluster_update,
			"requiresImport": testAccCassandraCluster_requiresImport,
		},
		"dataCenter": {
			"basic":  testAccCassandraDatacenter_basic,
			"update": testAccCassandraDatacenter_update,
		},
	})
}
