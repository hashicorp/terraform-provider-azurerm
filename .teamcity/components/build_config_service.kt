import jetbrains.buildServer.configs.kotlin.*

class ServiceDetails(name: String, val displayName: String, val environment: String, val vcsRootId: String) {
    val packageName = name

    fun buildConfiguration(providerName : String, nightlyTestsEnabled: Boolean, startHour: Int, parallelism: Int, daysOfWeek: String, daysOfMonth: String, timeout: Int, disableTriggers: Boolean) : BuildType {
        return BuildType {
            // TC needs a consistent ID for dynamically generated packages
            id(uniqueID(providerName))

            name = "%s - Acceptance Tests".format(displayName)

            vcs {
                root(rootId = AbsoluteId(vcsRootId))
                cleanCheckout = true
            }

            steps {
                configureGoEnv()
                downloadTerraformBinary()
                downloadVCRCassettes(packageName)
                runAcceptanceTests(packageName)
                uploadVCRCassettes(packageName)
                postTestResultsToGitHubPullRequest()
            }

            failureConditions {
                errorMessage = true
                executionTimeoutMin = 60 * timeout
            }

            features {
                golang()
                buildCacheFeature()
            }

            params {
                terraformAcceptanceTestParameters(parallelism, "TestAcc", timeout)
                terraformAcceptanceTestsFlag()
                terraformCoreBinaryTesting()
                terraformShouldPanicForSchemaErrors()
                readOnlySettings()
                workingDirectory(packageName)
                goCache()
            }

            triggers {
                runNightly(nightlyTestsEnabled, startHour, daysOfWeek, daysOfMonth, disableTriggers)
            }
        }
    }

    fun uniqueID(provider : String) : String {
        return "%s_SERVICE_%s_%s".format(provider.uppercase(), environment.uppercase(), packageName.uppercase())
    }
}
