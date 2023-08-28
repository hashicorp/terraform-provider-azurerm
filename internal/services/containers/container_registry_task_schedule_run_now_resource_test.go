// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerRegistryTaskScheduleResource struct {
	githubRepo
}

func TestAccContainerRegistryTaskSchedule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_task_schedule_run_now", "test")

	preCheckGithubRepo(t)

	r := ContainerRegistryTaskScheduleResource{
		githubRepo: githubRepo{
			url:   os.Getenv("ARM_TEST_ACR_TASK_GITHUB_REPO_URL"),
			token: os.Getenv("ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN"),
		},
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, r.dockerTaskStep),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, r.fileTaskStep),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, r.encodedTaskStep),
		},
		data.ImportStep(),
	})
}

func (r ContainerRegistryTaskScheduleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	ret := false
	return &ret, nil
}

func (r ContainerRegistryTaskScheduleResource) basic(data acceptance.TestData, tpl func(data acceptance.TestData) string) string {
	template := tpl(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_task_schedule_run_now" "test" {
  container_registry_task_id = azurerm_container_registry_task.test.id
}
`, template)
}

func (r ContainerRegistryTaskScheduleResource) dockerTaskStep(data acceptance.TestData) string {
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

func (r ContainerRegistryTaskScheduleResource) fileTaskStep(data acceptance.TestData) string {
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

func (r ContainerRegistryTaskScheduleResource) encodedTaskStep(data acceptance.TestData) string {
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

func (r ContainerRegistryTaskScheduleResource) template(data acceptance.TestData) string {
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
