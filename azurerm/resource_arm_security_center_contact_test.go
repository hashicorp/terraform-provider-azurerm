package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSecurityCenterContact_basic(t *testing.T) {
	resourceName := "azurerm_security_center_contact.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSecurityCenterContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSecurityCenterContact_template("email1@example.com", "+1-555-555-5555", true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterContactExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "email", "email1@example.com"),
					resource.TestCheckResourceAttr(resourceName, "phone", "+1-555-555-5555"),
					resource.TestCheckResourceAttr(resourceName, "alert_notifications", "true"),
					resource.TestCheckResourceAttr(resourceName, "alerts_to_admins", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSecurityCenterContact_update(t *testing.T) {
	resourceName := "azurerm_security_center_contact.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSecurityCenterContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSecurityCenterContact_template("email1@example.com", "+1-555-555-5555", true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterContactExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "email", "email1@example.com"),
					resource.TestCheckResourceAttr(resourceName, "phone", "+1-555-555-5555"),
					resource.TestCheckResourceAttr(resourceName, "alert_notifications", "true"),
					resource.TestCheckResourceAttr(resourceName, "alerts_to_admins", "true"),
				),
			},
			{
				Config: testAccAzureRMSecurityCenterContact_template("email2@example.com", "+1-555-678-6789", false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterContactExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "email", "email2@example.com"),
					resource.TestCheckResourceAttr(resourceName, "phone", "+1-555-678-6789"),
					resource.TestCheckResourceAttr(resourceName, "alert_notifications", "false"),
					resource.TestCheckResourceAttr(resourceName, "alerts_to_admins", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMSecurityCenterContactExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).securityCenterContactsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
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

func testCheckAzureRMSecurityCenterContactDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).securityCenterContactsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext
	for _, res := range s.RootModule().Resources {
		if res.Type != "azurerm_security_center_contact" {
			continue
		}
		resp, err := client.Get(ctx, securityCenterContactName)
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

func testAccAzureRMSecurityCenterContact_template(email, phone string, notifications, adminAlerts bool) string {
	return fmt.Sprintf(`
resource "azurerm_security_center_contact" "test" {
    email = "%s"
    phone = "%s"

    alert_notifications = %t
    alerts_to_admins    = %t
}
`, email, phone, notifications, adminAlerts)
}
