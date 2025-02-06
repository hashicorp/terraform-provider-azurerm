// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type ExadataInfraDataSource struct{}

func TestExadataInfraDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExadataInfraDataSource{}.ResourceType(), "test")
	r := ExadataInfraDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("compute_count").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("shape").Exists(),
				check.That(data.ResourceName).Key("storage_count").Exists(),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
			),
		},
	})
}

func (d ExadataInfraDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_exadata_infrastructure" "test" {
  name                = azurerm_oracle_exadata_infrastructure.test.name
  resource_group_name = azurerm_oracle_exadata_infrastructure.test.resource_group_name
}
`, ExadataInfraResource{}.basic(data))
}
