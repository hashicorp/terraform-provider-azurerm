package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryCacheRuleResource struct{}

func TestAccContainerRegistryCacheRuleName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 2,
		},
		{
			Value:    "four",
			ErrCount: 1,
		},
		{
			Value:    "5five",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 0,
		},
		{
			Value:    "hello-world-foo-bar-12345",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloWorld",
			ErrCount: 0,
		},
		{
			Value:    "helloworld12",
			ErrCount: 0,
		},
		{
			Value:    "hello@world",
			ErrCount: 1,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd3324120",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd33241202",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqfjjfewsqwcdw21ddwqwd3324120",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validate.ContainerRegistryCacheRuleName(tc.Value, "azurerm_container_registry_cache_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Container Registry Cache Rule Name to trigger a validation error: %v", errors)
		}
	}
}

func TestAccContainerRegistryCacheRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_cache_rule", "test")
	r := ContainerRegistryCacheRuleResource{}

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

func TestAccContainerRegistryCacheRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_cache_rule", "test")
	r := ContainerRegistryCacheRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_container_registry_cache_rule"),
		},
	})
}

func (t ContainerRegistryCacheRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := cacherules.ParseCacheRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.ContainerRegistryClient_v2023_07_01.CacheRules.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ContainerRegistryCacheRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
	features {}
}

resource "azurerm_resource_group" "test" {
	name     = "accTestRG-acr-cache-rule-%d"
	location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_container_registry_cache_rule" "test" {
  name                = "testacc-cr-cache-rule-%d"
  resource_group_name = azurerm_resource_group.test.name
  registry            = azurerm_container_registry.test.name
  target_repo         = "target"
  source_repo         = "docker.io/hello-world"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ContainerRegistryCacheRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_cache_rule" "import" {
  name                = azurerm_container_registry_cache_rule.test.name
  resource_group_name = azurerm_container_registry_cache_rule.test.resource_group_name
}
`, r.basic(data))
}
