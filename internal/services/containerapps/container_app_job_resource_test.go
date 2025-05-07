// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerAppJobResource struct{}

func (r ContainerAppJobResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := jobs.ParseJobID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ContainerApps.JobClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func TestAccContainerAppJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

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

func TestAccContainerAppJob_withWorkloadProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withWorkloadProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppJob_withSystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

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

func TestAccContainerAppJob_withUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppJob_withSystemAndUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSystemAndUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppJob_withIdentityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSystemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withUserAssignedIdentity(data),
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

func TestAccContainerAppJob_eventTrigger(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventTrigger(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppJob_manualTrigger(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualTrigger(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppJob_scheduleTrigger(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scheduleTrigger(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppJob_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppJob_withKeyVaultSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withKeyVaultSecret(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppJob_withKeyVaultSecretVersioningUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withKeyVaultSecret(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withKeyVaultSecretVersionless(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppJob_withKeyVaultSecretIdentityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_job", "test")
	r := ContainerAppJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withKeyVaultSecretUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withKeyVaultSecretSystemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withKeyVaultSecretSystemAndUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ContainerAppJobResource) eventTrigger(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  event_trigger_config {
    parallelism = 4
    scale {
      max_executions              = 10
      min_executions              = 1
      polling_interval_in_seconds = 10
      rules {
        authentication {
          secret_name       = "my-secret"
          trigger_parameter = "my-trigger-parameter"
        }
        metadata = {
          topic_name = "my-topic"
        }
        name             = "servicebuscalingrule"
        custom_rule_type = "azure-servicebus"
      }
    }
  }

  template {
    container {
      image = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name  = "testcontainerappsjob0"
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
      cpu    = 0.5
      memory = "1Gi"
    }
  }
}
`, template, data.RandomInteger)
}

func (r ContainerAppJobResource) manualTrigger(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  template {
    container {
      image = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name  = "testcontainerappsjob0"
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

      cpu    = 0.5
      memory = "1Gi"
    }
  }
}
`, template, data.RandomInteger)
}

func (r ContainerAppJobResource) scheduleTrigger(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  replica_timeout_in_seconds = 1800
  replica_retry_limit        = 10
  schedule_trigger_config {
    cron_expression          = "*/1 * * * *"
    parallelism              = 1
    replica_completion_count = 1
  }
  template {
    volume {
      name         = "appsettings-volume"
      storage_type = "EmptyDir"
    }
    container {
      image = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name  = "testcontainerappsjob0"
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
      cpu    = 0.5
      memory = "1Gi"
      volume_mounts {
        path = "/appsettings"
        name = "appsettings-volume"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r ContainerAppJobResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  template {
    container {
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name   = "testcontainerappsjob0"
      cpu    = 0.5
      memory = "1Gi"
    }
  }
}
`, template, data.RandomInteger)
}

func (r ContainerAppJobResource) withKeyVaultSecret(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctest-kv-%[3]s"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "Create",
      "Get",
    ]
    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id
    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[3]s"
  value        = "test-secret"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  template {
    container {
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name   = "testcontainerappsjob0"
      cpu    = 0.5
      memory = "1Gi"
      env {
        name        = "key-vault-secret"
        secret_name = "key-vault-secret"
      }
    }
  }

  secret {
    name                = "key-vault-secret"
    identity            = azurerm_user_assigned_identity.test.id
    key_vault_secret_id = azurerm_key_vault_secret.test.id
  }
}
`, template, data.RandomInteger, data.RandomString)
}

func (r ContainerAppJobResource) withKeyVaultSecretVersionless(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctest-kv-%[3]s"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "Create",
      "Get",
    ]
    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id
    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[3]s"
  value        = "test-secret"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  template {
    container {
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name   = "testcontainerappsjob0"
      cpu    = 0.5
      memory = "1Gi"
      env {
        name        = "key-vault-secret"
        secret_name = "key-vault-secret"
      }
    }
  }

  secret {
    name                = "key-vault-secret"
    identity            = azurerm_user_assigned_identity.test.id
    key_vault_secret_id = azurerm_key_vault_secret.test.versionless_id
  }
}
`, template, data.RandomInteger, data.RandomString)
}

func (r ContainerAppJobResource) withKeyVaultSecretUserIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctest-kv-%[3]s"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[3]s"
  value        = "test-secret"
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [
    azurerm_key_vault_access_policy.self_key_vault_admin
  ]
}

resource "azurerm_key_vault_access_policy" "self_key_vault_admin" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
  ]

  secret_permissions = [
    "Set",
    "Get",
    "Delete",
    "Purge",
    "Recover"
  ]
}

resource "azurerm_key_vault_access_policy" "mi_key_vault_secrets" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_container_app_job.test.identity[0].principal_id

  secret_permissions = [
    "Get",
  ]
}

resource "azurerm_key_vault_access_policy" "user_mi_key_vault_secrets" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  secret_permissions = [
    "Get",
  ]
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  template {
    container {
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name   = "testcontainerappsjob0"
      cpu    = 0.5
      memory = "1Gi"
      env {
        name        = "key-vault-secret"
        secret_name = "key-vault-secret"
      }
    }
  }

  secret {
    name                = "key-vault-secret"
    identity            = azurerm_user_assigned_identity.test.id
    key_vault_secret_id = azurerm_key_vault_secret.test.id
  }

  depends_on = [
    azurerm_key_vault_access_policy.user_mi_key_vault_secrets
  ]
}
`, template, data.RandomInteger, data.RandomString)
}

func (r ContainerAppJobResource) withKeyVaultSecretSystemAndUserIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctest-kv-%[3]s"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[3]s"
  value        = "test-secret"
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [
    azurerm_key_vault_access_policy.self_key_vault_admin
  ]
}

resource "azurerm_key_vault_access_policy" "self_key_vault_admin" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
  ]

  secret_permissions = [
    "Set",
    "Get",
    "Delete",
    "Purge",
    "Recover"
  ]
}

resource "azurerm_key_vault_access_policy" "mi_key_vault_secrets" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_container_app_job.test.identity[0].principal_id

  secret_permissions = [
    "Get",
  ]
}

resource "azurerm_key_vault_access_policy" "user_mi_key_vault_secrets" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  secret_permissions = [
    "Get",
  ]
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  template {
    container {
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name   = "testcontainerappsjob0"
      cpu    = 0.5
      memory = "1Gi"
      env {
        name        = "key-vault-secret"
        secret_name = "key-vault-secret"
      }
      env {
        name        = "key-vault-secret-system"
        secret_name = "key-vault-secret-system"
      }
    }
  }

  secret {
    name                = "key-vault-secret"
    identity            = azurerm_user_assigned_identity.test.id
    key_vault_secret_id = azurerm_key_vault_secret.test.id
  }

  secret {
    name                = "key-vault-secret-system"
    identity            = "System"
    key_vault_secret_id = azurerm_key_vault_secret.test.id
  }

  depends_on = [
    azurerm_key_vault_access_policy.user_mi_key_vault_secrets
  ]
}
`, template, data.RandomInteger, data.RandomString)
}

func (r ContainerAppJobResource) withKeyVaultSecretSystemIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctest-kv-%[3]s"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[3]s"
  value        = "test-secret"
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [
    azurerm_key_vault_access_policy.self_key_vault_admin
  ]
}

resource "azurerm_key_vault_access_policy" "self_key_vault_admin" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
  ]

  secret_permissions = [
    "Set",
    "Get",
    "Delete",
    "Purge",
    "Recover"
  ]
}

resource "azurerm_key_vault_access_policy" "mi_key_vault_secrets" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_container_app_job.test.identity[0].principal_id

  secret_permissions = [
    "Get",
  ]
}

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  identity {
    type = "SystemAssigned"
  }

  template {
    container {
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name   = "testcontainerappsjob0"
      cpu    = 0.5
      memory = "1Gi"
      env {
        name        = "key-vault-secret-system"
        secret_name = "key-vault-secret-system"
      }
    }
  }

  secret {
    name                = "key-vault-secret-system"
    identity            = "System"
    key_vault_secret_id = azurerm_key_vault_secret.test.id
  }
}
`, template, data.RandomInteger, data.RandomString)
}

func (r ContainerAppJobResource) withUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  template {
    container {
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name   = "testcontainerappsjob0"
      cpu    = 0.5
      memory = "1Gi"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppJobResource) withSystemIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  identity {
    type = "SystemAssigned"
  }

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  template {
    container {
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name   = "testcontainerappsjob0"
      cpu    = 0.5
      memory = "1Gi"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppJobResource) withSystemAndUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  template {
    container {
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name   = "testcontainerappsjob0"
      cpu    = 0.5
      memory = "1Gi"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppJobResource) withWorkloadProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

locals {
  workload_profiles = tolist(azurerm_container_app_environment.test.workload_profile)
}

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id
  workload_profile_name        = local.workload_profiles.0.name

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  template {
    container {
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name   = "testcontainerappsjob0"
      cpu    = 0.5
      memory = "1Gi"
    }
  }
}
`, r.templateWorkloadProfile(data), data.RandomInteger)
}

func (r ContainerAppJobResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  identity {
    type = "SystemAssigned"
  }

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }
  secret {
    name  = "registry-password"
    value = azurerm_container_registry.test.admin_password
  }
  registry {
    server               = azurerm_container_registry.test.login_server
    username             = azurerm_container_registry.test.admin_username
    password_secret_name = "registry-password"
  }

  template {
    volume {
      name          = azurerm_container_app_environment_storage.test.name
      storage_type  = "AzureFile"
      storage_name  = azurerm_container_app_environment_storage.test.name
      mount_options = "dir_mode=0777,file_mode=0666"
    }
    container {
      args = [
        "-c",
        "while true; do echo hello; sleep 10;done",
      ]
      command = [
        "/bin/sh",
      ]
      image = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name  = "testcontainerappsjob0"
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

      cpu    = 0.5
      memory = "1Gi"
      volume_mounts {
        path = "/appsettings"
        name = azurerm_container_app_environment_storage.test.name
      }
    }

    init_container {
      name   = "init-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"
      volume_mounts {
        name = azurerm_container_app_environment_storage.test.name
        path = "/appsettings"
      }
    }
  }
  tags = {
    ENV = "test"
  }
}
`, ContainerAppResource{}.templatePlusExtras(data), data.RandomInteger)
}

func (r ContainerAppJobResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_container_app_job" "test" {
  name                         = "acctest-cajob%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  identity {
    type = "SystemAssigned"
  }

  replica_timeout_in_seconds = 20
  replica_retry_limit        = 20

  manual_trigger_config {
    parallelism              = 5
    replica_completion_count = 2
  }

  secret {
    name  = "registry-password"
    value = azurerm_container_registry.test.admin_password
  }

  secret {
    name  = "foo"
    value = "bar"
  }

  registry {
    server               = azurerm_container_registry.test.login_server
    username             = azurerm_container_registry.test.admin_username
    password_secret_name = "registry-password"
  }

  template {
    volume {
      name          = azurerm_container_app_environment_storage.test.name
      storage_type  = "AzureFile"
      storage_name  = azurerm_container_app_environment_storage.test.name
      mount_options = "dir_mode=0777,file_mode=0666"
    }
    container {
      args = [
        "-c",
        "while true; do echo hello; sleep 10;done",
      ]
      command = [
        "/bin/sh",
      ]
      image = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      name  = "testcontainerappsjob0"
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

        initial_delay           = 10
        interval_seconds        = 40
        timeout                 = 3
        failure_count_threshold = 2
      }
      startup_probe {
        transport               = "TCP"
        port                    = 5000
        timeout                 = 5
        failure_count_threshold = 3
      }

      cpu    = 0.25
      memory = "0.5Gi"
      volume_mounts {
        path = "/appsettings"
        name = azurerm_container_app_environment_storage.test.name
      }
    }

    init_container {
      name   = "init-cont-%[2]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"
      volume_mounts {
        name = azurerm_container_app_environment_storage.test.name
        path = "/appsettings"
      }
    }
  }
  tags = {
    ENV = "test"
  }
}
`, ContainerAppResource{}.templatePlusExtras(data), data.RandomInteger)
}

func (r ContainerAppJobResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-CAJob%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-LAW%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ContainerAppJobResource) templateWorkloadProfile(data acceptance.TestData) string {
	return ContainerAppEnvironmentResource{}.complete(data)
}
