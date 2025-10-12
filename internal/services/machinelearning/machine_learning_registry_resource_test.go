// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/registrymanagement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MachineLearningRegistry struct{}

func TestAccMachineLearningRegistry_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_registry", "test")
	r := MachineLearningRegistry{}

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

func (r MachineLearningRegistry) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	registryClient := client.MachineLearning.RegistryManagement
	id, err := registrymanagement.ParseRegistryID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := registryClient.RegistriesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Machine Learning Data Store %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r MachineLearningRegistry) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_machine_learning_registry" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "accmlreg-%[2]d"
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func (r MachineLearningRegistry) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_machine_learning_registry" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "accmlreg-%[2]d"
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func (r MachineLearningRegistry) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
}

`, data.RandomInteger, data.Locations.Primary)
}
