/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

package tests

import PostTestResultsToGitHubPullRequest
import jetbrains.buildServer.configs.kotlin.BuildStep
import jetbrains.buildServer.configs.kotlin.BuildSteps
import org.junit.Assert.assertEquals
import org.junit.Test

class BuildComponentsTests {
    @Test
    fun postTestResultsToGitHubPullRequestShouldRunOnFailure() {
        val steps = BuildSteps().apply {
            PostTestResultsToGitHubPullRequest()
        }

        assertEquals(BuildStep.ExecutionMode.RUN_ON_FAILURE, steps.items.single().executionMode)
    }
}
