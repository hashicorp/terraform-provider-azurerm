---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_workspace_network_outbound_rule_service_tag"
description: |-
  Manages an Azure Machine Learning Workspace Network Outbound Rule Service Tag.
---
# azurerm_machine_learning_workspace_network_outbound_rule_service_tag

Manages an Azure Machine Learning Workspace Network Outbound Rule Service Tag.


## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "workspace-example-ai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_key_vault" "example" {
  name                = "workspaceexamplekeyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"
}

resource "azurerm_storage_account" "example" {
  name                     = "workspacestorageaccount"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_machine_learning_workspace" "example" {
  name                    = "example-workspace"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  application_insights_id = azurerm_application_insights.example.id
  key_vault_id            = azurerm_key_vault.example.id
  storage_account_id      = azurerm_storage_account.example.id

  managed_network {
    isolation_mode = "AllowOnlyApprovedOutbound"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_machine_learning_workspace_network_outbound_rule_service_tag" "example" {
  name         = "example-outboundrule"
  workspace_id = azurerm_machine_learning_workspace.example.id
  service_tag  = "AppService"
  protocol     = "TCP"
  port_ranges  = "443"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Machine Learning Workspace Network Outbound Rule Service Tag. Changing this forces a new resource to be created.

* `workspace_id` - (Required) Specifies the ID of the Machine Learning Workspace. Changing this forces a new resource to be created.

* `service_tag` - (Required) Specifies the fully qualified domain name to allow for outbound traffic. Possible values are `AppConfiguration`,`AppService`,`AzureActiveDirectory`,`AzureAdvancedThreatProtection`,`AzureArcInfrastructure`,`AzureAttestation`,`AzureBackup`,`AzureBotService`,`AzureContainerRegistry`,`AzureCosmosDB`,`AzureDataLake`,`AzureDevSpaces`,`AzureInformationProtection`,`AzureIoTHub`,`AzureKeyVault`,`AzureManagedGrafana`,`AzureMonitor`,`AzureOpenDatasets`,`AzurePlatformDNS`,`AzurePlatformIMDS`,`AzurePlatformLKM`,`AzureResourceManager`,`AzureSignalR`,`AzureSiteRecovery`,`AzureSpringCloud`,`AzureStack`,`AzureUpdateDelivery`,`DataFactoryManagement`,`EventHub`,`GuestAndHybridManagement`,`M365ManagementActivityApi`,`M365ManagementActivityApi`,`MicrosoftAzureFluidRelay`,`MicrosoftCloudAppSecurity`,`MicrosoftContainerRegistry`,`PowerPlatformInfra`,`ServiceBus`,`Sql`,`Storage`,`WindowsAdminCenter`,`AppServiceManagement`,`AutonomousDevelopmentPlatform`,`AzureActiveDirectoryDomainServices`,`AzureCloud`,`AzureConnectors`,`AzureContainerAppsService`,`AzureDatabricks`,`AzureDeviceUpdate`,`AzureEventGrid`,`AzureFrontDoor.Frontend`,`AzureFrontDoor.Backend`,`AzureFrontDoor.FirstParty`,`AzureHealthcareAPIs`,`AzureLoadBalancer`,`AzureMachineLearning`,`AzureSphere`,`AzureWebPubSub`,`BatchNodeManagement`,`ChaosStudio`,`CognitiveServicesFrontend`,`CognitiveServicesManagement`,`DataFactory`,`Dynamics365ForMarketingEmail`,`Dynamics365BusinessCentral`,`EOPExternalPublishedIPs`,`Internet`,`LogicApps`,`Marketplace`,`MicrosoftDefenderForEndpoint`,`PowerBI`,`PowerQueryOnline`,`ServiceFabric`,`SqlManagement`,`StorageSyncService`,`WindowsVirtualDesktop` and `VirtualNetwork`. 

* `protocol` - (Required) Specifies the network protocol. Possible values are `*`, `TCP`, `UDP` and `ICMP`

* `port_ranges` - (Required) Specifies which ports traffic will be allowed by this rule. You can specify a single port (e.g. ` 80`) ,  a port range  (e.g. `1024-655535`) or a comma-separated list of single ports and/or port ranges(e.g. `80,1024-655535`). `*` can be used to allow traffic on any port.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Machine Learning Workspace Network Outbound Rule Service Tag .

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Workspace Network Outbound Rule Service Tag.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Workspace Network Outbound Rule Service Tag.
* `update` - (Defaults to 30 minutes) Used when updating the Machine Learning Workspace Network Outbound Rule Service Tag.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Workspace Network Outbound Rule Service Tag.

## Import

Machine Learning Workspace Network Outbound Rule Service Tag can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_workspace_network_outbound_rule_service_tag.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/outboundRules/rule1
```
