import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot

object providerRepository : GitVcsRoot({
    name = "terraform-provider-azurerm"
    url = "https://github.com/hashicorp/terraform-provider-azurerm.git"
    agentCleanPolicy = AgentCleanPolicy.ALWAYS
    agentCleanFilesPolicy = AgentCleanFilesPolicy.ALL_UNTRACKED
    branchSpec = "+:*"
    branch = "refs/heads/main"
    authMethod = anonymous()
})
