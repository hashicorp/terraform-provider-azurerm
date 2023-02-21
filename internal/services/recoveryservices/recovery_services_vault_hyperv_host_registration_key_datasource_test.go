package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type RecoveryServicesVaultHyperVHostRegistrationKeyDataSource struct{}

func TestAccRecoveryServicesVaultHyperVHostRegistrationKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_recovery_services_vault_hyperv_host_registration_key", "test")
	r := RecoveryServicesVaultHyperVHostRegistrationKeyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("xml_content").Exists(),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("management_cert").Exists(),
				check.That(data.ResourceName).Key("aad_tenant_id").Exists(),
				check.That(data.ResourceName).Key("aad_authority").Exists(),
				check.That(data.ResourceName).Key("service_principal_client_id").Exists(),
				check.That(data.ResourceName).Key("aad_vault_audience").Exists(),
				check.That(data.ResourceName).Key("aad_management_endpoint").Exists(),
				check.That(data.ResourceName).Key("vault_private_endpoint_enabled").Exists(),
				check.That(data.ResourceName).Key("validate_to").Exists(),
			),
		},
	})
}

func (RecoveryServicesVaultHyperVHostRegistrationKeyDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_recovery_services_vault_hyperv_host_registration_key" "test" {
	site_recovery_services_vault_hyperv_site_id = azurerm_site_recovery_services_vault_hyperv_site.test.id
}
`, HyperVSiteResource{}.basic(data))
}
