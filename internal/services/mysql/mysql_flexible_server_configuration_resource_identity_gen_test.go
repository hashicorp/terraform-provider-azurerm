// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mysql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccMysqlFlexibleServerConfiguration_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_configuration", "test")
	r := MysqlFlexibleServerConfigurationResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":      {},
		"flexible_server_name": {},
		"name":                 {},
		"resource_group_name":  {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.characterSetServer(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_mysql_flexible_server_configuration.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_mysql_flexible_server_configuration.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_mysql_flexible_server_configuration.test", tfjsonpath.New("flexible_server_name"), tfjsonpath.New("server_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_mysql_flexible_server_configuration.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_mysql_flexible_server_configuration.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
