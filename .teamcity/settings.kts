import AzureRM
import ClientConfiguration
import jetbrains.buildServer.configs.kotlin.v2019_2.*

version = "2020.1"

var clientId = DslContext.getParameter("clientId", "")
var clientSecret = DslContext.getParameter("clientSecret", "")
var subscriptionId = DslContext.getParameter("subscriptionId", "")
var subscriptionIdAlt = DslContext.getParameter("subscriptionIdAlt", "")
var tenantId = DslContext.getParameter("tenantId", "")
var environment = DslContext.getParameter("environment", "public")
var clientIdAlt = DslContext.getParameter("clientIdAlt", "")
var clientSecretAlt = DslContext.getParameter("clientSecretAlt", "")

var clientConfig = ClientConfiguration(clientId, clientSecret, subscriptionId, tenantId, clientIdAlt, clientSecretAlt, subscriptionIdAlt )

project(AzureRM(environment, clientConfig))
