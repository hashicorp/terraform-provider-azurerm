import jetbrains.buildServer.configs.kotlin.v2019_2.*

class pullRequest(displayName: String, environment: String) {
    val displayName = displayName
    val environment = environment

    fun buildConfiguration(providerName : String) : BuildType {
        return BuildType {
            // TC needs a consistent ID for dynamically generated packages
            id(uniqueID(providerName))

            name = displayName

            vcs {
                root(providerRepository)
                cleanCheckout = true
            }

            steps {
                var packageName = "\"%SERVICES%\""

                ConfigureGoEnv()
                DownloadTerraformBinary()
                RunAcceptanceTestsForPullRequest(providerName, packageName)
            }

            failureConditions {
                errorMessage = true
            }

            features {
                Golang()
            }

            params {
                TerraformAcceptanceTestParameters(defaultParallelism, "TestAcc", "12")
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
