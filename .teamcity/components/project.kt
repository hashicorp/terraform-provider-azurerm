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

        var buildConfigs = buildConfigurationsForServices(services, providerName, environment, configuration)
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
        var buildConfig = service.buildConfiguration(providerName, runNightly, testConfig.startHour, testConfig.parallelism, testConfig.daysOfWeek, testConfig.daysOfMonth, testConfig.timeout)

        buildConfig.params.ConfigureAzureSpecificTestParameters(environment, config, locationsToUse,  testConfig.useAltSubscription, testConfig.useDevTestSubscription)

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

class testConfiguration(parallelism: Int = defaultParallelism, startHour: Int = defaultStartHour, daysOfWeek: String = defaultDaysOfWeek, daysOfMonth: String = defaultDaysOfMonth, timeout: Int = defaultTimeout, useAltSubscription: Boolean = false, useDevTestSubscription: Boolean = false, locationOverride: LocationConfiguration = LocationConfiguration("","","", false)) {
    var parallelism = parallelism
    var startHour = startHour
    var daysOfWeek = daysOfWeek
    var daysOfMonth = daysOfMonth
    var timeout = timeout
    var useAltSubscription = useAltSubscription
    var useDevTestSubscription = useDevTestSubscription
    var locationOverride = locationOverride
}