// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/containerapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppResource struct{}

func TestAccContainerAppResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_withSystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSystemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_withUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_withSystemAndUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSystemAndUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_withIdentityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSystemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withSystemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_basicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccContainerAppResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "rev1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_completeVolumeEmptyDir(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeEmptyDir(data, "rev1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_completeWithNoDaprAppPort(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeWithNoDaprAppPort(data, "rev1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_completeWithVNet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeWithVnet(data, "rev1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_completeWithSidecar(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeWithSidecar(data, "rev1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "rev1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data, "rev2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_removeDaprAppPort(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "rev1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate_withNoDaprAppPort(data, "rev2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("dapr.0.app_port").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppResource_secretRemoveShouldFail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeUpdate(data, "rev1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.complete(data, "rev2"),
			ExpectError: regexp.MustCompile("cannot remove secrets from Container Apps at this time"),
		},
	})
}

func TestAccContainerAppResource_secretRemoveWithAddShouldFail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app", "test")
	r := ContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeUpdate(data, "rev1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.completeChangedSecret(data, "rev2"),
			ExpectError: regexp.MustCompile("previously configured secret"),
		},
	})
}

func (r ContainerAppResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := containerapps.ParseContainerAppID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ContainerApps.ContainerAppClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ContainerAppResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppResource) withSystemIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Single"

  identity {
    type = "SystemAssigned"
  }

  template {
    container {
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppResource) withUserIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Single"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  template {
    container {
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppResource) withSystemAndUserIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Single"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  template {
    container {
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.5
      memory = "1Gi"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "import" {
  name                         = azurerm_container_app.test.name
  resource_group_name          = azurerm_container_app.test.resource_group_name
  container_app_environment_id = azurerm_container_app.test.container_app_environment_id
  revision_mode                = azurerm_container_app.test.revision_mode

  template {
    container {
      name   = azurerm_container_app.test.template.0.container.0.name
      image  = azurerm_container_app.test.template.0.container.0.image
      cpu    = azurerm_container_app.test.template.0.container.0.cpu
      memory = azurerm_container_app.test.template.0.container.0.memory
    }
  }
}
`, r.basic(data))
}

func (r ContainerAppResource) complete(data acceptance.TestData, revisionSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"

      readiness_probe {
        transport = "HTTP"
        port      = 5000
      }

      liveness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay           = 5
        interval_seconds        = 20
        timeout                 = 2
        failure_count_threshold = 1
      }

      startup_probe {
        transport = "TCP"
        port      = 5000
      }

      volume_mounts {
        name = azurerm_container_app_environment_storage.test.name
        path = "/tmp/testdata"
      }
    }

    volume {
      name         = azurerm_container_app_environment_storage.test.name
      storage_type = "AzureFile"
      storage_name = azurerm_container_app_environment_storage.test.name
    }

    min_replicas = 2
    max_replicas = 3

    revision_suffix = "%[3]s"
  }

  ingress {
    allow_insecure_connections = true
    external_enabled           = true
    target_port                = 5000
    transport                  = "http"
    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }

  registry {
    server               = azurerm_container_registry.test.login_server
    username             = azurerm_container_registry.test.admin_username
    password_secret_name = "registry-password"
  }

  secret {
    name  = "registry-password"
    value = azurerm_container_registry.test.admin_password
  }

  dapr {
    app_id       = "acctest-cont-%[2]d"
    app_port     = 5000
    app_protocol = "http"
  }

  tags = {
    foo     = "Bar"
    accTest = "1"
  }
}
`, r.templatePlusExtras(data), data.RandomInteger, revisionSuffix)
}

func (r ContainerAppResource) completeEmptyDir(data acceptance.TestData, revisionSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"

      readiness_probe {
        transport = "HTTP"
        port      = 5000
      }

      liveness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay           = 5
        interval_seconds        = 20
        timeout                 = 2
        failure_count_threshold = 1
      }

      startup_probe {
        transport = "TCP"
        port      = 5000
      }

      volume_mounts {
        name = azurerm_container_app_environment_storage.test.name
        path = "/tmp/testdata"
      }
    }

    volume {
      name         = azurerm_container_app_environment_storage.test.name
      storage_type = "EmptyDir"
    }

    min_replicas = 2
    max_replicas = 3

    revision_suffix = "%[3]s"
  }

  ingress {
    allow_insecure_connections = true
    external_enabled           = true
    target_port                = 5000
    transport                  = "http"
    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }

  registry {
    server               = azurerm_container_registry.test.login_server
    username             = azurerm_container_registry.test.admin_username
    password_secret_name = "registry-password"
  }

  secret {
    name  = "registry-password"
    value = azurerm_container_registry.test.admin_password
  }

  dapr {
    app_id       = "acctest-cont-%[2]d"
    app_port     = 5000
    app_protocol = "http"
  }

  tags = {
    foo     = "Bar"
    accTest = "1"
  }
}
`, r.templatePlusExtras(data), data.RandomInteger, revisionSuffix)
}

func (r ContainerAppResource) completeWithNoDaprAppPort(data acceptance.TestData, revisionSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"

      readiness_probe {
        transport = "HTTP"
        port      = 5000
      }

      liveness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay           = 5
        interval_seconds        = 20
        timeout                 = 2
        failure_count_threshold = 1
      }

      startup_probe {
        transport = "TCP"
        port      = 5000
      }

      volume_mounts {
        name = azurerm_container_app_environment_storage.test.name
        path = "/tmp/testdata"
      }
    }

    volume {
      name         = azurerm_container_app_environment_storage.test.name
      storage_type = "AzureFile"
      storage_name = azurerm_container_app_environment_storage.test.name
    }

    min_replicas = 2
    max_replicas = 3

    revision_suffix = "%[3]s"
  }

  ingress {
    allow_insecure_connections = true
    external_enabled           = true
    target_port                = 5000
    transport                  = "http"
    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }

  registry {
    server               = azurerm_container_registry.test.login_server
    username             = azurerm_container_registry.test.admin_username
    password_secret_name = "registry-password"
  }

  secret {
    name  = "registry-password"
    value = azurerm_container_registry.test.admin_password
  }

  dapr {
    app_id = "acctest-cont-%[2]d"
  }

  tags = {
    foo     = "Bar"
    accTest = "1"
  }
}
`, r.templatePlusExtras(data), data.RandomInteger, revisionSuffix)
}

func (r ContainerAppResource) completeWithVnet(data acceptance.TestData, revisionSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"

      readiness_probe {
        transport = "HTTP"
        port      = 5000
      }

      liveness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay           = 5
        timeout                 = 2
        failure_count_threshold = 1
      }

      startup_probe {
        transport = "TCP"
        port      = 5000
      }

      volume_mounts {
        name = azurerm_container_app_environment_storage.test.name
        path = "/tmp/testdata"
      }
    }

    volume {
      name         = azurerm_container_app_environment_storage.test.name
      storage_type = "AzureFile"
      storage_name = azurerm_container_app_environment_storage.test.name
    }

    min_replicas = 2
    max_replicas = 3

    revision_suffix = "%[3]s"
  }

  ingress {
    allow_insecure_connections = true
    target_port                = 5000
    transport                  = "http"
    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }

  registry {
    server               = azurerm_container_registry.test.login_server
    username             = azurerm_container_registry.test.admin_username
    password_secret_name = "registry-password"
  }

  secret {
    name  = "registry-password"
    value = azurerm_container_registry.test.admin_password
  }

  tags = {
    foo     = "Bar"
    accTest = "1"
  }
}
`, r.templateWithVnet(data), data.RandomInteger, revisionSuffix)
}

func (r ContainerAppResource) completeWithSidecar(data acceptance.TestData, revisionSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-sidecar-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"

      readiness_probe {
        transport = "HTTP"
        port      = 5000
      }

      liveness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay           = 5
        timeout                 = 2
        failure_count_threshold = 1
      }

      startup_probe {
        transport = "TCP"
        port      = 5000
      }

      volume_mounts {
        name = azurerm_container_app_environment_storage.test.name
        path = "/tmp/testdata"
      }
    }

    container {
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"

      readiness_probe {
        transport = "HTTP"
        port      = 5000
      }

      liveness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay           = 5
        timeout                 = 2
        failure_count_threshold = 1
      }

      startup_probe {
        transport = "TCP"
        port      = 5000
      }

      volume_mounts {
        name = azurerm_container_app_environment_storage.test.name
        path = "/tmp/testdata"
      }
    }

    volume {
      name         = azurerm_container_app_environment_storage.test.name
      storage_type = "AzureFile"
      storage_name = azurerm_container_app_environment_storage.test.name
    }

    min_replicas = 2
    max_replicas = 3

    revision_suffix = "%[3]s"
  }

  ingress {
    allow_insecure_connections = true
    target_port                = 5000
    transport                  = "http"
    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }

  registry {
    server               = azurerm_container_registry.test.login_server
    username             = azurerm_container_registry.test.admin_username
    password_secret_name = "registry-password"
  }

  secret {
    name  = "registry-password"
    value = azurerm_container_registry.test.admin_password
  }

  tags = {
    foo     = "Bar"
    accTest = "1"
  }
}
`, r.templatePlusExtras(data), data.RandomInteger, revisionSuffix)
}

func (r ContainerAppResource) completeChangedSecret(data acceptance.TestData, revisionSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Multiple"

  template {
    container {
      name  = "acctest-cont-%[2]d"
      image = "jackofallops/azure-containerapps-python-acctest:v0.0.1"

      cpu    = 0.5
      memory = "1Gi"

      readiness_probe {
        transport               = "HTTP"
        port                    = 5000
        path                    = "/uptime"
        timeout                 = 2
        failure_count_threshold = 1
        success_count_threshold = 1

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }
      }

      liveness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay           = 5
        timeout                 = 2
        failure_count_threshold = 3
      }

      startup_probe {
        transport               = "TCP"
        port                    = 5000
        timeout                 = 5
        failure_count_threshold = 1
      }
    }

    min_replicas = 1
    max_replicas = 4

    revision_suffix = "%[3]s"
  }

  ingress {
    allow_insecure_connections = true
    external_enabled           = true
    target_port                = 5000
    transport                  = "auto"

    traffic_weight {
      latest_revision = true
      percentage      = 20
    }

    traffic_weight {
      revision_suffix = "rev1"
      percentage      = 80
    }
  }

  registry {
    server               = azurerm_container_registry.test.login_server
    username             = azurerm_container_registry.test.admin_username
    password_secret_name = "registry-password"
  }

  secret {
    name  = "registry-password"
    value = azurerm_container_registry.test.admin_password
  }

  secret {
    name  = "pickle"
    value = "morty"
  }

  dapr {
    app_id       = "acctest-cont-%[2]d"
    app_port     = 5000
    app_protocol = "http"
  }

  tags = {
    foo     = "Bar"
    accTest = "1"
  }
}
`, r.templatePlusExtras(data), data.RandomInteger, revisionSuffix)
}

func (r ContainerAppResource) completeUpdate(data acceptance.TestData, revisionSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Multiple"

  template {
    container {
      name  = "acctest-cont-%[2]d"
      image = "jackofallops/azure-containerapps-python-acctest:v0.0.1"

      cpu    = 0.5
      memory = "1Gi"

      readiness_probe {
        transport               = "HTTP"
        port                    = 5000
        path                    = "/uptime"
        timeout                 = 2
        failure_count_threshold = 1
        success_count_threshold = 1

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }
      }

      liveness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay           = 5
        timeout                 = 2
        failure_count_threshold = 3
      }

      startup_probe {
        transport               = "TCP"
        port                    = 5000
        timeout                 = 5
        failure_count_threshold = 1
      }

      volume_mounts {
        name = azurerm_container_app_environment_storage.test.name
        path = "/tmp/testdata"
      }
    }

    volume {
      name         = azurerm_container_app_environment_storage.test.name
      storage_type = "AzureFile"
      storage_name = azurerm_container_app_environment_storage.test.name
    }

    min_replicas = 1
    max_replicas = 4

    revision_suffix = "%[3]s"
  }

  ingress {
    allow_insecure_connections = true
    external_enabled           = true
    target_port                = 5000
    transport                  = "auto"

    traffic_weight {
      latest_revision = true
      percentage      = 20
    }

    traffic_weight {
      revision_suffix = "rev1"
      percentage      = 80
    }
  }

  registry {
    server               = azurerm_container_registry.test.login_server
    username             = azurerm_container_registry.test.admin_username
    password_secret_name = "registry-password"
  }

  secret {
    name  = "registry-password"
    value = azurerm_container_registry.test.admin_password
  }

  secret {
    name  = "rick"
    value = "morty"
  }

  dapr {
    app_id       = "acctest-cont-%[2]d"
    app_port     = 5000
    app_protocol = "http"
  }

  tags = {
    foo     = "Bar"
    accTest = "1"
  }
}
`, r.templatePlusExtras(data), data.RandomInteger, revisionSuffix)
}

func (r ContainerAppResource) completeUpdate_withNoDaprAppPort(data acceptance.TestData, revisionSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Multiple"

  template {
    container {
      name  = "acctest-cont-%[2]d"
      image = "jackofallops/azure-containerapps-python-acctest:v0.0.1"

      cpu    = 0.5
      memory = "1Gi"

      readiness_probe {
        transport               = "HTTP"
        port                    = 5000
        path                    = "/uptime"
        timeout                 = 2
        failure_count_threshold = 1
        success_count_threshold = 1

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }
      }

      liveness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay           = 5
        timeout                 = 2
        failure_count_threshold = 3
      }

      startup_probe {
        transport               = "TCP"
        port                    = 5000
        timeout                 = 5
        failure_count_threshold = 1
      }

      volume_mounts {
        name = azurerm_container_app_environment_storage.test.name
        path = "/tmp/testdata"
      }
    }

    volume {
      name         = azurerm_container_app_environment_storage.test.name
      storage_type = "AzureFile"
      storage_name = azurerm_container_app_environment_storage.test.name
    }

    min_replicas = 1
    max_replicas = 4

    revision_suffix = "%[3]s"
  }

  ingress {
    allow_insecure_connections = true
    external_enabled           = true
    target_port                = 5000
    transport                  = "auto"

    traffic_weight {
      latest_revision = true
      percentage      = 20
    }

    traffic_weight {
      revision_suffix = "rev1"
      percentage      = 80
    }
  }

  registry {
    server               = azurerm_container_registry.test.login_server
    username             = azurerm_container_registry.test.admin_username
    password_secret_name = "registry-password"
  }

  secret {
    name  = "registry-password"
    value = azurerm_container_registry.test.admin_password
  }

  secret {
    name  = "rick"
    value = "morty"
  }

  dapr {
    app_id = "acctest-cont-%[2]d"
  }

  tags = {
    foo     = "Bar"
    accTest = "1"
  }
}
`, r.templatePlusExtras(data), data.RandomInteger, revisionSuffix)
}

func (ContainerAppResource) template(data acceptance.TestData) string {
	return ContainerAppEnvironmentResource{}.basic(data)
}

func (ContainerAppResource) templateWithVnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"
  admin_enabled       = true

  network_rule_set = []
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%[3]s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_share" "test" {
  name                 = "testshare%[3]s"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 1
}

resource "azurerm_container_app_environment_storage" "test" {
  name                         = "testacc-caes-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  account_name                 = azurerm_storage_account.test.name
  access_key                   = azurerm_storage_account.test.primary_access_key
  share_name                   = azurerm_storage_share.test.name
  access_mode                  = "ReadWrite"
}
`, ContainerAppEnvironmentResource{}.complete(data), data.RandomInteger, data.RandomString)
}

func (ContainerAppResource) templatePlusExtras(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"
  admin_enabled       = true

  network_rule_set = []
}

resource "azurerm_storage_share" "test" {
  name                 = "testshare%[3]s"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 1
}

resource "azurerm_container_app_environment_storage" "test" {
  name                         = "testacc-caes-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  account_name                 = azurerm_storage_account.test.name
  access_key                   = azurerm_storage_account.test.primary_access_key
  share_name                   = azurerm_storage_share.test.name
  access_mode                  = "ReadWrite"
}
`, ContainerAppEnvironmentDaprComponentResource{}.complete(data), data.RandomInteger, data.RandomString)
}
