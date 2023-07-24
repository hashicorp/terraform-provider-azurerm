package paloalto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/certificateobjectlocalrulestack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRulestackOutboundUnTrustCertificateResource struct{}

func TestAccLocalRulestackOutboundUnTrustCertificateResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_rulestack_outbound_untrust_certificate_association", "test")
	r := LocalRulestackOutboundUnTrustCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LocalRulestackOutboundUnTrustCertificateResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := certificateobjectlocalrulestack.ParseLocalRulestackCertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PaloAlto.CertificatesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r LocalRulestackOutboundUnTrustCertificateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_local_rulestack_outbound_untrust_certificate_association" "test" {
  rulestack_id   = azurerm_palo_alto_local_rule_stack.test.id
  certificate_id = azurerm_palo_alto_local_rule_stack_certificate.test.id
}

`, r.template(data))
}

func (r LocalRulestackOutboundUnTrustCertificateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-PAN-%[1]d"
  location = "%[2]s"
}

resource "azurerm_palo_alto_local_rule_stack" "test" {
  name                = "testAcc-palrs-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
}

resource "azurerm_palo_alto_local_rule_stack_certificate" "test" {
  name          = "testacc-palc-%[1]d"
  rule_stack_id = azurerm_palo_alto_local_rule_stack.test.id

  certificate_signer_id = "https://example.com/not-a-real-url"
}
`, data.RandomInteger, data.Locations.Primary)
}
