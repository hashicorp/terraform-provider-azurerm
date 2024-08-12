// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"testing"
)

type DBServersDataSource struct{}

func TestDBServersDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracledatabase_db_servers", "test")
	r := DBServersDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("db_servers.0.compartment_id").Exists(),
			),
		},
	})
}

func (d DBServersDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

data "azurerm_oracledatabase_db_servers" "test" {
	resource_group_name = "PowerShellTestRg"
	cloud_exadata_infrastructure_name = "OFake_PowerShellTestExaInfra"
}
`, d.template(data))
}

func (d DBServersDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azurerm" {
	features {}
}

data "azurerm_client_config" "current" {}

`, ExadataInfraResource{}.basic(data))
}
