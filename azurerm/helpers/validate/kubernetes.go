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

func KubernetesCidr() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?$`),
		"CIDR must start with IPV4 address and/or slash, number of bits (0-32) as prefix. Example: 127.0.0.1/8.")
}

func KubernetesDNSServiceIP() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`),
		"DNSServiceIP must follow IPV4 address format.")
}

func KubernetesAgentPoolName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-z]{1}[a-z0-9]{0,11}$"),
		"Agent Pool names must start with a lowercase letter, have max length of 12, and only have characters a-z0-9.",
	)
}

func KubernetesDnsPrefix() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{0,43}[a-zA-Z0-9]$"),
		"The DNS name must contain between 3 and 45 characters. The name can contain only letters, numbers, and hyphens. The name must start with a letter and must end with a letter or a number.",
	)
}
