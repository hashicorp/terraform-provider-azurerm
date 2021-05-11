// specifies the default hour (UTC) at which tests should be triggered, if enabled
var defaultStartHour = 0

// specifies the default level of parallelism per-service-package
var defaultParallelism = 20

// specifies the default version of Terraform Core which should be used for testing
var defaultTerraformCoreVersion = "0.15.3"

var locations = mapOf(
        "public" to LocationConfiguration("westeurope", "eastus2", "francecentral", false),
        "germany" to LocationConfiguration("germanynortheast", "germanycentral", "", false)
)

// specifies the list of Azure Environments where tests should be run nightly
var runNightly = mapOf(
        "public" to true
)

// specifies a list of services which should be run with a custom test configuration
var serviceTestConfigurationOverrides = mapOf(
        // these tests all conflict with one another
        "authorization" to testConfiguration(1, defaultStartHour),

        //Blueprints are constrained on the number of targets available - these execute quickly, so can be serialised
        "blueprints" to testConfiguration(1, defaultStartHour),

        // The AKS API has a low rate limit
        "containers" to testConfiguration(5, defaultStartHour),

        // Data Lake has a low quota
        "datalake" to testConfiguration(2, defaultStartHour),

        // HPC Cache has a 4 instance per subscription quota as of early 2021
        "hpccache" to testConfiguration(3, defaultStartHour),

        // HSM has low quota and potentially slow recycle time
        "hsm" to testConfiguration(1, defaultStartHour),

        // Log Analytics Clusters have a max deployments of 2 - parallelism set to 1 or `importTest` fails
        "loganalytics" to testConfiguration(1, defaultStartHour),

        // netapp has a max of 20 accounts per subscription so lets limit it to 10 to account for broken ones
        "netapp" to testConfiguration(10, defaultStartHour),

        // servicebus quotas are limited and we experience failures if tests
        // execute too quickly as we run out of namespaces in the sub
        "servicebus" to testConfiguration(10, defaultStartHour),

        // SignalR only allows provisioning one "Free" instance at a time,
        // which is used in multiple tests
        "signalr" to testConfiguration(1, defaultStartHour),

        // Spring Cloud only allows a max of 10 provisioned
        "springcloud" to testConfiguration(5, defaultStartHour),

        // Currently have a quota of 10 nodes, 3 nodes required per test so lets limit it to 3
        "vmware" to testConfiguration(3, defaultStartHour)
)
