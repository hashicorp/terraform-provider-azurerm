import jetbrains.buildServer.configs.kotlin.v2019_2.BuildType
import jetbrains.buildServer.configs.kotlin.v2019_2.Project

const val providerName = "azurerm"

fun AzureRM(environment: String, configuration : ClientConfiguration) : Project {
    return Project{
        vcsRoot(providerRepository)

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

        var service = serviceDetails(serviceName, displayName, environment)
        var buildConfig = service.buildConfiguration(providerName, runNightly, testConfig.startHour, testConfig.parallelism, testConfig.daysOfWeek, testConfig.daysOfMonth)

        buildConfig.params.ConfigureAzureSpecificTestParameters(environment, config, locationsToUse,  testConfig.useAltSubscription, testConfig.useDevTestSubscription)

        list.add(buildConfig)
    }

    return list
}

fun pullRequestBuildConfiguration(environment: String, configuration: ClientConfiguration) : BuildType {
    var locationsForEnv = locations.get(environment)!!
    var pullRequest = pullRequest("! Run Pull Request", environment)
    var buildConfiguration = pullRequest.buildConfiguration(providerName)
    buildConfiguration.params.ConfigureAzureSpecificTestParameters(environment, configuration, locationsForEnv)
    return buildConfiguration
}

class testConfiguration(parallelism: Int = defaultParallelism, startHour: Int = defaultStartHour, daysOfWeek: String = defaultDaysOfWeek, daysOfMonth: String = defaultDaysOfMonth, useAltSubscription: Boolean = false, useDevTestSubscription: Boolean = false, locationOverride: LocationConfiguration = LocationConfiguration("","","", false)) {
    var parallelism = parallelism
    var startHour = startHour
    var daysOfWeek = daysOfWeek
    var daysOfMonth = daysOfMonth
    var useAltSubscription = useAltSubscription
    var useDevTestSubscription = useDevTestSubscription
    var locationOverride = locationOverride
}