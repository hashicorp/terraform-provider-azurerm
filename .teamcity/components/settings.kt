// specifies the default hour (UTC) at which tests should be triggered, if enabled
var defaultStartHour = 0

// specifies the default level of parallelism per-service-package
var defaultParallelism = 20

// specifies the default version of Terraform Core which should be used for testing
var defaultTerraformCoreVersion = "0.14.5"

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
        // Spring Cloud only allows a max of 10 provisioned
        "appplatform" to testConfiguration(5, defaultStartHour),

        // these tests all conflict with one another
        "authorization": to testConfiguration(1, defaultStartHour),

        // The AKS API has a low rate limit
        "containers" to testConfiguration(5, defaultStartHour),

        // Data Lake has a low quota
        "datalake" to testConfiguration(2, defaultStartHour),

        // HSM has low quota and potentially slow recycle time
        "hsm" to testConfiguration(1, defaultStartHour),

        // Log Analytics Clusters have a max deployments of 2 - parallelism set to 1 or `importTest` fails
        "loganalytics" to testConfiguration(1, defaultStartHour),

        // servicebus quotas are limited and we experience failures if tests
        // execute too quickly as we run out of namespaces in the sub
        "servicebus" to testConfiguration(10, defaultStartHour),

        // SignalR only allows provisioning one "Free" instance at a time,
        // which is used in multiple tests
        "signalr" to testConfiguration(1, defaultStartHour)
)
