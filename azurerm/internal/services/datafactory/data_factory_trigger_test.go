package datafactory_test

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func findTriggerByName(name string, s *terraform.State) (*string, *string, *string, error) {
	// Ensure we have enough information in state to look up in API
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return nil, nil, nil, fmt.Errorf("Not found: %s", name)
	}

	triggerName := rs.Primary.Attributes["name"]
	resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
	dataFactoryName := rs.Primary.Attributes["data_factory_name"]
	if !hasResourceGroup {
		return nil, nil, nil, fmt.Errorf("Bad: no resource group found in state for Data Factory: %s", name)
	}

	return &triggerName, &dataFactoryName, &resourceGroup, nil
}

func testCheckAzureRMDataFactoryTriggerStarts(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.TriggersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		name, dataFactoryName, resourceGroup, err := findTriggerByName(name, s)
		if err != nil {
			return err
		}

		future, err := client.Start(ctx, *resourceGroup, *dataFactoryName, *name)
		if err != nil {
			return err
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return err
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryTriggerStops(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.TriggersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		name, dataFactoryName, resourceGroup, err := findTriggerByName(name, s)
		if err != nil {
			return err
		}

		future, err := client.Stop(ctx, *resourceGroup, *dataFactoryName, *name)
		if err != nil {
			return err
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return err
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryTriggerExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.TriggersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		name, dataFactoryName, resourceGroup, err := findTriggerByName(name, s)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, *resourceGroup, *dataFactoryName, *name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on dataFactory.TriggersClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Data Factory Trigger %q (data factory name: %q / resource group: %q) does not exist", *name, *dataFactoryName, *resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryTriggerDestroy(resource_type string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.TriggersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resource_type {
				continue
			}

			name := rs.Primary.Attributes["name"]
			resourceGroup := rs.Primary.Attributes["resource_group_name"]
			dataFactoryName := rs.Primary.Attributes["data_factory_name"]

			resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")

			if err != nil {
				return nil
			}

			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("Data Factory Trigger still exists:\n%#v", resp.Properties)
			}
		}

		return nil
	}
}
