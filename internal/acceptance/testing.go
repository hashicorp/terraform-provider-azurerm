// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

func GetAuthConfig(t *testing.T) *auth.Credentials {
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance test skipped unless env '%s' set", resource.EnvTfAcc)
		return nil
	}

	var (
		ctx = context.TODO()

		env *environments.Environment
		err error

		envName      = EnvironmentName()
		metadataHost = os.Getenv("ARM_METADATA_HOSTNAME")
	)

	if metadataHost != "" {
		if env, err = environments.FromEndpoint(ctx, fmt.Sprintf("https://%s", metadataHost)); err != nil {
			t.Fatalf("building test client: %+v", err)
			return nil
		}
	} else if env, err = environments.FromName(envName); err != nil {
		t.Fatalf("building test client: %+v", err)
		return nil
	}

	return &auth.Credentials{
		Environment: *env,
		ClientID:    os.Getenv("ARM_CLIENT_ID"),
		TenantID:    os.Getenv("ARM_TENANT_ID"),

		ClientCertificatePath:     os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"),
		ClientCertificatePassword: os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"),
		ClientSecret:              os.Getenv("ARM_CLIENT_SECRET"),

		EnableAuthenticatingUsingClientCertificate: true,
		EnableAuthenticatingUsingClientSecret:      true,
		EnableAuthenticatingUsingAzureCLI:          false,
		EnableAuthenticatingUsingManagedIdentity:   false,
		EnableAuthenticationUsingOIDC:              false,
		EnableAuthenticationUsingGitHubOIDC:        false,
	}
}

func RequiresImportError(resourceName string) *regexp.Regexp {
	message := `to\s+be\s+managed\s+via\s+Terraform\s+this\s+resource\s+needs\s+to\s+be\s+imported\s+into\s+the\s+State\.\s+Please\s+see\s+the\s+resource\s+documentation\s+for\s+%q\s+for\s+more\s+information`
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}

func RequiresImportAssociationError(resourceName string) *regexp.Regexp {
	message := `to\s+be\s+managed\s+via\s+Terraform\s+this\s+association\s+needs\s+to\s+be\s+imported\s+into\s+the\s+State\.\s+Please\s+see\s+the\s+resource\s+documentation\s+for\s+%q\s+for\s+more\s+information`
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}
