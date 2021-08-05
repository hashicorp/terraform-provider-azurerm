package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SecurityCenterContactResource struct {
}

func TestAccSecurityCenterContact(t *testing.T) {
	// there is only *one* read contact, if tests will conflict if run at the same time
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"contact": {
			"basic":          testAccSecurityCenterContact_basic,
			"update":         testAccSecurityCenterContact_update,
			"requiresImport": testAccSecurityCenterContact_requiresImport,
			"phoneOptional":  testAccSecurityCenterContact_phoneOptional,
		},
	})
}

func testAccSecurityCenterContact_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_contact", "test")
	r := SecurityCenterContactResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.template("basic@example.com", "+1-555-555-5555", true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("email").HasValue("basic@example.com"),
				check.That(data.ResourceName).Key("phone").HasValue("+1-555-555-5555"),
				check.That(data.ResourceName).Key("alert_notifications").HasValue("true"),
				check.That(data.ResourceName).Key("alerts_to_admins").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func testAccSecurityCenterContact_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_contact", "test")
	r := SecurityCenterContactResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.template("require@example.com", "+1-555-555-5555", true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("email").HasValue("require@example.com"),
				check.That(data.ResourceName).Key("phone").HasValue("+1-555-555-5555"),
				check.That(data.ResourceName).Key("alert_notifications").HasValue("true"),
				check.That(data.ResourceName).Key("alerts_to_admins").HasValue("true"),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImportCfg("email1@example.com", "+1-555-555-5555", true, true)
		}),
	})
}

func testAccSecurityCenterContact_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_contact", "test")
	r := SecurityCenterContactResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.template("update@example.com", "+1-555-555-5555", true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("email").HasValue("update@example.com"),
				check.That(data.ResourceName).Key("phone").HasValue("+1-555-555-5555"),
				check.That(data.ResourceName).Key("alert_notifications").HasValue("true"),
				check.That(data.ResourceName).Key("alerts_to_admins").HasValue("true"),
			),
		},
		{
			Config: r.template("updated@example.com", "+1-555-678-6789", false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("email").HasValue("updated@example.com"),
				check.That(data.ResourceName).Key("phone").HasValue("+1-555-678-6789"),
				check.That(data.ResourceName).Key("alert_notifications").HasValue("false"),
				check.That(data.ResourceName).Key("alerts_to_admins").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func testAccSecurityCenterContact_phoneOptional(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_contact", "test")
	r := SecurityCenterContactResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.templateWithoutPhone("basic@example.com", true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("email").HasValue("basic@example.com"),
				check.That(data.ResourceName).Key("phone").HasValue(""),
				check.That(data.ResourceName).Key("alert_notifications").HasValue("true"),
				check.That(data.ResourceName).Key("alerts_to_admins").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func (SecurityCenterContactResource) Exists(ctx context.Context, clients *clients.Client, _ *pluginsdk.InstanceState) (*bool, error) {
	contactName := "default1"

	resp, err := clients.SecurityCenter.ContactsClient.Get(ctx, contactName)
	if err != nil {
		return nil, fmt.Errorf("reading Security Center Subscription Contact (%s): %+v", contactName, err)
	}

	return utils.Bool(resp.ContactProperties != nil), nil
}

func (SecurityCenterContactResource) template(email, phone string, notifications, adminAlerts bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_contact" "test" {
  email = "%s"
  phone = "%s"

  alert_notifications = %t
  alerts_to_admins    = %t
}
`, email, phone, notifications, adminAlerts)
}

func (SecurityCenterContactResource) templateWithoutPhone(email string, notifications, adminAlerts bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_contact" "test" {
  email = "%s"

  alert_notifications = %t
  alerts_to_admins    = %t
}
`, email, notifications, adminAlerts)
}

func (r SecurityCenterContactResource) requiresImportCfg(email, phone string, notifications, adminAlerts bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_contact" "import" {
  email = azurerm_security_center_contact.test.email
  phone = azurerm_security_center_contact.test.phone

  alert_notifications = azurerm_security_center_contact.test.alert_notifications
  alerts_to_admins    = azurerm_security_center_contact.test.alerts_to_admins
}
`, r.template(email, phone, notifications, adminAlerts))
}
