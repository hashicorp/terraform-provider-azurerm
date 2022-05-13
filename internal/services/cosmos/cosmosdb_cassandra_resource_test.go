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
			"requiresImport": testAccCassandraCluster_requiresImport,
			"tags":           testAccCassandraCluster_tags,
		},
		"dataCenter": {
			"basic":          testAccCassandraDatacenter_basic,
			"requiresImport": testAccCassandraDatacenter_update,
		},
	})
}
