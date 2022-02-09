package confidentialledger_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ConfidentialLedgerDataSource struct{}

func TestAccConfidentialLedgerDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_confidential_ledger", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: ConfidentialLedgerResource{}.standard(data),
		},
		{
			Config: ConfidentialLedgerDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(ConfidentialLedgerResource{}),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("sku").Exists(),
				check.That(data.ResourceName).Key("primary_read_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_read_key.0.id").Exists(),
				check.That(data.ResourceName).Key("primary_read_key.0.secret").Exists(),
				check.That(data.ResourceName).Key("primary_write_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_write_key.0.id").Exists(),
				check.That(data.ResourceName).Key("primary_write_key.0.secret").Exists(),
				check.That(data.ResourceName).Key("secondary_read_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_read_key.0.id").Exists(),
				check.That(data.ResourceName).Key("secondary_read_key.0.secret").Exists(),
				check.That(data.ResourceName).Key("secondary_write_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_write_key.0.id").Exists(),
				check.That(data.ResourceName).Key("secondary_write_key.0.secret").Exists(),
			),
		},
	})
}

func (ConfidentialLedgerDataSource) basic(data acceptance.TestData) string {
	template := ConfidentialLedgerResource{}.standard(data)
	return fmt.Sprintf(`
%s

data "azurerm_confidential_ledger" "test" {
  name                = azurerm_confidential_ledger.test.name
  resource_group_name = azurerm_confidential_ledger.test.resource_group_name
}
`, template)
}
