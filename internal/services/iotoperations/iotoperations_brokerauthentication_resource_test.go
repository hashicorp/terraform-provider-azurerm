package iotoperations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBrokerAuthentication_basic(t *testing.T) {
	resourceName := "azurerm_broker_authentication.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { /* add pre-checks if needed */ },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBrokerAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBrokerAuthenticationConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-broker-auth"),
					resource.TestCheckResourceAttr(resourceName, "resource_group_name", "test-rg"),
					resource.TestCheckResourceAttr(resourceName, "instance_name", "test-instance"),
					resource.TestCheckResourceAttr(resourceName, "broker_name", "test-broker"),
				),
			},
		},
	})
}

func testAccBrokerAuthenticationConfig_basic() string {
	return fmt.Sprintf(`
resource "azurerm_broker_authentication" "test" {
  name                = "test-broker-auth"
  resource_group_name = "test-rg"
  instance_name       = "test-instance"
  broker_name         = "test-broker"

  authentication_methods {
    method = "ServiceAccountToken"
    # Add other nested fields as needed
  }

  extended_location {
    name = "qmbrfwcpwwhggszhrdjv"
    type = "CustomLocation"
  }
}
`)
}
