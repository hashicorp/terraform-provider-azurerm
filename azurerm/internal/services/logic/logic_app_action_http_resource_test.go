package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

type LogicAppActionHttpResource struct {
}

func TestAccLogicAppActionHttp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_http", "test")
	r := LogicAppActionHttpResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppActionHttp_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_http", "test")
	r := LogicAppActionHttpResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_logic_app_action_http"),
		},
	})
}

func TestAccLogicAppActionHttp_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_http", "test")
	r := LogicAppActionHttpResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.headers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppActionHttp_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_http", "test")
	r := LogicAppActionHttpResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			// delete it
			Config: r.template(data),
		},
		{
			Config:             r.basic(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccLogicAppActionHttp_runAfter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_http", "test")
	r := LogicAppActionHttpResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.runAfterCondition(data, "Succeeded"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.runAfterCondition(data, "Failed"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (LogicAppActionHttpResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	return actionExists(ctx, clients, state)
}

func (r LogicAppActionHttpResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "test" {
  name         = "action%d"
  logic_app_id = azurerm_logic_app_workflow.test.id
  method       = "GET"
  uri          = "http://example.com/hello"
}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppActionHttpResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "import" {
  name         = azurerm_logic_app_action_http.test.name
  logic_app_id = azurerm_logic_app_action_http.test.logic_app_id
  method       = azurerm_logic_app_action_http.test.method
  uri          = azurerm_logic_app_action_http.test.uri
}
`, r.basic(data))
}

func (r LogicAppActionHttpResource) headers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "test" {
  name         = "action%d"
  logic_app_id = azurerm_logic_app_workflow.test.id
  method       = "GET"
  uri          = "http://example.com/hello"

  headers = {
    "Hello"     = "World"
    "Something" = "New"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppActionHttpResource) runAfterCondition(data acceptance.TestData, condition string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "testp1" {
  name         = "action%dp1"
  logic_app_id = azurerm_logic_app_workflow.test.id
  method       = "GET"
  uri          = "http://example.com/hello"
}

resource "azurerm_logic_app_action_http" "testp2" {
  name         = "action%dp2"
  logic_app_id = azurerm_logic_app_workflow.test.id
  method       = "GET"
  uri          = "http://example.com/hello"
}

resource "azurerm_logic_app_action_http" "test" {
  name         = "action%d"
  logic_app_id = azurerm_logic_app_workflow.test.id
  method       = "GET"
  uri          = "http://example.com/hello"
  run_after {
    action_name   = azurerm_logic_app_action_http.testp1.name
    action_result = "%s"
  }
  run_after {
    action_name   = azurerm_logic_app_action_http.testp2.name
    action_result = "%s"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, condition, condition)
}

func (LogicAppActionHttpResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
