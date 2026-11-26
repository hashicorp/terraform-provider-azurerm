// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_db_servers_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type ExadataInfraResource struct{}

func (a ExadataInfraResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_oracle_exadata_infrastructure" "test" {
  name                = "OFakeacctest%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  compute_count       = "2"
  display_name        = "OFakeacctest%[2]d"
  shape               = "Exadata.X9M"
  storage_count       = "3"
  zones               = local.zones
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExadataInfraResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  zones = ["3"]
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
