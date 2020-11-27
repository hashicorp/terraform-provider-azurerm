package batch_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BatchCertificateResource struct {
}

func TestAccBatchCertificate_Pfx(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_certificate", "test")
	r := BatchCertificateResource{}
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	certificateID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-42c107874fd0e4a9583292a2f1098e8fe4b2edda", subscriptionID, data.RandomInteger, data.RandomString)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.pfx(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").HasValue(certificateID),
				check.That(data.ResourceName).Key("format").HasValue("Pfx"),
				check.That(data.ResourceName).Key("thumbprint").HasValue("42c107874fd0e4a9583292a2f1098e8fe4b2edda"),
				check.That(data.ResourceName).Key("thumbprint_algorithm").HasValue("sha1"),
			),
		},
		data.ImportStep("certificate", "password"),
	})
}

func TestAccBatchCertificate_PfxWithoutPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_certificate", "test")
	r := BatchCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.pfxWithoutPassword(data),
			ExpectError: regexp.MustCompile("Password is required"),
		},
	})
}

func TestAccBatchCertificate_Cer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_certificate", "test")
	r := BatchCertificateResource{}
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	certificateID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-312d31a79fa0cef49c00f769afc2b73e9f4edf34", subscriptionID, data.RandomInteger, data.RandomString)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.cer(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").HasValue(certificateID),
				check.That(data.ResourceName).Key("format").HasValue("Cer"),
				check.That(data.ResourceName).Key("thumbprint").HasValue("312d31a79fa0cef49c00f769afc2b73e9f4edf34"),
				check.That(data.ResourceName).Key("thumbprint_algorithm").HasValue("sha1"),
			),
		},
		data.ImportStep("certificate"),
	})
}

func TestAccBatchCertificate_CerWithPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_certificate", "test")
	r := BatchCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.cerwithPassword(data),
			ExpectError: regexp.MustCompile("Password must not be specified"),
		},
	})
}

func (t BatchCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.CertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Batch.CertificateClient.Get(ctx, id.Name, id.BatchAccountName, id.ResourceGroup)
	if err != nil {
		return nil, fmt.Errorf("retrieving Batch Certificate %q (Account Name %q / Resource Group %q) does not exist", id.Name, id.BatchAccountName, id.ResourceGroup)
	}

	return utils.Bool(resp.CertificateProperties != nil), nil
}

func (BatchCertificateResource) pfx(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate.pfx")
  format               = "Pfx"
  password             = "terraform"
  thumbprint           = "42c107874fd0e4a9583292a2f1098e8fe4b2edda"
  thumbprint_algorithm = "SHA1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (BatchCertificateResource) pfxWithoutPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate.pfx")
  format               = "Pfx"
  thumbprint           = "42c107874fd0e4a9583292a2f1098e8fe4b2edda"
  thumbprint_algorithm = "SHA1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (BatchCertificateResource) cer(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate.cer")
  format               = "Cer"
  thumbprint           = "312d31a79fa0cef49c00f769afc2b73e9f4edf34"
  thumbprint_algorithm = "SHA1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (BatchCertificateResource) cerwithPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate.cer")
  format               = "Cer"
  password             = "should not have a password for Cer"
  thumbprint           = "312d31a79fa0cef49c00f769afc2b73e9f4edf34"
  thumbprint_algorithm = "SHA1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
