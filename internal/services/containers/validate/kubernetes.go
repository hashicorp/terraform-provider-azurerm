// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func KubernetesAdminUserName(i interface{}, k string) (warnings []string, errors []error) {
	adminUserName, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	re := regexp.MustCompile(`^[A-Za-z][-A-Za-z\d_]*$`)
	if re != nil && !re.MatchString(adminUserName) {
		errors = append(errors, fmt.Errorf("the %q must begin with a letter, contain only letters, numbers, underscores and hyphens, got %q", k, adminUserName))
	}

	return warnings, errors
}

func KubernetesAgentPoolName(i interface{}, k string) (warnings []string, errors []error) {
	agentPoolName, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	re := regexp.MustCompile(`^[a-z]{1}[a-z\d]{0,11}$`)
	if re != nil && !re.MatchString(agentPoolName) {
		errors = append(errors, fmt.Errorf("the %q must begin with a lowercase letter, contain only lowercase letters and numbers and be between 1 and 12 characters in length, got %q", k, agentPoolName))
	}

	return warnings, errors
}

func KubernetesDNSPrefix(i interface{}, k string) (warnings []string, errors []error) {
	dnsPrefix, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	errMsg := fmt.Sprintf("the %q must begin and end with a letter or number, contain only letters, numbers, and hyphens and be between 1 and 54 characters in length, got", k)

	if len(dnsPrefix) < 2 {
		re := regexp.MustCompile(`^[a-zA-Z\d]`)
		if re != nil && !re.MatchString(dnsPrefix) {
			errors = append(errors, fmt.Errorf("%s %q", errMsg, dnsPrefix))
		}
	} else {
		re := regexp.MustCompile(`^[a-zA-Z\d][-a-zA-Z\d]{0,52}[a-zA-Z\d]$`)
		if re != nil && !re.MatchString(dnsPrefix) {
			errors = append(errors, fmt.Errorf("%s %q", errMsg, dnsPrefix))
		}
	}

	return warnings, errors
}

func KubernetesNodeTaint(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	parts := strings.Split(v, ":")
	if len(parts) != 2 {
		errors = append(errors, fmt.Errorf("%q must be in the format `key[=value]:effect`, got %q", k, v))
		return
	}
	if parts[0] == "" {
		errors = append(errors, fmt.Errorf("the key portion of %q cannot be empty, got %q", k, v))
	}
	switch parts[1] {
	case "NoSchedule", "PreferNoSchedule", "NoExecute":
	default:
		errors = append(errors, fmt.Errorf("the effect of %q must be one of `NoSchedule`, `PreferNoSchedule` or `NoExecute`, got %q", k, parts[1]))
	}

	return
}

func KubernetesGitRepositoryUrl() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
		}

		if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") || strings.HasPrefix(v, "git@") || strings.HasPrefix(v, "ssh://") {
			return nil, nil
		}

		return nil, []error{fmt.Errorf("expected %q to start with `http://`, `https://`, `git@` or `ssh://`", k)}
	}
}
