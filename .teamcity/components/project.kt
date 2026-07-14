import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Project

const val providerName = "azurerm"
var enableTestTriggersGlobally = false

fun AzureRM(environment: String, configuration : ClientConfiguration) : Project {
    return Project{
        // @tombuildsstuff: this temporary flag enables/disables all triggers, allowing a migration between CI servers
        enableTestTriggersGlobally = configuration.enableTestTriggersGlobally

        var pullRequestBuildConfig = pullRequestBuildConfiguration(environment, configuration)
        buildType(pullRequestBuildConfig)

        var cacheBuildConfig = buildConfigurationForCache(environment, configuration)
        buildType(cacheBuildConfig)

        var buildConfigs = buildConfigurationsForServices(services, providerName, environment, configuration)
        buildConfigs.forEach { buildConfiguration ->
            buildType(buildConfiguration)
        }
        if (configuration.createBetaProject) {
            subProject(AzureRMBetaVersion(environment, configuration))
        }
    }
}

fun AzureRMBetaVersion(environment: String, configuration : ClientConfiguration) : Project {
    return Project{
        id("AzureRMBetaVersion")
        name = "Beta"

        enableTestTriggersGlobally = configuration.enableTestTriggersGlobally

        var buildConfigs = buildConfigurationsForServicesBetaVersion(services, providerName, environment, configuration)
        buildConfigs.forEach { buildConfiguration ->
            buildType(buildConfiguration)
        }
    }
}

fun buildConfigurationsForServices(services: Map<String, String>, providerName : String, environment: String, config : ClientConfiguration): List<BuildType> {
    var list = ArrayList<BuildType>()
    var locationsForEnv = locations[environment]!!

    services.forEach { (serviceName, displayName) ->
        var defaultTestConfig = testConfiguration()
        var testConfig = serviceTestConfigurationOverrides.getOrDefault(serviceName, defaultTestConfig)
        var locationsToUse = if (testConfig.locationOverride.primary != "") testConfig.locationOverride else locationsForEnv
        var runNightly = runNightly.getOrDefault(environment, false)

        var service = serviceDetails(serviceName, displayName, environment, config.vcsRootId)
        var buildConfig = service.buildConfiguration(providerName, runNightly, testConfig.startHour, testConfig.parallelism, testConfig.daysOfWeek, testConfig.daysOfMonth, testConfig.timeout, testConfig.disableTriggers, false)

        buildConfig.params.ConfigureAzureSpecificTestParameters(environment, config, locationsToUse, testConfig.useAltSubscription, testConfig.useDevTestSubscription, enableBetaVersion = false)

        list.add(buildConfig)
    }

    return list
}

fun buildConfigurationsForServicesBetaVersion(services: Map<String, String>, providerName : String, environment: String, config : ClientConfiguration): List<BuildType> {
    var list = ArrayList<BuildType>()
    var locationsForEnv = locations[environment]!!

    services.forEach { (serviceName, displayName) ->
        var defaultTestConfig = testConfiguration()
        var testConfig = serviceTestConfigurationOverrides.getOrDefault(serviceName, defaultTestConfig)
        var locationsToUse = if (testConfig.locationOverride.primary != "") testConfig.locationOverride else locationsForEnv
        var runNightly = runNightly.getOrDefault(environment, false)

        var service = serviceDetails(serviceName, displayName, environment, config.vcsRootId)
        var buildConfig = service.buildConfiguration(providerName, runNightly, testConfig.startHour, testConfig.parallelism, testConfig.weeklyTestDay, testConfig.daysOfMonth, testConfig.timeout, testConfig.disableTriggers, true)

        buildConfig.params.ConfigureAzureSpecificTestParameters(environment, config, locationsToUse, testConfig.useAltSubscription, testConfig.useDevTestSubscription, enableBetaVersion = true)

        list.add(buildConfig)
    }

    return list
}

fun pullRequestBuildConfiguration(environment: String, config: ClientConfiguration) : BuildType {
    var locationsForEnv = locations.get(environment)!!
    var pullRequest = pullRequest("! Run Pull Request", environment, config.vcsRootId)
    var buildConfiguration = pullRequest.buildConfiguration(providerName)
    buildConfiguration.params.ConfigureAzureSpecificTestParameters(environment, config, locationsForEnv)
    return buildConfiguration
}

fun buildConfigurationForCache(environment: String, config: ClientConfiguration) : BuildType {
    return buildCacheConfiguration(environment, config.vcsRootId).buildConfiguration(providerName)
}

class testConfiguration(parallelism: Int = defaultParallelism, startHour: Int = defaultStartHour, daysOfWeek: String = defaultDaysOfWeek, weeklyTestDay: String = defaultWeeklyDay, daysOfMonth: String = defaultDaysOfMonth, timeout: Int = defaultTimeout, useAltSubscription: Boolean = false, useDevTestSubscription: Boolean = false, locationOverride: LocationConfiguration = LocationConfiguration("","","", false), terraformCoreOverride: String = defaultTerraformCoreVersion, disableTriggers: Boolean = false) {
    var parallelism = parallelism
    var startHour = startHour
    var daysOfWeek = daysOfWeek
    var weeklyTestDay = weeklyTestDay
    var daysOfMonth = daysOfMonth
    var timeout = timeout
    var useAltSubscription = useAltSubscription
    var useDevTestSubscription = useDevTestSubscription
    var locationOverride = locationOverride
    var terraformCoreOverride = terraformCoreOverride
    var disableTriggers = disableTriggers
}
