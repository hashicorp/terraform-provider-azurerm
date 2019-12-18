package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azuread/azuread"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	azureProvider := Provider().(*schema.Provider)

	testAccProvider = azureProvider
	acceptance.AzureProvider = azureProvider

	// NOTE: these /cannot/ be simplified into a single shared variable (tried, it causes a nil-slice)
	testAccProviders = map[string]terraform.ResourceProvider{
		"azurerm": testAccProvider,
		"azuread": azuread.Provider().(*schema.Provider),
	}
	acceptance.SupportedProviders = map[string]terraform.ResourceProvider{
		"azurerm": testAccProvider,
		"azuread": azuread.Provider().(*schema.Provider),
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestDataSourcesSupportCustomTimeouts(t *testing.T) {
	// this is required until 2.0
	os.Setenv("ARM_PROVIDER_CUSTOM_TIMEOUTS", "true")

	provider := Provider().(*schema.Provider)
	for dataSourceName, dataSource := range provider.DataSourcesMap {
		t.Run(fmt.Sprintf("DataSource/%s", dataSourceName), func(t *testing.T) {
			t.Logf("[DEBUG] Testing Data Source %q..", dataSourceName)

			if dataSource.Timeouts == nil {
				t.Fatalf("Data Source %q has no timeouts block defined!", dataSourceName)
			}

			// Azure's bespoke enough we want to be explicit about the timeouts for each value
			if dataSource.Timeouts.Default != nil {
				t.Fatalf("Data Source %q defines a Default timeout when it shouldn't!", dataSourceName)
			}

			// Data Sources must have a Read
			if dataSource.Timeouts.Read == nil {
				t.Fatalf("Data Source %q doesn't define a Read timeout", dataSourceName)
			}

			// but shouldn't have anything else
			if dataSource.Timeouts.Create != nil {
				t.Fatalf("Data Source %q defines a Create timeout when it shouldn't!", dataSourceName)
			}

			if dataSource.Timeouts.Update != nil {
				t.Fatalf("Data Source %q defines a Update timeout when it shouldn't!", dataSourceName)
			}

			if dataSource.Timeouts.Delete != nil {
				t.Fatalf("Data Source %q defines a Delete timeout when it shouldn't!", dataSourceName)
			}
		})
	}
}

func TestResourcesSupportCustomTimeouts(t *testing.T) {
	// this is required until 2.0
	os.Setenv("ARM_PROVIDER_CUSTOM_TIMEOUTS", "true")

	provider := Provider().(*schema.Provider)
	for resourceName, resource := range provider.ResourcesMap {
		t.Run(fmt.Sprintf("Resource/%s", resourceName), func(t *testing.T) {
			t.Logf("[DEBUG] Testing Resource %q..", resourceName)

			if resource.Timeouts == nil {
				t.Fatalf("Resource %q has no timeouts block defined!", resourceName)
			}

			// Azure's bespoke enough we want to be explicit about the timeouts for each value
			if resource.Timeouts.Default != nil {
				t.Fatalf("Resource %q defines a Default timeout when it shouldn't!", resourceName)
			}

			// every Resource has to have a Create, Read & Destroy timeout
			if resource.Timeouts.Create == nil && resource.Create != nil {
				t.Fatalf("Resource %q defines a Create method but no Create Timeout", resourceName)
			}
			if resource.Timeouts.Delete == nil && resource.Delete != nil {
				t.Fatalf("Resource %q defines a Delete method but no Delete Timeout", resourceName)
			}
			if resource.Timeouts.Read == nil {
				t.Fatalf("Resource %q doesn't define a Read timeout", resourceName)
			}

			// Optional
			if resource.Timeouts.Update == nil && resource.Update != nil {
				t.Fatalf("Resource %q defines a Update method but no Update Timeout", resourceName)
			}
		})
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}
