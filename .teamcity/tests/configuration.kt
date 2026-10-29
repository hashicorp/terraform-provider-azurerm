/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

package tests

import AzureRM
import AzureRMBetaVersion
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import serviceTestConfigurationOverrides
import useTeamCityGoTest

class ConfigurationTests {
    @Test
    fun buildShouldFailOnError() {
        val project = AzureRM("public", TestConfiguration())
        project.buildTypes.forEach { bt ->
            if (bt.type == jetbrains.buildServer.configs.kotlin.BuildTypeSettings.Type.COMPOSITE) {
                return@forEach
            }
            assertTrue("Build '${bt.id}' should fail on errors!", bt.failureConditions.errorMessage)
        }
    }

    @Test
    fun buildShouldHaveGoTestFeature() {
        val project = AzureRM("public", TestConfiguration())
        project.buildTypes.forEach{ bt ->
            if (bt.type == jetbrains.buildServer.configs.kotlin.BuildTypeSettings.Type.COMPOSITE || bt.id.toString().contains("CACHE")) {
                return@forEach
            }
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
            if (bt.type == jetbrains.buildServer.configs.kotlin.BuildTypeSettings.Type.COMPOSITE) {
                return@forEach
            }
            assertTrue("Build '${bt.id}' should fail on errors!", bt.failureConditions.errorMessage)
        }
    }

    @Test
    fun buildShouldHaveGoTestFeatureBetaVersion() {
        val project = AzureRMBetaVersion("public", TestConfiguration())
        project.buildTypes.forEach{ bt ->
            if (bt.type == jetbrains.buildServer.configs.kotlin.BuildTypeSettings.Type.COMPOSITE || bt.id.toString().contains("CACHE")) {
                return@forEach
            }
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
    fun betaVersionBuildsShouldSetSixPointZeroFlag() {
        val config = TestConfiguration()
        val project = AzureRMBetaVersion("public", config)
        project.buildTypes.forEach { bt ->
            if (bt.type == jetbrains.buildServer.configs.kotlin.BuildTypeSettings.Type.COMPOSITE) {
                return@forEach
            }
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
    fun standardBuildsShouldNotSetSixPointZeroFlag() {
        val config = TestConfiguration()
        val project = AzureRM("public", config)
        project.buildTypes.forEach { bt ->
            if (bt.type == jetbrains.buildServer.configs.kotlin.BuildTypeSettings.Type.COMPOSITE) {
                return@forEach
            }
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

    @Test
    fun buildChainShouldBeCompositeAndHaveDependencies() {
        val config = TestConfiguration()
        val project = AzureRM("public", config)

        val buildChain = project.buildTypes.find { it.name.contains("All Services Build Chain") }
        assertTrue("Build chain must exist", buildChain != null)
        assertTrue("Build chain must be a COMPOSITE build", buildChain!!.type == jetbrains.buildServer.configs.kotlin.BuildTypeSettings.Type.COMPOSITE)
        assertTrue("Build chain must have snapshot dependencies", buildChain.dependencies.items.isNotEmpty())

        val isNightlyParam = buildChain.params.findRawParam("reverse.dep.*.env.IS_NIGHTLY_RUN")?.value
        assertTrue("Build chain must pass down reverse.dep.*.env.IS_NIGHTLY_RUN parameter", isNightlyParam == "true")
    }

    @Test
    fun betaBuildChainShouldBeCompositeAndHaveDependencies() {
        val config = TestConfiguration()
        val project = AzureRMBetaVersion("public", config)

        val buildChain = project.buildTypes.find { it.name.contains("All Services Build Chain") }
        assertTrue("Beta Build chain must exist", buildChain != null)
        assertTrue("Beta Build chain must be a COMPOSITE build", buildChain!!.type == jetbrains.buildServer.configs.kotlin.BuildTypeSettings.Type.COMPOSITE)
        assertTrue("Beta Build chain must have snapshot dependencies", buildChain.dependencies.items.isNotEmpty())

        val isNightlyParam = buildChain.params.findRawParam("reverse.dep.*.env.IS_NIGHTLY_RUN")?.value
        assertTrue("Beta Build chain must pass down reverse.dep.*.env.IS_NIGHTLY_RUN parameter", isNightlyParam == "true")
    }

    @Test
    fun serviceConfigurationOverridesPropagated() {
        val config = TestConfiguration()
        val project = AzureRM("public", config)

        serviceTestConfigurationOverrides.forEach { (serviceName, testConfig) ->
            val bt = project.buildTypes.find { it.id.toString().endsWith("_SERVICE_PUBLIC_${serviceName.uppercase()}") }
            if (bt != null) {
                assertEquals(
                    "Service '$serviceName' should have PARALLELISM set to ${testConfig.parallelism}",
                    testConfig.parallelism.toString(),
                    bt.params.findRawParam("PARALLELISM")?.value
                )
                assertEquals(
                    "Service '$serviceName' should have TIMEOUT set to ${testConfig.timeout}",
                    testConfig.timeout.toString(),
                    bt.params.findRawParam("TIMEOUT")?.value
                )
                assertEquals(
                    "Service '$serviceName' should have executionTimeoutMin set to ${60 * testConfig.timeout}",
                    60 * testConfig.timeout,
                    bt.failureConditions.executionTimeoutMin
                )
                assertEquals(
                    "Service '$serviceName' should have maxRunningBuilds set to ${testConfig.maxConcurrentBuilds}",
                    testConfig.maxConcurrentBuilds,
                    bt.maxRunningBuilds
                )
                assertEquals(
                    "Service '$serviceName' should have maxRunningBuildsPerBranch set to *:${testConfig.maxConcurrentBuildsPerBranch}",
                    "*:%d".format(testConfig.maxConcurrentBuildsPerBranch),
                    bt.maxRunningBuildsPerBranch
                )
                assertEquals(
                    "Service '$serviceName' should have DAYS_OF_WEEK set to ${testConfig.daysOfWeek}",
                    testConfig.daysOfWeek,
                    bt.params.findRawParam("DAYS_OF_WEEK")?.value
                )
                assertEquals(
                    "Service '$serviceName' should have DISABLE_TRIGGERS set to ${testConfig.disableTriggers}",
                    testConfig.disableTriggers.toString(),
                    bt.params.findRawParam("DISABLE_TRIGGERS")?.value
                )

                if (testConfig.locationOverride.primary != "") {
                    assertEquals(
                        "Service '$serviceName' should have ARM_TEST_LOCATION overridden to ${testConfig.locationOverride.primary}",
                        testConfig.locationOverride.primary,
                        bt.params.findRawParam("env.ARM_TEST_LOCATION")?.value
                    )
                    assertEquals(
                        "Service '$serviceName' should have ARM_TEST_LOCATION_ALT overridden to ${testConfig.locationOverride.secondary}",
                        testConfig.locationOverride.secondary,
                        bt.params.findRawParam("env.ARM_TEST_LOCATION_ALT")?.value
                    )
                    assertEquals(
                        "Service '$serviceName' should have ARM_TEST_LOCATION_ALT2 overridden to ${testConfig.locationOverride.tertiary}",
                        testConfig.locationOverride.tertiary,
                        bt.params.findRawParam("env.ARM_TEST_LOCATION_ALT2")?.value
                    )
                    assertEquals(
                        "Service '$serviceName' should have ARM_PROVIDER_DYNAMIC_TEST overridden to ${testConfig.locationOverride.rotate}",
                        testConfig.locationOverride.rotate.toString(),
                        bt.params.findRawParam("env.ARM_PROVIDER_DYNAMIC_TEST")?.value
                    )
                }

                if (testConfig.useAltSubscription) {
                    assertEquals(
                        "Service '$serviceName' should have ARM_SUBSCRIPTION_ID swapped to alt",
                        config.subscriptionIdAlt,
                        bt.params.findRawParam("env.ARM_SUBSCRIPTION_ID")?.value
                    )
                    assertEquals(
                        "Service '$serviceName' should have ARM_SUBSCRIPTION_ID_ALT swapped to main",
                        config.subscriptionId,
                        bt.params.findRawParam("env.ARM_SUBSCRIPTION_ID_ALT")?.value
                    )
                } else if (testConfig.useDevTestSubscription) {
                    assertEquals(
                        "Service '$serviceName' should have ARM_SUBSCRIPTION_ID swapped to devtest",
                        config.subscriptionIdDevTest,
                        bt.params.findRawParam("env.ARM_SUBSCRIPTION_ID")?.value
                    )
                    assertEquals(
                        "Service '$serviceName' should have ARM_SUBSCRIPTION_ID_ALT swapped to main",
                        config.subscriptionId,
                        bt.params.findRawParam("env.ARM_SUBSCRIPTION_ID_ALT")?.value
                    )
                } else {
                    assertEquals(
                        "Service '$serviceName' should have main ARM_SUBSCRIPTION_ID",
                        config.subscriptionId,
                        bt.params.findRawParam("env.ARM_SUBSCRIPTION_ID")?.value
                    )
                    assertEquals(
                        "Service '$serviceName' should have main ARM_SUBSCRIPTION_ID_ALT",
                        config.subscriptionIdAlt,
                        bt.params.findRawParam("env.ARM_SUBSCRIPTION_ID_ALT")?.value
                    )
                }
            }
        }
    }

    @Test
    fun betaVersionWeeklyTestDayPropagated() {
        val config = TestConfiguration()
        val project = AzureRMBetaVersion("public", config)

        serviceTestConfigurationOverrides.forEach { (serviceName, testConfig) ->
            val bt = project.buildTypes.find { it.id.toString().endsWith("_SERVICE_PUBLIC_${serviceName.uppercase()}") }
            if (bt != null) {
                assertEquals(
                    "Beta service '$serviceName' should have DAYS_OF_WEEK set to ${testConfig.weeklyTestDay}",
                    testConfig.weeklyTestDay,
                    bt.params.findRawParam("DAYS_OF_WEEK")?.value
                )
            }
        }
    }
}
