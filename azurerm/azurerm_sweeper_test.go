package azurerm

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfigForRegion(region string) (*ArmClient, error) {
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	clientID := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantID := os.Getenv("ARM_TENANT_ID")
	environment := os.Getenv("ARM_ENVIRONMENT")

	if environment == "" {
		environment = "public"
	}

	if subscriptionID == "" || clientID == "" || clientSecret == "" || tenantID == "" {
		return nil, fmt.Errorf("ARM_SUBSCRIPTION_ID, ARM_CLIENT_ID, ARM_CLIENT_SECRET and ARM_TENANT_ID must be set for acceptance tests")
	}

	config := &Config{
		SubscriptionID:           subscriptionID,
		ClientID:                 clientID,
		ClientSecret:             clientSecret,
		TenantID:                 tenantID,
		Environment:              environment,
		SkipProviderRegistration: false,
	}

	return config.getArmClient()
}

func shouldSweepAcceptanceTestResource(name string) bool {
	loweredName := strings.ToLower(name)
	return strings.HasPrefix(loweredName, "acctest")
}

func resourceLocatedInRegion(resourceLocation string, region string) bool {
	normalisedResourceLocation := azureRMNormalizeLocation(resourceLocation)
	normalisedRegion := azureRMNormalizeLocation(region)
	return normalisedResourceLocation == normalisedRegion
}