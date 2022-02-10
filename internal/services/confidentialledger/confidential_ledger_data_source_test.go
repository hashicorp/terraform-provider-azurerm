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
			Config: ConfidentialLedgerResource{}.combinedServicePrincipals(data),
		},
		{
			Config: ConfidentialLedgerDataSource{}.combinedServicePrincipals(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(ConfidentialLedgerResource{}),
				check.That(data.ResourceName).Key("aad_based_security_principals.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.0.ledger_role_name").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.1.principal_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.1.tenant_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.1.ledger_role_name").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.2.principal_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.2.tenant_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.2.ledger_role_name").Exists(),
				check.That(data.ResourceName).Key("cert_based_security_principals.0.cert").Exists(),
				check.That(data.ResourceName).Key("cert_based_security_principals.0.ledger_role_name").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("ledger_type").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("tags").Exists(),
			),
		},
	})
}

func (ConfidentialLedgerDataSource) combinedServicePrincipals(data acceptance.TestData) string {
	template := ConfidentialLedgerResource{}.combinedServicePrincipals(data)
	return fmt.Sprintf(`
%s

data "azurerm_confidential_ledger" "test" {
  name                = azurerm_confidential_ledger.test.name
  resource_group_name = azurerm_confidential_ledger.test.resource_group_name
}
`, template)
}
