package datadog_test

import (
    "testing"
    "context"
)

type DatadogMonitorResource struct{}

func TestAccDatadogMonitor_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
    r := DatadogMonitorResource{}
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

func TestAccDatadogMonitor_requiresImport(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
    r := DatadogMonitorResource{}
    data.ResourceTest(t, r, []resource.TestStep{
        {
            Config: r.basic(data),
            Check: resource.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.RequiresImportErrorStep(r.requiresImport),
    })
}

func TestAccDatadogMonitor_complete(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
    r := DatadogMonitorResource{}
    data.ResourceTest(t, r, []resource.TestStep{
        {
            Config: r.complete(data),
            Check: resource.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
    })
}

func TestAccDatadogMonitor_update(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
    r := DatadogMonitorResource{}
    data.ResourceTest(t, r, []resource.TestStep{
        {
            Config: r.basic(data),
            Check: resource.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
        {
            Config: r.complete(data),
            Check: resource.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
        {
            Config: r.basic(data),
            Check: resource.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
    })
}

func (r DatadogMonitorResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
    id, err := parse.DatadogMonitorID(state.ID)
    if err != nil {
        return nil, err
    }
    resp, err := client.Datadog.MonitorClient.Get(ctx, id.ResourceGroup, id.Name)
    if err != nil {
        if utils.ResponseWasNotFound(resp.Response) {
            return utils.Bool(false), nil
        }
        return nil, fmt.Errorf("retrieving Datadog Monitor %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
    }
    return utils.Bool(true), nil
}

func (r DatadogMonitorResource) template(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-datadog-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DatadogMonitorResource) basic(data acceptance.TestData) string {
    template := r.template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_datadog_monitor" "test" {
  name                = "acctest-dm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func (r DatadogMonitorResource) requiresImport(data acceptance.TestData) string {
    config := r.basic(data)
    return fmt.Sprintf(`
%s

resource "azurerm_datadog_monitor" "import" {
  name                = azurerm_datadog_monitor.test.name
  resource_group_name = azurerm_datadog_monitor.test.resource_group_name
  location            = azurerm_datadog_monitor.test.location
}
`, config)
}

func (r DatadogMonitorResource) complete(data acceptance.TestData) string {
    template := r.template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_datadog_monitor" "test" {
  name                = "acctest-dm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datadog_organization_properties {
    api_key           = ""
    application_key   = ""
    enterprise_app_id = ""
    linking_auth_code = ""
    linking_client_id = ""
    redirect_uri      = ""
  }

  identity {
    type = ""
  }

  sku {
    name = ""
  }

  user_info {
    name          = ""
    email_address = ""
    phone_number  = ""
  }
  monitoring_status = false
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
