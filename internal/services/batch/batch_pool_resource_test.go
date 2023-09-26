// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2023-05-01/pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BatchPoolResource struct{}

func TestAccBatchPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("18.04-lts"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_acceleratedNetworkingEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.acceleratedNetworkingEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_configuration.0.accelerated_networking_enabled").HasValue("true"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_identityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("18.04-lts"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("18.04-lts"),
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("18.04-lts"),
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fixedScale_complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("max_tasks_per_node").HasValue("2"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("18.04-lts"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.node_deallocation_method").HasValue("Terminate"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("2"),
				check.That(data.ResourceName).Key("fixed_scale.0.resize_timeout").HasValue("PT15M"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_low_priority_nodes").HasValue("0"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		data.ImportStep("stop_pending_resize_operation", "fixed_scale.0.node_deallocation_method"),
	})
}

func TestAccBatchPool_autoScale_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoScale_complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("18.04-lts"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("auto_scale.0.evaluation_interval").HasValue("PT15M"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		data.ImportStep("stop_pending_resize_operation", "fixed_scale.0.node_deallocation_method"),
	})
}

func TestAccBatchPool_completeUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fixedScale_complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("18.04-lts"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("2"),
				check.That(data.ResourceName).Key("fixed_scale.0.resize_timeout").HasValue("PT15M"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_low_priority_nodes").HasValue("0"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
			),
		},
		data.ImportStep("stop_pending_resize_operation", "fixed_scale.0.node_deallocation_method"),
		{
			Config: r.autoScale_complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("18.04-lts"),
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.startTask_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("18.04-lts"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.resize_timeout").HasValue("PT15M"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_low_priority_nodes").HasValue("0"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.0.task_retry_maximum").HasValue("5"),
				check.That(data.ResourceName).Key("start_task.0.common_environment_properties.%").HasValue("2"),
				check.That(data.ResourceName).Key("start_task.0.common_environment_properties.env").HasValue("TEST"),
				check.That(data.ResourceName).Key("start_task.0.common_environment_properties.bu").HasValue("Research&Dev"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.0.auto_user.#").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.0.auto_user.0.scope").HasValue("Task"),
				check.That(data.ResourceName).Key("start_task.0.user_identity.0.auto_user.0.elevation_level").HasValue("NonAdmin"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_startTask_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.startTask_complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_task.0.container.0.registry.0.user_name").HasValue("myUserName"),
				check.That(data.ResourceName).Key("start_task.0.container.0.registry.0.registry_server").HasValue("myContainerRegistry.azurecr.io"),
				check.That(data.ResourceName).Key("start_task.0.container.0.registry.0.user_name").HasValue("myUserName"),
				check.That(data.ResourceName).Key("start_task.0.container.0.run_options").HasValue("cat /proc/cpuinfo"),
				check.That(data.ResourceName).Key("start_task.0.container.0.image_name").HasValue("centos7"),
				check.That(data.ResourceName).Key("start_task.0.container.0.working_directory").HasValue("ContainerImageDefault"),
			),
		},
		data.ImportStep("stop_pending_resize_operation",
			"container_configuration.0.container_registries.0.password",
			"start_task.0.container.0.registry.0.password"),
	})
}

func TestAccBatchPool_startTask_userIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.startTask_userIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.certificates(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("2"),
				check.That(data.ResourceName).Key("certificate.0.id").HasValue(certificate0ID),
				check.That(data.ResourceName).Key("certificate.0.store_location").HasValue("CurrentUser"),
				check.That(data.ResourceName).Key("certificate.0.store_name").HasValue(""),
				check.That(data.ResourceName).Key("certificate.0.visibility.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.1.id").HasValue(certificate1ID),
				check.That(data.ResourceName).Key("certificate.1.store_location").HasValue("CurrentUser"),
				check.That(data.ResourceName).Key("certificate.1.store_name").HasValue(""),
				check.That(data.ResourceName).Key("certificate.1.visibility.#").HasValue("2"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_validateResourceFileWithoutSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.validateResourceFileWithoutSource(data),
			ExpectError: regexp.MustCompile("exactly one of auto_storage_container_name, storage_container_url and http_url must be specified"),
		},
	})
}

func TestAccBatchPool_containerWithUser(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.containerConfigurationWithRegistryUser(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container_configuration.0.type").HasValue("DockerCompatible"),
				check.That(data.ResourceName).Key("container_configuration.0.container_image_names.#").HasValue("1"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.#").HasValue("1"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.registry_server").HasValue("myContainerRegistry.azurecr.io"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.user_name").HasValue("myUserName"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.password").HasValue("myPassword"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.user_assigned_identity_id").IsEmpty(),
			),
		},
		data.ImportStep(
			"stop_pending_resize_operation",
			"container_configuration.0.container_registries.0.password",
		),
	})
}

func TestAccBatchPool_containerWithUAMI(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.containerConfigurationWithRegistryUAMI(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container_configuration.0.type").HasValue("DockerCompatible"),
				check.That(data.ResourceName).Key("container_configuration.0.container_image_names.#").HasValue("1"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.#").HasValue("1"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.registry_server").HasValue("myContainerRegistry.azurecr.io"),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.user_name").IsEmpty(),
				check.That(data.ResourceName).Key("container_configuration.0.container_registries.0.user_assigned_identity_id").IsSet(),
			),
		},
		data.ImportStep(
			"stop_pending_resize_operation",
		),
	})
}

func TestAccBatchPool_validateResourceFileWithMultipleSources(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.validateResourceFileWithMultipleSources(data),
			ExpectError: regexp.MustCompile("exactly one of auto_storage_container_name, storage_container_url and http_url must be specified"),
		},
	})
}

func TestAccBatchPool_validateResourceFileBlobPrefixWithoutAutoStorageContainerUrl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.validateResourceFileBlobPrefixWithoutAutoStorageContainerName(data),
			ExpectError: regexp.MustCompile("auto_storage_container_name or storage_container_url must be specified when using blob_prefix"),
		},
	})
}

func TestAccBatchPool_validateResourceFileHttpURLWithoutFilePath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.validateResourceFileHttpURLWithoutFilePath(data),
			ExpectError: regexp.MustCompile("file_path must be specified when using http_url"),
		},
	})
}

func TestAccBatchPool_validateResourceFileWithIdentityReference(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.validateResourceFileWithIdentityReference(data),
			Check:  acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r)),
		},
	})
}

func TestAccBatchPool_customImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customImageConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("max_tasks_per_node").HasValue("2"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("storage_image_reference.#").HasValue("1"),
				check.That(data.ResourceName).Key("storage_image_reference.0.publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("storage_image_reference.0.sku").HasValue("18.04-lts"),
				check.That(data.ResourceName).Key("storage_image_reference.0.offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("auto_scale.#").HasValue("0"),
				check.That(data.ResourceName).Key("fixed_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("fixed_scale.0.target_dedicated_nodes").HasValue("1"),
				check.That(data.ResourceName).Key("start_task.#").HasValue("0"),
				check.That(data.ResourceName).Key("network_configuration.#").HasValue("1"),
				check.That(data.ResourceName).Key("network_configuration.0.dynamic_vnet_assignment_scope").HasValue("none"),
				check.That(data.ResourceName).Key("network_configuration.0.subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_configuration.0.public_ips.#").HasValue("1"),
				check.That(data.ResourceName).Key("network_configuration.0.endpoint_configuration.0.network_security_group_rules.0.source_port_ranges.0").HasValue("*"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_fixedScaleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fixedScale_complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fixed_scale.0.node_deallocation_method").HasValue("Terminate"),
			),
		},
		data.ImportStep("stop_pending_resize_operation", "fixed_scale.0.node_deallocation_method"),
		{
			Config: r.fixedScale_completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("stop_pending_resize_operation", "fixed_scale.0.node_deallocation_method"),
	})
}

func TestAccBatchPool_mountConfigurationAzureBlobFileSystem(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mountConfigurationAzureBlobFileSystem(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount.0.azure_blob_file_system.#").HasValue("1"),
				check.That(data.ResourceName).Key("mount.0.azure_blob_file_system.0.relative_mount_path").HasValue("/mnt/"),
			),
		},
		data.ImportStep(
			"stop_pending_resize_operation",
			"mount.0.azure_blob_file_system.0.account_key",
		),
	})
}

func TestAccBatchPool_mountConfigurationAzureBlobFileSystemWithUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mountConfigurationAzureBlobFileSystemWithUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount.0.azure_blob_file_system.#").HasValue("1"),
				check.That(data.ResourceName).Key("mount.0.azure_blob_file_system.0.relative_mount_path").HasValue("/mnt/"),
				check.That(data.ResourceName).Key("mount.0.azure_blob_file_system.0.identity_id").Exists(),
			),
		},
		data.ImportStep(
			"stop_pending_resize_operation",
			"mount.0.azure_blob_file_system.0.account_key",
		),
	})
}

func TestAccBatchPool_mountConfigurationAzureFileShare(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mountConfigurationAzureFileShare(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount.0.azure_file_share.#").HasValue("1"),
				check.That(data.ResourceName).Key("mount.0.azure_file_share.0.relative_mount_path").HasValue("/mnt/"),
			),
		},
		data.ImportStep(
			"stop_pending_resize_operation",
			"mount.0.azure_file_share.0.account_key",
		),
	})
}

func TestAccBatchPool_mountConfigurationCIFS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mountConfigurationCIFS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount.0.cifs_mount.#").HasValue("1"),
				check.That(data.ResourceName).Key("mount.0.cifs_mount.0.user_name").HasValue("myUserName"),
				check.That(data.ResourceName).Key("mount.0.cifs_mount.0.password").HasValue("myPassword"),
				check.That(data.ResourceName).Key("mount.0.cifs_mount.0.source").HasValue("https://testaccount.file.core.windows.net/"),
				check.That(data.ResourceName).Key("mount.0.cifs_mount.0.relative_mount_path").HasValue("/mnt/"),
				check.That(data.ResourceName).Key("mount.0.cifs_mount.0.mount_options").HasValue("sampleops"),
			),
		},
		data.ImportStep(
			"stop_pending_resize_operation",
			"mount.0.cifs_mount.0.password",
		),
	})
}

func TestAccBatchPool_mountConfigurationNFS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mountConfigurationNFS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount.0.nfs_mount.#").HasValue("1"),
				check.That(data.ResourceName).Key("mount.0.nfs_mount.0.source").HasValue("https://testaccount.file.core.windows.net/"),
				check.That(data.ResourceName).Key("mount.0.nfs_mount.0.relative_mount_path").HasValue("/mnt/"),
				check.That(data.ResourceName).Key("mount.0.nfs_mount.0.mount_options").HasValue("sampleops"),
			),
		},
		data.ImportStep(
			"stop_pending_resize_operation",
		),
	})
}

func TestAccBatchPool_targetNodeCommunicationMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.targetNodeCommunicationMode(data, "Default"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
		{
			Config: r.targetNodeCommunicationMode(data, "Simplified"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_diskSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.diskSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vm_size").HasValue("STANDARD_A1"),
				check.That(data.ResourceName).Key("node_agent_sku_id").HasValue("batch.node.ubuntu 18.04"),
				check.That(data.ResourceName).Key("data_disks.0.lun").HasValue("20"),
				check.That(data.ResourceName).Key("data_disks.0.caching").HasValue("None"),
				check.That(data.ResourceName).Key("data_disks.0.disk_size_gb").HasValue("1"),
				check.That(data.ResourceName).Key("data_disks.0.storage_account_type").HasValue("Standard_LRS"),
				check.That(data.ResourceName).Key("os_disk_placement").HasValue("CacheDisk"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_extensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("extensions.0.name").HasValue("KeyVaultForLinux"),
				check.That(data.ResourceName).Key("extensions.0.publisher").HasValue("Microsoft.Azure.KeyVault"),
				check.That(data.ResourceName).Key("extensions.0.type").HasValue("KeyVaultForLinux"),
				check.That(data.ResourceName).Key("extensions.0.type_handler_version").HasValue("2.0"),
				check.That(data.ResourceName).Key("extensions.0.auto_upgrade_minor_version").HasValue("true"),
				check.That(data.ResourceName).Key("extensions.0.automatic_upgrade_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("extensions.0.settings_json").HasValue("{}"),
				check.That(data.ResourceName).Key("extensions.0.protected_settings").HasValue("sensitive"),
				check.That(data.ResourceName).Key("extensions.0.provision_after_extensions.0").HasValue("newProv1"),
			),
		},
	})
}

func TestAccBatchPool_interNodeCommunicationWithTaskSchedulingPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.interNodeCommunicationWithTaskSchedulingPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("inter_node_communication").HasValue("Disabled"),
				check.That(data.ResourceName).Key("task_scheduling_policy.0.node_fill_type").HasValue("Pack"),
			),
		},
		data.ImportStep("stop_pending_resize_operation"),
	})
}

func TestAccBatchPool_linuxUserAccounts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxUserAccounts(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("user_accounts.0.name").HasValue("username1"),
				check.That(data.ResourceName).Key("user_accounts.0.password").HasValue("<ExamplePassword>"),
				check.That(data.ResourceName).Key("user_accounts.0.elevation_level").HasValue("Admin"),
				check.That(data.ResourceName).Key("user_accounts.0.linux_user_configuration.0.ssh_private_key").HasValue("sshprivatekeyvalue"),
				check.That(data.ResourceName).Key("user_accounts.0.linux_user_configuration.0.uid").HasValue("1234"),
				check.That(data.ResourceName).Key("user_accounts.0.linux_user_configuration.0.gid").HasValue("4567"),
			),
		},
		data.ImportStep(
			"stop_pending_resize_operation",
			"user_accounts.0.password",
			"user_accounts.0.linux_user_configuration.0.ssh_private_key",
		),
	})
}

func TestAccBatchPool_windowsUserAccountsWithAdditionalConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_pool", "test")
	r := BatchPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsUserAccountsWithConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("license_type").HasValue("Windows_Server"),
				check.That(data.ResourceName).Key("node_placement.0.policy").HasValue("Regional"),
				check.That(data.ResourceName).Key("disk_encryption.0.disk_encryption_target").HasValue("TemporaryDisk"),
				check.That(data.ResourceName).Key("user_accounts.0.name").HasValue("username1"),
				check.That(data.ResourceName).Key("user_accounts.0.password").HasValue("<ExamplePassword>"),
				check.That(data.ResourceName).Key("user_accounts.0.elevation_level").HasValue("Admin"),
				check.That(data.ResourceName).Key("user_accounts.0.windows_user_configuration.0.login_mode").HasValue("Interactive"),
				check.That(data.ResourceName).Key("windows.0.enable_automatic_updates").HasValue("true"),
			),
		},
		data.ImportStep(
			"stop_pending_resize_operation",
			"user_accounts.0.password",
		),
	})
}

func (t BatchPoolResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := pool.ParsePoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Batch.PoolClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s", *id)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (BatchPoolResource) fixedScale_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-batch-%d"
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
  name                                = "testaccbatch%s"
  resource_group_name                 = azurerm_resource_group.test.name
  location                            = azurerm_resource_group.test.location
  pool_allocation_mode                = "BatchService"
  storage_account_id                  = azurerm_storage_account.test.id
  storage_account_authentication_mode = "StorageKeys"

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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"

  fixed_scale {
    node_deallocation_method  = "Terminate"
    target_dedicated_nodes    = 2
    resize_timeout            = "PT15M"
    target_low_priority_nodes = 0
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
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
  name     = "acctestRG-batch-%d"
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
  name                                = "testaccbatch%s"
  resource_group_name                 = azurerm_resource_group.test.name
  location                            = azurerm_resource_group.test.location
  pool_allocation_mode                = "BatchService"
  storage_account_id                  = azurerm_storage_account.test.id
  storage_account_authentication_mode = "StorageKeys"

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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"

  fixed_scale {
    target_dedicated_nodes = 3
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
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
  name     = "acctestRG-batch-%d"
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
  name                                = "testaccbatch%s"
  resource_group_name                 = azurerm_resource_group.test.name
  location                            = azurerm_resource_group.test.location
  pool_allocation_mode                = "BatchService"
  storage_account_id                  = azurerm_storage_account.test.id
  storage_account_authentication_mode = "StorageKeys"

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
  node_agent_sku_id             = "batch.node.ubuntu 18.04"
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
    sku       = "18.04-lts"
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
  name     = "acctestRG-batch-%d"
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (BatchPoolResource) acceleratedNetworkingEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-batch-%d"
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
  node_agent_sku_id   = "batch.node.windows amd64"
  vm_size             = "Standard_D1_v2"

  fixed_scale {
    target_dedicated_nodes = 2
  }

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-datacenter-smalldisk"
    version   = "latest"
  }

  network_configuration {
    accelerated_networking_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (BatchPoolResource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-batch-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
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
    sku       = "18.04-lts"
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
  name     = "acctestRG-batch-%d"
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }

  start_task {
    command_line       = "echo 'Hello World from $env'"
    wait_for_success   = true
    task_retry_maximum = 5
    common_environment_properties = {
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
      http_url  = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/main/README.md"
      file_path = "README.md"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (BatchPoolResource) startTask_complete(data acceptance.TestData) string {
	template := BatchPoolResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 20.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "microsoft-azure-batch"
    offer     = "ubuntu-server-container"
    sku       = "20-04-lts"
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

  start_task {
    command_line       = "echo 'Hello World from $env'"
    wait_for_success   = true
    task_retry_maximum = 5
    common_environment_properties = {
      env = "TEST"
      bu  = "Research&Dev"
    }

    container {
      run_options = "cat /proc/cpuinfo"
      image_name  = "centos7"
      registry {
        registry_server = "myContainerRegistry.azurecr.io"
        user_name       = "myUserName"
        password        = "myPassword"
      }
      working_directory = "ContainerImageDefault"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }

    resource_file {
      storage_container_url = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/main/README.md"
      file_path             = "README.md"
    }
  }
}
`, template, data.RandomString, data.RandomString)
}

func (BatchPoolResource) startTask_userIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-batch-%[1]d"
  location = "%[2]s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }

  start_task {
    command_line       = "echo 'Hello World from $env'"
    wait_for_success   = true
    task_retry_maximum = 5

    common_environment_properties = {
      env = "TEST"
      bu  = "Research&Dev"
    }

    user_identity {
      user_name = "adminuser"
    }

    resource_file {
      storage_container_url = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/main/README.md"
      file_path             = "README.md"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (BatchPoolResource) validateResourceFileWithoutSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestbatch%d"
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }

  start_task {
    command_line       = "echo 'Hello World from $env'"
    task_retry_maximum = 1
    wait_for_success   = true

    common_environment_properties = {
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
  name     = "acctestbatch%d"
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }

  start_task {
    command_line       = "echo 'Hello World from $env'"
    task_retry_maximum = 1
    wait_for_success   = true

    common_environment_properties = {
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
  name     = "acctestbatch%d"
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }

  start_task {
    command_line       = "echo 'Hello World from $env'"
    task_retry_maximum = 1
    wait_for_success   = true

    common_environment_properties = {
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
  name     = "acctestbatch%d"
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }

  start_task {
    command_line       = "echo 'Hello World from $env'"
    task_retry_maximum = 1
    wait_for_success   = true

    common_environment_properties = {
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

func (BatchPoolResource) validateResourceFileWithIdentityReference(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "useridentity%s"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }

  start_task {
    command_line       = "echo 'Hello World from $env'"
    task_retry_maximum = 1
    wait_for_success   = true

    common_environment_properties = {
      env = "TEST"
      bu  = "Research&Dev"
    }

    user_identity {
      user_name = "testUserIndentity"
    }

    resource_file {
      http_url                  = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/main/README.md"
      file_path                 = "README.md"
      user_assigned_identity_id = azurerm_user_assigned_identity.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}

func (BatchPoolResource) certificates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestbatch%d"
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
  certificate          = filebase64("testdata/batch_certificate_password.pfx")
  format               = "Pfx"
  password             = "terraform"
  thumbprint           = "42c107874fd0e4a9583292a2f1098e8fe4b2edda"
  thumbprint_algorithm = "SHA1"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
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

func (BatchPoolResource) containerConfigurationWithRegistryUser(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestbatch%d"
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
  node_agent_sku_id   = "batch.node.ubuntu 20.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "microsoft-azure-batch"
    offer     = "ubuntu-server-container"
    sku       = "20-04-lts"
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

func (BatchPoolResource) containerConfigurationWithRegistryUAMI(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestbatch%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "testaccuami%d"
}

resource "azurerm_container_registry" "test" {
  name                = "testregistry%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
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
  node_agent_sku_id   = "batch.node.ubuntu 20.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "microsoft-azure-batch"
    offer     = "ubuntu-server-container"
    sku       = "20-04-lts"
    version   = "latest"
  }

  container_configuration {
    type                  = "DockerCompatible"
    container_image_names = ["centos7"]
    container_registries {
      registry_server           = "myContainerRegistry.azurecr.io"
      user_assigned_identity_id = azurerm_user_assigned_identity.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomString, data.RandomString)
}

func (BatchPoolResource) customImageConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-batch-%d"
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
  address_prefixes     = ["10.0.2.0/24"]
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

  allow_nested_items_to_be_public = true

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
    sku       = "18.04-lts"
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

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"

  fixed_scale {
    target_dedicated_nodes = 2
  }

  storage_image_reference {
    id = azurerm_shared_image.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString)
}

func (BatchPoolResource) networkConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d-batchpool"
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
  address_prefixes     = ["10.0.2.0/24"]
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }

  network_configuration {
    dynamic_vnet_assignment_scope    = "none"
    public_address_provisioning_type = "UserManaged"
    public_ips                       = [azurerm_public_ip.test.id]
    subnet_id                        = azurerm_subnet.test.id

    endpoint_configuration {
      name                = "SSH"
      protocol            = "TCP"
      backend_port        = 22
      frontend_port_range = "4000-4100"

      network_security_group_rules {
        access                = "Deny"
        priority              = 1001
        source_address_prefix = "*"
        source_port_ranges    = ["*"]
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (BatchPoolResource) mountConfigurationAzureBlobFileSystem(data acceptance.TestData) string {
	template := BatchPoolResource{}.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_storage_account" "test" {
  name                     = "accbatchsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
resource "azurerm_storage_container" "test" {
  name                  = "accbatchsc%s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"
  mount {
    azure_blob_file_system {
      account_name        = azurerm_storage_account.test.name
      container_name      = azurerm_storage_container.test.name
      account_key         = azurerm_storage_account.test.primary_access_key
      relative_mount_path = "/mnt/"
    }
  }
  fixed_scale {
    target_dedicated_nodes = 1
  }
  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }
}
`, template, data.RandomString, data.RandomString, data.RandomString, data.RandomString)
}

func (BatchPoolResource) mountConfigurationAzureBlobFileSystemWithUserAssignedIdentity(data acceptance.TestData) string {
	template := BatchPoolResource{}.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_storage_account" "test" {
  name                     = "accbatchsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "accbatchsc%s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "testidentity%s"
}

resource "azurerm_role_assignment" "blob_contributor" {
  principal_id         = azurerm_user_assigned_identity.test.principal_id
  role_definition_name = "Storage Blob Data Contributor"
  scope                = azurerm_storage_account.test.id
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  display_name        = "Test Acc Pool Auto"
  vm_size             = "Standard_A1"
  node_agent_sku_id   = "batch.node.ubuntu 20.04"

  fixed_scale {
    target_dedicated_nodes = 0
  }

  storage_image_reference {
    publisher = "microsoft-azure-batch"
    offer     = "ubuntu-server-container"
    sku       = "20-04-lts"
    version   = "latest"
  }

  mount {
    azure_blob_file_system {
      account_name        = azurerm_storage_account.test.name
      container_name      = azurerm_storage_container.test.name
      relative_mount_path = "/mnt/"
      identity_id         = azurerm_user_assigned_identity.test.id
    }
  }
}
`, template, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString)
}

func (BatchPoolResource) mountConfigurationAzureFileShare(data acceptance.TestData) string {
	template := BatchPoolResource{}.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_storage_account" "test" {
  name                     = "accbatchsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
resource "azurerm_storage_container" "test" {
  name                  = "accbatchsc%s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"
  mount {
    azure_file_share {
      account_name        = azurerm_storage_account.test.name
      account_key         = azurerm_storage_account.test.primary_access_key
      azure_file_url      = "https://testaccount.file.core.windows.net/"
      relative_mount_path = "/mnt/"
    }
  }
  fixed_scale {
    target_dedicated_nodes = 1
  }
  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }
}
`, template, data.RandomString, data.RandomString, data.RandomString, data.RandomString)
}

func (BatchPoolResource) mountConfigurationCIFS(data acceptance.TestData) string {
	template := BatchPoolResource{}.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"
  mount {
    cifs_mount {
      user_name           = "myUserName"
      password            = "myPassword"
      source              = "https://testaccount.file.core.windows.net/"
      relative_mount_path = "/mnt/"
      mount_options       = "sampleops"
    }
  }
  fixed_scale {
    target_dedicated_nodes = 1
  }
  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }
}
`, template, data.RandomString, data.RandomString)
}

func (BatchPoolResource) mountConfigurationNFS(data acceptance.TestData) string {
	template := BatchPoolResource{}.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"
  mount {
    nfs_mount {
      source              = "https://testaccount.file.core.windows.net/"
      relative_mount_path = "/mnt/"
      mount_options       = "sampleops"
    }
  }
  fixed_scale {
    target_dedicated_nodes = 1
  }
  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }
}
`, template, data.RandomString, data.RandomString)
}

func (BatchPoolResource) targetNodeCommunicationMode(data acceptance.TestData, targetNodeCommunicationMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-batch-%d"
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
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"

  target_node_communication_mode = "%s"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, targetNodeCommunicationMode)
}

func (BatchPoolResource) extensions(data acceptance.TestData) string {
	template := BatchPoolResource{}.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"
  extensions {
    name                       = "KeyVaultForLinux"
    publisher                  = "Microsoft.Azure.KeyVault"
    type                       = "KeyVaultForLinux"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true
    automatic_upgrade_enabled  = true
    settings_json              = "{}"
    protected_settings         = "sensitive"
    provision_after_extensions = ["newProv1"]
  }
  fixed_scale {
    target_dedicated_nodes = 1
  }
  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }
}
`, template, data.RandomString, data.RandomString)
}

func (BatchPoolResource) diskSettings(data acceptance.TestData) string {
	template := BatchPoolResource{}.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"
  data_disks {
    lun                  = 20
    caching              = "None"
    disk_size_gb         = 1
    storage_account_type = "Standard_LRS"
  }
  os_disk_placement = "CacheDisk"
  fixed_scale {
    target_dedicated_nodes = 1
  }
  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }
}
`, template, data.RandomString, data.RandomString)
}

func (BatchPoolResource) interNodeCommunicationWithTaskSchedulingPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-batch-%d"
  location = "%s"
}
resource "azurerm_batch_account" "test" {
  name                          = "testaccbatch%s"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  public_network_access_enabled = false
}
resource "azurerm_batch_pool" "test" {
  name                     = "testaccpool%s"
  resource_group_name      = azurerm_resource_group.test.name
  account_name             = azurerm_batch_account.test.name
  node_agent_sku_id        = "batch.node.ubuntu 18.04"
  vm_size                  = "Standard_A1"
  inter_node_communication = "Disabled"
  task_scheduling_policy {
    node_fill_type = "Pack"
  }
  fixed_scale {
    target_dedicated_nodes = 1
  }
  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (BatchPoolResource) linuxUserAccounts(data acceptance.TestData) string {
	template := BatchPoolResource{}.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 18.04"
  vm_size             = "Standard_A1"
  fixed_scale {
    target_dedicated_nodes = 1
  }
  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-lts"
    version   = "latest"
  }
  user_accounts {
    name            = "username1"
    password        = "<ExamplePassword>"
    elevation_level = "Admin"
    linux_user_configuration {
      ssh_private_key = "sshprivatekeyvalue"
      uid             = 1234
      gid             = 4567
    }
  }
}
`, template, data.RandomString, data.RandomString)
}

func (BatchPoolResource) windowsUserAccountsWithConfig(data acceptance.TestData) string {
	template := BatchPoolResource{}.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_batch_pool" "test" {
  name                = "testaccpool%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.windows amd64"
  vm_size             = "Standard_A1"
  fixed_scale {
    target_dedicated_nodes = 1
  }
  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter"
    version   = "latest"
  }
  license_type = "Windows_Server"
  node_placement {
    policy = "Regional"
  }
  disk_encryption {
    disk_encryption_target = "TemporaryDisk"
  }
  windows {
    enable_automatic_updates = true
  }
  user_accounts {
    name            = "username1"
    password        = "<ExamplePassword>"
    elevation_level = "Admin"
    windows_user_configuration {
      login_mode = "Interactive"
    }
  }
}
`, template, data.RandomString, data.RandomString)
}

func (BatchPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-batch-%d"
  location = "%s"
}
resource "azurerm_network_security_group" "test" {
  name                = "testnsg-batch-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
resource "azurerm_virtual_network" "test" {
  name                = "testvn-batch-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}
resource "azurerm_subnet" "testsubnet" {
  name                 = "testsn-%s"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}
resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.testsubnet.id
  network_security_group_id = azurerm_network_security_group.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}
