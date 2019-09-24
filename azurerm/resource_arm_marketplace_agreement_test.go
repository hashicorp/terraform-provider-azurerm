package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMarketplaceAgreement(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":          testAccAzureRMMarketplaceAgreement_basic,
			"requiresImport": testAccAzureRMMarketplaceAgreement_requiresImport,
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

func testAccAzureRMMarketplaceAgreement_basic(t *testing.T) {
	resourceName := "azurerm_marketplace_agreement.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMarketplaceAgreementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMarketplaceAgreement_basicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMarketplaceAgreementExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "license_text_link"),
					resource.TestCheckResourceAttrSet(resourceName, "privacy_policy_link"),
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

func testAccAzureRMMarketplaceAgreement_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_marketplace_agreement.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMarketplaceAgreementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMarketplaceAgreement_basicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMarketplaceAgreementExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMMarketplaceAgreement_requiresImportConfig(),
				ExpectError: testRequiresImportError("azurerm_marketplace_agreement"),
			},
		},
	})
}

func testCheckAzureRMMarketplaceAgreementExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		offer := rs.Primary.Attributes["offer"]
		plan := rs.Primary.Attributes["plan"]
		publisher := rs.Primary.Attributes["publisher"]

		client := testAccProvider.Meta().(*ArmClient).compute.MarketplaceAgreementsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, publisher, offer, plan)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Marketplace Agreement for Publisher %q / Offer %q / Plan %q does not exist", publisher, offer, plan)
			}
			return fmt.Errorf("Bad: Get on MarketplaceAgreementsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMarketplaceAgreementDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_marketplace_agreement" {
			continue
		}

		offer := rs.Primary.Attributes["offer"]
		plan := rs.Primary.Attributes["plan"]
		publisher := rs.Primary.Attributes["publisher"]

		client := testAccProvider.Meta().(*ArmClient).compute.MarketplaceAgreementsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, publisher, offer, plan)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Marketplace Agreement still exists:\n%#v", resp)
			}
		}
	}

	return nil
}

func testAccAzureRMMarketplaceAgreement_basicConfig() string {
	return fmt.Sprintf(`
resource "azurerm_marketplace_agreement" "test" {
  publisher = "barracudanetworks"
  offer     = "waf"
  plan      = "hourly"
}
`)
}

func testAccAzureRMMarketplaceAgreement_requiresImportConfig() string {
	template := testAccAzureRMMarketplaceAgreement_basicConfig()
	return fmt.Sprintf(`
%s

resource "azurerm_marketplace_agreement" "import" {
  publisher = azurerm_marketplace_agreement.test.publisher
  offer     = azurerm_marketplace_agreement.test.offer
  plan      = azurerm_marketplace_agreement.test.plan
}
`, template)
}
