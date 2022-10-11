package containerapps_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/containerapps"
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
		// TODO - Uncomment the following stages when https://github.com/Azure/azure-rest-api-specs/issues/19285 / https://github.com/microsoft/azure-container-apps/issues/395 are resolved and secrets can be managed?
		// {
		// 	Config: r.complete(data, "rev3"),
		// 	Check: acceptance.ComposeTestCheckFunc(
		// 		check.That(data.ResourceName).ExistsInAzure(r),
		// 	),
		// },
		// data.ImportStep(),
		// {
		// 	Config: r.completeUpdate2(data, "rev4"),
		// 	Check: acceptance.ComposeTestCheckFunc(
		// 		check.That(data.ResourceName).ExistsInAzure(r),
		// 	),
		// },
		// data.ImportStep(),
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
			ExpectError: regexp.MustCompile("previously configured secret \"rick\" was removed. Removing secrets is not supported at this time"),
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
	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}

	return pointer.To(true), nil
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
      name   = container_app_environment.test.template.0.container.0.name
      image  = container_app_environment.test.template.0.container.0.image
      cpu    = container_app_environment.test.template.0.container.0.cpu
      memory = container_app_environment.test.template.0.container.0.memory
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
      //args    = ["HOSTNAME", "APPLICATION_PORT"]
      //command = ["node"]

      readiness_probe {
        transport = "http"
        port      = 5000
      }

      liveness_probe {
        transport = "http"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay     = 5
        interval          = 20
        timeout           = 2
        failure_threshold = 1
      }

      startup_probe {
        transport = "tcp"
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
    }

    min_replicas = 2
    max_replicas = 3

    revision_suffix = "%[3]s"
  }

  ingress {
    allow_insecure_connections = true
    is_external                = true
    target_port                = 5000
    transport                  = "http"
    traffic_weight {
      latest_revision = true
      weight          = 100
    }
  }

  registry {
    server                    = azurerm_container_registry.test.login_server
    username                  = azurerm_container_registry.test.admin_username
    password_secret_reference = "registry-password"
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
      //args    = ["HOSTNAME", "APPLICATION_PORT"]
      //command = ["node"]

      readiness_probe {
        transport = "http"
        port      = 5000
      }

      liveness_probe {
        transport = "http"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay     = 5
        interval          = 20
        timeout           = 2
        failure_threshold = 1
      }

      startup_probe {
        transport = "tcp"
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
      weight          = 100
    }
  }

  registry {
    server                    = azurerm_container_registry.test.login_server
    username                  = azurerm_container_registry.test.admin_username
    password_secret_reference = "registry-password"
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
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.5
      memory = "1Gi"
      //args    = ["HOSTNAME"] // TODO - Add a container Image where args and command can be used
      //command = ["node"]

      readiness_probe {
        transport         = "http"
        port              = 5000
        path              = "/uptime"
        timeout           = 2
        failure_threshold = 1
        success_threshold = 1

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }
      }

      liveness_probe {
        transport = "http"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay     = 5
        interval          = 20
        timeout           = 2
        failure_threshold = 3
      }

      startup_probe {
        transport         = "tcp"
        port              = 5000
        timeout           = 5
        failure_threshold = 1
      }

      //volume_mounts {
      // name = "testVol"
      // path = "/tmp/testdata"
      //}
    }

    //volume {
    // name = "testVol"
    // storage_type = "EmptyDir"
    //}

    min_replicas = 1
    max_replicas = 4

    revision_suffix = "%[3]s"
  }

  ingress {
    allow_insecure_connections = true
    is_external                = true
    target_port                = 5000
    transport                  = "auto"

    traffic_weight {
      latest_revision = true
      weight          = 20
    }

    traffic_weight {
      revision_suffix = "rev1"
      weight          = 80
    }
  }

  registry {
    server                    = azurerm_container_registry.test.login_server
    username                  = azurerm_container_registry.test.admin_username
    password_secret_reference = "registry-password"
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
      name   = "acctest-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.5
      memory = "1Gi"
      //args    = ["HOSTNAME"] // TODO - Add a container Image where args and command can be used
      //command = ["node"]

      readiness_probe {
        transport         = "http"
        port              = 5000
        path              = "/uptime"
        timeout           = 2
        failure_threshold = 1
        success_threshold = 1

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }
      }

      liveness_probe {
        transport = "http"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay     = 5
        interval          = 20
        timeout           = 2
        failure_threshold = 3
      }

      startup_probe {
        transport         = "tcp"
        port              = 5000
        timeout           = 5
        failure_threshold = 1
      }

      //volume_mounts {
      // name = "testVol"
      // path = "/tmp/testdata"
      //}
    }

    //volume {
    // name = "testVol"
    // storage_type = "EmptyDir"
    //}

    min_replicas = 1
    max_replicas = 4

    revision_suffix = "%[3]s"
  }

  ingress {
    allow_insecure_connections = true
    is_external                = true
    target_port                = 5000
    transport                  = "auto"

    traffic_weight {
      latest_revision = true
      weight          = 20
    }

    traffic_weight {
      revision_suffix = "rev1"
      weight          = 80
    }
  }

  registry {
    server                    = azurerm_container_registry.test.login_server
    username                  = azurerm_container_registry.test.admin_username
    password_secret_reference = "registry-password"
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

func (r ContainerAppResource) completeUpdate2(data acceptance.TestData, revisionSuffix string) string {
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
      cpu    = 1.0
      memory = "2Gi"
    }
    revision_suffix = "%[3]s"
  }

  secret {
    name  = "doesnt-matter"
    value = "anything"
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
