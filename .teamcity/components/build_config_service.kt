import jetbrains.buildServer.configs.kotlin.*

class serviceDetails(name: String, displayName: String, environment: String, vcsRootId : String) {
    val packageName = name
    val displayName = displayName
    val environment = environment
    val vcsRootId = vcsRootId

    fun buildConfiguration(providerName : String, nightlyTestsEnabled: Boolean, startHour: Int, parallelism: Int, daysOfWeek: String, daysOfMonth: String, timeout: Int, disableTriggers: Boolean, runWithBetaVersion: Boolean) : BuildType {
        return BuildType {
            // TC needs a consistent ID for dynamically generated packages
            id(uniqueID(providerName, runWithBetaVersion))

            name = "%s - Acceptance Tests".format(displayName)

            vcs {
                root(rootId = AbsoluteId(vcsRootId))
                cleanCheckout = true
            }

            steps {
                SetBuildStartTime()
                ConfigureGoEnv()
                DownloadTerraformBinary()
                RunAcceptanceTests(packageName)
                PostTestResultsToGitHubPullRequest()
            }

            failureConditions {
                errorMessage = true
                executionTimeoutMin = 60 * timeout
            }

            features {
                Golang()
                BuildCacheFeature()
            }

            params {
                TerraformAcceptanceTestParameters(parallelism, "TestAcc", timeout)
                TerraformAcceptanceTestsFlag()
                TerraformCoreBinaryTesting()
                TerraformShouldPanicForSchemaErrors()
                ReadOnlySettings()
                WorkingDirectory(packageName)
                GoCache()
                BuildStartTime()
                SetForBetaVersion(runWithBetaVersion)
            }

            triggers {
                RunNightly(nightlyTestsEnabled, startHour, daysOfWeek, daysOfMonth, disableTriggers)
            }
        }
    }

    fun uniqueID(provider : String, runWithBetaVersion: Boolean ) : String {
        if (runWithBetaVersion) {
            return "%s_BETA_VERSION_SERVICE_%s_%s".format(provider.uppercase(), environment.uppercase(), packageName.uppercase())
        }
        return "%s_SERVICE_%s_%s".format(provider.uppercase(), environment.uppercase(), packageName.uppercase())
    }
}
