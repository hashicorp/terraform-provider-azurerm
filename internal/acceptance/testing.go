// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"fmt"
	"os"
	"regexp"
	"testing"
)

func PreCheck(t *testing.T) {
	variables := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_SUBSCRIPTION_ID",
		"ARM_TENANT_ID",
		"ARM_TEST_LOCATION",
		"ARM_TEST_LOCATION_ALT",
		"ARM_TEST_LOCATION_ALT2",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

func EnvironmentName() string {
	envName, exists := os.LookupEnv("ARM_ENVIRONMENT")
	if !exists {
		envName = "public"
	}

	return envName
}

func RequiresImportError(resourceName string) *regexp.Regexp {
	message := `to\s+be\s+managed\s+via\s+Terraform\s+this\s+resource\s+needs\s+to\s+be\s+imported\s+into\s+the\s+State\.\s+Please\s+see\s+the\s+resource\s+documentation\s+for\s+%q\s+for\s+more\s+information`
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}

func RequiresImportAssociationError(resourceName string) *regexp.Regexp {
	message := `to\s+be\s+managed\s+via\s+Terraform\s+this\s+association\s+needs\s+to\s+be\s+imported\s+into\s+the\s+State\.\s+Please\s+see\s+the\s+resource\s+documentation\s+for\s+%q\s+for\s+more\s+information`
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}
