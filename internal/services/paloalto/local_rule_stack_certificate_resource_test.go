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

type LocalRulestackCertificateResource struct{}

func TestAccPaloAltoLocalRulestackCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rule_stack_certificate", "test")

	r := LocalRulestackCertificateResource{}

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

func TestAccPaloAltoLocalRulestackCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rule_stack_certificate", "test")

	r := LocalRulestackCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError(data.ResourceName),
		},
	})
}

func TestAccPaloAltoLocalRulestackCertificate_completeSelfSigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rule_stack_certificate", "test")

	r := LocalRulestackCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeSelfSigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPaloAltoLocalRulestackCertificate_selfSignedUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rule_stack_certificate", "test")

	r := LocalRulestackCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeSelfSigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeSelfSignedUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPaloAltoLocalRulestackCertificate_authoritySigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rule_stack_certificate", "test")

	r := LocalRulestackCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeAuthoritySigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LocalRulestackCertificateResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r LocalRulestackCertificateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s 

resource "azurerm_palo_alto_local_rule_stack_certificate" "test" {
  name          = "testacc-palc-%[2]d"
  rule_stack_id = azurerm_palo_alto_local_rule_stack.test.id
  self_signed   = true
}

`, r.template(data), data.RandomInteger)
}

func (r LocalRulestackCertificateResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s 

resource "azurerm_palo_alto_local_rule_stack_certificate" "test" {
  name          = azurerm_palo_alto_local_rule_stack_certificate.test.name
  rule_stack_id = azurerm_palo_alto_local_rule_stack_certificate.test.rule_stack_id
  self_signed   = azurerm_palo_alto_local_rule_stack_certificate.test.self_signed
}

`, r.basic(data))
}

func (r LocalRulestackCertificateResource) completeSelfSigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s 

resource "azurerm_palo_alto_local_rule_stack_certificate" "test" {
  name          = "testacc-palc-%[2]d"
  rule_stack_id = azurerm_palo_alto_local_rule_stack.test.id
  self_signed   = true

  audit_comment = "Acceptance test audit comment - %[2]d"
  description   = "Acceptance test Desc - %[2]d"
}

`, r.template(data), data.RandomInteger)
}

func (r LocalRulestackCertificateResource) completeSelfSignedUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s 

resource "azurerm_palo_alto_local_rule_stack_certificate" "test" {
  name          = "testacc-palc-%[2]d"
  rule_stack_id = azurerm_palo_alto_local_rule_stack.test.id
  self_signed   = true

  audit_comment = "Updated acceptance test audit comment - %[2]d"
  description   = "Updated acceptance test Desc - %[2]d"
}

`, r.template(data), data.RandomInteger)
}

func (r LocalRulestackCertificateResource) completeAuthoritySigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s 

resource "azurerm_palo_alto_local_rule_stack_certificate" "test" {
  name          = "testacc-palc-%[2]d"
  rule_stack_id = azurerm_palo_alto_local_rule_stack.test.id
  
  certificate_signer_id = "https://example.com/not-a-real-url"

  audit_comment = "Acceptance test audit comment - %[2]d"
  description   = "Acceptance test Desc - %[2]d"
}

`, r.template(data), data.RandomInteger)
}

func (r LocalRulestackCertificateResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
}
