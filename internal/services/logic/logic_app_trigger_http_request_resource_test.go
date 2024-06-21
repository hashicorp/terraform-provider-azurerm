// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogicAppTriggerHttpRequestResource struct{}

func TestAccLogicAppTriggerHttpRequest_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")
	r := LogicAppTriggerHttpRequestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fullSchema(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.method(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("method").HasValue("PUT"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerHttpRequest_callbackUrl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")
	r := LogicAppTriggerHttpRequestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.method(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("callback_url").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerHttpRequest_relativePath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")
	r := LogicAppTriggerHttpRequestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.relativePath(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func TestAccLogicAppTriggerHttpRequest_workflowWithISE(t *testing.T) {
	t.Skip("skip as Integration Service Environment is being deprecated")

	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")
	r := LogicAppTriggerHttpRequestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.workflowWithISE(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("schema").HasValue("{}"),
			),
		},
		data.ImportStep(),
	})
}

func (LogicAppTriggerHttpRequestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
`, r.basic(data))
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

func (r LogicAppTriggerHttpRequestResource) workflowWithISE(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = azurerm_logic_app_workflow.test.id
  schema       = "{}"
}
`, LogicAppWorkflowResource{}.integrationServiceEnvironment(data))
}
