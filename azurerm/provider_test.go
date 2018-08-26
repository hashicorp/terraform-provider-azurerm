package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/authentication"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"azurerm": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	variables := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_SUBSCRIPTION_ID",
		"ARM_TENANT_ID",
		"ARM_TEST_LOCATION",
		"ARM_TEST_LOCATION_ALT",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

func testLocation() string {
	return os.Getenv("ARM_TEST_LOCATION")
}

func testAltLocation() string {
	return os.Getenv("ARM_TEST_LOCATION_ALT")
}

func testArmEnvironmentName() string {
	envName, exists := os.LookupEnv("ARM_ENVIRONMENT")
	if !exists {
		envName = "public"
	}

	return envName
}

func testArmEnvironment() (*azure.Environment, error) {
	envName := testArmEnvironmentName()

	// detect cloud from environment
	env, envErr := azure.EnvironmentFromName(envName)
	if envErr != nil {
		// try again with wrapped value to support readable values like german instead of AZUREGERMANCLOUD
		wrapped := fmt.Sprintf("AZURE%sCLOUD", envName)
		var innerErr error
		if env, innerErr = azure.EnvironmentFromName(wrapped); innerErr != nil {
			return nil, envErr
		}
	}

	return &env, nil
}

func testGetAzureConfig(t *testing.T) *authentication.Config {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Integration test skipped unless env '%s' set", resource.TestEnvVar))
		return nil
	}

	environment := testArmEnvironmentName()

	// we deliberately don't use the main config - since we care about
	config := authentication.Config{
		SubscriptionID:           os.Getenv("ARM_SUBSCRIPTION_ID"),
		ClientID:                 os.Getenv("ARM_CLIENT_ID"),
		TenantID:                 os.Getenv("ARM_TENANT_ID"),
		ClientSecret:             os.Getenv("ARM_CLIENT_SECRET"),
		Environment:              environment,
		SkipProviderRegistration: false,
	}
	return &config
}

func TestAccAzureRMResourceProviderRegistration(t *testing.T) {
	config := testGetAzureConfig(t)
	if config == nil {
		return
	}

	armClient, err := getArmClient(config)
	if err != nil {
		t.Fatalf("Error building ARM Client: %+v", err)
	}

	client := armClient.providersClient
	ctx := testAccProvider.StopContext()
	providerList, err := client.List(ctx, nil, "")
	if err != nil {
		t.Fatalf("Unable to list provider registration status, it is possible that this is due to invalid "+
			"credentials or the service principal does not have permission to use the Resource Manager API, Azure "+
			"error: %s", err)
	}

	err = registerAzureResourceProvidersWithSubscription(ctx, providerList.Values(), client)
	if err != nil {
		t.Fatalf("Error registering Resource Providers: %+v", err)
	}

	needingRegistration := determineAzureResourceProvidersToRegister(providerList.Values())
	if len(needingRegistration) > 0 {
		t.Fatalf("'%d' Resource Providers are still Pending Registration: %s", len(needingRegistration), spew.Sprint(needingRegistration))
	}
}
