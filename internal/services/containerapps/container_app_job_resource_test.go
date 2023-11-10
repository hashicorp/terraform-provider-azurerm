package containerapps_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"testing"
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

func TestFAccContainerAppJob_scheduleTrigger(t *testing.T) {
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

func (r ContainerAppJobResource) eventTrigger(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_container_app_job" "test" {
  name                = "acctest-cajob%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id
  
  configuration {
	trigger_type = "Event"
	replica_timeout = 10
 	replica_retry_limit = 10
	event_trigger_config {
	  parallelism = 4
	  replica_completion_count = 1
	  scale {
	    max_executions = 10
	    min_executions = 1
		polling_interval = 10
	    rules {
		  metadata = {
		    topic_name = "my-topic"
		  }
		  name = "servicebuscalingrule"
		  type = "azure-servicebus"
		}
	  }
	}
  }

  template {
    containers {
	  image = "repo/testcontainerAppsJob0:v1"
  	  name = "testcontainerappsjob0"
	  probes {
	    http_get {
          http_headers {
			name = "testheader"
			value = "testvalue"
		  }
		  path = "/testpath"
		  port = 8080
		}
		initial_delay_seconds = 10
		period_seconds = 10
		type = "Liveness"
      }
	  resources {
	    cpu = 0.5
		memory = "1Gi"
	  }
	}

	init_containers {
	  args = ["testarg"]
	  command = ["testcommand"]
	  image = "repo/testcontainerAppsJob0:v1"
	  name = "testcontainerappsjob0"
	  resources {
	    cpu = 0.5
		memory = "1Gi"
	  }
	}
  }	
}
`, template, data.RandomInteger)
}

func (r ContainerAppJobResource) manualTrigger(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_container_app_job" "test" {
  name = "acctest-cajob%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  configuration {
	trigger_type = "Manual"
	replica_timeout = 10
	replica_retry_limit = 10
	manual_trigger_config {
	  parallelism = 4
	  replica_completion_count = 1
	}
  }
  
  template {
	containers {
	  image = "repo/testcontainerAppsJob0:v1"
	  name = "testcontainerappsjob0"
	  probes {
	    http_get {
		  http_headers {
		    name = "testheader"
			value = "testvalue"
		  }
		  path = "/testpath"
		  port = 8080
		  host = "testhost"
		  scheme = "HTTPS"
		}
		initial_delay_seconds = 10
		period_seconds = 10
		type = "Liveness"
		failure_threshold = 1
		success_threshold = 1
		timeout_seconds = 1
	  }
	  resources {
		cpu = 0.5
		memory = "1Gi"
	  }
	}

    init_containers {
	  args = ["testarg"]
	  command = ["testcommand"]
	  image = "repo/testcontainerAppsJob0:v1"
	  name = "testcontainerappsjob0"
	  resources {
	    cpu = 0.5
		memory = "1Gi"
	  }
	}
  }
}
`, template, data.RandomInteger)
}

func (r ContainerAppJobResource) scheduleTrigger(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_container_app_job" "test" {
  name                = "acctest-cajob%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  configuration {
	trigger_type = "Schedule"
	replica_timeout = 1800
	replica_retry_limit = 0
	schedule_trigger_config {
	  cron_expression = "*/1 * * * *"
	  parallelism = 1
	  replica_completion_count = 1
	}
  }
  template {
	volumes {
	  name = "appsettings-volume"
	  storage_type = "EmptyDir"
	}
	containers {
	  image = "repo/testcontainerAppsJob0:v1"
	  name = "testcontainerappsjob0"
	  probes {
	    tcp_socket {
		  host = "testhost"
		  port = 8080
		}
		initial_delay_seconds = 5
		timeout_seconds = 1
		success_threshold = 1
		failure_threshold = 1
		period_seconds = 3
		type = "Liveness"
	  }
	  resources {
		cpu = 0.5
		memory = "1Gi"
	  }
	  volume_mounts {
		mount_path = "/appsettings"
	    volume_name = "appsettings-volume"
	  }
	}
  }
}
`, template, data.RandomInteger)
}

func (r ContainerAppJobResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_container_app_job" "test" {
  name                = "acctest-cajob%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id

  configuration {
	trigger_type = "Manual"
	replica_timeout = 10
	replica_retry_limit = 10
    manual_trigger_config {
	  parallelism = 4
	  replica_completion_count = 1
	}
  }
  
  template {
    containers {
	  image = "repo/testcontainerAppsJob0:v1"
  	  name = "testcontainerappsjob0"
	  probes {
	    http_get {
          http_headers {
			name = "testheader"
			value = "testvalue"
		  }
		  path = "/testpath"
		  port = 8080
		}
		initial_delay_seconds = 10
		period_seconds = 10
		type = "Liveness"
      }
	  resources {
	    cpu = 0.5
		memory = "1Gi"
	  }
	}

	init_containers {
	  args = ["testarg"]
	  command = ["testcommand"]
	  image = "repo/testcontainerAppsJob0:v1"
	  name = "testcontainerappsjob0"
	  resources {
	    cpu = 0.5
		memory = "1Gi"
	  }
	}
  }
}
`, template, data.RandomInteger)
}

func (r ContainerAppJobResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
	features {}
}

resource "azurerm_resource_group" "test" {
	name    = "acctest-CAJob%[1]d"
	location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name = "acctest-LAW%[1]d"
  location = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku = "PerGB2018"
  retention_in_days = 30
}

resource "azurerm_container_app_environment" "test" {
  name                = "acctest-CAEnv%[1]d"
  location = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}
