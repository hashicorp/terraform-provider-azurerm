package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerGroupResource struct{}

func TestAccContainerGroup_SystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.SystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("0"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
			),
		},
		data.ImportStep("identity.0.principal_id"),
	})
}

func TestAccContainerGroup_ProbeHttpGet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ProbeHttpGet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("identity.0.principal_id"),
	})
}

func TestAccContainerGroup_ProbeExec(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ProbeExec(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("identity.0.principal_id"),
	})
}

func TestAccContainerGroup_SystemAssignedIdentityNoNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.SystemAssignedIdentityNoNetwork(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("0"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
			),
		},
		data.ImportStep("identity.0.principal_id", "ip_address_type"),
	})
}

func TestAccContainerGroup_UserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.UserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerGroup_UserAssignedIdentityWithVirtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.UserAssignedIdentityWithVirtualNetwork(data),
		},
		data.ImportStep("network_profile_id"),
	})
}

func TestAccContainerGroup_multipleAssignedIdentities(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.MultipleAssignedIdentities(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerGroup_imageRegistryCredentials(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageRegistryCredentials(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageRegistryCredentials(data),
			Check: acceptance.ComposeTestCheckFunc(
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
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logTypeUnset(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func TestAccContainerGroup_exposedPort(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.exposedPort(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
				check.That(data.ResourceName).Key("os_type").HasValue("Linux"),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("2"),
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
			),
		},
		{
			Config: r.linuxBasicUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("2"),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("2"),
			),
		},
	})
}

func TestAccContainerGroup_exposedPortUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.exposedPort(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("exposed_port.#").HasValue("1"),
			),
		},
		{
			Config: r.exposedPortUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("2"),
				check.That(data.ResourceName).Key("exposed_port.#").HasValue("2"),
			),
		},
	})
}

func TestAccContainerGroup_linuxBasicTagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
			),
		},
		{
			Config: r.linuxBasicTagsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetwork(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckNoResourceAttr(data.ResourceName, "dns_label_name"),
				acceptance.TestCheckNoResourceAttr(data.ResourceName, "identity"),
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

func TestAccContainerGroup_virtualNetworkParallel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetworkParallel(data, 4),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName+".0").ExistsInAzure(r),
				check.That(data.ResourceName+".1").ExistsInAzure(r),
				check.That(data.ResourceName+".2").ExistsInAzure(r),
				check.That(data.ResourceName+".3").ExistsInAzure(r),
			),
		},
	})
}

func TestAccContainerGroup_SystemAssignedIdentityVirtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.SystemAssignedIdentityVirtualNetwork(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckNoResourceAttr(data.ResourceName, "dns_label_name"),
				acceptance.TestCheckNoResourceAttr(data.ResourceName, "identity"),
				check.That(data.ResourceName).Key("container.#").HasValue("1"),
				check.That(data.ResourceName).Key("os_type").HasValue("Linux"),
				check.That(data.ResourceName).Key("container.0.ports.#").HasValue("1"),
				check.That(data.ResourceName).Key("ip_address_type").HasValue("Private"),
				check.That(data.ResourceName).Key("network_profile_id").Exists(),
				check.That(data.ResourceName).Key("dns_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("0"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
			),
		},
		data.ImportStep("identity.0.principal_id", "network_profile_id"),
	})
}

func TestAccContainerGroup_windowsBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPrivateEmpty(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gitRepoVolume(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerGroup_emptyDirVolume(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emptyDirVolume(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerGroup_emptyDirVolumeShared(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emptyDirVolumeShared(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("ip_address_type"),
	})
}

func TestAccContainerGroup_emptyDirVolumeSharedWithInitContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emptyDirVolumeSharedWithInitContainer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("ip_address_type"),
	})
}

func TestAccContainerGroup_withInitContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withInitContainer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("ip_address_type"),
	})
}

func TestAccContainerGroup_secretVolume(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.secretVolume(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("container.0.volume.0.secret"),
	})
}

func TestAccContainerGroup_encryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group", "test")
	r := ContainerGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
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

func (ContainerGroupResource) ProbeHttpGet(data acceptance.TestData) string {
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
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 80
      protocol = "TCP"
    }
    readiness_probe {
      http_get {
        path   = "/"
        port   = 443
        scheme = "Http"
      }
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

func (ContainerGroupResource) ProbeExec(data acceptance.TestData) string {
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
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
    cpu    = "0.5"
    memory = "0.5"
    liveness_probe {
      exec = ["cat", "/tmp/healthy"]
    }
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

func (ContainerGroupResource) SystemAssignedIdentityNoNetwork(data acceptance.TestData) string {
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
  ip_address_type     = "None"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
    cpu    = "0.5"
    memory = "0.5"
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
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 80
      protocol = "TCP"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
func (ContainerGroupResource) UserAssignedIdentityWithVirtualNetwork(data acceptance.TestData) string {
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
  address_prefixes     = ["10.1.0.0/24"]

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

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Private"
  os_type             = "Linux"
  network_profile_id  = azurerm_network_profile.test.id
  container {
    name   = "hw"
    image  = "ubuntu:20.04"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 80
      protocol = "TCP"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
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
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 80
      protocol = "TCP"
    }
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
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

func (ContainerGroupResource) exposedPort(data acceptance.TestData) string {
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
  ip_address_type     = "Public"
  os_type             = "Linux"

  exposed_port {
    port     = 80
    protocol = "TCP"
  }

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 80
      protocol = "TCP"
    }
    ports {
      port     = 5443
      protocol = "UDP"
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
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
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
  name                = azurerm_container_group.test.name
  location            = azurerm_container_group.test.location
  resource_group_name = azurerm_container_group.test.resource_group_name
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
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
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
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
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
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
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
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
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
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
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
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
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
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
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "0.5"
  }

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerGroupResource) exposedPortUpdated(data acceptance.TestData) string {
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
  ip_address_type     = "Public"
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
    image  = "ubuntu:20.04"
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
  address_prefixes     = ["10.1.0.0/24"]

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
    image  = "ubuntu:20.04"
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

func (ContainerGroupResource) virtualNetworkParallel(data acceptance.TestData, count int) string {
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
  address_prefixes     = ["10.1.0.0/24"]

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
  count               = %d
  name                = "acctestcontainergroup-${count.index}-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Private"
  network_profile_id  = azurerm_network_profile.test.id
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port = 80
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, count, data.RandomInteger)
}

func (ContainerGroupResource) SystemAssignedIdentityVirtualNetwork(data acceptance.TestData) string {
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
  address_prefixes     = ["10.1.0.0/24"]

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
    image  = "ubuntu:20.04"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port = 80
    }
  }

  dns_config {
    nameservers = ["reddog.microsoft.com", "somecompany.somedomain"]
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
  ip_address_type     = "Public"
  os_type             = "Windows"

  container {
    name   = "windowsservercore"
    image  = "mcr.microsoft.com/windows/servercore/iis:20210810-windowsservercore-ltsc2019"
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
  ip_address_type     = "Public"
  dns_name_label      = "acctestcontainergroup-%d"
  os_type             = "Windows"
  restart_policy      = "Never"

  container {
    name   = "windowsservercore"
    image  = "mcr.microsoft.com/windows/servercore/iis:20210810-windowsservercore-ltsc2019"
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

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name = "acctestss-%d"

  storage_account_name = azurerm_storage_account.test.name

  quota = 50
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Public"
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
      share_name = azurerm_storage_share.test.name

      storage_account_name = azurerm_storage_account.test.name
      storage_account_key  = azurerm_storage_account.test.primary_access_key
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Public"
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Public"
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

func (ContainerGroupResource) emptyDirVolumeSharedWithInitContainer(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroupemptyshared-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "None"
  os_type             = "Linux"
  restart_policy      = "Never"

  init_container {
    name     = "init"
    image    = "busybox"
    commands = ["touch", "/sharedempty/file.txt"]

    volume {
      name       = "logs"
      mount_path = "/sharedempty"
      read_only  = false
      empty_dir  = true
    }
  }

  container {
    name   = "reader"
    image  = "ubuntu:20.04"
    cpu    = "1"
    memory = "1.5"

    volume {
      name       = "logs"
      mount_path = "/sharedempty"
      read_only  = false
      empty_dir  = true
    }

    commands = ["/bin/bash", "-c", "timeout 30 watch --interval 1 --errexit \"! cat /sharedempty/file.txt\""]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerGroupResource) emptyDirVolumeShared(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroupemptyshared-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "None"
  os_type             = "Linux"
  restart_policy      = "Never"

  container {
    name     = "writer"
    image    = "ubuntu:20.04"
    cpu      = "1"
    memory   = "1.5"
    commands = ["touch", "/sharedempty/file.txt"]

    volume {
      name       = "logs"
      mount_path = "/sharedempty"
      read_only  = false
      empty_dir  = true
    }
  }

  container {
    name   = "reader"
    image  = "ubuntu:20.04"
    cpu    = "1"
    memory = "1.5"

    volume {
      name       = "logs"
      mount_path = "/sharedempty"
      read_only  = false
      empty_dir  = true
    }

    commands = ["/bin/bash", "-c", "timeout 30 watch --interval 1 --errexit \"! cat /sharedempty/file.txt\""]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerGroupResource) withInitContainer(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroupemptyshared-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "None"
  os_type             = "Linux"
  restart_policy      = "Never"

  init_container {
    name     = "init"
    image    = "busybox"
    commands = ["echo", "hello from init"]
  }

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
    cpu    = "1"
    memory = "1.5"

    commands = ["echo", "hello from ubuntu"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
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

func (t ContainerGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ContainerGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.GroupsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Container Group (%s): %+v", id.String(), err)
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
  ip_address_type     = "Public"
  dns_name_label      = "jerome-aci-label"
  os_type             = "Linux"

  container {
    name   = "hello-world"
    image  = "ubuntu:20.04"
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

func (ContainerGroupResource) encryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acc-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.test.id
  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "List",
    "Purge",
    "Update",
  ]

  secret_permissions = [
    "Get",
    "Delete",
    "Set",
  ]
  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%[1]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
  depends_on = [azurerm_key_vault_access_policy.terraform]
}

data "azuread_service_principal" "test" {
  display_name = "Azure Container Instance Service"
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey"
  ]

  tenant_id  = data.azurerm_client_config.current.tenant_id
  object_id  = data.azuread_service_principal.test.object_id
  depends_on = [azurerm_key_vault_access_policy.terraform]
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "ubuntu:20.04"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port     = 80
      protocol = "TCP"
    }
  }
  key_vault_key_id = azurerm_key_vault_key.test.id
  depends_on       = [azurerm_key_vault_access_policy.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
