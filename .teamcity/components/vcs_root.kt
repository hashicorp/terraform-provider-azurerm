import jetbrains.buildServer.configs.kotlin.v2019_2.vcs.GitVcsRoot

object providerRepository : GitVcsRoot({
    name = "terraform-provider-azurerm"
    url = "https://github.com/terraform-providers/terraform-provider-azurerm.git"
    agentCleanPolicy = AgentCleanPolicy.ALWAYS
    agentCleanFilesPolicy = AgentCleanFilesPolicy.ALL_UNTRACKED
    branchSpec = "+:*"
    branch = "refs/heads/master"
    authMethod = anonymous()
})
