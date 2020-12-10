package securitycenter_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccSecurityCenter_contact(t *testing.T) {
	// there is only *one* read contact, if tests will conflict if run at the same time
	testCases := map[string]map[string]func(t *testing.T){
		"contact": {
			"basic":          testAccSecurityCenterContact_basic,
			"update":         testAccSecurityCenterContact_update,
			"requiresImport": testAccSecurityCenterContact_requiresImport,
			"phoneOptional":  testAccSecurityCenterContact_phoneOptional,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccSecurityCenterContact_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_contact", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSecurityCenterContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterContact_template("basic@example.com", "+1-555-555-5555", true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityCenterContactExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "email", "basic@example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "phone", "+1-555-555-5555"),
					resource.TestCheckResourceAttr(data.ResourceName, "alert_notifications", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "alerts_to_admins", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccSecurityCenterContact_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_contact", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSecurityCenterContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterContact_template("require@example.com", "+1-555-555-5555", true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityCenterContactExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "email", "require@example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "phone", "+1-555-555-5555"),
					resource.TestCheckResourceAttr(data.ResourceName, "alert_notifications", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "alerts_to_admins", "true"),
				),
			},
			data.RequiresImportErrorStep(func(data acceptance.TestData) string {
				return testAccSecurityCenterContact_requiresImportCfg("email1@example.com", "+1-555-555-5555", true, true)
			}),
		},
	})
}

func testAccSecurityCenterContact_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_contact", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSecurityCenterContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterContact_template("update@example.com", "+1-555-555-5555", true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityCenterContactExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "email", "update@example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "phone", "+1-555-555-5555"),
					resource.TestCheckResourceAttr(data.ResourceName, "alert_notifications", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "alerts_to_admins", "true"),
				),
			},
			{
				Config: testAccSecurityCenterContact_template("updated@example.com", "+1-555-678-6789", false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityCenterContactExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "email", "updated@example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "phone", "+1-555-678-6789"),
					resource.TestCheckResourceAttr(data.ResourceName, "alert_notifications", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "alerts_to_admins", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccSecurityCenterContact_phoneOptional(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_contact", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSecurityCenterContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterContact_templateWithoutPhone("basic@example.com", true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityCenterContactExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "email", "basic@example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "phone", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "alert_notifications", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "alerts_to_admins", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckSecurityCenterContactExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.ContactsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		contactName := rs.Primary.Attributes["securityContacts"]

		resp, err := client.Get(ctx, contactName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Security Center Subscription Contact %q was not found: %+v", contactName, err)
			}

			return fmt.Errorf("Bad: GetContact: %+v", err)
		}

		return nil
	}
}

func testCheckSecurityCenterContactDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.ContactsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	for _, res := range s.RootModule().Resources {
		if res.Type != "azurerm_security_center_contact" {
			continue
		}
		resp, err := client.Get(ctx, "default1")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}
		return fmt.Errorf("security center contact still exists")
	}
	return nil
}

func testAccSecurityCenterContact_template(email, phone string, notifications, adminAlerts bool) string {
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

func testAccSecurityCenterContact_templateWithoutPhone(email string, notifications, adminAlerts bool) string {
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

func testAccSecurityCenterContact_requiresImportCfg(email, phone string, notifications, adminAlerts bool) string {
	template := testAccSecurityCenterContact_template(email, phone, notifications, adminAlerts)
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_contact" "import" {
  email = azurerm_security_center_contact.test.email
  phone = azurerm_security_center_contact.test.phone

  alert_notifications = azurerm_security_center_contact.test.alert_notifications
  alerts_to_admins    = azurerm_security_center_contact.test.alerts_to_admins
}
`, template)
}
