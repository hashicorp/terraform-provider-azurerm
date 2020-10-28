package batch_test

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccBatchPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.target_dedicated_nodes", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.#", "0"),
				),
			},
			data.ImportStep("stop_pending_resize_operation"),
		},
	})
}

func TestAccBatchPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.target_dedicated_nodes", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.#", "0"),
				),
			},
			{
				Config:      testAccBatchPool_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_batch_pool"),
			},
		},
	})
}

func TestAccBatchPool_fixedScale_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchPool_fixedScale_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(data.ResourceName, "max_tasks_per_node", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_scale.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.target_dedicated_nodes", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.resize_timeout", "PT15M"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.target_low_priority_nodes", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.#", "0"),
				),
			},
			data.ImportStep("stop_pending_resize_operation"),
		},
	})
}

func TestAccBatchPool_autoScale_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchPool_autoScale_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_scale.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_scale.0.evaluation_interval", "PT15M"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.#", "0"),
				),
			},
			data.ImportStep("stop_pending_resize_operation"),
		},
	})
}

func TestAccBatchPool_completeUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchPool_fixedScale_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_scale.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.target_dedicated_nodes", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.resize_timeout", "PT15M"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.target_low_priority_nodes", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.#", "0"),
				),
			},
			data.ImportStep("stop_pending_resize_operation"),
			{
				Config: testAccBatchPool_autoScale_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_scale.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_scale.0.evaluation_interval", "PT15M"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.#", "0"),
				),
			},
			data.ImportStep("stop_pending_resize_operation"),
		},
	})
}

func TestAccBatchPoolStartTask_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchPoolStartTask_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_scale.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.target_dedicated_nodes", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.resize_timeout", "PT15M"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.target_low_priority_nodes", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.0.max_task_retry_count", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.0.environment.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.0.environment.env", "TEST"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.0.environment.bu", "Research&Dev"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.0.user_identity.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.0.user_identity.0.auto_user.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.0.user_identity.0.auto_user.0.scope", "Task"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.0.user_identity.0.auto_user.0.elevation_level", "NonAdmin"),
				),
			},
			data.ImportStep("stop_pending_resize_operation"),
		},
	})
}

func TestAccBatchPool_certificates(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	certificate0ID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-312d31a79fa0cef49c00f769afc2b73e9f4edf34", subscriptionID, data.RandomInteger, data.RandomString)
	certificate1ID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-42c107874fd0e4a9583292a2f1098e8fe4b2edda", subscriptionID, data.RandomInteger, data.RandomString)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchPoolCertificates(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.id", certificate0ID),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.store_location", "CurrentUser"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.store_name", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.visibility.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.visibility.3294600504", "StartTask"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.1.id", certificate1ID),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.1.store_location", "CurrentUser"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.1.store_name", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.1.visibility.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.1.visibility.3294600504", "StartTask"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.1.visibility.4077195354", "RemoteUser"),
				),
			},
			data.ImportStep("stop_pending_resize_operation"),
		},
	})
}

func TestAccBatchPool_validateResourceFileWithoutSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccBatchPoolValidateResourceFileWithoutSource(data),
				ExpectError: regexp.MustCompile("Exactly one of auto_storage_container_name, storage_container_url and http_url must be specified"),
			},
		},
	})
}

func TestAccBatchPool_container(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchPoolContainerConfiguration(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "container_configuration.0.type", "DockerCompatible"),
					resource.TestCheckResourceAttr(data.ResourceName, "container_configuration.0.container_image_names.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container_configuration.0.container_registries.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "container_configuration.0.container_registries.0.registry_server", "myContainerRegistry.azurecr.io"),
					resource.TestCheckResourceAttr(data.ResourceName, "container_configuration.0.container_registries.0.user_name", "myUserName"),
					resource.TestCheckResourceAttr(data.ResourceName, "container_configuration.0.container_registries.0.password", "myPassword"),
				),
			},
			data.ImportStep(
				"stop_pending_resize_operation",
				"container_configuration.0.container_registries.0.password",
			),
		},
	})
}

func TestAccBatchPool_validateResourceFileWithMultipleSources(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccBatchPoolValidateResourceFileWithMultipleSources(data),
				ExpectError: regexp.MustCompile("Exactly one of auto_storage_container_name, storage_container_url and http_url must be specified"),
			},
		},
	})
}

func TestAccBatchPool_validateResourceFileBlobPrefixWithoutAutoStorageContainerUrl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccBatchPoolValidateResourceFileBlobPrefixWithoutAutoStorageContainerName(data),
				ExpectError: regexp.MustCompile("auto_storage_container_name or storage_container_url must be specified when using blob_prefix"),
			},
		},
	})
}

func TestAccBatchPool_validateResourceFileHttpURLWithoutFilePath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccBatchPoolValidateResourceFileHttpURLWithoutFilePath(data),
				ExpectError: regexp.MustCompile("file_path must be specified when using http_url"),
			},
		},
	})
}

func TestAccBatchPool_customImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchPoolCustomImageConfiguration(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(data.ResourceName, "max_tasks_per_node", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_scale.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.target_dedicated_nodes", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.resize_timeout", "PT15M"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.target_low_priority_nodes", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.#", "0"),
				),
			},
			data.ImportStep("stop_pending_resize_operation"),
		},
	})
}

func TestAccBatchPool_frontEndPortRanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchPool_networkConfiguration(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.sku", "16.04.0-LTS"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_image_reference.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_scale.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "fixed_scale.0.target_dedicated_nodes", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_task.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_configuration.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_configuration.0.subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_configuration.0.public_ips.#", "1"),
				),
			},
			data.ImportStep("stop_pending_resize_operation"),
		},
	})
}

func TestAccBatchPool_fixedScaleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchPool_fixedScale_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
				),
			},
			data.ImportStep("stop_pending_resize_operation"),
			{
				Config: testAccBatchPool_fixedScale_completeUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchPoolExists(data.ResourceName),
				),
			},
			data.ImportStep("stop_pending_resize_operation"),
		},
	})
}

func testCheckBatchPoolExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Batch.PoolClient

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id, err := parse.BatchPoolID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on batchPoolClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Batch pool %q (account: %q, resource group: %q) does not exist", id.Name, id.AccountName, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckBatchPoolDestroy(s *terraform.State) error {
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Batch.PoolClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_batch_pool" {
			continue
		}

		id, err := parse.BatchPoolID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := conn.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testAccBatchPool_fixedScale_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchpool"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
  storage_account_id   = azurerm_storage_account.test.id

  tags = {
    env = "test"
  }
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  display_name        = "Test Acc Pool"
  vm_size             = "Standard_A1"
  max_tasks_per_node  = 2
  node_agent_sku_id   = "batch.node.ubuntu 16.04"

  fixed_scale {
    target_dedicated_nodes = 2
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  metadata = {
    tagName = "Example tag"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}

func testAccBatchPool_fixedScale_completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchpool"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
  storage_account_id   = azurerm_storage_account.test.id

  tags = {
    env = "test"
  }
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  display_name        = "Test Acc Pool"
  vm_size             = "Standard_A1"
  max_tasks_per_node  = 2
  node_agent_sku_id   = "batch.node.ubuntu 16.04"

  fixed_scale {
    target_dedicated_nodes = 3
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  metadata = {
    tagName = "Example tag"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}

func testAccBatchPool_autoScale_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchpool"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
  storage_account_id   = azurerm_storage_account.test.id

  tags = {
    env = "test"
  }
}

resource "azurerm_batch_pool" "test" {
  name                          = "testaccpool%s"
  resource_group_name           = azurerm_resource_group.test.name
  account_name                  = azurerm_batch_account.test.name
  display_name                  = "Test Acc Pool"
  vm_size                       = "Standard_A1"
  node_agent_sku_id             = "batch.node.ubuntu 16.04"
  stop_pending_resize_operation = true

  auto_scale {
    evaluation_interval = "PT15M"

    formula = <<EOF
      startingNumberOfVMs = 1;
      maxNumberofVMs = 25;
      pendingTaskSamplePercent = $PendingTasks.GetSamplePercent(180 * TimeInterval_Second);
      pendingTaskSamples = pendingTaskSamplePercent < 70 ? startingNumberOfVMs : avg($PendingTasks.GetSample(180 * TimeInterval_Second));
      $TargetDedicatedNodes=min(maxNumberofVMs, pendingTaskSamples);
EOF

  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}

func testAccBatchPool_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchpool"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccBatchPool_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_batch_pool" "import" {
  name                = azurerm_batch_pool.test.name
  resource_group_name = azurerm_batch_pool.test.resource_group_name
  account_name        = azurerm_batch_pool.test.account_name
  node_agent_sku_id   = azurerm_batch_pool.test.node_agent_sku_id
  vm_size             = azurerm_batch_pool.test.vm_size

  fixed_scale {
    target_dedicated_nodes = azurerm_batch_pool.test.fixed_scale[0].target_dedicated_nodes
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}
`, testAccBatchPool_basic(data))
}

func testAccBatchPoolStartTask_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchpool"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
      env = "TEST"
      bu  = "Research&Dev"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }

    resource_file {
      http_url  = "https://raw.githubusercontent.com/terraform-providers/terraform-provider-azurerm/master/README.md"
      file_path = "README.md"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccBatchPoolValidateResourceFileWithoutSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
      env = "TEST"
      bu  = "Research&Dev"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }

    resource_file {
      # no valid values for sources
      auto_storage_container_name = ""
      file_mode                   = "0770"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccBatchPoolValidateResourceFileWithMultipleSources(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
      env = "TEST"
      bu  = "Research&Dev"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }

    resource_file {
      auto_storage_container_name = "test"
      http_url                    = "test"
      file_path                   = "README.md"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccBatchPoolValidateResourceFileBlobPrefixWithoutAutoStorageContainerName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
      env = "TEST"
      bu  = "Research&Dev"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }

    resource_file {
      http_url    = "test"
      blob_prefix = "test"
      file_path   = "README.md"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccBatchPoolValidateResourceFileHttpURLWithoutFilePath(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
      env = "TEST"
      bu  = "Research&Dev"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }

    resource_file {
      http_url  = "test"
      file_path = ""
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccBatchPoolCertificates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_certificate" "testcer" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate.cer")
  format               = "Cer"
  thumbprint           = "312d31a79fa0cef49c00f769afc2b73e9f4edf34" # deliberately using lowercase here as verification
  thumbprint_algorithm = "SHA1"
}

resource "azurerm_batch_certificate" "testpfx" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate.pfx")
  format               = "Pfx"
  password             = "terraform"
  thumbprint           = "42c107874fd0e4a9583292a2f1098e8fe4b2edda"
  thumbprint_algorithm = "SHA1"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  certificate {
    id             = azurerm_batch_certificate.testcer.id
    store_location = "CurrentUser"
    visibility     = ["StartTask"]
  }

  certificate {
    id             = azurerm_batch_certificate.testpfx.id
    store_location = "CurrentUser"
    visibility     = ["StartTask", "RemoteUser"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccBatchPoolContainerConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testregistry%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "microsoft-azure-batch"
    offer     = "ubuntu-server-container"
    sku       = "16-04-lts"
    version   = "latest"
  }

  container_configuration {
    type                  = "DockerCompatible"
    container_image_names = ["centos7"]
    container_registries {
      registry_server = "myContainerRegistry.azurecr.io"
      user_name       = "myUserName"
      password        = "myPassword"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}

func testAccBatchPoolCustomImageConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchpool"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name               = "internal"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "acctestpip%d"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "Dev"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_virtual_machine" "testsource" {
  name                  = "acctestvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testsource.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "30"
  }

  os_profile {
    computer_name  = "acctest-%d"
    admin_username = "tfuser"
    admin_password = "P@ssW0RD7890"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_image" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = azurerm_virtual_machine.testsource.storage_os_disk[0].vhd_uri
    size_gb  = 30
    caching  = "None"
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"

  tags = {
    env = "test"
  }
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  display_name        = "Test Acc Pool"
  vm_size             = "Standard_A1"
  max_tasks_per_node  = 2
  node_agent_sku_id   = "batch.node.ubuntu 16.04"

  fixed_scale {
    target_dedicated_nodes = 2
  }

  storage_image_reference {
    id = azurerm_image.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString)
}

func testAccBatchPool_networkConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%[1]d-batchpool"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name               = "internal"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "acctest-publicip-%[1]d"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%[3]s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%[3]s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }

  network_configuration {
    subnet_id  = azurerm_subnet.test.id
    public_ips = [azurerm_public_ip.test.id]

    endpoint_configuration {
      name                = "SSH"
      protocol            = "TCP"
      backend_port        = 22
      frontend_port_range = "4000-4100"

      network_security_group_rules {
        access                = "Deny"
        priority              = 1001
        source_address_prefix = "*"
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
