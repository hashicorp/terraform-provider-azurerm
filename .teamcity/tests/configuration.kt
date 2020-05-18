package tests

import AzureRM
import org.junit.Assert.assertTrue
import org.junit.Test
import useTeamCityGoTest

class ConfigurationTests {
    @Test
    fun buildShouldFailOnError() {
        val project = AzureRM("public", TestConfiguration())
        project.buildTypes.forEach { bt ->
            assertTrue("Build '${bt.id}' should fail on errors!", bt.failureConditions.errorMessage)
        }
    }

    @Test
    fun buildShouldHaveGoTestFeature() {
        val project = AzureRM("public", TestConfiguration())
        project.buildTypes.forEach{ bt ->
            var exists = false
            bt.features.items.forEach { f ->
                if (f.type == "golang") {
                    exists = true
                }
            }

            if (useTeamCityGoTest) {
                assertTrue("Build %s doesn't have Go Test Json enabled".format(bt.name), exists)
            }
        }
    }

    @Test
    fun buildShouldHaveTrigger() {
        val project = AzureRM("public", TestConfiguration())
        var exists = false
        project.buildTypes.forEach{ bt ->
            bt.triggers.items.forEach { t ->
                if (t.type == "schedulingTrigger") {
                    exists = true
                }
            }
        }
        assertTrue("The Build Configuration should have a Trigger", exists)
    }
}
