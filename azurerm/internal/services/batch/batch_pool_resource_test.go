package batch_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BatchPoolResource struct {
}

func TestAccBatchPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 16.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("16.04.0-LTS"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 16.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("16.04.0-LTS"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_batch_pool"),
		},
	})
}

func TestAccBatchPool_fixedScale_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.fixedScale_complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("max_tasks_per_node").HasValue("2"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 16.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("16.04.0-LTS"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("2"),
				check.That(data.ResourceName).Key("fixed_scale.0.resize_timeout").HasValue("PT15M"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_low_priority_nodes").HasValue("0"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_autoScale_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.autoScale_complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 16.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("16.04.0-LTS"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("auto_scale.0.evaluation_interval").HasValue("PT15M"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_completeUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.fixedScale_complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 16.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("16.04.0-LTS"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("2"),
				check.That(data.ResourceName).Key("fixed_scale.0.resize_timeout").HasValue("PT15M"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_low_priority_nodes").HasValue("0"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
		{
			Config: r.autoScale_complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 16.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("16.04.0-LTS"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("auto_scale.0.evaluation_interval").HasValue("PT15M"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_startTask_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.startTask_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 16.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("16.04.0-LTS"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.resize_timeout").HasValue("PT15M"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_low_priority_nodes").HasValue("0"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.0.max_task_retry_count").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.0.environment.%").HasValue("2"),
				check.That(data.ResourceName).Key("start_task.0.environment.env").HasValue("TEST"),
				check.That(data.ResourceName).Key("start_task.0.environment.bu").HasValue("Research&Dev"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.0.auto_user.#").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.0.auto_user.0.scope").HasValue("Task"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.0.auto_user.0.elevation_level").HasValue("NonAdmin"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_certificates(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	certificate0ID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-312d31a79fa0cef49c00f769afc2b73e9f4edf34", subscriptionID, data.RandomInteger, data.RandomString)
	certificate1ID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-42c107874fd0e4a9583292a2f1098e8fe4b2edda", subscriptionID, data.RandomInteger, data.RandomString)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.certificates(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 16.04"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("2"),
				check.That(data.ResourceName).Key("certificate.0.id").HasValue(certificate0ID),
				check.That(data.ResourceName).Key("certificate.0.store_location").HasValue("CurrentUser"),
				check.That(data.ResourceName).Key("certificate.0.store_name").HasValue(""),
				check.That(data.ResourceName).Key("certificate.0.visibility.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.visibility.3294600504").HasValue("StartTask"),
				check.That(data.ResourceName).Key("certificate.1.id").HasValue(certificate1ID),
				check.That(data.ResourceName).Key("certificate.1.store_location").HasValue("CurrentUser"),
				check.That(data.ResourceName).Key("certificate.1.store_name").HasValue(""),
				check.That(data.ResourceName).Key("certificate.1.visibility.#").HasValue("2"),
				check.That(data.ResourceName).Key("certificate.1.visibility.3294600504").HasValue("StartTask"),
				check.That(data.ResourceName).Key("certificate.1.visibility.4077195354").HasValue("RemoteUser"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_validateResourceFileWithoutSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.validateResourceFileWithoutSource(data),
			ExpectError: regexp.MustCompile("Exactly one of auto_storage_container_name, storage_container_url and http_url must be specified"),
		},
	})
}

func TestAccBatchPool_container(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.containerConfiguration(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container_configuration.0.type").HasValue("DockerCompatible"),
				check.That(data.ResourceName).Key("container_configuration.0.container_image_names.#").HasValue("1"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.#").HasValue("1"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.registry_server").HasValue("myContainerRegistry.azurecr.io"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.user_name").HasValue("myUserName"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.password").HasValue("myPassword"),
			),
		},
		data.ImportStep(
			"stop_pending_resize_operation",
			"container_configuration.0.container_registries.0.password",
		),
	})
}

func TestAccBatchPool_validateResourceFileWithMultipleSources(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.validateResourceFileWithMultipleSources(data),
			ExpectError: regexp.MustCompile("Exactly one of auto_storage_container_name, storage_container_url and http_url must be specified"),
		},
	})
}

func TestAccBatchPool_validateResourceFileBlobPrefixWithoutAutoStorageContainerUrl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.validateResourceFileBlobPrefixWithoutAutoStorageContainerName(data),
			ExpectError: regexp.MustCompile("auto_storage_container_name or storage_container_url must be specified when using blob_prefix"),
		},
	})
}

func TestAccBatchPool_validateResourceFileHttpURLWithoutFilePath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.validateResourceFileHttpURLWithoutFilePath(data),
			ExpectError: regexp.MustCompile("file_path must be specified when using http_url"),
		},
	})
}

func TestAccBatchPool_customImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.customImageConfiguration(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("max_tasks_per_node").HasValue("2"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 16.04"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("2"),
				check.That(data.ResourceName).Key("fixed_scale.0.resize_timeout").HasValue("PT15M"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_low_priority_nodes").HasValue("0"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_frontEndPortRanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkConfiguration(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 16.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("16.04.0-LTS"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
				check.That(data.ResourceName).Key("network_configuration.#").HasValue("1"),
				check.That(data.ResourceName).Key("network_configuration.0.subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_configuration.0.public_ips.#").HasValue("1"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_fixedScaleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.fixedScale_complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
		{
			Config: r.fixedScale_completeUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func (t BatchPoolResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.PoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Batch.PoolClient.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Batch Pool %q (Account Name %q / Resource Group %q) does not exist", id.Name, id.BatchAccountName, id.ResourceGroup)
	}

	return utils.Bool(resp.PoolProperties != nil), nil
}

func (BatchPoolResource) fixedScale_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
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

func (BatchPoolResource) fixedScale_completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
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

func (BatchPoolResource) autoScale_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
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

func (BatchPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
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

func (BatchPoolResource) requiresImport(data acceptance.TestData) string {
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
`, BatchPoolResource{}.basic(data))
}

func (BatchPoolResource) startTask_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
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

func (BatchPoolResource) validateResourceFileWithoutSource(data acceptance.TestData) string {
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

func (BatchPoolResource) validateResourceFileWithMultipleSources(data acceptance.TestData) string {
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

func (BatchPoolResource) validateResourceFileBlobPrefixWithoutAutoStorageContainerName(data acceptance.TestData) string {
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

func (BatchPoolResource) validateResourceFileHttpURLWithoutFilePath(data acceptance.TestData) string {
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

func (BatchPoolResource) certificates(data acceptance.TestData) string {
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

func (BatchPoolResource) containerConfiguration(data acceptance.TestData) string {
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

func (BatchPoolResource) customImageConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
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

func (BatchPoolResource) networkConfiguration(data acceptance.TestData) string {
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
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
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
