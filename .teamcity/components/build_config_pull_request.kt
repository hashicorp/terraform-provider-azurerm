import jetbrains.buildServer.configs.kotlin.AbsoluteId
import jetbrains.buildServer.configs.kotlin.BuildType

class pullRequest(displayName: String, environment: String, vcsRootId : String) {
    val displayName = displayName
    val environment = environment
    val vcsRootId = vcsRootId

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
                var packageName = "\"%SERVICES%\""

                ConfigureGoEnv()
                DownloadTerraformBinary()
                RunAcceptanceTestsForPullRequest(packageName)
            }

            failureConditions {
                errorMessage = true
            }

            features {
                Golang()
            }

            params {
                TerraformAcceptanceTestParameters(defaultParallelism, "TestAcc", defaultTimeout)
                TerraformAcceptanceTestsFlag()
                TerraformShouldPanicForSchemaErrors()
                TerraformCoreBinaryTesting()
                ReadOnlySettings()

                text("SERVICES", "portal")
            }
        }
    }

    fun uniqueID(provider : String) : String {
        return "%s_PR_%s".format(provider.toUpperCase(), environment.toUpperCase())
    }
}
