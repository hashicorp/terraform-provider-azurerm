import jetbrains.buildServer.configs.kotlin.*

class nextRelease(name: String, displayName: String, environment: String, vcsRootId : String) {
    val packageName = name
    val displayName = displayName
    val environment = environment
    val vcsRootId = vcsRootId

    fun buildConfiguration(providerName : String, nightlyTestsEnabled: Boolean, startHour: Int, parallelism: Int, weeklyTestDay: String, daysOfMonth: String, timeout: Int, disableTriggers: Boolean) : BuildType {
        return BuildType {
            // TC needs a consistent ID for dynamically generated packages
            id(uniqueID(providerName))

            name = "%s - Next Release Acceptance Tests".format(displayName)

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
                SetForNextRelease()
            }

            triggers {
                RunNightly(nightlyTestsEnabled, startHour, weeklyTestDay, daysOfMonth, disableTriggers)
            }
        }
    }

    fun uniqueID(provider : String) : String {
        return "%s_NEXT_RELEASE_SERVICE_%s_%s".format(provider.uppercase(), environment.uppercase(), packageName.uppercase())
    }
}
