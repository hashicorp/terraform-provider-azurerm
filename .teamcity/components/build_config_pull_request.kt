import jetbrains.buildServer.configs.kotlin.AbsoluteId
import jetbrains.buildServer.configs.kotlin.BuildType

class PullRequest(val displayName: String, val environment: String, val vcsRootId: String) {

    fun buildConfiguration(providerName : String) : BuildType {
        return BuildType {
            // TC needs a consistent ID for dynamically generated packages
            id(uniqueID(providerName))

            name = displayName

            vcs {
                root(rootId = AbsoluteId(vcsRootId))
                cleanCheckout = true
            }

            steps {
                val packageName = "\"%SERVICES%\""

                configureGoEnv()
                downloadTerraformBinary()
                downloadVCRCassettes(packageName)
                runAcceptanceTestsForPullRequest(packageName)
                uploadVCRCassettes(packageName)
                postTestResultsToGitHubPullRequest()
            }

            failureConditions {
                errorMessage = true
            }

            features {
                golang()
                buildCacheFeature()
            }

            params {
                terraformAcceptanceTestParameters(defaultParallelism, "TestAcc", defaultTimeout)
                terraformAcceptanceTestsFlag()
                terraformShouldPanicForSchemaErrors()
                terraformCoreBinaryTesting()
                readOnlySettings()
                goCache()

                text("SERVICES", "portal")
            }
        }
    }

    fun uniqueID(provider : String) : String {
        return "%s_PR_%s".format(provider.uppercase(), environment.uppercase())
    }
}
