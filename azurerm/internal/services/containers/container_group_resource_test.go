package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ContainerGroupResource struct {
}

func TestAccContainerGroup_SystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.SystemAssignedIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("0"),
				resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
			),
		},
		data.ImportStep("identity.0.principal_id"),
	})
}

func TestAccContainerGroup_UserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.UserAssignedIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerGroup_multipleAssignedIdentities(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.MultipleAssignedIdentities(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
			),
		},
		data.ImportStep("identity.0.principal_id"),
	})
}

func TestAccContainerGroup_imageRegistryCredentials(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.imageRegistryCredentials(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("image_registry_credential.#").HasValue("2"),
				check.That(data.ResourceName).Key("image_registry_credential.0.server").HasValue("hub.docker.com"),
				check.That(data.ResourceName).Key("image_registry_credential.0.username").HasValue("yourusername"),
				check.That(data.ResourceName).Key("image_registry_credential.0.password").HasValue("yourpassword"),
				check.That(data.ResourceName).Key("image_registry_credential.1.server").HasValue("mine.acr.io"),
				check.That(data.ResourceName).Key("image_registry_credential.1.username").HasValue("acrusername"),
				check.That(data.ResourceName).Key("image_registry_credential.1.password").HasValue("acrpassword"),
			),
		},
		data.ImportStep(
			"image_registry_credential.0.password",
			"image_registry_credential.1.password",
		),
	})
}

func TestAccContainerGroup_imageRegistryCredentialsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.imageRegistryCredentials(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("image_registry_credential.#").HasValue("2"),
				check.That(data.ResourceName).Key("image_registry_credential.0.server").HasValue("hub.docker.com"),
				check.That(data.ResourceName).Key("image_registry_credential.0.username").HasValue("yourusername"),
				check.That(data.ResourceName).Key("image_registry_credential.0.password").HasValue("yourpassword"),
				check.That(data.ResourceName).Key("image_registry_credential.1.server").HasValue("mine.acr.io"),
				check.That(data.ResourceName).Key("image_registry_credential.1.username").HasValue("acrusername"),
				check.That(data.ResourceName).Key("image_registry_credential.1.password").HasValue("acrpassword"),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("1"),
			),
		},
		{
			Config: r.imageRegistryCredentialsUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("image_registry_credential.#").HasValue("1"),
				check.That(data.ResourceName).Key("image_registry_credential.0.server").HasValue("hub.docker.com"),
				check.That(data.ResourceName).Key("image_registry_credential.0.username").HasValue("updatedusername"),
				check.That(data.ResourceName).Key("image_registry_credential.0.password").HasValue("updatedpassword"),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("1"),
			),
		},
	})
}

func TestAccContainerGroup_logTypeUnset(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.logTypeUnset(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.#").HasValue("1"),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.log_type").HasValue(""),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.metadata.%").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.workspace_id").Exists(),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.workspace_key").Exists(),
			),
		},
		data.ImportStep("diagnostics.0.log_analytics.0.workspace_key"),
	})
}

func TestAccContainerGroup_linuxBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
				check.That(data.ResourceName).Key("os_type").HasValue("Linux"),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("1"),
			),
		},
		data.ImportStep(
			"image_registry_credential.0.password",
			"image_registry_credential.1.password",
		),
	})
}

func TestAccContainerGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_container_group"),
		},
	})
}

func TestAccContainerGroup_linuxBasicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
			),
		},
		{
			Config: r.linuxBasicUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("2"),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("2"),
			),
		},
	})
}

func TestAccContainerGroup_linuxBasicTagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
			),
		},
		{
			Config: r.linuxBasicTagsUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
				check.That(data.ResourceName).Key("tags.OS").HasValue("Linux"),
			),
		},
	})
}

func TestAccContainerGroup_linuxComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.linuxComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.commands.#").HasValue("3"),
				check.That(data.ResourceName).Key("container.0.commands.0").HasValue("/bin/bash"),
				check.That(data.ResourceName).Key("container.0.commands.1").HasValue("-c"),
				check.That(data.ResourceName).Key("container.0.commands.2").HasValue("ls"),
				check.That(data.ResourceName).Key("container.0.environment_variables.%").HasValue("2"),
				check.That(data.ResourceName).Key("container.0.environment_variables.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("container.0.environment_variables.foo1").HasValue("bar1"),
				check.That(data.ResourceName).Key("container.0.secure_environment_variables.%").HasValue("2"),
				check.That(data.ResourceName).Key("container.0.secure_environment_variables.secureFoo").HasValue("secureBar"),
				check.That(data.ResourceName).Key("container.0.secure_environment_variables.secureFoo1").HasValue("secureBar1"),
				check.That(data.ResourceName).Key("container.0.gpu.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.gpu.0.count").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.gpu.0.sku").HasValue("K80"),
				check.That(data.ResourceName).Key("container.0.volume.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.volume.0.mount_path").HasValue("/aci/logs"),
				check.That(data.ResourceName).Key("container.0.volume.0.name").HasValue("logs"),
				check.That(data.ResourceName).Key("container.0.volume.0.read_only").HasValue("false"),
				check.That(data.ResourceName).Key("os_type").HasValue("Linux"),
				check.That(data.ResourceName).Key("restart_policy").HasValue("OnFailure"),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.#").HasValue("1"),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.log_type").HasValue("ContainerInsights"),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.metadata.%").HasValue("1"),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.workspace_id").Exists(),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.workspace_key").Exists(),
				check.That(data.ResourceName).Key("container.0.readiness_probe.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.exec.#").HasValue("2"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.exec.0").HasValue("cat"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.exec.1").HasValue("/tmp/healthy"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.http_get.#").HasValue("0"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.initial_delay_seconds").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.period_seconds").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.failure_threshold").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.success_threshold").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.timeout_seconds").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.failure_threshold").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.http_get.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.http_get.0.path").HasValue("/"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.http_get.0.port").HasValue("443"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.http_get.0.scheme").HasValue("Http"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.initial_delay_seconds").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.period_seconds").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.success_threshold").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.timeout_seconds").HasValue("1"),
			),
		},
		data.ImportStep(
			"container.0.volume.0.storage_account_key",
			"container.0.secure_environment_variables.%",
			"container.0.secure_environment_variables.secureFoo",
			"container.0.secure_environment_variables.secureFoo1",
			"diagnostics.0.log_analytics.0.workspace_key",
		),
	})
}

func TestAccContainerGroup_virtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.virtualNetwork(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestCheckNoResourceAttr(data.ResourceName, "dns_label_name"),
				resource.TestCheckNoResourceAttr(data.ResourceName, "identity"),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
				check.That(data.ResourceName).Key("os_type").HasValue("Linux"),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("1"),
				check.That(data.ResourceName).Key("ip_address_type").HasValue("Private"),
				check.That(data.ResourceName).Key("network_profile_id").Exists(),
				check.That(data.ResourceName).Key("dns_config.#").HasValue("1"),
			),
		},
	})
}

func TestAccContainerGroup_windowsBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.windowsBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
				check.That(data.ResourceName).Key("os_type").HasValue("Windows"),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerGroup_windowsComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.windowsComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.commands.#").HasValue("3"),
				check.That(data.ResourceName).Key("container.0.commands.0").HasValue("cmd.exe"),
				check.That(data.ResourceName).Key("container.0.commands.1").HasValue("echo"),
				check.That(data.ResourceName).Key("container.0.commands.2").HasValue("hi"),
				check.That(data.ResourceName).Key("container.0.environment_variables.%").HasValue("2"),
				check.That(data.ResourceName).Key("container.0.environment_variables.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("container.0.environment_variables.foo1").HasValue("bar1"),
				check.That(data.ResourceName).Key("container.0.secure_environment_variables.%").HasValue("2"),
				check.That(data.ResourceName).Key("container.0.secure_environment_variables.secureFoo").HasValue("secureBar"),
				check.That(data.ResourceName).Key("container.0.secure_environment_variables.secureFoo1").HasValue("secureBar1"),
				check.That(data.ResourceName).Key("os_type").HasValue("Windows"),
				check.That(data.ResourceName).Key("restart_policy").HasValue("Never"),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.#").HasValue("1"),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.log_type").HasValue("ContainerInsights"),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.metadata.%").HasValue("1"),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.workspace_id").Exists(),
				check.That(data.ResourceName).Key("diagnostics.0.log_analytics.0.workspace_key").Exists(),
				check.That(data.ResourceName).Key("container.0.readiness_probe.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.exec.#").HasValue("2"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.exec.0").HasValue("cat"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.exec.1").HasValue("/tmp/healthy"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.http_get.#").HasValue("0"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.initial_delay_seconds").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.period_seconds").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.failure_threshold").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.success_threshold").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.readiness_probe.0.timeout_seconds").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.failure_threshold").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.http_get.#").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.http_get.0.path").HasValue("/"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.http_get.0.port").HasValue("443"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.http_get.0.scheme").HasValue("Http"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.initial_delay_seconds").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.period_seconds").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.success_threshold").HasValue("1"),
				check.That(data.ResourceName).Key("container.0.liveness_probe.0.timeout_seconds").HasValue("1"),
			),
		},
		data.ImportStep(
			"container.0.secure_environment_variables.%",
			"container.0.secure_environment_variables.secureFoo",
			"container.0.secure_environment_variables.secureFoo1",
			"diagnostics.0.log_analytics.0.workspace_key",
		),
	})
}

func TestAccContainerGroup_withPrivateEmpty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withPrivateEmpty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"container.0.secure_environment_variables.PRIVATE_VALUE",
		),
	})
}

func TestAccContainerGroup_gitRepoVolume(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.gitRepoVolume(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerGroup_emptyDirVolume(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.emptyDirVolume(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerGroup_secretVolume(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.secretVolume(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("container.0.volume.0.secret"),
	})
}

func (ContainerGroupResource) SystemAssignedIdentity(data acceptance.TestData) string {
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

func (ContainerGroupResource) UserAssignedIdentity(data acceptance.TestData) string {
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

func (ContainerGroupResource) MultipleAssignedIdentities(data acceptance.TestData) string {
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

func (ContainerGroupResource) linuxBasic(data acceptance.TestData) string {
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

  exposed_port {
    port     = 80
    protocol = "TCP"
  }

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

func (ContainerGroupResource) linuxBasicTagsUpdated(data acceptance.TestData) string {
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
      port     = 80
      protocol = "TCP"
    }
  }

  tags = {
    environment = "Testing"
    OS          = "Linux"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ContainerGroupResource) requiresImport(data acceptance.TestData) string {
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
`, r.linuxBasic(data))
}

func (ContainerGroupResource) imageRegistryCredentials(data acceptance.TestData) string {
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

func (ContainerGroupResource) imageRegistryCredentialsUpdated(data acceptance.TestData) string {
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

func (ContainerGroupResource) logTypeUnset(data acceptance.TestData) string {
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

func (ContainerGroupResource) linuxBasicUpdated(data acceptance.TestData) string {
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

  exposed_port {
    port = 80
  }

  exposed_port {
    port     = 5443
    protocol = "UDP"
  }

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

func (ContainerGroupResource) virtualNetwork(data acceptance.TestData) string {
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

func (ContainerGroupResource) windowsBasic(data acceptance.TestData) string {
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

func (ContainerGroupResource) windowsComplete(data acceptance.TestData) string {
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

func (ContainerGroupResource) linuxComplete(data acceptance.TestData) string {
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

func (ContainerGroupResource) gitRepoVolume(data acceptance.TestData) string {
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

func (ContainerGroupResource) emptyDirVolume(data acceptance.TestData) string {
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
      empty_dir  = true
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

func (ContainerGroupResource) secretVolume(data acceptance.TestData) string {
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

    volume {
      name       = "config"
      mount_path = "/var/config"

      secret = {
        mysecret1 = "TXkgZmlyc3Qgc2VjcmV0IEZPTwo="
        mysecret2 = "TXkgc2Vjb25kIHNlY3JldCBCQVIK"
      }
    }
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t ContainerGroupResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["containerGroups"]

	resp, err := clients.Containers.GroupsClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading Container Group (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ContainerGroupResource) withPrivateEmpty(data acceptance.TestData) string {
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
