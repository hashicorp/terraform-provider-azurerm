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

type LogicAppTriggerHttpRequestResource struct {
}

func TestAccLogicAppTriggerHttpRequest_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")
	r := LogicAppTriggerHttpRequestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("schema").HasValue("{}"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerHttpRequest_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")
	r := LogicAppTriggerHttpRequestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_logic_app_trigger_http_request"),
		},
	})
}

func TestAccLogicAppTriggerHttpRequest_fullSchema(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")
	r := LogicAppTriggerHttpRequestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.fullSchema(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("schema").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerHttpRequest_method(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")
	r := LogicAppTriggerHttpRequestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.method(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("method").HasValue("PUT"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerHttpRequest_relativePath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")
	r := LogicAppTriggerHttpRequestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.relativePath(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("method").HasValue("POST"),
				check.That(data.ResourceName).Key("relative_path").HasValue("customers/{id}"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerHttpRequest_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")
	r := LogicAppTriggerHttpRequestResource{}

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
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(LogicAppWorkflowResource{}),
			),
		},
		{
			Config:             r.basic(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func (LogicAppTriggerHttpRequestResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	return triggerExists(ctx, clients, state)
}

func (r LogicAppTriggerHttpRequestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = azurerm_logic_app_workflow.test.id
  schema       = "{}"
}
`, r.template(data))
}

func (r LogicAppTriggerHttpRequestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "import" {
  name         = azurerm_logic_app_trigger_http_request.test.name
  logic_app_id = azurerm_logic_app_trigger_http_request.test.logic_app_id
  schema       = azurerm_logic_app_trigger_http_request.test.schema
}
`, r.template(data))
}

func (r LogicAppTriggerHttpRequestResource) fullSchema(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = azurerm_logic_app_workflow.test.id

  schema = <<SCHEMA
{
    "type": "object",
    "properties": {
        "hello": {
            "type": "string"
        }
    }
}
SCHEMA

}
`, r.template(data))
}

func (r LogicAppTriggerHttpRequestResource) method(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = azurerm_logic_app_workflow.test.id
  schema       = "{}"
  method       = "PUT"
}
`, r.template(data))
}

func (r LogicAppTriggerHttpRequestResource) relativePath(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name          = "some-http-trigger"
  logic_app_id  = azurerm_logic_app_workflow.test.id
  schema        = "{}"
  method        = "POST"
  relative_path = "customers/{id}"
}
`, r.template(data))
}

func (LogicAppTriggerHttpRequestResource) template(data acceptance.TestData) string {
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
