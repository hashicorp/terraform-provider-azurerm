// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccMssqlDatabaseSensitivityLabel_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_sensitivity_label", "test")
	r := MssqlDatabaseSensitivityLabelResource{}

	checkedFields := map[string]struct{}{
		"sensitivity_label_source": {},
		"column_name":              {},
		"schema_name":              {},
		"table_name":               {},
		"database_name":            {},
		"resource_group_name":      {},
		"server_name":              {},
		"subscription_id":          {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_mssql_database_sensitivity_label.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_mssql_database_sensitivity_label.test", tfjsonpath.New("sensitivity_label_source"), knownvalue.StringExact("current")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_mssql_database_sensitivity_label.test", tfjsonpath.New("column_name"), tfjsonpath.New("column_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_mssql_database_sensitivity_label.test", tfjsonpath.New("schema_name"), tfjsonpath.New("schema_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_mssql_database_sensitivity_label.test", tfjsonpath.New("table_name"), tfjsonpath.New("table_name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_mssql_database_sensitivity_label.test", tfjsonpath.New("database_name"), tfjsonpath.New("database_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_mssql_database_sensitivity_label.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("database_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_mssql_database_sensitivity_label.test", tfjsonpath.New("server_name"), tfjsonpath.New("database_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_mssql_database_sensitivity_label.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("database_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
