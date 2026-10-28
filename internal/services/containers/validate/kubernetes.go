// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func KubernetesAdminUserName(i any, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[A-Za-z][-A-Za-z\d_]*$`), "must begin with a letter, contain only letters, numbers, underscores and hyphens")(i, k)
}

func KubernetesAgentPoolName(i any, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-z]{1}[a-z\d]{0,38}$`), "must begin with a lowercase letter, contain only lowercase letters and numbers and be between 1 and 39 characters in length")(i, k)
}

func KubernetesAgentPoolNameWithFeatures(name, key string, maxLength int) error {
	if maxLength < 12 || maxLength > 39 {
		return fmt.Errorf("`features.kubernetes_cluster.node_pool_name_max_length` must be between 12 and 39")
	}
	_, errs := validation.StringMatch(
		regexp.MustCompile(fmt.Sprintf(`^[a-z]{1}[a-z\d]{0,%d}$`, maxLength-1)),
		fmt.Sprintf("must begin with a lowercase letter, contain only lowercase letters and numbers and be between 1 and %d characters in length", maxLength),
	)(name, key)
	return errors.Join(errs...)
}

func KubernetesClusterName(i any, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]$|^[a-zA-Z0-9][-_a-zA-Z0-9]{0,61}[a-zA-Z0-9]$`), "name must start and end with a letter or number, and can only contain letters, numbers, hyphens, and underscores, and be between 1 and 63 characters in length")(i, k)
}

func KubernetesDNSPrefix(i any, k string) (warnings []string, errors []error) {
	return validation.StringMatch(
		regexp.MustCompile(`^[a-zA-Z\d]$|^[a-zA-Z\d][-a-zA-Z\d]{0,52}[a-zA-Z\d]$`),
		"must begin and end with a letter or number, contain only letters, numbers, and hyphens and be between 1 and 54 characters in length",
	)(i, k)
}

func KubernetesGitRepositoryUrl() pluginsdk.SchemaValidateFunc {
	return validation.StringStartsWithOneOf("http://", "https://", "git@", "ssh://")
}
