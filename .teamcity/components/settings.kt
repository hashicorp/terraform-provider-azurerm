// specifies the default hour (UTC) at which tests should be triggered, if enabled
var defaultStartHour = 0

// specifies the default level of parallelism per-service-package
var defaultParallelism = 20

// specifies the default version of Terraform Core which should be used for testing
var defaultTerraformCoreVersion = "1.0.1"

// This represents a cron view of days of the week, Monday - Friday.
const val defaultDaysOfWeek = "2,3,4,5,6"

// Cron value for any day of month
const val defaultDaysOfMonth = "*"

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
        "authorization" to testConfiguration(parallelism = 1),

        //Blueprints are constrained on the number of targets available - these execute quickly, so can be serialised
        "blueprints" to testConfiguration(parallelism = 1),

        // "cognitive" is expensive - Monday, Wednesday, Friday
        "cognitive" to testConfiguration(daysOfWeek = "2,4,6"),

        // The AKS API has a low rate limit
        "containers" to testConfiguration(parallelism = 5),

        // Data Lake has a low quota
        "datalake" to testConfiguration(parallelism = 2),

        // "hdinsight" is super expensive
        "hdinsight" to testConfiguration(daysOfWeek = "2,4,6"),

        // HPC Cache has a 4 instance per subscription quota as of early 2021
        "hpccache" to testConfiguration(parallelism = 3, daysOfWeek = "2,4,6"),

        // HSM has low quota and potentially slow recycle time, Only run on Mondays
        "hsm" to testConfiguration(parallelism = 1, daysOfWeek = "1"),

        // Log Analytics Clusters have a max deployments of 2 - parallelism set to 1 or `importTest` fails
        "loganalytics" to testConfiguration(parallelism = 1),

        // netapp has a max of 20 accounts per subscription so lets limit it to 10 to account for broken ones, run Monday, Wednesday, Friday
        "netapp" to testConfiguration(parallelism = 10, daysOfWeek = "2,4,6"),

        // redisenterprise is costly - Monday, Wednesday, Friday
        "redisenterprise" to testConfiguration(daysOfWeek = "2,4,6"),

        // servicebus quotas are limited and we experience failures if tests
        // execute too quickly as we run out of namespaces in the sub
        "servicebus" to testConfiguration(parallelism = 10),

        // SignalR only allows provisioning one "Free" instance at a time,
        // which is used in multiple tests
        "signalr" to testConfiguration(parallelism = 1),

        // Spring Cloud only allows a max of 10 provisioned
        "springcloud" to testConfiguration(parallelism = 5),

        // Currently have a quota of 10 nodes, 3 nodes required per test so lets limit it to 3
        "vmware" to testConfiguration(parallelism = 3)
)
