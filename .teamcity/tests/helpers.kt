package tests

import ClientConfiguration

fun TestConfiguration() : ClientConfiguration {
    return ClientConfiguration("clientId", "clientSecret", "subscriptionId", "subscriptionIdAlt", "tenantId", "clientIdAlt", "clientSecretAlt")
}