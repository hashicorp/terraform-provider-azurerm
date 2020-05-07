import jetbrains.buildServer.configs.kotlin.v2019_2.*
import jetbrains.buildServer.configs.kotlin.v2019_2.buildFeatures.GolangFeature
import jetbrains.buildServer.configs.kotlin.v2019_2.buildSteps.ScriptBuildStep
import jetbrains.buildServer.configs.kotlin.v2019_2.triggers.schedule

// NOTE: in time this could be pulled out into a separate Kotlin package

// unfortunately TeamCity's Go Test Json parser appears to be broken
// as such for the moment let's use the old method with a feature-flag
const val useTeamCityGoTest = true

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

fun BuildSteps.RunAcceptanceTests(providerName : String, packageName: String) {
    if (useTeamCityGoTest) {
        var servicePath = "./%s/internal/services/%s/...".format(providerName, packageName)
        step(ScriptBuildStep {
            name = "Run Tests"
            scriptContent = "go test -v \"$servicePath\" -timeout=\"%TIMEOUT%h\" -test.parallel=\"%PARALLELISM%\" -run=\"%TEST_PREFIX%\" -json"
        })
    } else {
        step(ScriptBuildStep {
            name = "Install tombuildsstuff/teamcity-go-test-json"
            scriptContent = "wget https://github.com/tombuildsstuff/teamcity-go-test-json/releases/download/v0.2.0/teamcity-go-test-json_linux_amd64 && chmod +x teamcity-go-test-json_linux_amd64"
        })

        var servicePath = "./%s/internal/services/%s/...".format(providerName, packageName)
        step(ScriptBuildStep {
            name = "Run Tests"
            scriptContent = "GOFLAGS=\"-mod=vendor\" ./teamcity-go-test-json_linux_amd64 -scope \"$servicePath\" -prefix \"%TEST_PREFIX%\" -count=1 -parallelism=%PARALLELISM% -timeout %TIMEOUT%"
        })
    }
}

fun ParametrizedWithType.TerraformAcceptanceTestParameters(parallelism : Int, prefix : String, timeout: String) {
    text("PARALLELISM", "%d".format(parallelism))
    text("TEST_PREFIX", prefix)
    text("TIMEOUT", timeout)
}

fun ParametrizedWithType.ReadOnlySettings() {
    hiddenVariable("teamcity.ui.settings.readOnly", "true", "Requires build configurations be edited via Kotlin")
}

fun ParametrizedWithType.TerraformAcceptanceTestsFlag() {
    hiddenVariable("env.TF_ACC", "1", "Set to a value to run the Acceptance Tests")
}

fun ParametrizedWithType.TerraformShouldPanicForSchemaErrors() {
    hiddenVariable("env.TF_SCHEMA_PANIC_ON_ERROR", "1", "Panic if unknown/unmatched fields are set into the state")
}

fun ParametrizedWithType.hiddenVariable(name: String, value: String, description: String) {
    text(name, value, "", description, ParameterDisplay.HIDDEN)
}

fun ParametrizedWithType.hiddenPasswordVariable(name: String, value: String, description: String) {
    password(name, value, "", description, ParameterDisplay.HIDDEN)
}

fun Triggers.RunNightly(nightlyTestsEnabled: Boolean, startHour: Int) {
    schedule{
        enabled = nightlyTestsEnabled
        branchFilter = "+:refs/heads/master"

        schedulingPolicy = daily {
            hour = startHour
            timezone = "SERVER"
        }
    }
}
