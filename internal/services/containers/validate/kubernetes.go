// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func KubernetesAdminUserName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[A-Za-z][-A-Za-z\d_]*$`), "must begin with a letter, contain only letters, numbers, underscores and hyphens")(i, k)
}

func KubernetesAgentPoolName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-z]{1}[a-z\d]{0,11}$`), "must begin with a lowercase letter, contain only lowercase letters and numbers and be between 1 and 12 characters in length")(i, k)
}

func KubernetesClusterName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]$|^[a-zA-Z0-9][-_a-zA-Z0-9]{0,61}[a-zA-Z0-9]$`), "name must start and end with a letter or number, and can only contain letters, numbers, hyphens, and underscores, and be between 1 and 63 characters in length")(i, k)
}

func KubernetesDNSPrefix(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(
		regexp.MustCompile(`^[a-zA-Z\d]$|^[a-zA-Z\d][-a-zA-Z\d]{0,52}[a-zA-Z\d]$`),
		"must begin and end with a letter or number, contain only letters, numbers, and hyphens and be between 1 and 54 characters in length",
	)(i, k)
}

func KubernetesGitRepositoryUrl() pluginsdk.SchemaValidateFunc {
	return validation.StringStartsWithOneOf("http://", "https://", "git@", "ssh://")
}
