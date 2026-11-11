// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_custom_keys_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountProjectConnectionCustomKeysResource struct{}

func TestAccCognitiveAccountProjectConnectionCustomKeys_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_custom_keys", "test")
	r := CognitiveAccountProjectConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_keys"),
	})
}

func TestAccCognitiveAccountProjectConnectionCustomKeys_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_custom_keys", "test")
	r := CognitiveAccountProjectConnectionCustomKeysResource{}

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

func TestAccCognitiveAccountProjectConnectionCustomKeys_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_custom_keys", "test")
	r := CognitiveAccountProjectConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("metadata.ApiType").HasValue("Azure"),
			),
		},
		data.ImportStep("custom_keys", "metadata"),
	})
}

func TestAccCognitiveAccountProjectConnectionCustomKeys_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_custom_keys", "test")
	r := CognitiveAccountProjectConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("metadata.ApiType").HasValue("Azure"),
			),
		},
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("metadata.ApiType").HasValue("Custom"),
				check.That(data.ResourceName).Key("custom_keys.%").HasValue("1"),
			),
		},
		data.ImportStep("custom_keys", "metadata"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("metadata.%").HasValue("0"),
				check.That(data.ResourceName).Key("custom_keys.%").HasValue("2"),
			),
		},
		data.ImportStep("custom_keys"),
	})
}

func TestAccCognitiveAccountProjectConnectionCustomKeys_remoteA2A(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_custom_keys", "test")
	r := CognitiveAccountProjectConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.category(data, "RemoteA2A"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_keys"),
	})
}

func TestAccCognitiveAccountProjectConnectionCustomKeys_remoteTool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_custom_keys", "test")
	r := CognitiveAccountProjectConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.category(data, "RemoteTool"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_keys"),
	})
}

func TestAccCognitiveAccountProjectConnectionCustomKeys_importMismatchedAuthTypeFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_custom_keys", "test")
	r := CognitiveAccountProjectConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mixedConnections(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			ResourceName: data.ResourceName,
			ImportState:  true,
			ImportStateIdFunc: func(s *terraform.State) (string, error) {
				rs, ok := s.RootModule().Resources["azurerm_cognitive_account_project_connection_entra_id.wrong"]
				if !ok {
					return "", fmt.Errorf("resource `%s` not found in state", "azurerm_cognitive_account_project_connection_entra_id.wrong")
				}
				return rs.Primary.ID, nil
			},
			ExpectError: regexp.MustCompile("cannot be managed by"),
		},
	})
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := projectconnectionresource.ParseProjectConnectionID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.Cognitive.ProjectConnectionResourceClient.ProjectConnectionsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(helpers.CognitiveAccountProjectConnectionHasExpectedAuthType(resp.Model, projectconnectionresource.ConnectionAuthTypeCustomKeys)), nil
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-pcck-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cognitive_account" "test" {
  name                       = "acctest-cogacc-%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctestaiservices-%[1]d"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_project" "test" {
  name                 = "acctest-proj-%[1]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  location             = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) basic(data acceptance.TestData) string {
	return r.category(data, "CustomKeys")
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) category(data acceptance.TestData, category string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_project_connection_custom_keys" "test" {
  name                         = "acctest-conn-%[2]d"
  cognitive_account_project_id = azurerm_cognitive_account_project.test.id
  category                     = "%[3]s"
  target                       = azurerm_cognitive_account.test.endpoint

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.test.primary_access_key
    secondaryKey = azurerm_cognitive_account.test.secondary_access_key
  }
}
`, r.template(data), data.RandomInteger, category)
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_project_connection_custom_keys" "test" {
  name                         = "acctest-conn-%[2]d"
  cognitive_account_project_id = azurerm_cognitive_account_project.test.id
  category                     = "CustomKeys"
  target                       = azurerm_cognitive_account.test.endpoint

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.test.primary_access_key
    secondaryKey = azurerm_cognitive_account.test.secondary_access_key
  }

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.test.id
    Location   = azurerm_cognitive_account.test.location
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_project_connection_custom_keys" "test" {
  name                         = "acctest-conn-%[2]d"
  cognitive_account_project_id = azurerm_cognitive_account_project.test.id
  category                     = "CustomKeys"
  target                       = "https://updated.example.com/"

  custom_keys = {
    apiKey = azurerm_cognitive_account.test.secondary_access_key
  }

  metadata = {
    ApiType = "Custom"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_project_connection_custom_keys" "import" {
  name                         = azurerm_cognitive_account_project_connection_custom_keys.test.name
  cognitive_account_project_id = azurerm_cognitive_account_project_connection_custom_keys.test.cognitive_account_project_id
  category                     = azurerm_cognitive_account_project_connection_custom_keys.test.category
  target                       = azurerm_cognitive_account_project_connection_custom_keys.test.target
  custom_keys                  = azurerm_cognitive_account_project_connection_custom_keys.test.custom_keys
}
`, r.basic(data))
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) mixedConnections(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_project_connection_entra_id" "wrong" {
  name                         = "acctest-wrong-%[2]d"
  cognitive_account_project_id = azurerm_cognitive_account_project.test.id
  category                     = "Databricks"
  target                       = "https://workspace-%[2]d.databricks.net/"
}
`, r.basic(data), data.RandomInteger)
}
