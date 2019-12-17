package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMContainerGroup_SystemAssignedIdentity(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMContainerGroup_SystemAssignedIdentity(ri, acceptance.Location())
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.identity_ids.#", "0"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", validate.UUIDRegExp),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"identity.0.principal_id",
				},
			},
		},
	})
}

func TestAccAzureRMContainerGroup_UserAssignedIdentity(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(14)
	config := testAccAzureRMContainerGroup_UserAssignedIdentity(ri, acceptance.Location(), rs)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "UserAssigned"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.identity_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.principal_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMContainerGroup_multipleAssignedIdentities(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(14)
	config := testAccAzureRMContainerGroup_MultipleAssignedIdentities(ri, acceptance.Location(), rs)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned, UserAssigned"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.identity_ids.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", validate.UUIDRegExp),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"identity.0.principal_id",
				},
			},
		},
	})
}

func TestAccAzureRMContainerGroup_imageRegistryCredentials(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()

	config := testAccAzureRMContainerGroup_imageRegistryCredentials(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.server", "hub.docker.com"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.username", "yourusername"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.password", "yourpassword"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.server", "mine.acr.io"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.username", "acrusername"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.password", "acrpassword"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"image_registry_credential.0.password",
					"image_registry_credential.1.password",
				},
			},
		},
	})
}

func TestAccAzureRMContainerGroup_imageRegistryCredentialsUpdate(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()

	config := testAccAzureRMContainerGroup_imageRegistryCredentials(ri, acceptance.Location())
	updated := testAccAzureRMContainerGroup_imageRegistryCredentialsUpdated(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.server", "hub.docker.com"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.username", "yourusername"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.password", "yourpassword"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.server", "mine.acr.io"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.username", "acrusername"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.password", "acrpassword"),
					resource.TestCheckResourceAttr(resourceName, "container.0.port", "5443"),
					resource.TestCheckResourceAttr(resourceName, "container.0.protocol", "UDP"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.server", "hub.docker.com"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.username", "updatedusername"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.password", "updatedpassword"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_logTypeUnset(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMContainerGroup_logTypeUnset(ri, acceptance.Location())
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "diagnostics.0.log_analytics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics.0.log_analytics.0.log_type", ""),
					resource.TestCheckResourceAttr(resourceName, "diagnostics.0.log_analytics.0.metadata.%", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics.0.log_analytics.0.workspace_id"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics.0.log_analytics.0.workspace_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"diagnostics.0.log_analytics.0.workspace_key",
				},
			},
		},
	})
}

func TestAccAzureRMContainerGroup_linuxBasic(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()

	config := testAccAzureRMContainerGroup_linuxBasic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "container.0.port", "80"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"image_registry_credential.0.password",
					"image_registry_credential.1.password",
				},
			},
		},
	})
}

func TestAccAzureRMContainerGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerGroup_linuxBasic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMContainerGroup_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_container_group"),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_linuxBasicUpdate(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()

	config := testAccAzureRMContainerGroup_linuxBasic(ri, acceptance.Location())
	updatedConfig := testAccAzureRMContainerGroup_linuxBasicUpdated(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_linuxComplete(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()

	config := testAccAzureRMContainerGroup_linuxComplete(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.command", "/bin/bash -c ls"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.0", "/bin/bash"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.1", "-c"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.2", "ls"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.secureFoo", "secureBar"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.secureFoo1", "secureBar1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.gpu.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.gpu.0.count", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.gpu.0.sku", "K80"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.0.mount_path", "/aci/logs"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.0.name", "logs"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.0.read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "restart_policy", "OnFailure"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics.0.log_analytics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics.0.log_analytics.0.log_type", "ContainerInsights"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics.0.log_analytics.0.metadata.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics.0.log_analytics.0.workspace_id"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics.0.log_analytics.0.workspace_key"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.exec.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.exec.0", "cat"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.exec.1", "/tmp/healthy"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.http_get.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.initial_delay_seconds", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.period_seconds", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.failure_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.success_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.timeout_seconds", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.failure_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.http_get.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.http_get.0.path", "/"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.http_get.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.http_get.0.scheme", "Http"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.initial_delay_seconds", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.period_seconds", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.success_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.timeout_seconds", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"container.0.volume.0.storage_account_key",
					"container.0.secure_environment_variables.%",
					"container.0.secure_environment_variables.secureFoo",
					"container.0.secure_environment_variables.secureFoo1",
					"diagnostics.0.log_analytics.0.workspace_key",
				},
			},
		},
	})
}

func TestAccAzureRMContainerGroup_virtualNetwork(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMContainerGroup_virtualNetwork(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckNoResourceAttr(resourceName, "dns_label_name"),
					resource.TestCheckNoResourceAttr(resourceName, "identity"),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "container.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "ip_address_type", "Private"),
					resource.TestCheckResourceAttrSet(resourceName, "network_profile_id"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_windowsBasic(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()

	config := testAccAzureRMContainerGroup_windowsBasic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMContainerGroup_windowsComplete(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := tf.AccRandTimeInt()

	config := testAccAzureRMContainerGroup_windowsComplete(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.command", "cmd.exe echo hi"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.0", "cmd.exe"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.1", "echo"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.2", "hi"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.secureFoo", "secureBar"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.secureFoo1", "secureBar1"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(resourceName, "restart_policy", "Never"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics.0.log_analytics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics.0.log_analytics.0.log_type", "ContainerInsights"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics.0.log_analytics.0.metadata.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics.0.log_analytics.0.workspace_id"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics.0.log_analytics.0.workspace_key"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.exec.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.exec.0", "cat"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.exec.1", "/tmp/healthy"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.http_get.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.initial_delay_seconds", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.period_seconds", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.failure_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.success_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.readiness_probe.0.timeout_seconds", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.failure_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.http_get.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.http_get.0.path", "/"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.http_get.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.http_get.0.scheme", "Http"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.initial_delay_seconds", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.period_seconds", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.success_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.liveness_probe.0.timeout_seconds", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"container.0.secure_environment_variables.%",
					"container.0.secure_environment_variables.secureFoo",
					"container.0.secure_environment_variables.secureFoo1",
					"diagnostics.0.log_analytics.0.workspace_key",
				},
			},
		},
	})
}

func testAccAzureRMContainerGroup_SystemAssignedIdentity(ri int, location string) string {
	return fmt.Sprintf(`
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
    port   = 80
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "Testing"
  }
}
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_UserAssignedIdentity(ri int, location string, rString string) string {
	return fmt.Sprintf(`
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
    port   = 80
  }

  identity {
    type         = "UserAssigned"
    identity_ids = ["${azurerm_user_assigned_identity.test.id}"]
  }

  tags = {
    environment = "Testing"
  }
}
`, ri, location, rString, ri)
}

func testAccAzureRMContainerGroup_MultipleAssignedIdentities(ri int, location string, rString string) string {
	return fmt.Sprintf(`
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
    port   = 80
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = ["${azurerm_user_assigned_identity.test.id}"]
  }

  tags = {
    environment = "Testing"
  }
}
`, ri, location, rString, ri)
}

func testAccAzureRMContainerGroup_linuxBasic(ri int, location string) string {
	return fmt.Sprintf(`
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
    port   = 80
  }

  tags = {
    environment = "Testing"
  }
}
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_requiresImport(rInt int, location string) string {
	template := testAccAzureRMContainerGroup_linuxBasic(rInt, location)
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
    port   = "80"
  }

  tags = {
    environment = "Testing"
  }
}
`, template)
}

func testAccAzureRMContainerGroup_imageRegistryCredentials(ri int, location string) string {
	return fmt.Sprintf(`
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
    name     = "hw"
    image    = "microsoft/aci-helloworld:latest"
    cpu      = "0.5"
    memory   = "0.5"
    port     = 5443
    protocol = "UDP"
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
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_imageRegistryCredentialsUpdated(ri int, location string) string {
	return fmt.Sprintf(`
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
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_logTypeUnset(ri int, location string) string {
	return fmt.Sprintf(`
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
    port   = 80
  }

  diagnostics {
    log_analytics {
      workspace_id  = "${azurerm_log_analytics_workspace.test.workspace_id}"
      workspace_key = "${azurerm_log_analytics_workspace.test.primary_shared_key}"
    }
  }

  tags = {
    environment = "Testing"
  }
}
`, ri, location, ri, ri)
}

func testAccAzureRMContainerGroup_linuxBasicUpdated(ri int, location string) string {
	return fmt.Sprintf(`
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
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_virtualNetwork(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "testvnet"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  container_network_interface {
    name = "testcnic"

    ip_configuration {
      name      = "testipconfig"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "Private"
  network_profile_id  = "${azurerm_network_profile.test.id}"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    port   = 80
  }

  tags = {
    environment = "Testing"
  }
}
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_windowsBasic(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_windowsComplete(ri int, location string) string {
	return fmt.Sprintf(`
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

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, ri, location, ri, ri, ri)
}

func testAccAzureRMContainerGroup_linuxComplete(ri int, location string) string {
	return fmt.Sprintf(`
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

  resource_group_name  = "${azurerm_resource_group.test.name}"
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
`, ri, location, ri, ri, ri, ri, ri)
}

func testCheckAzureRMContainerGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Containers.GroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
