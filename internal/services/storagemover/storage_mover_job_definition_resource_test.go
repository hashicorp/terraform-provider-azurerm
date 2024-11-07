// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/jobdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageMoverJobDefinitionTestResource struct{}

func TestAccStorageMoverJobDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_job_definition", "test")
	r := StorageMoverJobDefinitionTestResource{}
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

func TestAccStorageMoverJobDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_job_definition", "test")
	r := StorageMoverJobDefinitionTestResource{}
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

func TestAccStorageMoverJobDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_job_definition", "test")
	r := StorageMoverJobDefinitionTestResource{}
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

func TestAccStorageMoverJobDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_job_definition", "test")
	r := StorageMoverJobDefinitionTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageMoverJobDefinitionTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := jobdefinitions.ParseJobDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.StorageMover.JobDefinitionsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r StorageMoverJobDefinitionTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_storage_mover_agent" "test" {
  name                     = "acctest-sa-%[2]d"
  storage_mover_id         = azurerm_storage_mover.test.id
  arc_virtual_machine_id   = data.azurerm_arc_machine.test.id
  arc_virtual_machine_uuid = data.azurerm_arc_machine.test.vm_uuid
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}

resource "azurerm_storage_account" "test" {
  name                            = "accsa%[4]s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "test" {
  name                  = "acccontainer%[4]s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_mover_target_endpoint" "test" {
  name                   = "acctest-smte-%[2]d"
  storage_mover_id       = azurerm_storage_mover.test.id
  storage_account_id     = azurerm_storage_account.test.id
  storage_container_name = azurerm_storage_container.test.name
}

resource "azurerm_storage_mover_source_endpoint" "test" {
  name             = "acctest-smse-%[2]d"
  storage_mover_id = azurerm_storage_mover.test.id
  host             = "192.168.0.1"
}

resource "azurerm_storage_mover_project" "test" {
  name             = "acctest-sp-%[2]d"
  storage_mover_id = azurerm_storage_mover.test.id
}
`, StorageMoverAgentTestResource{}.template(data), data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageMoverJobDefinitionTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%s

resource "azurerm_storage_mover_job_definition" "test" {
  name                     = "acctest-sjd-%d"
  storage_mover_project_id = azurerm_storage_mover_project.test.id
  agent_name               = azurerm_storage_mover_agent.test.name
  copy_mode                = "Additive"
  source_name              = azurerm_storage_mover_source_endpoint.test.name
  target_name              = azurerm_storage_mover_target_endpoint.test.name
}
`, template, data.RandomInteger)
}

func (r StorageMoverJobDefinitionTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_mover_job_definition" "import" {
  name                     = azurerm_storage_mover_job_definition.test.name
  storage_mover_project_id = azurerm_storage_mover_job_definition.test.storage_mover_project_id
  agent_name               = azurerm_storage_mover_job_definition.test.agent_name
  copy_mode                = azurerm_storage_mover_job_definition.test.copy_mode
  source_name              = azurerm_storage_mover_job_definition.test.source_name
  target_name              = azurerm_storage_mover_job_definition.test.target_name
}
`, config)
}

func (r StorageMoverJobDefinitionTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%s

resource "azurerm_storage_mover_job_definition" "test" {
  name                     = "acctest-sjd-%d"
  storage_mover_project_id = azurerm_storage_mover_project.test.id
  agent_name               = azurerm_storage_mover_agent.test.name
  copy_mode                = "Additive"
  source_name              = azurerm_storage_mover_source_endpoint.test.name
  source_sub_path          = "/"
  target_name              = azurerm_storage_mover_target_endpoint.test.name
  target_sub_path          = "/"
  description              = "Example Job Definition Description"
}
`, template, data.RandomInteger)
}

func (r StorageMoverJobDefinitionTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%s

resource "azurerm_storage_mover_job_definition" "test" {
  name                     = "acctest-sjd-%[2]d"
  storage_mover_project_id = azurerm_storage_mover_project.test.id
  agent_name               = azurerm_storage_mover_agent.test.name
  copy_mode                = "Additive"
  source_name              = azurerm_storage_mover_source_endpoint.test.name
  source_sub_path          = "/"
  target_name              = azurerm_storage_mover_target_endpoint.test.name
  target_sub_path          = "/"
  description              = "Update example Job Definition Description"
}
`, template, data.RandomInteger)
}
