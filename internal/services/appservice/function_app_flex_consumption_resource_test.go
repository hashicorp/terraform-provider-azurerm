// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FunctionAppFlexConsumptionResource struct{}

func TestAccFunctionAppFlexConsumption_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_connectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_stickySettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.stickySettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_connectionStringUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.connectionStringUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.connectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_appSettingsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appSettingsAddKvps(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appSettingsRemoveKvps(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimePython(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.10"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimeNode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.runtimeNode(data, "20"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimeJava(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimeJavaUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.java(data, "17"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.java(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimeDotNet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "8.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimePowerShell(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.powerShell(data, "7.4"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_systemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageUserAssignedIdentity1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_userAssignedIdentityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageUserAssignedIdentity1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.storageUserAssignedIdentity2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func (r FunctionAppFlexConsumptionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseFunctionAppID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r FunctionAppFlexConsumptionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-tf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) connectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-tf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) stickySettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-tf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  app_settings = {
    foo    = "bar"
    secret = "sauce"
    third  = "degree"
  }

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }

  connection_string {
    name  = "Second"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  connection_string {
    name  = "Third"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  sticky_settings {
    app_setting_names       = ["foo", "secret"]
    connection_string_names = ["First", "Third"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) connectionStringUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-tf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  connection_string {
    name  = "AnotherExample"
    value = "some-other-connection-string"
    type  = "Custom"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-tf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  site_config {
    app_command_line                       = "whoami"
    api_definition_url                     = "https://example.com/azure_function_app_def.json"
    application_insights_connection_string = azurerm_application_insights.test.connection_string

    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
      "hostingstart.html",
    ]

    http2_enabled = true

    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    load_balancing_mode      = "LeastResponseTime"
    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"

    scm_ip_restriction {
      ip_address = "10.20.20.20/32"
      name       = "test-scm-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    scm_ip_restriction {
      ip_address = "fd80::/64"
      name       = "test-scm-restriction-v6"
      priority   = 124
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    websockets_enabled                = true
    health_check_path                 = "/health-check"
    health_check_eviction_time_in_min = 7
    worker_count                      = 3

    minimum_tls_version     = "1.2"
    scm_minimum_tls_version = "1.2"

    cors {
      allowed_origins = [
        "https://www.contoso.com",
        "www.contoso.com",
      ]

      support_credentials = true
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) appSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-tf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  app_settings = {
    "tftest" : "tftestvalue"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) appSettingsAddKvps(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-tf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  app_settings = {
    "tftest" : "tftestvalue",
    "tftestkvp1" : "tftestkvpvalue1"
    "tftestkvp2" : "tftestkvpvalue2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) appSettingsRemoveKvps(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-tf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  app_settings = {
    "tftest" : "tftestvalue",
    "tftestkvp1" : "tftestkvpvalue1"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) python(data acceptance.TestData, pythonVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = azurerm_storage_container.test.id
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "python"
  runtime_version             = "%s"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, pythonVersion)
}

func (r FunctionAppFlexConsumptionResource) java(data acceptance.TestData, javaVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "java"
  runtime_version             = "%s"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, javaVersion)
}

func (r FunctionAppFlexConsumptionResource) dotNet(data acceptance.TestData, dotNetVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "dotnet-isolated"
  runtime_version             = "%s"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, dotNetVersion)
}

func (r FunctionAppFlexConsumptionResource) powerShell(data acceptance.TestData, powerShellVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "powershell"
  runtime_version             = "%s"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, powerShellVersion)
}

func (r FunctionAppFlexConsumptionResource) storageSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "SystemAssignedIdentity"
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) storageUserAssignedIdentity1(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s
resource "azurerm_user_assigned_identity" "test1" {
  name                = "acct-uai1-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test1.id
  runtime_name                      = "node"
  runtime_version                   = "20"
  maximum_instance_count            = 50
  instance_memory_in_mb             = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) storageUserAssignedIdentity2(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acct-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                      = "node"
  runtime_version                   = "20"
  maximum_instance_count            = 50
  instance_memory_in_mb             = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) runtimeNode(data acceptance.TestData, nodeVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acct-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                      = "node"
  runtime_version                   = "%s"
  maximum_instance_count            = 50
  instance_memory_in_mb             = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, nodeVersion)
}

func (FunctionAppFlexConsumptionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-LFA-%[1]d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestblobforfc"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2021-04-01"
  expiry = "2024-03-30"

  permissions {
    read    = false
    write   = true
    delete  = false
    list    = false
    add     = false
    create  = false
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "FC1"
}
`, data.RandomInteger, "eastus2", data.RandomString, data.RandomInteger) // location needs to be hardcoded for the moment because flex isn't available in all regions yet and appservice already has location overrides in TC
}
