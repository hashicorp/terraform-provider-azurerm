package hpccache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type HPCCacheBlobNFSTargetResource struct {
}

func TestAccHPCCacheBlobNFSTarget_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_nfs_target", "test")
	r := HPCCacheBlobNFSTargetResource{}

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

func TestAccHPCCacheBlobNFSTarget_accessPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_nfs_target", "test")
	r := HPCCacheBlobNFSTargetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.accessPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.accessPolicyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCacheBlobNFSTarget_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_nfs_target", "test")
	r := HPCCacheBlobNFSTargetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.namespace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCacheBlobNFSTarget_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_nfs_target", "test")
	r := HPCCacheBlobNFSTargetResource{}

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

func (HPCCacheBlobNFSTargetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.StorageTargetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HPCCache.StorageTargetsClient.Get(ctx, id.ResourceGroup, id.CacheName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving HPC Cache Blob Target (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r HPCCacheBlobNFSTargetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_blob_nfs_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = jsondecode(azurerm_resource_group_template_deployment.storage-containers.output_content).id.value
  namespace_path       = "/p1"
  usage_model          = "READ_HEAVY_INFREQ"
}
`, r.template(data), data.RandomString)
}

func (r HPCCacheBlobNFSTargetResource) namespace(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_blob_nfs_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = jsondecode(azurerm_resource_group_template_deployment.storage-containers.output_content).id.value
  namespace_path       = "/p2"
  usage_model          = "READ_HEAVY_INFREQ"
}
`, r.template(data), data.RandomString)
}

func (r HPCCacheBlobNFSTargetResource) accessPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_access_policy" "test" {
  name         = "p1"
  hpc_cache_id = azurerm_hpc_cache.test.id
  access_rule {
    scope  = "default"
    access = "rw"
  }

  # This is not needed in Terraform v0.13, whilst needed in v0.14.
  # Once https://github.com/hashicorp/terraform/issues/28193 is fixed, we can remove this lifecycle block.
  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_hpc_cache_blob_nfs_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = jsondecode(azurerm_resource_group_template_deployment.storage-containers.output_content).id.value
  namespace_path       = "/p1"
  access_policy_name   = azurerm_hpc_cache_access_policy.test.name
  usage_model          = "READ_HEAVY_INFREQ"
}
`, r.template(data), data.RandomString)
}

func (r HPCCacheBlobNFSTargetResource) accessPolicyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_access_policy" "test" {
  name         = "p2"
  hpc_cache_id = azurerm_hpc_cache.test.id
  access_rule {
    scope  = "default"
    access = "rw"
  }
  # This is necessary to make the Terraform apply order works correctly.
  # Without CBD: azurerm_hpc_cache_access_policy-p1 (destroy) -> azurerm_hpc_cache_blob_nfs_target (update) -> azurerm_hpc_cache_access_policy-p2 (create)
  # 			 (the 1st step wil fail as the access policy is under used by the blob target)
  # With CBD   : azurerm_hpc_cache_access_policy-p2 (create) -> azurerm_hpc_cache_blob_nfs_target (update) -> azurerm_hpc_cache_access_policy-p1 (delete)
  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_hpc_cache_blob_nfs_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = jsondecode(azurerm_resource_group_template_deployment.storage-containers.output_content).id.value
  namespace_path       = "/blob_storage1"
  access_policy_name   = azurerm_hpc_cache_access_policy.test.name
  usage_model          = "READ_HEAVY_INFREQ"
}
`, r.template(data), data.RandomString)
}

func (r HPCCacheBlobNFSTargetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_blob_nfs_target" "import" {
  name                 = azurerm_hpc_cache_blob_nfs_target.test.name
  resource_group_name  = azurerm_hpc_cache_blob_nfs_target.test.resource_group_name
  cache_name           = azurerm_hpc_cache_blob_nfs_target.test.cache_name
  storage_container_id = azurerm_hpc_cache_blob_nfs_target.test.storage_container_id
  namespace_path       = azurerm_hpc_cache_blob_nfs_target.test.namespace_path
  usage_model          = azurerm_hpc_cache_blob_nfs_target.test.usage_model
}
`, r.basic(data))
}

func (HPCCacheBlobNFSTargetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VN-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

data "azuread_service_principal" "test" {
  display_name = "HPC Cache Resource Provider"
}

resource "azurerm_storage_account" "test" {
  name                      = "accteststorgacc%[3]s"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  account_tier              = "Standard"
  account_kind              = "StorageV2"
  account_replication_type  = "LRS"
  is_hns_enabled            = true
  nfsv3_enabled             = true
  enable_https_traffic_only = false
  network_rules {
    default_action             = "Deny"
    virtual_network_subnet_ids = [azurerm_subnet.test.id]
  }
}

# Due to https://github.com/terraform-providers/terraform-provider-azurerm/issues/2977 and the fact
# that the NFSv3 enabled storage account can't allow public network access - otherwise the NFSv3 protocol will fail,
# we have to use the ARM template to deploy the storage container as a workaround.
# Once the issue above got resolved, we can instead use the azurerm_storage_container resource.
resource "azurerm_resource_group_template_deployment" "storage-containers" {
  name                = "acctest-strgctn-deployment-%[1]d"
  resource_group_name = azurerm_storage_account.test.resource_group_name
  deployment_mode     = "Incremental"

  parameters_content = jsonencode({
    name = {
      value = "acctest-strgctn-hpc-%[1]d"
    }
  })

  template_content = <<EOF
{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "name": {
      "type": "String"
    }
  },
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts/blobServices/containers",
      "apiVersion": "2019-06-01",
      "name": "[concat('${azurerm_storage_account.test.name}/', 'default/', parameters('name'))]",
      "location": "${azurerm_storage_account.test.location}",
      "properties": {}
    }
  ],

  "outputs": {
    "id": {
      "type": "String",
      "value": "[resourceId('Microsoft.Storage/storageAccounts/blobServices/containers', '${azurerm_storage_account.test.name}', 'default', parameters('name'))]"
    }
  }
}
EOF
}

resource "azurerm_role_assignment" "test_storage_account_contrib" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Account Contributor"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
