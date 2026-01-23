// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databasemigration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestAccDatabaseMigrationProject_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	r := DatabaseMigrationProjectResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basicForResourceIdentity(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValue("azurerm_database_migration_project.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_database_migration_project.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_database_migration_project.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_database_migration_project.test", tfjsonpath.New("service_name"), tfjsonpath.New("service_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
