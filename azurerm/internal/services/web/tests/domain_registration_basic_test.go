package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

var accTestDomainRegistrationName = "acctest-hashi-1234567890.com"

func TestAccDomainRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_domain_registration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkDomainRegistrationIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testDomainRegistration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRegistrationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testDomainRegistration_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-domain-%[1]d"
  location = "%[2]s"
}

resource "azurerm_dns_zone" "test" {
  name                = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_domain_registration" "test" {
  name                = azurerm_dns_zone.test.name
  resource_group_name = azurerm_resource_group.test.name

  admin_contact {
    first_name   = "Admin"
    last_name    = "Contact"
    email        = "admincontacttest@%[3]s"
    phone_number = "555-555-5555"

    mailing_address {
      address_1   = "1 Test Avenue"
      city        = "TestCity1"
      country     = "USA"
      postal_code = "90210"
      state       = "CA"
    }
  }

  billing_contact {
    first_name   = "Bill"
    middle_name  = "Ing"
    last_name    = "Contact"
    email        = "billcontacttest@%[3]s"
    phone_number = "555-555-5555"
    job_title    = "Financial Controller"

    mailing_address {
      address_1   = "2 Test Avenue"
      city        = "TestCity2"
      country     = "USA"
      postal_code = "90210"
      state       = "CA"
    }
  }

  registrant_contact {
    first_name   = "Reg"
    middle_name  = "Istrant"
    last_name    = "Contact"
    email        = "regcontacttest@%[3]s"
    phone_number = "555-555-5555"
    job_title    = "Director"

    mailing_address {
      address_1   = "3 Test Avenue"
      address_2   = "Test Address 3"
      city        = "TestCity3"
      country     = "USA"
      postal_code = "90210"
      state       = "CA"
    }
  }

  technical_contact {
    first_name   = "Tech"
    last_name    = "Contact"
    email        = "techcontacttest@%[3]s"
    phone_number = "555-555-5555"

    mailing_address {
      address_1   = "4 Test Avenue"
      city        = "TestCity4"
      country     = "USA"
      postal_code = "90210"
      state       = "CA"
    }
  }

  dns_zone_id = azurerm_dns_zone.test.id
  privacy     = true

}
`, data.RandomInteger, data.Locations.Primary, accTestDomainRegistrationName)
}
