package azurerm

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	clientID := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantID := os.Getenv("ARM_TENANT_ID")
	testLocation := os.Getenv("ARM_TEST_LOCATION")

	if subscriptionID == "" || clientID == "" || clientSecret == "" || tenantID == "" || testLocation == "" {
		t.Fatal("ARM_SUBSCRIPTION_ID, ARM_CLIENT_ID, ARM_CLIENT_SECRET, ARM_TENANT_ID and ARM_TEST_LOCATION must be set for acceptance tests")
	}
}

func testLocation() string {
	return os.Getenv("ARM_TEST_LOCATION")
}
