package validate

import (
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func KubernetesAdminUserName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[A-Za-z][-A-Za-z0-9_]*$`),
		"AdminUserName must start with alphabet and/or continue with alphanumeric characters, underscores, hyphens.")
}

func KubernetesAgentPoolName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-z]{1}[a-z0-9]{0,11}$"),
		"Agent Pool names must start with a lowercase letter, have max length of 12, and only have characters a-z0-9.",
	)
}

func KubernetesDNSPrefix() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{0,43}[a-zA-Z0-9]$"),
		"The DNS Prefix must contain between 2 and 45 characters. The name can contain only letters, numbers, and hyphens. The name must start with a letter and must end with an alphanumeric character.",
	)
}
