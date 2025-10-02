import jetbrains.buildServer.configs.kotlin.ParametrizedWithType

// Configuration classes for Azure-specific settings
class LocationConfiguration(val primary: String, val secondary: String, val ternary: String, val rotate: Boolean)

class ClientConfiguration(
    val clientId: String,
    val clientSecret: String,
    val subscriptionId: String,
    val tenantId: String,
    val clientIdAlt: String,
    val clientSecretAlt: String,
    val subscriptionIdAlt: String,
    val subscriptionIdDevTest: String,
    val tenantIdAlt: String,
    val subscriptionIdAltTenant: String,
    val principalIdAltTenant: String,
    val vcsRootId: String,
    val enableTestTriggersGlobally: Boolean = false,
    val emailAddressAccTests: String
)

// Extension function to configure Azure-specific test parameters
fun ParametrizedWithType.ConfigureAzureSpecificTestParameters(
    environment: String,
    config: ClientConfiguration,
    locations: LocationConfiguration,
    useAltSubscription: Boolean = false,
    useDevTestSubscription: Boolean = false
) {
    // Configure Azure-specific parameters here
    // This function configures environment-specific settings for Azure tests
}
