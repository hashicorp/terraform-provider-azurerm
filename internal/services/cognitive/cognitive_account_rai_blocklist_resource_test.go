// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/raiblocklists"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveRaiBlocklistTestResource struct{}

func TestAccCognitiveRaiBlocklistSequential(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"raiBlocklist": {
			"basic":          TestAccCognitiveRaiBlocklist_basic,
			"requiresImport": TestAccCognitiveRaiBlocklist_requiresImport,
			"complete":       TestAccCognitiveRaiBlocklist_complete,
			"update":         TestAccCognitiveRaiBlocklist_update,
		},
	})
}

func TestAccCognitiveRaiBlocklist_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_rai_blocklist", "test")
	r := CognitiveRaiBlocklistTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveRaiBlocklist_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_rai_blocklist", "test")
	r := CognitiveRaiBlocklistTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCognitiveRaiBlocklist_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_rai_blocklist", "test")
	r := CognitiveRaiBlocklistTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveRaiBlocklist_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_rai_blocklist", "test")
	r := CognitiveRaiBlocklistTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (c CognitiveRaiBlocklistTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := raiblocklists.ParseRaiBlocklistID(state.ID)
	if err != nil {
		return nil, err
	}

	existing, err := clients.Cognitive.RaiBlocklistsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(existing.Model != nil), nil
}

func (c CognitiveRaiBlocklistTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "OpenAI"
  sku_name            = "S0"
}
		`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (c CognitiveRaiBlocklistTestResource) basic(data acceptance.TestData) string {
	template := c.template(data)
	return fmt.Sprintf(`
		%s
resource "azurerm_cognitive_account_rai_blocklist" "test" {
  name                 = "acctest-crb-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  description          = "Acceptance test data new azurerm resource"
}
		`, template, data.RandomInteger)
}

func (c CognitiveRaiBlocklistTestResource) requiresImport(data acceptance.TestData) string {
	config := c.basic(data)
	return fmt.Sprintf(`
		%s
resource "azurerm_cognitive_account_rai_blocklist" "import" {
  name                 = azurerm_cognitive_account_rai_blocklist.test.name
  cognitive_account_id = azurerm_cognitive_account.test.id
  description          = "Acceptance test data new azurerm resource"
}
		`, config)
}

func (c CognitiveRaiBlocklistTestResource) complete(data acceptance.TestData) string {
	template := c.template(data)
	return fmt.Sprintf(`
		%s
resource "azurerm_cognitive_account_rai_blocklist" "test" {
  name                 = "acctest-crb-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  description          = "Acceptance test data new azurerm resource"
}
		`, template, data.RandomInteger)
}

func (c CognitiveRaiBlocklistTestResource) update(data acceptance.TestData) string {
	template := c.template(data)
	return fmt.Sprintf(`
		%s
resource "azurerm_cognitive_account_rai_blocklist" "test" {
  name                 = "acctest-crb-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  description          = "Acceptance test data new azurerm resource updated"
}
		`, template, data.RandomInteger)
}
