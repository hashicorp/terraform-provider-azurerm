import jetbrains.buildServer.configs.kotlin.ParametrizedWithType

class ClientConfiguration(var clientId: String,
                          var clientSecret: String,
                          val subscriptionId : String,
                          val tenantId : String,
                          val clientIdAlt: String,
                          val clientSecretAlt: String,
                          val subscriptionIdAlt : String,
                          val subscriptionIdDevTest : String,
                          val tenantIdAlt : String,
                          val subscriptionIdAltTenant : String,
                          val principalIdAltTenant : String,
                          val vcsRootId : String,
                          val enableTestTriggersGlobally : Boolean) {
}

class LocationConfiguration(var primary : String, var secondary : String, var tertiary : String, var rotate : Boolean) {
}

fun ParametrizedWithType.ConfigureAzureSpecificTestParameters(environment: String, config: ClientConfiguration, locationsForEnv: LocationConfiguration, useAltSubscription: Boolean = false, useDevTestSubscription: Boolean = false) {
    hiddenPasswordVariable("env.ARM_CLIENT_ID", config.clientId, "The Client ID of the Application used for Testing")
    hiddenPasswordVariable("env.ARM_CLIENT_ID_ALT", config.clientIdAlt, "The Client ID of the Alternate Application used for Testing")
    hiddenPasswordVariable("env.ARM_CLIENT_SECRET", config.clientSecret, "The Client Secret of the Application used for Testing")
    hiddenPasswordVariable("env.ARM_CLIENT_SECRET_ALT", config.clientSecretAlt, "The Client Secret of the Alternate Application used for Testing")
    hiddenVariable("env.ARM_ENVIRONMENT", environment, "The Azure Environment in which the tests are running")
    hiddenVariable("env.ARM_PROVIDER_DYNAMIC_TEST", "%b".format(locationsForEnv.rotate), "Should tests rotate between the supported regions?")
    if (useAltSubscription) {
        hiddenPasswordVariable("env.ARM_SUBSCRIPTION_ID", config.subscriptionIdAlt, "The ID of the Azure Subscription used for Testing (Alt swapped in)")
        hiddenPasswordVariable("env.ARM_SUBSCRIPTION_ID_ALT", config.subscriptionId, "The ID of the Alternate Azure Subscription used for Testing (Main swapped out)")
    } else if (useDevTestSubscription) {
        hiddenPasswordVariable("env.ARM_SUBSCRIPTION_ID", config.subscriptionIdDevTest, "The ID of the Azure Subscription used for Testing (DevTest swapped in)")
        hiddenPasswordVariable("env.ARM_SUBSCRIPTION_ID_ALT", config.subscriptionId, "The ID of the Alternate Azure Subscription used for Testing (Main swapped out)")
    } else {
        hiddenPasswordVariable("env.ARM_SUBSCRIPTION_ID", config.subscriptionId, "The ID of the Azure Subscription used for Testing")
        hiddenPasswordVariable("env.ARM_SUBSCRIPTION_ID_ALT", config.subscriptionIdAlt, "The ID of the Alternate Azure Subscription used for Testing")
    }
    hiddenPasswordVariable("env.ARM_TENANT_ID", config.tenantId, "The ID of the Azure Tenant used for Testing")
    hiddenPasswordVariable("env.ARM_TENANT_ID_ALT", config.tenantIdAlt, "The ID of the Secondary Azure Tenant used for Testing")
    hiddenPasswordVariable("env.ARM_SUBSCRIPTION_ID_ALT_TENANT", config.subscriptionIdAltTenant, "The ID of the Azure Subscription attached to the Secondary Azure Tenant used for Testing")
    hiddenPasswordVariable("env.ARM_PRINCIPAL_ID_ALT_TENANT", config.principalIdAltTenant, "The Object ID of the Service Principal in the Secondary Azure Tenant representing the Application used for Testing")
    hiddenVariable("env.ARM_TEST_LOCATION", locationsForEnv.primary, "The Primary region which should be used for testing")
    hiddenVariable("env.ARM_TEST_LOCATION_ALT", locationsForEnv.secondary, "The Secondary region which should be used for testing")
    hiddenVariable("env.ARM_TEST_LOCATION_ALT2", locationsForEnv.tertiary, "The Tertiary region which should be used for testing")
    hiddenVariable("env.ARM_FIVEPOINTZERO_BETA", "false", "Opt into the 5.0 beta")
}