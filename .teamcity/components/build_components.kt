import jetbrains.buildServer.configs.kotlin.*
import jetbrains.buildServer.configs.kotlin.buildFeatures.GolangFeature
import jetbrains.buildServer.configs.kotlin.buildSteps.ScriptBuildStep
import jetbrains.buildServer.configs.kotlin.triggers.schedule

// NOTE: in time this could be pulled out into a separate Kotlin package

// The native Go test runner (which TeamCity shells out to) will fail
// the entire test suite when a single test panics, which isn't ideal.
//
// Until that changes, we'll continue to use `teamcity-go-test` to run
// each test individually
const val useTeamCityGoTest = false

fun BuildFeatures.Golang() {
    if (useTeamCityGoTest) {
        feature(GolangFeature {
            testFormat = "json"
        })
    }
}

fun BuildSteps.ConfigureGoEnv() {
    step(ScriptBuildStep {
        name = "Configure Go Version"
        scriptContent = "goenv install -s \$(goenv local) && goenv rehash"
    })
}

fun BuildSteps.DownloadTerraformBinary() {
    // https://releases.hashicorp.com/terraform/0.12.28/terraform_0.12.28_linux_amd64.zip
    var terraformUrl = "https://releases.hashicorp.com/terraform/%env.TERRAFORM_CORE_VERSION%/terraform_%env.TERRAFORM_CORE_VERSION%_linux_amd64.zip"
    step(ScriptBuildStep {
        name = "Download Terraform Core v%env.TERRAFORM_CORE_VERSION%.."
        scriptContent = "mkdir -p tools && wget -O tf.zip %s && unzip tf.zip && mv terraform tools/".format(terraformUrl)
    })
}

fun servicePath(packageName: String) : String {
    return "./internal/services/%s".format(packageName)
}

fun BuildSteps.RunAcceptanceTests(packageName: String) {
    var packagePath = servicePath(packageName)
    var withTestsDirectoryPath = "##teamcity[setParameter name='SERVICE_PATH' value='%s/tests']".format(packagePath)

    // some packages use a ./tests folder, others don't - conditionally append that if needed
    step(ScriptBuildStep {
        name          = "Determine Working Directory for this Package"
        scriptContent = "if [ -d \"%s/tests\" ]; then echo \"%s\"; fi".format(packagePath, withTestsDirectoryPath)
    })

    if (useTeamCityGoTest) {
        step(ScriptBuildStep {
            name = "Run Tests"
            scriptContent = "go test -v \"%SERVICE_PATH%\" -timeout=\"%TIMEOUT%h\" -test.parallel=\"%PARALLELISM%\" -run=\"%TEST_PREFIX%\" -json"
        })
    } else {
        step(ScriptBuildStep {
            name = "Compile Test Binary"
            scriptContent = "go test -c -o test-binary"
            workingDir = "%SERVICE_PATH%"
        })

        step(ScriptBuildStep {
            // ./test-binary -test.list=TestAccAzureRMResourceGroup_ | teamcity-go-test -test ./test-binary -timeout 1s
            name = "Run via jen20/teamcity-go-test"
            scriptContent = "./test-binary -test.list=\"%TEST_PREFIX%\" | teamcity-go-test -test ./test-binary -parallelism \"%PARALLELISM%\" -timeout \"%TIMEOUT%h\""
            workingDir = "%SERVICE_PATH%"
        })
    }
}

fun BuildSteps.RunAcceptanceTestsForPullRequest(packageName: String) {
    var servicePath = "./internal/services/%s/...".format(packageName)
    if (useTeamCityGoTest) {
        step(ScriptBuildStep {
            name = "Run Tests"
            scriptContent = "go test -v \"$servicePath\" -timeout=\"%TIMEOUT%h\" -test.parallel=\"%PARALLELISM%\" -run=\"%TEST_PREFIX%\" -json"
        })
    } else {
        // Building a binary with teamcity-go-test doesn't work for multiple packages, so fallback to this
        step(ScriptBuildStep {
            name = "Install tombuildsstuff/teamcity-go-test-json"
            scriptContent = "wget https://github.com/tombuildsstuff/teamcity-go-test-json/releases/download/v0.2.0/teamcity-go-test-json_linux_amd64 && chmod +x teamcity-go-test-json_linux_amd64"
        })

        step(ScriptBuildStep {
            name = "Run Tests"
            scriptContent = "GOFLAGS=\"-mod=vendor\" ./teamcity-go-test-json_linux_amd64 -scope \"$servicePath\" -prefix \"%TEST_PREFIX%\" -count=1 -parallelism=%PARALLELISM% -timeout %TIMEOUT%"
        })
    }
}

fun ParametrizedWithType.TerraformAcceptanceTestParameters(parallelism : Int, prefix : String, timeout: Int) {
    text("PARALLELISM", "%d".format(parallelism))
    text("TEST_PREFIX", prefix)
    text("TIMEOUT", "%d".format(timeout))
}

fun ParametrizedWithType.ReadOnlySettings() {
    hiddenVariable("teamcity.ui.settings.readOnly", "true", "Requires build configurations be edited via Kotlin")
}

fun ParametrizedWithType.TerraformAcceptanceTestsFlag() {
    hiddenVariable("env.TF_ACC", "1", "Set to a value to run the Acceptance Tests")
}

fun ParametrizedWithType.TerraformCoreBinaryTesting() {
    text("env.TERRAFORM_CORE_VERSION", defaultTerraformCoreVersion, "The version of Terraform Core which should be used for testing")
    hiddenVariable("env.TF_ACC_TERRAFORM_PATH", "%system.teamcity.build.checkoutDir%/tools/terraform", "The path where the Terraform Binary is located")
}

fun ParametrizedWithType.TerraformShouldPanicForSchemaErrors() {
    hiddenVariable("env.TF_SCHEMA_PANIC_ON_ERROR", "1", "Panic if unknown/unmatched fields are set into the state")
}

fun ParametrizedWithType.WorkingDirectory(packageName: String) {
    text("SERVICE_PATH", servicePath(packageName), "", "The path at which to run - automatically updated", ParameterDisplay.HIDDEN)
}

fun ParametrizedWithType.hiddenVariable(name: String, value: String, description: String) {
    text(name, value, "", description, ParameterDisplay.HIDDEN)
}

fun ParametrizedWithType.hiddenPasswordVariable(name: String, value: String, description: String) {
    password(name, value, "", description, ParameterDisplay.HIDDEN)
}

fun Triggers.RunNightly(nightlyTestsEnabled: Boolean, startHour: Int, daysOfWeek: String, daysOfMonth: String, disableTriggers: Boolean = false) {
    // @tombuildsstuff: this temporary flag enables/disables all triggers, allowing a migration between CI servers
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
    }
}
