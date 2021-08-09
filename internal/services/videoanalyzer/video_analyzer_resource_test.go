package videoanalyzer_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/videoanalyzer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VideoAnalyzerResource struct {
}

func TestAccVideoAnalyzer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_video_analyzer", "test")
	r := VideoAnalyzerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("storage_account.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVideoAnalyzer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_video_analyzer", "test")
	r := VideoAnalyzerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("storage_account.#").HasValue("1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccVideoAnalyzer_multipleStorageAccounts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_video_analyzer", "test")
	r := VideoAnalyzerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.multipleAccounts(data),
			ExpectError: regexp.MustCompile("Error: Too many list items"),
		},
	})
}

func TestAccVideoAnalyzer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_video_analyzer", "test")
	r := VideoAnalyzerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func (VideoAnalyzerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VideoAnalyzerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.VideoAnalyzer.VideoAnalyzersClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Video Analyzer %s (resource group: %s): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.PropertiesType != nil), nil
}

func (r VideoAnalyzerResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_video_analyzer" "test" {
  name                = "acctestva%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id                        = azurerm_storage_account.first.id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  depends_on = [
    azurerm_role_assignment.contributor,
    azurerm_role_assignment.reader,
  ]
}
`, template, data.RandomString)
}

func (r VideoAnalyzerResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_video_analyzer" "import" {
  name                = azurerm_video_analyzer.test.name
  location            = azurerm_video_analyzer.test.location
  resource_group_name = azurerm_video_analyzer.test.resource_group_name

  storage_account {
    id                        = azurerm_storage_account.first.id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}
`, template)
}

func (VideoAnalyzerResource) multipleAccounts(data acceptance.TestData) string {
	template := VideoAnalyzerResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_video_analyzer" "test" {
  name                = "acctestva%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id                        = azurerm_storage_account.first.id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  storage_account {
    id                        = azurerm_storage_account.first.id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  depends_on = [
    azurerm_role_assignment.contributor,
    azurerm_role_assignment.reader,
  ]
}
`, template, data.RandomString)
}

func (r VideoAnalyzerResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_video_analyzer" "test" {
  name                = "acctestva%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id                        = azurerm_storage_account.first.id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  depends_on = [
    azurerm_role_assignment.contributor,
    azurerm_role_assignment.reader,
  ]

  tags = {
    label = "test"
  }
}
`, template, data.RandomString)
}

func (VideoAnalyzerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-video-analyzer-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "contributor" {
  scope                = azurerm_storage_account.first.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "reader" {
  scope                = azurerm_storage_account.first.id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_storage_account" "first" {
  name                     = "acctestsa1%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
