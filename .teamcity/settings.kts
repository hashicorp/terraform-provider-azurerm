import jetbrains.buildServer.configs.kotlin.*

version = "2025.11"

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
var createBetaProject = DslContext.getParameter("createBetaProject", "false").equals("true", ignoreCase = true)
var enableTestTriggersGlobally = DslContext.getParameter("enableTestTriggersGlobally", "true").equals("true", ignoreCase = true)
var emailAddressAccTests = DslContext.getParameter("emailAddressAccTests", "")
var gitHubRepo = DslContext.getParameter("gitHubRepo", "hashicorp/terraform-provider-azurerm")
var gitPat = DslContext.getParameter("gitPat", "")
var teamcityToken = DslContext.getParameter("teamcityToken", "")
var betaVersionEnvVar = DslContext.getParameter("betaVersionEnvVar", "env.ARM_FIVEPOINTZERO_BETA")
var labelSuccess = DslContext.getParameter("labelSuccess", "teamcity-passed")
var labelFailure = DslContext.getParameter("labelFailure", "teamcity-failed")
var labelOutdated = DslContext.getParameter("labelOutdated", "teamcity-outdated")
var labelNewFailure = DslContext.getParameter("labelNewFailure", "teamcity-new-failure")
var applyTestingLabelsEnabled = DslContext.getParameter("applyTestingLabelsEnabled", "true").equals("true", ignoreCase = true)


var clientConfig = ClientConfiguration(
    clientId,
    clientSecret,
    subscriptionId,
    tenantId,
    clientIdAlt,
    clientSecretAlt,
    subscriptionIdAlt,
    subscriptionIdDevTest,
    tenantIdAlt,
    subscriptionIdAltTenant,
    principalIdAltTenant,
    vcsRootId,
    createBetaProject,
    enableTestTriggersGlobally,
    emailAddressAccTests,
    gitHubRepo,
    gitPat,
    teamcityToken,
    betaVersionEnvVar,
    labelSuccess,
    labelFailure,
    labelOutdated,
    labelNewFailure,
    applyTestingLabelsEnabled
)

project(AzureRM(environment, clientConfig))
