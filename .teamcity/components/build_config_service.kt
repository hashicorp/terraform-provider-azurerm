import jetbrains.buildServer.configs.kotlin.v2019_2.*

class serviceDetails(name: String, displayName: String, environment: String) {
    val packageName = name
    val displayName = displayName
    val environment = environment

    fun buildConfiguration(providerName : String, nightlyTestsEnabled: Boolean, startHour: Int, parallelism: Int) : BuildType {
        return BuildType {
            // TC needs a consistent ID for dynamically generated packages
            id(uniqueID(providerName))

            name = "%s - Acceptance Tests".format(displayName)

            vcs {
                root(providerRepository)
                cleanCheckout = true
            }

            steps {
                ConfigureGoEnv()
                RunAcceptanceTests(providerName, packageName)
            }

            failureConditions {
                errorMessage = true
            }

            features {
                Golang()
            }

            params {
                TerraformAcceptanceTestParameters(parallelism, "TestAcc", "12")
                TerraformAcceptanceTestsFlag()
                TerraformShouldPanicForSchemaErrors()
                ReadOnlySettings()
            }

            triggers {
                RunNightly(nightlyTestsEnabled, startHour)
            }
        }
    }

    fun uniqueID(provider : String) : String {
        return "%s_SERVICE_%s_%s".format(provider.toUpperCase(), environment.toUpperCase(), packageName.toUpperCase())
    }
}
