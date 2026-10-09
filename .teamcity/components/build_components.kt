import jetbrains.buildServer.configs.kotlin.*
import java.io.File
import jetbrains.buildServer.configs.kotlin.buildFeatures.GolangFeature
import jetbrains.buildServer.configs.kotlin.buildSteps.ScriptBuildStep
import jetbrains.buildServer.configs.kotlin.triggers.schedule

// NOTE: in time this could be pulled out into a separate Kotlin package

// TeamCity's own Go support names each test after its package and doesn't group a test's output in
// the build log, so tests are reported by `internal/tools/teamcity-test-reporter` instead (see
// RunAcceptanceTests) - enabling this as well would report every test twice.
const val useTeamCityGoTest = false

fun BuildFeatures.Golang() {
    if (useTeamCityGoTest) {
        feature(GolangFeature {
            testFormat = "json"
        })
    }
}

// Ensure that daysOfWeek constraints in the overrides are honoured.
fun BuildSteps.CheckScheduleConstraints() {
    step(ScriptBuildStep {
        name = "Check Schedule Constraints"
        scriptContent = """
            #!/bin/bash
            # TeamCity days: 1=Sun, 2=Mon, 3=Tue, 4=Wed, 5=Thu, 6=Fri, 7=Sat
            DAY_NUM=${'$'}(( ${'$'}(date +%w) + 1 ))

            # manual/gh trigger should bypass the daysOfWeek overrides.
            if [[ "%env.IS_NIGHTLY_RUN%" == "true" ]]; then
                if [[ ! ",%DAYS_OF_WEEK%," =~ ",${'$'}{DAY_NUM}," && "%DAYS_OF_WEEK%" != "*" ]]; then
                    echo "Today is day ${'$'}{DAY_NUM}. This job is constrained to run only on days: %DAYS_OF_WEEK%."
                    echo "Skipping test execution for this service."
                    # Tell TeamCity to skip subsequent steps using a service message
                    echo "##teamcity[setParameter name='env.SCHEDULE_MATCHES' value='false']"
                else
                    echo "Schedule matches. Proceeding with tests."
                fi
            else
                echo "This is not a scheduled nightly run (likely triggered manually or via PR chat-ops)."
                echo "Bypassing schedule constraints."
            fi
        """.trimIndent()
    })
}

fun BuildSteps.SetBuildStartTime() {
    step(ScriptBuildStep {
        name = "Set Build Start Time"
        scriptContent = File("scripts/set_build_start_time.sh").readText()
        conditions {
            equals("env.SCHEDULE_MATCHES", "true")
        }
    })
}

fun BuildSteps.ConfigureGoEnv() {
    step(ScriptBuildStep {
        name = "Configure Go Version"
        scriptContent = "goenv install -s \$(goenv local) && goenv rehash"
        conditions {
            equals("env.SCHEDULE_MATCHES", "true")
        }
    })
}

// Downloads Terraform Core to the agent only when that version isn't already there.
fun BuildSteps.DownloadTerraformBinary() {
    step(ScriptBuildStep {
        name = "Download Terraform Core v%env.TERRAFORM_CORE_VERSION%.."
        scriptContent = File("scripts/download_terraform.sh").readText()
        conditions {
            equals("env.SCHEDULE_MATCHES", "true")
        }
    })
}

// Fetches the providers the tests will ask for into the agent's shared provider directory.
// Requires TerraformProviderMirror() in the build's params.
fun BuildSteps.DownloadTerraformProviders(packageName: String) {
    step(ScriptBuildStep {
        name = "Download Terraform Providers"
        // not every build defines SERVICE_PATH as a parameter, so the package under test is filled in here
        scriptContent = File("scripts/download_terraform_providers.sh").readText().replace("%SERVICE_PATH%", servicePath(packageName))
        conditions {
            equals("env.SCHEDULE_MATCHES", "true")
        }
    })
}

fun servicePath(packageName: String) : String {
    return "./internal/services/%s".format(packageName)
}

// Says what state the agent's Go cache is in and drops what nothing has used lately - see GoCache().
fun BuildSteps.PrepareGoCache(packageName: String) {
    step(ScriptBuildStep {
        name = "Prepare Go Cache"
        scriptContent = File("scripts/go_cache.sh").readText().replace("%SERVICE_PATH%", servicePath(packageName))
        conditions {
            equals("env.SCHEDULE_MATCHES", "true")
        }
    })
}

fun BuildSteps.RunAcceptanceTests(packageName: String) {
    step(ScriptBuildStep {
        name = "Run Tests"
        scriptContent = File("scripts/run_tests.sh").readText().replace("%SERVICE_PATH%", servicePath(packageName))
        conditions {
            equals("env.SCHEDULE_MATCHES", "true")
        }
    })
}

fun BuildSteps.RunAcceptanceTestsForPullRequest(packageName: String) {
    step(ScriptBuildStep {
        name = "Run Tests"
        scriptContent = File("scripts/run_tests.sh").readText().replace("%SERVICE_PATH%", servicePath(packageName))
        conditions {
            equals("env.SCHEDULE_MATCHES", "true")
        }
    })
}

fun BuildSteps.PostTestResultsToGitHubPullRequest() {
    step(ScriptBuildStep {
        name = "Post Test Results to GitHub Pull Request"
        scriptContent = File("scripts/post_github_comment.sh").readText()
        workingDir = "%SERVICE_PATH%"
        conditions {
            equals("env.SCHEDULE_MATCHES", "true")
        }
        executionMode = BuildStep.ExecutionMode.RUN_ON_FAILURE
    })
}

fun ParametrizedWithType.TerraformAcceptanceTestParameters(parallelism : Int, prefix : String, timeout: Int) {
    text("PARALLELISM", "%d".format(parallelism))
    text("TEST_PREFIX", prefix)
    text("TIMEOUT", "%d".format(timeout))
    text("POST_GITHUB_COMMENT", "false")
    text("TRACKING_ID", "0", "Tracking ID for comment management (typically PR commit SHA)")
}

fun ParametrizedWithType.ReadOnlySettings() {
    hiddenVariable("teamcity.ui.settings.readOnly", "true", "Requires build configurations be edited via Kotlin")
}

fun ParametrizedWithType.TerraformAcceptanceTestsFlag() {
    hiddenVariable("env.TF_ACC", "1", "Set to a value to run the Acceptance Tests")
}

// Where Terraform Core and the providers are kept between builds. This can't live in the agent's work
// directory: once a build finishes TeamCity deletes everything in there which isn't a checkout directory.
// The agent's persistent cache directory is left alone until the agent runs short of disk space.
const val terraformCacheDir = "%system.agent.persistent.cache%/terraform-cache"

fun ParametrizedWithType.TerraformCoreBinaryTesting() {
    text("env.TERRAFORM_CORE_VERSION", defaultTerraformCoreVersion, "The version of Terraform Core which should be used for testing")
    hiddenVariable("env.TF_ACC_TERRAFORM_PATH", "$terraformCacheDir/core/%env.TERRAFORM_CORE_VERSION%/terraform", "The path where the Terraform Binary is located - shared by every build on the agent")
}

fun ParametrizedWithType.TerraformProviderMirror() {
    hiddenVariable("env.TF_ACC_TERRAFORM_PROVIDER_CACHE_MIRROR_PATH", "$terraformCacheDir/providers", "The directory of provider binaries shared by every build on the agent, which tests link to rather than downloading their own")
    hiddenVariable("env.TF_ACC_TERRAFORM_PROVIDER_CACHE_CONFIG_FILE", "$terraformCacheDir/providers.tfrc", "The Terraform CLI config which has tests install providers from that directory when they are present")
    // TF_CLI_CONFIG_FILE is the name Terraform itself reads, so it has to be set for the tests to pick the config up
    hiddenVariable("env.TF_CLI_CONFIG_FILE", "%env.TF_ACC_TERRAFORM_PROVIDER_CACHE_CONFIG_FILE%", "Points Terraform at the provider cache CLI config")
}

fun ParametrizedWithType.TerraformShouldPanicForSchemaErrors() {
    hiddenVariable("env.TF_SCHEMA_PANIC_ON_ERROR", "1", "Panic if unknown/unmatched fields are set into the state")
}

fun ParametrizedWithType.WorkingDirectory(packageName: String) {
    text("SERVICE_PATH", servicePath(packageName), "", "The path at which to run - automatically updated", ParameterDisplay.HIDDEN)
}

fun ParametrizedWithType.BuildStartTime() {
    text("env.BUILD_START_TIME", "1777662664", "The time at which the build started")
}

// Each agent keeps its own Go caches between builds, so only the first build on an agent compiles everything.
// They can't live in the agent's work directory: once a build finishes TeamCity deletes everything in there
// which isn't a checkout directory. Nothing but Go and run_tests.sh touches the persistent cache directory.
fun ParametrizedWithType.GoCache() {
    text("env.GOMODCACHE", "%system.agent.persistent.cache%/go-cache/mod", "The location of the Go Module Cache")
    text("env.GOCACHE", "%system.agent.persistent.cache%/go-cache/build", "The location of the Go Cache")
}

fun ParametrizedWithType.hiddenVariable(name: String, value: String, description: String) {
    text(name, value, "", description, ParameterDisplay.HIDDEN)
}

fun ParametrizedWithType.hiddenPasswordVariable(name: String, value: String, description: String) {
    password(name, value, "", description, ParameterDisplay.HIDDEN)
}

fun Triggers.RunNightly(nightlyTestsEnabled: Boolean, startHour: Int, daysOfWeek: String, daysOfMonth: String, disableTriggers: Boolean = false) {
    if (!enableTestTriggersGlobally) {
        return
    }

    if (disableTriggers) {
        return
    }

    schedule{
        enabled = nightlyTestsEnabled
        branchFilter = "+:refs/heads/main"

        schedulingPolicy = cron {
            hours = startHour.toString()
            timezone = "SERVER"

            dayOfWeek = daysOfWeek
            dayOfMonth = daysOfMonth
        }
        withPendingChangesOnly = false
    }
}
