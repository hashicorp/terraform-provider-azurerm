// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// Your server name can contain only lowercase letters, numbers, and '-', but can't start or end with '-' or have more than 63 characters.
func ValidateMsSqlServerName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-z]([-0-9a-z]{0,61}[0-9a-z])?$`), "can contain only lowercase letters, numbers, and '-', but can't start or end with '-' or have more than 63 characters")(i, k)
}

// Your database name can't end with '.' or ' ', can't contain '<,>,*,%,&,:,\,/,?' or control characters, and can't have more than 128 characters.
func ValidateMsSqlDatabaseName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^<>*%&:\\\/?]{0,127}[^\s.<>*%&:\\\/?]$`), `can't end with '.' or ' ', can't contain '<,>,*,%,&,:,\,/,?' or control characters, and can't have more than 128 characters`)(i, k)
}

func ValidateMsSqlFailoverGroupName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-z]([-0-9a-z]{0,61}[0-9a-z])?$`), "can contain only lowercase letters, numbers, and '-', but can't start or end with '-'")(i, k)
}

// Following characters and any control characters are not allowed for resource name '%,&,\\\\,?,/'.\"
// The name can not end with characters: '. '
// TODO: unsure about length, was able to deploy one at 120
func ValidateMsSqlElasticPoolName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^&%\\\/?]{0,127}[^\s.&%\\\/?]$`), `can't end with '.' or ' ', can't contain '%,&,\,/,?' or control characters, and can't have more than 128 characters`)(i, k)
}

// Job Agent name must not contain any of ?<>*%&:\/? and must not end with a space or .
func ValidateMsSqlJobAgentName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^?<>*%&:\/?]{0,127}[^?<>*%&:\/?. ]$`), `must not contain any of ?<>*%&:\/?, must not end with a space or a period and can't have more than 128 characters`)(i, k)
}

// ValidateMsSqlDNSAliasName
// Server DNS Alias name cannot be empty or null. It can only be made
//
//	up of lowercase letters 'a'-'z', the numbers 0-9 and the hyphen. The hyphen
//	may not lead or trail in the name.
func ValidateMsSqlDNSAliasName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-z][-0-9a-z]{0,127}[0-9a-z]$`), "Server DNS Alias name cannot be empty or null. It can only be made up of lowercase letters 'a'-'z', the numbers 0-9 and the hyphen. The hyphen may not lead or trail in the name")(i, k)
}
