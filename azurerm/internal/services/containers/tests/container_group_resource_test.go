package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMContainerGroup_SystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_SystemAssignedIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.identity_ids.#", "0"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
				),
			},
			data.ImportStep("identity.0.principal_id"),
		},
	})
}

func TestAccAzureRMContainerGroup_UserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_UserAssignedIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "UserAssigned"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.identity_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.principal_id", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerGroup_multipleAssignedIdentities(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_MultipleAssignedIdentities(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned, UserAssigned"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.identity_ids.#", "1"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
				),
			},
			data.ImportStep("identity.0.principal_id"),
		},
	})
}

func TestAccAzureRMContainerGroup_imageRegistryCredentials(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_imageRegistryCredentials(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.0.server", "hub.docker.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.0.username", "yourusername"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.0.password", "yourpassword"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.1.server", "mine.acr.io"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.1.username", "acrusername"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.1.password", "acrpassword"),
				),
			},
			data.ImportStep(
				"image_registry_credential.0.password",
				"image_registry_credential.1.password",
			),
		},
	})
}

func TestAccAzureRMContainerGroup_imageRegistryCredentialsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_imageRegistryCredentials(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.0.server", "hub.docker.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.0.username", "yourusername"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.0.password", "yourpassword"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.1.server", "mine.acr.io"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.1.username", "acrusername"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.1.password", "acrpassword"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.ports.#", "1"),
				),
			},
			{
				Config: testAccAzureRMContainerGroup_imageRegistryCredentialsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.0.server", "hub.docker.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.0.username", "updatedusername"),
					resource.TestCheckResourceAttr(data.ResourceName, "image_registry_credential.0.password", "updatedpassword"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.ports.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_logTypeUnset(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_logTypeUnset(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics.0.log_analytics.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics.0.log_analytics.0.log_type", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics.0.log_analytics.0.metadata.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics.0.log_analytics.0.workspace_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics.0.log_analytics.0.workspace_key"),
				),
			},
			data.ImportStep("diagnostics.0.log_analytics.0.workspace_key"),
		},
	})
}

func TestAccAzureRMContainerGroup_linuxBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_linuxBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.ports.#", "1"),
				),
			},
			data.ImportStep(
				"image_registry_credential.0.password",
				"image_registry_credential.1.password",
			),
		},
	})
}

func TestAccAzureRMContainerGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_linuxBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMContainerGroup_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_container_group"),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_linuxBasicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_linuxBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "container.#", "1"),
				),
			},
			{
				Config: testAccAzureRMContainerGroup_linuxBasicUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "container.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.ports.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_linuxComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_linuxComplete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.ports.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.commands.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.commands.0", "/bin/bash"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.commands.1", "-c"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.commands.2", "ls"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.environment_variables.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.environment_variables.foo", "bar"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.environment_variables.foo1", "bar1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.secure_environment_variables.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.secure_environment_variables.secureFoo", "secureBar"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.secure_environment_variables.secureFoo1", "secureBar1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.gpu.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.gpu.0.count", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.gpu.0.sku", "K80"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.volume.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.volume.0.mount_path", "/aci/logs"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.volume.0.name", "logs"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.volume.0.read_only", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(data.ResourceName, "restart_policy", "OnFailure"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics.0.log_analytics.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics.0.log_analytics.0.log_type", "ContainerInsights"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics.0.log_analytics.0.metadata.%", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics.0.log_analytics.0.workspace_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics.0.log_analytics.0.workspace_key"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.exec.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.exec.0", "cat"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.exec.1", "/tmp/healthy"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.http_get.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.initial_delay_seconds", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.period_seconds", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.failure_threshold", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.success_threshold", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.timeout_seconds", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.failure_threshold", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.http_get.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.http_get.0.path", "/"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.http_get.0.port", "443"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.http_get.0.scheme", "Http"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.initial_delay_seconds", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.period_seconds", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.success_threshold", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.timeout_seconds", "1"),
				),
			},
			data.ImportStep(
				"container.0.volume.0.storage_account_key",
				"container.0.secure_environment_variables.%",
				"container.0.secure_environment_variables.secureFoo",
				"container.0.secure_environment_variables.secureFoo1",
				"diagnostics.0.log_analytics.0.workspace_key",
			),
		},
	})
}

func TestAccAzureRMContainerGroup_virtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_virtualNetwork(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckNoResourceAttr(data.ResourceName, "dns_label_name"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "identity"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.ports.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_address_type", "Private"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "dns_config.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_windowsBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_windowsBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.ports.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerGroup_windowsComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_windowsComplete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.ports.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.commands.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.commands.0", "cmd.exe"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.commands.1", "echo"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.commands.2", "hi"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.environment_variables.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.environment_variables.foo", "bar"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.environment_variables.foo1", "bar1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.secure_environment_variables.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.secure_environment_variables.secureFoo", "secureBar"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.secure_environment_variables.secureFoo1", "secureBar1"),
					resource.TestCheckResourceAttr(data.ResourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(data.ResourceName, "restart_policy", "Never"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics.0.log_analytics.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics.0.log_analytics.0.log_type", "ContainerInsights"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics.0.log_analytics.0.metadata.%", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics.0.log_analytics.0.workspace_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics.0.log_analytics.0.workspace_key"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.exec.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.exec.0", "cat"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.exec.1", "/tmp/healthy"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.http_get.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.initial_delay_seconds", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.period_seconds", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.failure_threshold", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.success_threshold", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.readiness_probe.0.timeout_seconds", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.failure_threshold", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.http_get.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.http_get.0.path", "/"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.http_get.0.port", "443"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.http_get.0.scheme", "Http"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.initial_delay_seconds", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.period_seconds", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.success_threshold", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container.0.liveness_probe.0.timeout_seconds", "1"),
				),
			},
			data.ImportStep(
				"container.0.secure_environment_variables.%",
				"container.0.secure_environment_variables.secureFoo",
				"container.0.secure_environment_variables.secureFoo1",
				"diagnostics.0.log_analytics.0.workspace_key",
			),
		},
	})
}

func TestAccAzureRMContainerGroup_withPrivateEmpty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_withPrivateEmpty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"container.0.secure_environment_variables.PRIVATE_VALUE",
			),
		},
	})
}

func TestAccAzureRMContainerGroup_gitRepoVolume(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_gitRepoVolume(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMContainerGroup_SystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 80
      protocol = "TCP"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMContainerGroup_UserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  name = "acctest%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 80
      protocol = "TCP"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = ["${azurerm_user_assigned_identity.test.id}"]
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMContainerGroup_MultipleAssignedIdentities(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  name = "acctest%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 80
      protocol = "TCP"
    }
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = ["${azurerm_user_assigned_identity.test.id}"]
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMContainerGroup_linuxBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 80
      protocol = "TCP"
    }
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMContainerGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMContainerGroup_linuxBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_group" "import" {
  name                = "${azurerm_container_group.test.name}"
  location            = "${azurerm_container_group.test.location}"
  resource_group_name = "${azurerm_container_group.test.resource_group_name}"
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 80
      protocol = "TCP"
    }
  }

  tags = {
    environment = "Testing"
  }
}
`, template)
}

func testAccAzureRMContainerGroup_imageRegistryCredentials(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 5443
      protocol = "UDP"
    }
  }

  image_registry_credential {
    server   = "hub.docker.com"
    username = "yourusername"
    password = "yourpassword"
  }

  image_registry_credential {
    server   = "mine.acr.io"
    username = "acrusername"
    password = "acrpassword"
  }

  container {
    name   = "sidecar"
    image  = "microsoft/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "0.5"
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMContainerGroup_imageRegistryCredentialsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"

    ports {
      port = 80
    }
  }

  image_registry_credential {
    server   = "hub.docker.com"
    username = "updatedusername"
    password = "updatedpassword"
  }

  container {
    name   = "sidecar"
    image  = "microsoft/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "0.5"
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMContainerGroup_logTypeUnset(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port = 80
    }
  }

  diagnostics {
    log_analytics {
      workspace_id  = azurerm_log_analytics_workspace.test.workspace_id
      workspace_key = azurerm_log_analytics_workspace.test.primary_shared_key
    }
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMContainerGroup_linuxBasicUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"

    ports {
      port = 80
    }

    ports {
      port     = 5443
      protocol = "UDP"
    }
  }

  container {
    name   = "sidecar"
    image  = "microsoft/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "0.5"
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMContainerGroup_virtualNetwork(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "testvnet"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_network_profile" "test" {
  name                = "testnetprofile"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  container_network_interface {
    name = "testcnic"

    ip_configuration {
      name      = "testipconfig"
      subnet_id = azurerm_subnet.test.id
    }
  }
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Private"
  network_profile_id  = azurerm_network_profile.test.id
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port = 80
    }
  }
  dns_config {
    nameservers    = ["reddog.microsoft.com", "somecompany.somedomain"]
    options        = ["one:option", "two:option", "red:option", "blue:option"]
    search_domains = ["default.svc.cluster.local."]
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMContainerGroup_windowsBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "public"
  os_type             = "windows"

  container {
    name   = "windowsservercore"
    image  = "microsoft/iis:windowsservercore"
    cpu    = "2.0"
    memory = "3.5"

    ports {
      port     = 80
      protocol = "TCP"
    }

    ports {
      port     = 443
      protocol = "TCP"
    }
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMContainerGroup_windowsComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "ContainerInsights"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/ContainerInsights"
  }
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "public"
  dns_name_label      = "acctestcontainergroup-%d"
  os_type             = "windows"
  restart_policy      = "Never"

  container {
    name   = "windowsservercore"
    image  = "microsoft/iis:windowsservercore"
    cpu    = "2.0"
    memory = "3.5"

    ports {
      port     = 80
      protocol = "TCP"
    }

    environment_variables = {
      foo  = "bar"
      foo1 = "bar1"
    }

    secure_environment_variables = {
      secureFoo  = "secureBar"
      secureFoo1 = "secureBar1"
    }

    readiness_probe {
      exec                  = ["cat", "/tmp/healthy"]
      initial_delay_seconds = 1
      period_seconds        = 1
      failure_threshold     = 1
      success_threshold     = 1
      timeout_seconds       = 1
    }

    liveness_probe {
      http_get {
        path   = "/"
        port   = 443
        scheme = "Http"
      }

      initial_delay_seconds = 1
      period_seconds        = 1
      failure_threshold     = 1
      success_threshold     = 1
      timeout_seconds       = 1
    }

    commands = ["cmd.exe", "echo", "hi"]
  }

  diagnostics {
    log_analytics {
      workspace_id  = azurerm_log_analytics_workspace.test.workspace_id
      workspace_key = azurerm_log_analytics_workspace.test.primary_shared_key
      log_type      = "ContainerInsights"

      metadata = {
        node-name = "acctestContainerGroup"
      }
    }
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMContainerGroup_linuxComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "ContainerInsights"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  workspace_resource_id = "${azurerm_log_analytics_workspace.test.id}"
  workspace_name        = "${azurerm_log_analytics_workspace.test.name}"

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/ContainerInsights"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name = "acctestss-%d"

  storage_account_name = "${azurerm_storage_account.test.name}"

  quota = 50
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  dns_name_label      = "acctestcontainergroup-%d"
  os_type             = "Linux"
  restart_policy      = "OnFailure"

  container {
    name   = "hf"
    image  = "seanmckenna/aci-hellofiles"
    cpu    = "1"
    memory = "1.5"

    ports {
      port     = 80
      protocol = "TCP"
    }

    gpu {
      count = 1
      sku   = "K80"
    }

    volume {
      name       = "logs"
      mount_path = "/aci/logs"
      read_only  = false
      share_name = "${azurerm_storage_share.test.name}"

      storage_account_name = "${azurerm_storage_account.test.name}"
      storage_account_key  = "${azurerm_storage_account.test.primary_access_key}"
    }

    environment_variables = {
      foo  = "bar"
      foo1 = "bar1"
    }

    secure_environment_variables = {
      secureFoo  = "secureBar"
      secureFoo1 = "secureBar1"
    }

    readiness_probe {
      exec                  = ["cat", "/tmp/healthy"]
      initial_delay_seconds = 1
      period_seconds        = 1
      failure_threshold     = 1
      success_threshold     = 1
      timeout_seconds       = 1
    }

    liveness_probe {
      http_get {
        path   = "/"
        port   = 443
        scheme = "Http"
      }

      initial_delay_seconds = 1
      period_seconds        = 1
      failure_threshold     = 1
      success_threshold     = 1
      timeout_seconds       = 1
    }

    commands = ["/bin/bash", "-c", "ls"]
  }

  diagnostics {
    log_analytics {
      workspace_id  = "${azurerm_log_analytics_workspace.test.workspace_id}"
      workspace_key = "${azurerm_log_analytics_workspace.test.primary_shared_key}"
      log_type      = "ContainerInsights"

      metadata = {
        node-name = "acctestContainerGroup"
      }
    }
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMContainerGroup_gitRepoVolume(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  dns_name_label      = "acctestcontainergroup-%d"
  os_type             = "Linux"
  restart_policy      = "OnFailure"

  container {
    name   = "hf"
    image  = "seanmckenna/aci-hellofiles"
    cpu    = "1"
    memory = "1.5"

    ports {
      port     = 80
      protocol = "TCP"
    }

    volume {
      name       = "logs"
      mount_path = "/aci/logs"
      read_only  = false

      git_repo {
        url       = "https://github.com/Azure-Samples/aci-helloworld"
        directory = "app"
        revision  = "d5ccfce"
      }
    }

    environment_variables = {
      foo  = "bar"
      foo1 = "bar1"
    }

    readiness_probe {
      exec                  = ["cat", "/tmp/healthy"]
      initial_delay_seconds = 1
      period_seconds        = 1
      failure_threshold     = 1
      success_threshold     = 1
      timeout_seconds       = 1
    }

    liveness_probe {
      http_get {
        path   = "/"
        port   = 443
        scheme = "Http"
      }

      initial_delay_seconds = 1
      period_seconds        = 1
      failure_threshold     = 1
      success_threshold     = 1
      timeout_seconds       = 1
    }

    commands = ["/bin/bash", "-c", "ls"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testCheckAzureRMContainerGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Containers.GroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Container Registry: %s", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Container Group %q (resource group: %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on containerGroupsClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMContainerGroupDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Containers.GroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_container_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Container Group still exists:\n%#v", resp)
			}

			return nil
		}
	}

	return nil
}

func testAccAzureRMContainerGroup_withPrivateEmpty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-containergroup-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "public"
  dns_name_label      = "jerome-aci-label"
  os_type             = "Linux"

  container {
    name   = "hello-world"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "1.5"

    ports {
      port     = 8000
      protocol = "TCP"
    }

    secure_environment_variables = {
      PRIVATE_EMPTY = ""
      PRIVATE_VALUE = "test"
    }

    environment_variables = {
      PUBLIC_EMPTY = ""
      PUBLIC_VALUE = "test"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
