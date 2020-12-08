package tests

import ClientConfiguration

fun TestConfiguration() : ClientConfiguration {
    return ClientConfiguration("clientId", "clientSecret", "subscriptionId", "tenantId", "clientIdAlt", "clientSecretAlt", "subscriptionIdAlt")
}