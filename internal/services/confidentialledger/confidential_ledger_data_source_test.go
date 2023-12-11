// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
			Config: ConfidentialLedgerDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(ConfidentialLedgerResource{}),
				check.That(data.ResourceName).Key("identity_service_endpoint").Exists(),
				check.That(data.ResourceName).Key("ledger_endpoint").Exists(),
				check.That(data.ResourceName).Key("ledger_type").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("azuread_based_service_principal.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("azuread_based_service_principal.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("azuread_based_service_principal.0.ledger_role_name").Exists(),
				check.That(data.ResourceName).Key("certificate_based_security_principal.0.ledger_role_name").Exists(),
				check.That(data.ResourceName).Key("certificate_based_security_principal.0.pem_public_key").Exists(),
			),
		},
	})
}

func (ConfidentialLedgerDataSource) basic(data acceptance.TestData) string {
	template := ConfidentialLedgerResource{}.certBased(data)
	return fmt.Sprintf(`
%s

data "azurerm_confidential_ledger" "test" {
  name                = azurerm_confidential_ledger.test.name
  resource_group_name = azurerm_confidential_ledger.test.resource_group_name
}
`, template)
}
