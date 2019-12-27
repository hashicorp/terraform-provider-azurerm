package tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventHubNamespaceAuthorizationRule_listen(t *testing.T) {
	testAccAzureRMEventHubNamespaceAuthorizationRule(t, true, false, false)
}

func TestAccAzureRMEventHubNamespaceAuthorizationRule_send(t *testing.T) {
	testAccAzureRMEventHubNamespaceAuthorizationRule(t, false, true, false)
}

func TestAccAzureRMEventHubNamespaceAuthorizationRule_listensend(t *testing.T) {
	testAccAzureRMEventHubNamespaceAuthorizationRule(t, true, true, false)
}

func TestAccAzureRMEventHubNamespaceAuthorizationRule_manage(t *testing.T) {
	testAccAzureRMEventHubNamespaceAuthorizationRule(t, true, true, true)
}

func testAccAzureRMEventHubNamespaceAuthorizationRule(t *testing.T, listen, send, manage bool) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceAuthorizationRule_base(data, listen, send, manage),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", strconv.FormatBool(listen)),
					resource.TestCheckResourceAttr(data.ResourceName, "send", strconv.FormatBool(send)),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", strconv.FormatBool(manage)),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventHubNamespaceAuthorizationRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceAuthorizationRule_base(data, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceAuthorizationRuleExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMEventHubNamespaceAuthorizationRule_requiresImport(data, true, true, true),
				ExpectError: acceptance.RequiresImportError("azurerm_eventhub_namespace_authorization_rule"),
			},
		},
	})
}

func TestAccAzureRMEventHubNamespaceAuthorizationRule_rightsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceAuthorizationRule_base(data, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "false"),
				),
			},
			{
				Config: testAccAzureRMEventHubNamespaceAuthorizationRule_base(data, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMEventHubNamespaceAuthorizationRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.NamespacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub_authorization_rule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testCheckAzureRMEventHubNamespaceAuthorizationRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.NamespacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Event Hub: %s", name)
		}

		resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Event Hub Namespace Authorization Rule %q (namespace %q / resource group: %q) does not exist", name, namespaceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on eventHubClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMEventHubNamespaceAuthorizationRule_base(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku = "Standard"
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "acctest-EHN-AR%[1]d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  listen = %[3]t
  send   = %[4]t
  manage = %[5]t
}
`, data.RandomInteger, data.Locations.Primary, listen, send, manage)
}

func testAccAzureRMEventHubNamespaceAuthorizationRule_requiresImport(data acceptance.TestData, listen, send, manage bool) string {
	template := testAccAzureRMEventHubNamespaceAuthorizationRule_base(data, listen, send, manage)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace_authorization_rule" "import" {
  name                = "${azurerm_eventhub_namespace_authorization_rule.test.name}"
  namespace_name      = "${azurerm_eventhub_namespace_authorization_rule.test.namespace_name}"
  resource_group_name = "${azurerm_eventhub_namespace_authorization_rule.test.resource_group_name}"
  listen              = "${azurerm_eventhub_namespace_authorization_rule.test.listen}"
  send                = "${azurerm_eventhub_namespace_authorization_rule.test.send}"
  manage              = "${azurerm_eventhub_namespace_authorization_rule.test.manage}"
}
`, template)
}
