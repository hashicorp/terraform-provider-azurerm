// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sensitivitylabels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MssqlDatabaseSensitivityLabelResource struct{}

func TestAccMssqlDatabaseSensitivityLabel_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_sensitivity_label", "test")
	r := MssqlDatabaseSensitivityLabelResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMssqlDatabaseSensitivityLabel_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_sensitivity_label", "test")
	r := MssqlDatabaseSensitivityLabelResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMssqlDatabaseSensitivityLabel_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_sensitivity_label", "test")
	r := MssqlDatabaseSensitivityLabelResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMssqlDatabaseSensitivityLabel_informationType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_sensitivity_label", "test")
	r := MssqlDatabaseSensitivityLabelResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.informationType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMssqlDatabaseSensitivityLabel_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_sensitivity_label", "test")
	r := MssqlDatabaseSensitivityLabelResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.informationType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MssqlDatabaseSensitivityLabelResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sensitivitylabels.ParseSensitivityLabelSourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.SensitivityLabelsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (MssqlDatabaseSensitivityLabelResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name        = "acctest-db-%[2]d"
  server_id   = azurerm_mssql_server.test.id
  sample_name = "AdventureWorksLT"
  sku_name    = "Basic"
}
`, MssqlDatabaseResource{}.template(data), data.RandomInteger)
}

func (r MssqlDatabaseSensitivityLabelResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database_sensitivity_label" "test" {
  column_name = "EmailAddress"
  database_id = azurerm_mssql_database.test.id
  schema_name = "SalesLT"
  table_name  = "Customer"
  label_name  = "Confidential"
}
`, r.template(data))
}

func (r MssqlDatabaseSensitivityLabelResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database_sensitivity_label" "import" {
  column_name = azurerm_mssql_database_sensitivity_label.test.column_name
  database_id = azurerm_mssql_database_sensitivity_label.test.database_id
  schema_name = azurerm_mssql_database_sensitivity_label.test.schema_name
  table_name  = azurerm_mssql_database_sensitivity_label.test.table_name
  label_name  = azurerm_mssql_database_sensitivity_label.test.label_name
}
`, r.basic(data))
}

func (r MssqlDatabaseSensitivityLabelResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database_sensitivity_label" "test" {
  column_name         = "EmailAddress"
  database_id         = azurerm_mssql_database.test.id
  schema_name         = "SalesLT"
  table_name          = "Customer"
  label_name          = "Highly Confidential"
  information_type    = "Email"
  information_type_id = "b7732f83-4fc9-4d68-bb8a-2f279b8ba46b"
  label_id            = "bf91e08c-f4f0-478a-b016-25164b2a65ff"
  rank                = "High"
}
`, r.template(data))
}

func (r MssqlDatabaseSensitivityLabelResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database_sensitivity_label" "test" {
  column_name         = "EmailAddress"
  database_id         = azurerm_mssql_database.test.id
  schema_name         = "SalesLT"
  table_name          = "Customer"
  label_name          = "Restricted"
  information_type    = "Contact Information"
  information_type_id = "e1f0e0f7-2b5a-4f52-85b9-8c929b75d35a"
  label_id            = "cf91e08c-f4f0-478a-b016-25164b2a65ff"
  rank                = "Critical"
}
`, r.template(data))
}

func (r MssqlDatabaseSensitivityLabelResource) informationType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database_sensitivity_label" "test" {
  column_name      = "EmailAddress"
  database_id      = azurerm_mssql_database.test.id
  schema_name      = "SalesLT"
  table_name       = "Customer"
  information_type = "Email"
}
`, r.template(data))
}
