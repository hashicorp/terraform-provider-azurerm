import jetbrains.buildServer.configs.kotlin.v2019_2.*

class serviceDetails(name: String, displayName: String, environment: String) {
    val packageName = name
    val displayName = displayName
    val environment = environment

    fun buildConfiguration(providerName : String, nightlyTestsEnabled: Boolean, startHour: Int, parallelism: Int, daysOfWeek: String, daysOfMonth: String) : BuildType {
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
                DownloadTerraformBinary()
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
                TerraformCoreBinaryTesting()
                TerraformShouldPanicForSchemaErrors()
                ReadOnlySettings()
                WorkingDirectory(providerName, packageName)
            }

            triggers {
                RunNightly(nightlyTestsEnabled, startHour, daysOfWeek, daysOfMonth)
            }
        }
    }

    fun uniqueID(provider : String) : String {
        return "%s_SERVICE_%s_%s".format(provider.toUpperCase(), environment.toUpperCase(), packageName.toUpperCase())
    }
}
