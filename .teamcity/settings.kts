import jetbrains.buildServer.configs.kotlin.*

version = "2023.11"

var clientId = DslContext.getParameter("clientId", "")
var clientSecret = DslContext.getParameter("clientSecret", "")
var subscriptionId = DslContext.getParameter("subscriptionId", "")
var subscriptionIdAlt = DslContext.getParameter("subscriptionIdAlt", "")
var subscriptionIdDevTest = DslContext.getParameter("subscriptionIdDevTest", "")
var tenantId = DslContext.getParameter("tenantId", "")
var environment = DslContext.getParameter("environment", "public")
var clientIdAlt = DslContext.getParameter("clientIdAlt", "")
var clientSecretAlt = DslContext.getParameter("clientSecretAlt", "")
var tenantIdAlt = DslContext.getParameter("tenantIdAlt", "")
var subscriptionIdAltTenant = DslContext.getParameter("subscriptionIdAltTenant", "")
var principalIdAltTenant = DslContext.getParameter("principalIdAltTenant", "")
var vcsRootId = DslContext.getParameter("vcsRootId", "TF_HashiCorp_AzureRM_Repository")
var enableTestTriggersGlobally = DslContext.getParameter("enableTestTriggersGlobally", "true").equals("true", ignoreCase = true)

var clientConfig = ClientConfiguration(clientId, clientSecret, subscriptionId, tenantId, clientIdAlt, clientSecretAlt, subscriptionIdAlt, subscriptionIdDevTest, tenantIdAlt, subscriptionIdAltTenant, principalIdAltTenant, vcsRootId, enableTestTriggersGlobally)

project(AzureRM(environment, clientConfig))
