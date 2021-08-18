package hdinsight_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestAccHDInsightCluster(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"hadoop": {
			"securityProfile": testAccHDInsightHadoopCluster_securityProfile,
		},
		"hbase": {
			"securityProfile": testAccHDInsightHBaseCluster_securityProfile,
		},
		"interactiveQuery": {
			"securityProfile": testAccHDInsightInteractiveQueryCluster_securityProfile,
		},
		"kafka": {
			"securityProfile": testAccHDInsightKafkaCluster_securityProfile,
		},
		"spark": {
			"securityProfile": testAccHDInsightSparkCluster_securityProfile,
		},
	})
}
