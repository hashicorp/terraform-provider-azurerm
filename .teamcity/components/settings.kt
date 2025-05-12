/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// specifies the default hour (UTC) at which tests should be triggered, if enabled
var defaultStartHour = 23

// specifies the default level of parallelism per-service-package
var defaultParallelism = 20

// specifies the default build timeout in hours
var defaultTimeout = 12

// specifies the default version of Terraform Core which should be used for testing
var defaultTerraformCoreVersion = "1.5.1"

// This represents a cron view of days of the week, Monday - Friday.
const val defaultDaysOfWeek = "2,3,4,5,6"

// Cron value for any day of month
const val defaultDaysOfMonth = "*"

var locations = mapOf(
        "public" to LocationConfiguration("westeurope", "eastus2", "westus2", true)
)

// specifies the list of Azure Environments where tests should be run nightly
var runNightly = mapOf(
        "public" to true
)

// specifies a list of services which should be run with a custom test configuration
var serviceTestConfigurationOverrides = mapOf(

        // Server is only available in certain locations
        "analysisservices" to testConfiguration(locationOverride = LocationConfiguration("westus", "northeurope", "southcentralus", true)),

        // App Service Plans for Linux are currently unavailable in WestUS2
        "appservice" to testConfiguration(startHour = 3, daysOfWeek = "2,4,6", locationOverride = LocationConfiguration("westeurope", "westus2", "eastus2", true)),

        // Arc Kubernetes Provisioned Cluster is only available in certain locations
        "arckubernetes" to testConfiguration(locationOverride = LocationConfiguration("australiaeast", "eastus", "westeurope", true)),

        // these tests all conflict with one another
        "authorization" to testConfiguration(parallelism = 1),

        // HCICluster is only available in certain locations
        "azurestackhci" to testConfiguration(locationOverride = LocationConfiguration("australiaeast", "eastus", "westeurope", true)),

        //Blueprints are constrained on the number of targets available - these execute quickly, so can be serialised
        "blueprints" to testConfiguration(parallelism = 1),

        // CDN is only available in certain locations
        "cdn" to testConfiguration(locationOverride = LocationConfiguration("centralus", "eastus2", "westeurope", true)),

        // Chaosstudio is only available in certain locations
        "chaosstudio" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "eastus", "westus", false)),

        // "cognitive" is expensive - Monday, Wednesday, Friday
        // cognitive is only available in certain locations
        "cognitive" to testConfiguration(daysOfWeek = "2,4,6", locationOverride = LocationConfiguration("westeurope", "eastus", "southcentralus", true)),

        // Cosmos is only available in certain locations
        "cosmos" to testConfiguration(locationOverride = LocationConfiguration("westus", "northeurope", "eastus2", true)),

        // Confidential Ledger
        "confidentialledger" to testConfiguration(locationOverride = LocationConfiguration("eastus","southcentralus","westeurope", false)),

        // Container App Managed Environments are limited to 20 per location, using 10 as they can take some time to clear
        // Enable rotation test to mitigate resource burden in a single region
        "containerapps" to testConfiguration(parallelism = 10, locationOverride = LocationConfiguration("eastus2","westus2","southcentralus", true)),

        // The AKS API has a low rate limit
        "containers" to testConfiguration(parallelism = 5, locationOverride = LocationConfiguration("eastus","westeurope","eastus2", false), timeout = 18),

        // Custom Providers is only available in certain locations
        "customproviders" to testConfiguration(locationOverride = LocationConfiguration("eastus", "westus2", "westeurope", true)),

        // Dashboard is only available in certain locations
        "dashboard" to testConfiguration(locationOverride = LocationConfiguration("eastus", "westus2", "eastus2", false)),

        // Datadog is available only in WestUS2 region
        "datadog" to testConfiguration(locationOverride = LocationConfiguration("westus2", "westus2", "centraluseuap", false)),

        // data factory uses NC class VMs which are not available in eastus2
        "datafactory" to testConfiguration(daysOfWeek = "2,4,6", locationOverride = LocationConfiguration("westeurope", "southeastasia", "westus2", false)),

        // Dev Center only available in some regions / has a quota of 5
        "devcenter" to testConfiguration(parallelism = 4, locationOverride = LocationConfiguration("westeurope", "uksouth", "canadacentral", false)),

        // "hdinsight" is super expensive - G class VM's are not available in westus2, quota only available in westeurope currently
        "hdinsight" to testConfiguration(daysOfWeek = "2,4,6", locationOverride = LocationConfiguration("westeurope", "southeastasia", "eastus2", false)),

        // Elastic can't provision many in parallel
        "elastic" to testConfiguration(parallelism = 1),

        // ElasticSAN has 5 instance quota per subscription per region
        "elasticsan" to testConfiguration(parallelism = 4),

        // HSM has low quota and potentially slow recycle time, Only run on Mondays
        "hsm" to testConfiguration(parallelism = 1, daysOfWeek = "1"),

        // IoT Central is only available in certain locations
        "iotcentral" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "southeastasia", "eastus2", false)),

        // IoT Hub Device Update is only available in certain locations
        "iothub" to testConfiguration(locationOverride = LocationConfiguration("eastus", "eastus2", "westus2", false)),

        // Lab Service is only available in certain locations
        "labservice" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "eastus", "westus", false)),

        // load balancer global tire Public IP is only available in
        "loadbalancer" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "eastus2", "westus", false)),

        // Log Analytics Clusters have a max deployments of 2 - parallelism set to 1 or `importTest` fails
        "loganalytics" to testConfiguration(parallelism = 1),

        // Logic uses app service which is only available in certain locations
        "logic" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false)),

        // Logz is only available in certain locations
        "logz" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "westus2", "eastus2", false)),

        // Maps is only available in certain locations
        "maps" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "westus2", "eastus", false)),

        // MobileNetwork is only available in certain locations
        "mobilenetwork" to testConfiguration(locationOverride = LocationConfiguration("eastus", "westeurope", "centraluseuap", false)),

        // Mongocluster free tier is currently only available in southindia
        "mongocluster" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "eastus2", "southindia", false)),

        // MSSQl uses app service which is only available in certain locations
        "mssql" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false)),

        // MSSQL Managed Instance creation can impact the service so limit the frequency and number of tests
        "mssqlmanagedinstance" to testConfiguration(parallelism = 4, daysOfWeek = "7", locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false), timeout = 18),

        // MySQL has quota available in certain locations
        "mysql" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false)),

        // netapp has a max of 10 accounts and the max capacity of pool is 25 TiB per subscription so lets limit it to 1 to account for broken ones, run Monday, Wednesday, Friday
        "netapp" to testConfiguration(parallelism = 1, daysOfWeek = "2,4,6", locationOverride = LocationConfiguration("westeurope", "eastus2", "westus2", false)),

        // network has increased timeout to accommodate the custom_ip_prefix resource
        "network" to testConfiguration(timeout = 24),

        // Run New Relic testcases in Canary Region to avoid generating pollution test data in Production Region, which will cause side effect in Service Partner's Database
        "newrelic" to testConfiguration(locationOverride = LocationConfiguration("centraluseuap", "eastus", "eastus", false)),

        // Network Function is only available in certain locations
        "networkfunction" to testConfiguration(locationOverride = LocationConfiguration("westus2", "eastus2", "westeurope", false)),

        // Network Regional Tire Public IP is only available in
        "network" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "eastus2", "westus", false)),

        // Orbital is only available in certain locations
        "orbital" to testConfiguration(locationOverride = LocationConfiguration("eastus", "southcentralus", "westus2", false)),

        "paloalto" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "eastus", "westus", false)),

        "policy" to testConfiguration(useAltSubscription = true),

        "postgres" to testConfiguration(locationOverride = LocationConfiguration("northeurope", "centralus", "westeurope", false)),

        // Private DNS Resolver is only available in certain locations
        "privatednsresolver" to testConfiguration(locationOverride = LocationConfiguration("eastus", "westus3", "westeurope", true)),

        // Purview Accounts are only available in certain locations
        "purview" to testConfiguration(locationOverride = LocationConfiguration("eastus", "southcentralus", "westus", true)),

        // redisenterprise is costly - Monday, Wednesday, Friday
        "redisenterprise" to testConfiguration(daysOfWeek = "2,4,6"),

        // servicebus quotas are limited and we experience failures if tests
        // execute too quickly as we run out of namespaces in the sub
        "servicebus" to testConfiguration(parallelism = 10),

        // Spring Cloud only allows a max of 10 provisioned
        "springcloud" to testConfiguration(parallelism = 5),

        // SQL has quota available in certain locations
        "sql" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false)),

        // HPC Cache has a 4 instance per subscription quota as of early 2021
        "storagecache" to testConfiguration(parallelism = 3, daysOfWeek = "2,4,6"),

        "storagemover" to testConfiguration(locationOverride = LocationConfiguration("eastus", "eastus2", "westus3", false)),

        // StreamAnalytics has quota available in certain locations
        "streamanalytics" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false)),

        // Synapse is only available in certain locations
        "synapse" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false)),

        // ServiceNetworking is available in certain locations
        "servicenetworking" to testConfiguration(locationOverride = LocationConfiguration("eastus","westus","westeurope", false)),

        // Currently, we have insufficient quota to actually run these, but there are a few nodes in West Europe, so we'll pin it there for now
        "vmware" to testConfiguration(parallelism = 3, locationOverride = LocationConfiguration("westeurope", "westus2", "eastus2", false)),

        // In general, Azure Voice Service is available in several different regions, but each subscription will only be allowlisted for specific regions(`westcentralus`, `westcentralus`, `westcentralus`).
        // Only the regions (`westcentralus`) is specified since the devtest subscription does not support creating resource group for the other two regions.
        "voiceservices" to testConfiguration(parallelism = 3, locationOverride = LocationConfiguration("westcentralus", "westcentralus", "westcentralus", false)),

        // Offset start hour to avoid collision with new App Service, reduce frequency of testing days
        "web" to testConfiguration(startHour = 3, daysOfWeek = "1,3,5", locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", true)),

        // Workloads has quota available in certain locations
        "workloads" to testConfiguration(parallelism = 1, locationOverride = LocationConfiguration("eastus", "westeurope", "francecentral", false))
)
