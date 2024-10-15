// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type CloudVmClusterDataSource struct{}

func TestCloudVmClusterDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.CloudVmClusterDataSource{}.ResourceType(), "test")
	r := CloudVmClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("cloud_exadata_infrastructure_id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("cpu_core_count").Exists(),
			),
		},
	})
}

func (d CloudVmClusterDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_cloud_vm_cluster" "test" {
  name                = azurerm_oracle_cloud_vm_cluster.test.name
  resource_group_name = azurerm_oracle_cloud_vm_cluster.test.resource_group_name
}
`, CloudVmClusterResource{}.basic(data))
}
