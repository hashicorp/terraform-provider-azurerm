package acceptance

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

func Environment() (*azure.Environment, error) {
	envName := EnvironmentName()
	metadataURL := os.Getenv("ARM_METADATA_URL")
	return authentication.AzureEnvironmentByNameFromEndpoint(context.TODO(), metadataURL, envName)
}

func GetAuthConfig(t *testing.T) *authentication.Config {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Integration test skipped unless env '%s' set", resource.TestEnvVar))
		return nil
	}

	environment := EnvironmentName()

	builder := authentication.Builder{
		SubscriptionID: os.Getenv("ARM_SUBSCRIPTION_ID"),
		ClientID:       os.Getenv("ARM_CLIENT_ID"),
		TenantID:       os.Getenv("ARM_TENANT_ID"),
		ClientSecret:   os.Getenv("ARM_CLIENT_SECRET"),
		Environment:    environment,
		MetadataHost:   os.Getenv("ARM_METADATA_HOST"),

		// we intentionally only support Client Secret auth for tests (since those variables are used all over)
		SupportsClientSecretAuth: true,
	}
	config, err := builder.Build()
	if err != nil {
		t.Fatalf("Error building ARM Client: %+v", err)
		return nil
	}

	return config
}

func RequiresImportError(resourceName string) *regexp.Regexp {
	message := "to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information."
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}
