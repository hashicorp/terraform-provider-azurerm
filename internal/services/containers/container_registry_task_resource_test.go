package containers_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryTaskResource struct {
	githubRepo
}

type githubRepo struct {
	url   string
	token string
}

func preCheckGithubRepo(t *testing.T) {
	// - ARM_TEST_ACR_TASK_GITHUB_REPO_URL represents the user forked repo from: https://github.com/Azure-Samples/acr-build-helloworld-node
	// - ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN represents the github personal token with the appropriate permissions per: https://docs.microsoft.com/en-us/azure/container-registry/container-registry-tutorial-build-task#create-a-github-personal-access-token
	// Checkout https://docs.microsoft.com/en-us/azure/container-registry/container-registry-tutorial-build-task for details.
	variables := []string{
		"ARM_TEST_ACR_TASK_GITHUB_REPO_URL",
		"ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Skipf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

func TestAccContainerRegistryTask_dockerStep(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_task", "test")

	preCheckGithubRepo(t)

	r := ContainerRegistryTaskResource{
		githubRepo: githubRepo{
			url:   os.Getenv("ARM_TEST_ACR_TASK_GITHUB_REPO_URL"),
			token: os.Getenv("ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN"),
		},
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dockerStepBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("docker_step.0.context_access_token"),
		{
			Config: r.dockerStepUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("docker_step.0.context_access_token", "docker_step.0.secret_arguments"),
		{
			Config: r.dockerStepBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("docker_step.0.context_access_token"),
	})
}

func TestAccContainerRegistryTask_fileTaskStep(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_task", "test")

	preCheckGithubRepo(t)

	r := ContainerRegistryTaskResource{
		githubRepo: githubRepo{
			url:   os.Getenv("ARM_TEST_ACR_TASK_GITHUB_REPO_URL"),
			token: os.Getenv("ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN"),
		},
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.fileTaskStepBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("file_step.0.context_access_token"),
		{
			Config: r.fileTaskStepUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("file_step.0.context_access_token", "file_step.0.secret_values"),
		{
			Config: r.fileTaskStepBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("file_step.0.context_access_token"),
	})
}

func TestAccContainerRegistryTask_encodedTaskStep(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_task", "test")

	preCheckGithubRepo(t)

	r := ContainerRegistryTaskResource{
		githubRepo: githubRepo{
			url:   os.Getenv("ARM_TEST_ACR_TASK_GITHUB_REPO_URL"),
			token: os.Getenv("ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN"),
		},
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.encodedTaskStepBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("encoded_step.0.context_access_token"),
		{
			Config: r.encodedTaskStepUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("encoded_step.0.context_access_token", "encoded_step.0.secret_values"),
		{
			Config: r.encodedTaskStepBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("encoded_step.0.context_access_token"),
	})
}

func TestAccContainerRegistryTask_dockerStepBaseImageTrigger(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_task", "test")

	preCheckGithubRepo(t)

	r := ContainerRegistryTaskResource{
		githubRepo: githubRepo{
			url:   os.Getenv("ARM_TEST_ACR_TASK_GITHUB_REPO_URL"),
			token: os.Getenv("ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN"),
		},
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dockerStepBaseImageTrigger(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("docker_step.0.context_access_token"),
		{
			Config: r.dockerStepBaseImageTriggerUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("docker_step.0.context_access_token", "base_image_trigger.0.update_trigger_endpoint"),
		{
			Config: r.dockerStepBaseImageTrigger(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("docker_step.0.context_access_token"),
	})
}

func TestAccContainerRegistryTask_dockerStepSourceTrigger(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_task", "test")

	preCheckGithubRepo(t)

	r := ContainerRegistryTaskResource{
		githubRepo: githubRepo{
			url:   os.Getenv("ARM_TEST_ACR_TASK_GITHUB_REPO_URL"),
			token: os.Getenv("ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN"),
		},
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dockerStepSourceTrigger(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"docker_step.0.context_access_token",
			"source_trigger.0.authentication.#",
			"source_trigger.0.authentication.0.%",
			"source_trigger.0.authentication.0.expire_in_seconds",
			"source_trigger.0.authentication.0.refresh_token",
			"source_trigger.0.authentication.0.scope",
			"source_trigger.0.authentication.0.token",
		),
		{
			Config: r.dockerStepSourceTriggerUpdateDockerStep(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"docker_step.0.context_access_token",
			"source_trigger.0.authentication.#",
			"source_trigger.0.authentication.0.%",
			"source_trigger.0.authentication.0.expire_in_seconds",
			"source_trigger.0.authentication.0.refresh_token",
			"source_trigger.0.authentication.0.scope",
			"source_trigger.0.authentication.0.token",
		),
		{
			Config: r.dockerStepSourceTriggerUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"docker_step.0.context_access_token",
			"source_trigger.0.authentication.#",
			"source_trigger.0.authentication.0.%",
			"source_trigger.0.authentication.0.expire_in_seconds",
			"source_trigger.0.authentication.0.refresh_token",
			"source_trigger.0.authentication.0.scope",
			"source_trigger.0.authentication.0.token",
		),
		{
			Config: r.dockerStepSourceTrigger(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"docker_step.0.context_access_token",
			"source_trigger.0.authentication.#",
			"source_trigger.0.authentication.0.%",
			"source_trigger.0.authentication.0.expire_in_seconds",
			"source_trigger.0.authentication.0.refresh_token",
			"source_trigger.0.authentication.0.scope",
			"source_trigger.0.authentication.0.token",
		),
	})
}

func TestAccContainerRegistryTask_dockerStepTimerTrigger(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_task", "test")

	preCheckGithubRepo(t)

	r := ContainerRegistryTaskResource{
		githubRepo: githubRepo{
			url:   os.Getenv("ARM_TEST_ACR_TASK_GITHUB_REPO_URL"),
			token: os.Getenv("ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN"),
		},
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dockerStepTimerTrigger(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("docker_step.0.context_access_token"),
		{
			Config: r.dockerStepTimerTriggerUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("docker_step.0.context_access_token"),
		{
			Config: r.dockerStepTimerTrigger(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("docker_step.0.context_access_token"),
	})
}

func TestAccContainerRegistryTask_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_task", "test")

	preCheckGithubRepo(t)

	r := ContainerRegistryTaskResource{
		githubRepo: githubRepo{
			url:   os.Getenv("ARM_TEST_ACR_TASK_GITHUB_REPO_URL"),
			token: os.Getenv("ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN"),
		},
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dockerStepBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"docker_step.0.context_access_token",
		),
		{
			Config: r.dockerStepSystemIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"docker_step.0.context_access_token",
		),
		{
			Config: r.dockerStepUserAssignedIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"docker_step.0.context_access_token",
		),
		{
			Config: r.dockerStepSystemUserAssignedIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"docker_step.0.context_access_token",
		),
	})
}

func TestAccContainerRegistryTask_fileTaskStepRegistryCredential(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_task", "test")

	preCheckGithubRepo(t)

	r := ContainerRegistryTaskResource{
		githubRepo: githubRepo{
			url:   os.Getenv("ARM_TEST_ACR_TASK_GITHUB_REPO_URL"),
			token: os.Getenv("ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN"),
		},
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.fileTaskStepBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("file_step.0.context_access_token"),
		{
			Config: r.fileTaskStepRegistryCredentialPassword(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"file_step.0.context_access_token",
			"registry_credential.0.custom.#",
			"registry_credential.0.custom.0.%",
			"registry_credential.0.custom.0.identity",
			"registry_credential.0.custom.0.login_server",
			"registry_credential.0.custom.0.password",
			"registry_credential.0.custom.0.username",
		),
		{
			Config: r.fileTaskStepRegistryCredentialIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"file_step.0.context_access_token",
			"registry_credential.0.custom.#",
			"registry_credential.0.custom.0.%",
			"registry_credential.0.custom.0.identity",
			"registry_credential.0.custom.0.login_server",
			"registry_credential.0.custom.0.password",
			"registry_credential.0.custom.0.username",
		),
		{
			Config: r.fileTaskStepBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("file_step.0.context_access_token"),
	})
}

func TestAccContainerRegistryTask_systemTask(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_task", "test")
	r := ContainerRegistryTaskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.systemTask(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistryTask_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_task", "test")

	preCheckGithubRepo(t)

	r := ContainerRegistryTaskResource{
		githubRepo: githubRepo{
			url:   os.Getenv("ARM_TEST_ACR_TASK_GITHUB_REPO_URL"),
			token: os.Getenv("ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN"),
		},
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dockerStepBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r ContainerRegistryTaskResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Containers.TasksClient

	id, err := parse.ContainerRegistryTaskID(state.ID)
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.TaskName); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r ContainerRegistryTaskResource) dockerStepBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) dockerStepUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile-app"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
    arguments = {
      REGISTRY_NAME = "some.azurecr.io"
    }
    secret_arguments = {
      secret = "secret"
    }
    push_enabled  = false
    cache_enabled = false
    target        = "some_target"
  }
  agent_setting {
    cpu = 2
  }
  enabled            = false
  timeout_in_seconds = 300
  tags = {
    env = "Test"
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) fileTaskStepBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  file_step {
    task_file_path       = "taskmulti.yaml"
    context_path         = "%s"
    context_access_token = "%s"
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) fileTaskStepUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  file_step {
    task_file_path       = "taskmulti-multiregistry.yaml"
    context_path         = "%s"
    context_access_token = "%s"
    values = {
      regDate = "mycontainerregistrydate.azurecr.io"
    }
    secret_values = {
      secret = "secret"
    }
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) encodedTaskStepBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  encoded_step {
    task_content         = <<EOF
FROM node:15-alpine

COPY . /src
RUN cd /src && npm install
EXPOSE 80
CMD ["node", "/src/server.js"]
EOF
    context_path         = "%s"
    context_access_token = "%s"
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) encodedTaskStepUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  encoded_step {
    task_content         = <<EOF
ARG REGISTRY_NAME
FROM $${REGISTRY_NAME}/baseimages/node:15-alpine

COPY . /src
RUN cd /src && npm install
EXPOSE 80
CMD ["node", "/src/server.js"]
EOF
    context_path         = "%s"
    context_access_token = "%s"
    values = {
      REGISTRY_NAME = "some.azurecr.io"
    }
    secret_values = {
      secret = "secret"
    }
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) dockerStepBaseImageTrigger(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
  base_image_trigger {
    name = "default"
    type = "Runtime"
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) dockerStepBaseImageTriggerUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
  base_image_trigger {
    name                        = "default-update"
    type                        = "All"
    enabled                     = false
    update_trigger_endpoint     = "https://foo.com"
    update_trigger_payload_type = "Default"
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) dockerStepSourceTrigger(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
  source_trigger {
    name           = "default"
    events         = ["commit"]
    source_type    = "Github"
    repository_url = "%s"
    branch         = "main"
    authentication {
      token_type = "PAT"
      token      = "%s"
    }
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) dockerStepSourceTriggerUpdateDockerStep(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld2:{{.Run.ID}}"]
  }
  source_trigger {
    name           = "default"
    events         = ["commit"]
    source_type    = "Github"
    repository_url = "%s"
    branch         = "main"
    authentication {
      token_type = "PAT"
      token      = "%s"
    }
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) dockerStepSourceTriggerUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
  source_trigger {
    name           = "default-update"
    events         = ["pullrequest"]
    source_type    = "Github"
    repository_url = "%s"
    branch         = "master"
    authentication {
      token_type = "PAT"
      token      = "%s"
    }
    enabled = false
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) dockerStepTimerTrigger(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
  timer_trigger {
    name     = "default"
    schedule = "0 21 * * *"
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) dockerStepTimerTriggerUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
  timer_trigger {
    name     = "default-update"
    schedule = "0 12 * * *"
    enabled  = false
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) dockerStepSystemIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  identity {
    type = "SystemAssigned"
  }
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
}
`, template, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) dockerStepUserAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "testacccrTask-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
}
`, template, data.RandomInteger, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) dockerStepSystemUserAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "testacccrTask-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
}
`, template, data.RandomInteger, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) fileTaskStepRegistryCredentialPassword(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry" "test2" {
  name                = "testacccrtask2%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"
}

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  platform {
    os = "Linux"
  }
  file_step {
    task_file_path       = "taskmulti-multiregistry.yaml"
    context_path         = "%s"
    context_access_token = "%s"
    values = {
      regDate = azurerm_container_registry.test2.login_server
    }
  }
  registry_credential {
    custom {
      login_server = azurerm_container_registry.test2.login_server
      username     = "%s"
      password     = "%s"
    }
  }
}
`, template, data.RandomInteger, data.RandomInteger, r.githubRepo.url, r.githubRepo.token, os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_CLIENT_SECRET"))
}

func (r ContainerRegistryTaskResource) fileTaskStepRegistryCredentialIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry" "test2" {
  name                = "testacccrtask2%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"
}

resource "azurerm_container_registry_task" "test" {
  name                  = "testacccrTask%d"
  container_registry_id = azurerm_container_registry.test.id
  identity {
    type = "SystemAssigned"
  }
  platform {
    os = "Linux"
  }
  file_step {
    task_file_path       = "taskmulti-multiregistry.yaml"
    context_path         = "%s"
    context_access_token = "%s"
    values = {
      regDate = azurerm_container_registry.test2.login_server
    }
  }
  registry_credential {
    custom {
      login_server = azurerm_container_registry.test2.login_server
      identity     = "[system]"
    }
  }
}
`, template, data.RandomInteger, data.RandomInteger, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) systemTask(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "test" {
  name                  = "quicktask"
  container_registry_id = azurerm_container_registry.test.id
  is_system_task        = true
}
`, template)
}

func (r ContainerRegistryTaskResource) requiresImport(data acceptance.TestData) string {
	template := r.dockerStepBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task" "import" {
  name                  = azurerm_container_registry_task.test.name
  container_registry_id = azurerm_container_registry_task.test.container_registry_id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "%s"
    context_access_token = "%s"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
}
`, template, r.githubRepo.url, r.githubRepo.token)
}

func (r ContainerRegistryTaskResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ACRTask-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccrtask%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
