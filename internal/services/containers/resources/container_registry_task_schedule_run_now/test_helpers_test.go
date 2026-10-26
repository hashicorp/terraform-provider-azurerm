// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_task_schedule_run_now_test

import (
	"os"
	"testing"
)

type githubRepo struct {
	url   string
	token string
}

func preCheckGithubRepo(t *testing.T) {
	// - ARM_TEST_ACR_TASK_GITHUB_REPO_URL represents the user forked repo from: https://github.com/Azure-Samples/acr-build-helloworld-node
	// - ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN represents the github personal token with the appropriate permissions per: https://docs.microsoft.com/azure/container-registry/container-registry-tutorial-build-task#create-a-github-personal-access-token
	// Checkout https://docs.microsoft.com/azure/container-registry/container-registry-tutorial-build-task for details.
	variables := []string{
		"ARM_TEST_ACR_TASK_GITHUB_REPO_URL",
		"ARM_TEST_ACR_TASK_GITHUB_USER_TOKEN",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Skipf("`%s` must be set for acceptance tests!", variable)
		}
	}
}
