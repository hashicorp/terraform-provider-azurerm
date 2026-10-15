/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

package tests

import CheckScheduleConstraints
import PostTestResultsToGitHubPullRequest
import jetbrains.buildServer.configs.kotlin.BuildStep
import jetbrains.buildServer.configs.kotlin.BuildSteps
import jetbrains.buildServer.configs.kotlin.buildSteps.ScriptBuildStep
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class BuildComponentsTests {
    @Test
    fun postTestResultsToGitHubPullRequestShouldRunOnFailure() {
        val steps = BuildSteps().apply {
            PostTestResultsToGitHubPullRequest()
        }

        assertEquals(BuildStep.ExecutionMode.RUN_ON_FAILURE, steps.items.single().executionMode)
    }

    @Test
    fun checkScheduleConstraintsShouldHandleDisableTriggers() {
        val steps = BuildSteps().apply {
            CheckScheduleConstraints()
        }
        val scriptStep = steps.items.single() as ScriptBuildStep
        val content = scriptStep.scriptContent ?: ""
        assertTrue(
            "CheckScheduleConstraints script should check DISABLE_TRIGGERS",
            content.contains("%DISABLE_TRIGGERS%")
        )
        assertTrue(
            "CheckScheduleConstraints script should set SCHEDULE_MATCHES to false when triggers are disabled",
            content.contains("SCHEDULE_MATCHES' value='false'")
        )
    }
}
