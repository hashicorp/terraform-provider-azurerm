// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dynatrace_tag_rules_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type MonitorsResource struct {
	dynatraceInfo dynatraceInfo
}

func (r MonitorsResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_dynatrace_monitor" "test" {
  name                     = "acctestacc%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  marketplace_subscription = "Active"

  identity {
    type = "SystemAssigned"
  }

  user {
    first_name   = "%s"
    last_name    = "%s"
    email        = "%s"
    phone_number = "%s"
    country      = "%s"
  }

  plan {
    usage_type    = "COMMITTED"
    billing_cycle = "MONTHLY"
    plan          = "azureportalintegration_privatepreview@TIDgmz7xq9ge3py"
  }

  tags = {
    environment = "Dev"
  }
}
`, template, data.RandomInteger, r.dynatraceInfo.UserFirstName, r.dynatraceInfo.UserLastName, r.dynatraceInfo.UserEmail, r.dynatraceInfo.UserPhoneNumber, r.dynatraceInfo.UserCountry)
}

func (r MonitorsResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, "eastus2euap", data.RandomString)
}

type dynatraceInfo struct {
	UserCountry     string
	UserEmail       string
	UserFirstName   string
	UserLastName    string
	UserPhoneNumber string
}
