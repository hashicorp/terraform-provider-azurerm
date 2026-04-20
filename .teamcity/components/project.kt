import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Project

const val providerName = "azurerm"
var enableTestTriggersGlobally = false

fun azureRM(environment: String, configuration : ClientConfiguration) : Project {
    return Project{
        // @tombuildsstuff: this temporary flag enables/disables all triggers, allowing a migration between CI servers
        enableTestTriggersGlobally = configuration.enableTestTriggersGlobally

        val pullRequestBuildConfig = pullRequestBuildConfiguration(environment, configuration)
        buildType(pullRequestBuildConfig)

        val cacheBuildConfig = buildConfigurationForCache(environment, configuration)
        buildType(cacheBuildConfig)

        val buildConfigs = buildConfigurationsForServices(services, providerName, environment, configuration)
        buildConfigs.forEach { buildConfiguration ->
            buildType(buildConfiguration)
        }
    }
}

fun buildConfigurationsForServices(services: Map<String, String>, providerName : String, environment: String, config : ClientConfiguration): List<BuildType> {
    val list = ArrayList<BuildType>()
    val locationsForEnv = locations[environment]!!

    services.forEach { (serviceName, displayName) ->
        val defaultTestConfig = testConfiguration()
        val testConfig = serviceTestConfigurationOverrides.getOrDefault(serviceName, defaultTestConfig)
        val locationsToUse = if (testConfig.locationOverride.primary != "") testConfig.locationOverride else locationsForEnv
        val runNightly = runNightly.getOrDefault(environment, false)

        val service = ServiceDetails(serviceName, displayName, environment, config.vcsRootId)
        val buildConfig = service.buildConfiguration(providerName, runNightly, testConfig.startHour, testConfig.parallelism, testConfig.daysOfWeek, testConfig.daysOfMonth, testConfig.timeout, testConfig.disableTriggers)

        buildConfig.params.configureAzureSpecificTestParameters(environment, config, locationsToUse,  testConfig.useAltSubscription, testConfig.useDevTestSubscription)

        list.add(buildConfig)
    }

    return list
}

fun pullRequestBuildConfiguration(environment: String, config: ClientConfiguration) : BuildType {
    val locationsForEnv = locations[environment]!!
    val pullRequest = PullRequest("! Run Pull Request", environment, config.vcsRootId)
    val buildConfiguration = pullRequest.buildConfiguration(providerName)
    buildConfiguration.params.configureAzureSpecificTestParameters(environment, config, locationsForEnv)
    return buildConfiguration
}

fun buildConfigurationForCache(environment: String, config: ClientConfiguration) : BuildType {
    return BuildCacheConfiguration(environment, config.vcsRootId).buildConfiguration(providerName)
}

class testConfiguration(var parallelism: Int = defaultParallelism,
                        var startHour: Int = defaultStartHour,
                        var daysOfWeek: String = defaultDaysOfWeek,
                        var daysOfMonth: String = defaultDaysOfMonth,
                        var timeout: Int = defaultTimeout,
                        var useAltSubscription: Boolean = false,
                        var useDevTestSubscription: Boolean = false,
                        var locationOverride: LocationConfiguration = LocationConfiguration("", "", "", false),
                        var terraformCoreOverride: String = defaultTerraformCoreVersion, // Note: never used/overridden, but we might want to for some reason?
                        var disableTriggers: Boolean = false
)
