/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

package tests

import AzureRM
import AzureRMBetaVersion
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

    @Test
    fun buildShouldFailOnErrorBetaVersion() {
        val project = AzureRMBetaVersion("public", TestConfiguration())
        project.buildTypes.forEach { bt ->
            assertTrue("Build '${bt.id}' should fail on errors!", bt.failureConditions.errorMessage)
        }
    }

    @Test
    fun buildShouldHaveGoTestFeatureBetaVersion() {
        val project = AzureRMBetaVersion("public", TestConfiguration())
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
    fun buildShouldHaveTriggerBetaVersion() {
        val project = AzureRMBetaVersion("public", TestConfiguration())
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

    @Test
    fun betaVersionBuildsShouldSetFivePointZeroFlag() {
        val config = TestConfiguration()
        val project = AzureRMBetaVersion("public", config)
        project.buildTypes.forEach { bt ->
            val betaFlag = bt.params.findRawParam(config.betaVersionEnvVar)?.value
            assertTrue(
                "Build '${bt.id}' should have ${config.betaVersionEnvVar} parameter set, but it was not found",
                betaFlag != null
            )
            assertTrue(
                "Build '${bt.id}' should set ${config.betaVersionEnvVar} to true, but was '$betaFlag'",
                betaFlag == "true"
            )
        }
    }

    @Test
    fun standardBuildsShouldNotSetFivePointZeroFlag() {
        val config = TestConfiguration()
        val project = AzureRM("public", config)
        project.buildTypes.forEach { bt ->
            // Skip cache and PR builds as they don't need the beta version flag
            if (bt.id.toString().contains("AZURERM_CACHE_PUBLIC") || bt.id.toString().contains("AZURERM_PR_PUBLIC")) {
                return@forEach
            }
            val betaFlag = bt.params.findRawParam(config.betaVersionEnvVar)?.value
            assertTrue(
                "Build '${bt.id}' should have ${config.betaVersionEnvVar} parameter set, but it was not found",
                betaFlag != null
            )
            assertTrue(
                "Build '${bt.id}' should set ${config.betaVersionEnvVar} to false, but was '$betaFlag'",
                betaFlag == "false"
            )
        }
    }
}
