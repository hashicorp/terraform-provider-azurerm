// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobcredentials"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlJobCredentialResource struct{}

func TestAccMsSqlJobCredential_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_credential", "test")
	r := MsSqlJobCredentialResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccMsSqlJobCredential_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_credential", "test")
	r := MsSqlJobCredentialResource{}

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

func TestAccMsSqlJobCredential_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_credential", "test")
	r := MsSqlJobCredentialResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccMsSqlJobCredential_writeOnlyPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_credential", "test")
	r := MsSqlJobCredentialResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.writeOnlyPassword(data, "secret", 1),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("password_wo_version"),
			{
				Config: r.writeOnlyPassword(data, "secretUpdate", 2),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("password_wo_version"),
		},
	})
}

func TestAccMsSqlJobCredential_updateToWriteOnlyPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_credential", "test")
	r := MsSqlJobCredentialResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("password"),
			{
				Config: r.writeOnlyPassword(data, "secret", 1),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("password", "password_wo_version"),
			{
				Config: r.basic(data),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("password"),
		},
	})
}

func (MsSqlJobCredentialResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := jobcredentials.ParseCredentialID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.JobCredentialsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("%s does not exist", *id)
		}
		return nil, fmt.Errorf("reading %s: %v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MsSqlJobCredentialResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_credential" "test" {
  name         = "acctestmssqljobcredential%[2]d"
  job_agent_id = azurerm_mssql_job_agent.test.id
  username     = "test"
  password     = "test"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlJobCredentialResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_credential" "import" {
  name         = azurerm_mssql_job_credential.test.name
  job_agent_id = azurerm_mssql_job_agent.test.id
  username     = "test"
  password     = "test"
}
`, r.basic(data))
}

func (r MsSqlJobCredentialResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_credential" "test" {
  name         = "acctestmssqljobcredential%[2]d"
  job_agent_id = azurerm_mssql_job_agent.test.id
  username     = "test1"
  password     = "test1"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlJobCredentialResource) writeOnlyPassword(data acceptance.TestData, secret string, version int) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_mssql_job_credential" "test" {
  name                = "acctestmssqljobcredential%[3]d"
  job_agent_id        = azurerm_mssql_job_agent.test.id
  username            = "test"
  password_wo         = ephemeral.azurerm_key_vault_secret.test.value
  password_wo_version = %[4]d
}
`, r.template(data), acceptance.WriteOnlyKeyVaultSecretTemplate(data, secret), data.RandomInteger, version)
}

func (MsSqlJobCredentialResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-jobcredential-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestmssqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dministr4t0r"
  administrator_login_password = "superSecur3!!!"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctestmssqldb%[1]d"
  server_id = azurerm_mssql_server.test.id
  collation = "SQL_Latin1_General_CP1_CI_AS"
  sku_name  = "S1"
}

resource "azurerm_mssql_job_agent" "test" {
  name        = "acctestmssqljobagent%[1]d"
  location    = azurerm_resource_group.test.location
  database_id = azurerm_mssql_database.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}
