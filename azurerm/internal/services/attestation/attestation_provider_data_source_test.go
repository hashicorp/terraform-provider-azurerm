package attestation_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AttestationProviderDataSource struct {
}

func TestAccAttestationProviderDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_attestation_provider", "test")
	randStr := strings.ToLower(acceptance.RandString(10))

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AttestationProviderDataSource{}.basic(data, randStr),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(AttestationProviderResource{}),
			),
		},
	})
}

func (AttestationProviderDataSource) basic(data acceptance.TestData, randStr string) string {
	config := AttestationProviderResource{}.basic(data, randStr)
	return fmt.Sprintf(`
%s

data "azurerm_attestation_provider" "test" {
  name                = azurerm_attestation_provider.test.name
  resource_group_name = azurerm_attestation_provider.test.resource_group_name
}
`, config)
}
