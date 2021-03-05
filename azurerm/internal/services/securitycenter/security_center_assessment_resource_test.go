package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SecurityCenterAssessmentResource struct{}

func TestAccSecurityCenterAssessment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment", "test")
	r := SecurityCenterAssessmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterAssessment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment", "test")
	r := SecurityCenterAssessmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSecurityCenterAssessment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment", "test")
	r := SecurityCenterAssessmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterAssessment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_assessment", "test")
	r := SecurityCenterAssessmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SecurityCenterAssessmentResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	assessmentClient := client.SecurityCenter.AssessmentsClient
	id, err := parse.AssessmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := assessmentClient.Get(ctx, id.TargetResourceID, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Azure Security Center Assessment %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.AssessmentProperties != nil), nil
}

func (r SecurityCenterAssessmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_assessment" "test" {
  assessment_policy_id = azurerm_security_center_assessment_policy.test.id
  target_resource_id   = azurerm_linux_virtual_machine_scale_set.test.id

  status {
    code = "Healthy"
  }
}
`, r.template(data))
}

func (r SecurityCenterAssessmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_assessment" "import" {
  assessment_policy_id = azurerm_security_center_assessment.test.assessment_policy_id
  target_resource_id   = azurerm_security_center_assessment.test.target_resource_id

  status {
    code = azurerm_security_center_assessment.test.status.0.code
  }
}
`, r.basic(data))
}

func (r SecurityCenterAssessmentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_assessment" "test" {
  assessment_policy_id = azurerm_security_center_assessment_policy.test.id
  target_resource_id   = azurerm_linux_virtual_machine_scale_set.test.id

  status {
    code        = "Unhealthy"
    cause       = "un healthy"
    description = "description for acctest"
  }

  additional_data = {
    "Env" : "Test",
    "Foo" : "Bar"
  }
}
`, r.template(data))
}

func (r SecurityCenterAssessmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-SecurityCenter-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}

resource "azurerm_security_center_assessment_policy" "test" {
  display_name = "Test Display Name"
  severity     = "Medium"
  description  = "Test Description"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
