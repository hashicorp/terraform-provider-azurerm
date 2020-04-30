import AzureRM
import ClientConfiguration
import jetbrains.buildServer.configs.kotlin.v2019_2.*

version = "2019.2"

var clientId = DslContext.getParameter("clientId", "")
var clientSecret = DslContext.getParameter("clientSecret", "")
var subscriptionId = DslContext.getParameter("subscriptionId", "")
var tenantId = DslContext.getParameter("tenantId", "")
var environment = DslContext.getParameter("environment", "public")

var clientConfig = ClientConfiguration(clientId, clientSecret, subscriptionId, tenantId)

project(AzureRM(environment, clientConfig))