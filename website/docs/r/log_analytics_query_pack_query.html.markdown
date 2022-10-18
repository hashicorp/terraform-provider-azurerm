---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_query_pack_query"
description: |-
  Manages a Log Analytics Query Pack Query.
---

# azurerm_log_analytics_query_pack_query

Manages a Log Analytics Query Pack Query.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_query_pack" "example" {
  name                = "example-laqp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_log_analytics_query_pack_query" "example" {
  name          = "19952bc3-0bf9-49eb-b713-6b80e7a41847"
  query_pack_id = azurerm_log_analytics_query_pack.example.id
  body          = "let newExceptionsTimeRange = 1d;\nlet timeRangeToCheckBefore = 7d;\nexceptions\n| where timestamp < ago(timeRangeToCheckBefore)\n| summarize count() by problemId\n| join kind= rightanti (\nexceptions\n| where timestamp >= ago(newExceptionsTimeRange)\n| extend stack = tostring(details[0].rawStack)\n| summarize count(), dcount(user_AuthenticatedId), min(timestamp), max(timestamp), any(stack) by problemId  \n) on problemId \n| order by  count_ desc\n"
  display_name  = "Exceptions - New in the last 24 hours"
}
```

## Arguments Reference

The following arguments are supported:

* `query_pack_id` - (Required) The ID of the Log Analytics Query Pack. Changing this forces a new resource to be created.

* `body` - (Required) The body of the Log Analytics Query Pack Query.

* `display_name` - (Required) The unique display name for the query within the Log Analytics Query Pack.

* `name` - (Optional) An unique UUID/GUID which identifies this Log Analytics Query Pack Query - one will be generated if not specified. Changing this forces a new resource to be created.

* `description` - (Optional) The description of the Log Analytics Query Pack Query.

* `categories` - (Optional) A list of the related categories for the function. Possible values are `applications`, `audit`, `container`, `databases`, `desktopanalytics`, `management`, `monitor`, `network`, `resources`, `security`, `virtualmachines`, `windowsvirtualdesktop` and `workloads`.

* `additional_settings_json` - (Optional) The additional properties that can be set for the Log Analytics Query Pack Query.

* `resource_types` - (Optional) A list of the related resource types for the function. Possible values are `default`, `microsoft.aad/domainservices`, `microsoft.aadiam/tenants`, `microsoft.agfoodplatform/farmbeats`, `microsoft.analysisservices/servers`, `microsoft.apimanagement/service`, `microsoft.appconfiguration/configurationstores`, `microsoft.appplatform/spring`, `microsoft.attestation/attestationproviders`, `microsoft.authorization/tenants`, `microsoft.automation/automationaccounts`, `microsoft.autonomousdevelopmentplatform/accounts`, `microsoft.azurestackhci/virtualmachines`, `microsoft.batch/batchaccounts`, `microsoft.blockchain/blockchainmembers`, `microsoft.botservice/botservices`, `microsoft.cache/redis`, `microsoft.cdn/profiles`, `microsoft.cognitiveservices/accounts`, `microsoft.communication/communicationservices`, `microsoft.compute/virtualmachines`, `microsoft.compute/virtualmachinescalesets`, `microsoft.connectedcache/cachenodes`, `microsoft.connectedvehicle/platformaccounts`, `microsoft.conenctedvmwarevsphere/virtualmachines`, `microsoft.containerregistry/registries`, `microsoft.containerservice/managedclusters`, `microsoft.d365customerinsights/instances`, `microsoft.dashboard/grafana`, `microsoft.databricks/workspaces`, `microsoft.datacollaboration/workspaces`, `microsoft.datafactory/factories`, `microsoft.datalakeanalytics/accounts`, `microsoft.datalakestore/accounts`, `microsoft.datashare/accounts`, `microsoft.dbformariadb/servers`, `microsoft.dbformysql/servers`, `microsoft.dbforpostgresql/flexibleservers`, `microsoft.dbforpostgresql/servers`, `microsoft.dbforpostgresql/serversv2`, `microsoft.digitaltwins/digitaltwinsinstances`, `microsoft.documentdb/cassandraclusters`, `microsoft.documentdb/databaseaccounts`, `microsoft.desktopvirtualization/applicationgroups`, `microsoft.desktopvirtualization/hostpools`, `microsoft.desktopvirtualization/workspaces`, `microsoft.devices/iothubs`, `microsoft.devices/provisioningservices`, `microsoft.dynamics/fraudprotection/purchase`, `microsoft.eventgrid/domains`, `microsoft.eventgrid/topics`, `microsoft.eventgrid/partnernamespaces`, `microsoft.eventgrid/partnertopics`, `microsoft.eventgrid/systemtopics`, `microsoft.eventhub/namespaces`, `microsoft.experimentation/experimentworkspaces`, `microsoft.hdinsight/clusters`, `microsoft.healthcareapis/services`, `microsoft.informationprotection/datasecuritymanagement`, `microsoft.intune/operations`, `microsoft.insights/autoscalesettings`, `microsoft.insights/components`, `microsoft.insights/workloadmonitoring`, `microsoft.keyvault/vaults`, `microsoft.kubernetes/connectedclusters`, `microsoft.kusto/clusters`, `microsoft.loadtestservice/loadtests`, `microsoft.logic/workflows`, `microsoft.machinelearningservices/workspaces`, `microsoft.media/mediaservices`, `microsoft.netapp/netappaccounts/capacitypools`, `microsoft.network/applicationgateways`, `microsoft.network/azurefirewalls`, `microsoft.network/bastionhosts`, `microsoft.network/expressroutecircuits`, `microsoft.network/frontdoors`, `microsoft.network/loadbalancers`, `microsoft.network/networkinterfaces`, `microsoft.network/networksecuritygroups`, `microsoft.network/networksecurityperimeters`, `microsoft.network/networkwatchers/connectionmonitors`, `microsoft.network/networkwatchers/trafficanalytics`, `microsoft.network/publicipaddresses`, `microsoft.network/trafficmanagerprofiles`, `microsoft.network/virtualnetworks`, `microsoft.network/virtualnetworkgateways`, `microsoft.network/vpngateways`, `microsoft.networkfunction/azuretrafficcollectors`, `microsoft.openenergyplatform/energyservices`, `microsoft.openlogisticsplatform/workspaces`, `microsoft.operationalinsights/workspaces`, `microsoft.powerbi/tenants`, `microsoft.powerbi/tenants/workspaces`, `microsoft.powerbidedicated/capacities`, `microsoft.purview/accounts`, `microsoft.recoveryservices/vaults`, `microsoft.resources/azureactivity`, `microsoft.scvmm/virtualmachines`, `microsoft.search/searchservices`, `microsoft.security/antimalwaresettings`, `microsoft.securityinsights/amazon`, `microsoft.securityinsights/anomalies`, `microsoft.securityinsights/cef`, `microsoft.securityinsights/datacollection`, `microsoft.securityinsights/dnsnormalized`, `microsoft.securityinsights/mda`, `microsoft.securityinsights/mde`, `microsoft.securityinsights/mdi`, `microsoft.securityinsights/mdo`, `microsoft.securityinsights/networksessionnormalized`, `microsoft.securityinsights/office365`, `microsoft.securityinsights/purview`, `microsoft.securityinsights/securityinsights`, `microsoft.securityinsights/securityinsights/mcas`, `microsoft.securityinsights/tvm`, `microsoft.securityinsights/watchlists`, `microsoft.servicebus/namespaces`, `microsoft.servicefabric/clusters`, `microsoft.signalrservice/signalr`, `microsoft.signalrservice/webpubsub`, `microsoft.sql/managedinstances`, `microsoft.sql/servers`, `microsoft.sql/servers/databases`, `microsoft.storage/storageaccounts`, `microsoft.storagecache/caches`, `microsoft.streamanalytics/streamingjobs`, `microsoft.synapse/workspaces`, `microsoft.timeseriesinsights/environments`, `microsoft.videoindexer/accounts`, `microsoft.web/sites`, `microsoft.workloadmonitor/monitors`, `resourcegroup` and `subscription`.

* `solutions` - (Optional) A list of the related Log Analytics solutions for the function. Possible values are `AADDomainServices`, `ADAssessment`, `ADAssessmentPlus`, `ADReplication`, `ADSecurityAssessment`, `AlertManagement`, `AntiMalware`, `ApplicationInsights`, `AzureAssessment`, `AzureSecurityOfThings`, `AzureSentinelDSRE`, `AzureSentinelPrivatePreview`, `BehaviorAnalyticsInsights`, `ChangeTracking`, `CompatibilityAssessment`, `ContainerInsights`, `Containers`, `CustomizedWindowsEventsFiltering`, `DeviceHealthProd`, `DnsAnalytics`, `ExchangeAssessment`, `ExchangeOnlineAssessment`, `IISAssessmentPlus`, `InfrastructureInsights`, `InternalWindowsEvent`, `LogManagement`, `Microsoft365Analytics`, `NetworkMonitoring`, `SCCMAssessmentPlus`, `SCOMAssessment`, `SCOMAssessmentPlus`, `Security`, `SecurityCenter`, `SecurityCenterFree`, `SecurityInsights`, `ServiceMap`, `SfBAssessment`, `SfBOnlineAssessment`, `SharePointOnlineAssessment`, `SPAssessment`, `SQLAdvancedThreatProtection`, `SQLAssessment`, `SQLAssessmentPlus`, `SQLDataClassification`, `SQLThreatDetection`, `SQLVulnerabilityAssessment`, `SurfaceHub`, `Updates`, `VMInsights`, `WEFInternalUat`, `WEF_10x`, `WEF_10xDSRE`, `WaaSUpdateInsights`, `WinLog`, `WindowsClientAssessmentPlus`, `WindowsEventForwarding`, `WindowsFirewall`, `WindowsServerAssessment`, `WireData` and `WireData2`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Log Analytics Query Pack Query.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Log Analytics Query Pack Query.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Query Pack Query.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Query Pack Query.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Query Pack Query.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Query Pack Query.

## Import

Log Analytics Query Pack Queries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_query_pack_query.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.OperationalInsights/queryPacks/queryPack1/queries/15b49e87-8555-4d92-8a7b-2014b469a9df
```
