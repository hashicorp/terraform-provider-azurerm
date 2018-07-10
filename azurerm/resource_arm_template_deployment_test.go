package azurerm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMTemplateDeployment_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMTemplateDeployment_basicMultiple(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTemplateDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTemplateDeploymentExists("azurerm_template_deployment.test"),
				),
			},
		},
	})
}

func TestAccAzureRMTemplateDeployment_disappears(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMTemplateDeployment_basicSingle(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTemplateDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTemplateDeploymentExists("azurerm_template_deployment.test"),
					testCheckAzureRMTemplateDeploymentDisappears("azurerm_template_deployment.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMTemplateDeployment_withParams(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMTemplateDeployment_withParams(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTemplateDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTemplateDeploymentExists("azurerm_template_deployment.test"),
					resource.TestCheckResourceAttr("azurerm_template_deployment.test", "outputs.testOutput", "Output Value"),
				),
			},
		},
	})
}

func TestAccAzureRMTemplateDeployment_withParamsBody(t *testing.T) {
	ri := acctest.RandInt()
	config := testaccAzureRMTemplateDeployment_withParamsBody(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTemplateDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTemplateDeploymentExists("azurerm_template_deployment.test"),
					resource.TestCheckResourceAttr("azurerm_template_deployment.test", "outputs.testOutput", "Output Value"),
				),
			},
		},
	})

}

func TestAccAzureRMTemplateDeployment_withOutputs(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMTemplateDeployment_withOutputs(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTemplateDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTemplateDeploymentExists("azurerm_template_deployment.test"),
					resource.TestCheckOutput("tfIntOutput", "-123"),
					resource.TestCheckOutput("tfStringOutput", "Standard_GRS"),

					// these values *should* be 'true' and 'false' but,
					// due to a bug in the way terraform represents bools at various times these are for now 0 and 1
					// see https://github.com/hashicorp/terraform/issues/13512#issuecomment-295389523
					// at a later date these may return the expected 'true' / 'false' and should be changed back
					resource.TestCheckOutput("tfFalseOutput", "0"),
					resource.TestCheckOutput("tfTrueOutput", "1"),
					resource.TestCheckResourceAttr("azurerm_template_deployment.test", "outputs.stringOutput", "Standard_GRS"),
				),
			},
		},
	})
}

func TestAccAzureRMTemplateDeployment_withError(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMTemplateDeployment_withError(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTemplateDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Code=\"DeploymentFailed\""),
			},
		},
	})
}

func testCheckAzureRMTemplateDeploymentExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for template deployment: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).deploymentsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on deploymentsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: TemplateDeployment %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMTemplateDeploymentDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		deploymentName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for template deployment: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).deploymentsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, err := client.Delete(ctx, resourceGroup, deploymentName)
		if err != nil {
			return fmt.Errorf("Failed deleting Deployment %q (Resource Group %q): %+v", deploymentName, resourceGroup, err)
		}

		return waitForTemplateDeploymentToBeDeleted(ctx, client, resourceGroup, deploymentName)
	}
}

func testCheckAzureRMTemplateDeploymentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).deploymentsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_template_deployment" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Template Deployment still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMTemplateDeployment_basicSingle(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
  }

  resource "azurerm_template_deployment" "test" {
    name = "acctesttemplate-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    template_body = <<DEPLOY
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "variables": {
    "location": "[resourceGroup().location]",
    "publicIPAddressType": "Dynamic",
    "apiVersion": "2015-06-15",
    "dnsLabelPrefix": "[concat('terraform-tdacctest', uniquestring(resourceGroup().id))]"
  },
  "resources": [
     {
      "type": "Microsoft.Network/publicIPAddresses",
      "apiVersion": "[variables('apiVersion')]",
      "name": "acctestpip-%d",
      "location": "[variables('location')]",
      "properties": {
        "publicIPAllocationMethod": "[variables('publicIPAddressType')]",
        "dnsSettings": {
          "domainNameLabel": "[variables('dnsLabelPrefix')]"
        }
      }
    }
  ]
}
DEPLOY
    deployment_mode = "Complete"
  }
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTemplateDeployment_basicMultiple(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
  }

  resource "azurerm_template_deployment" "test" {
    name = "acctesttemplate-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    template_body = <<DEPLOY
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "storageAccountType": {
      "type": "string",
      "defaultValue": "Standard_LRS",
      "allowedValues": [
        "Standard_LRS",
        "Standard_GRS",
        "Standard_ZRS"
      ],
      "metadata": {
        "description": "Storage Account type"
      }
    }
  },
  "variables": {
    "location": "[resourceGroup().location]",
    "storageAccountName": "[concat(uniquestring(resourceGroup().id), 'storage')]",
    "publicIPAddressName": "[concat('myPublicIp', uniquestring(resourceGroup().id))]",
    "publicIPAddressType": "Dynamic",
    "apiVersion": "2015-06-15",
    "dnsLabelPrefix": "[concat('terraform-tdacctest', uniquestring(resourceGroup().id))]"
  },
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts",
      "name": "[variables('storageAccountName')]",
      "apiVersion": "[variables('apiVersion')]",
      "location": "[variables('location')]",
      "properties": {
        "accountType": "[parameters('storageAccountType')]"
      }
    },
    {
      "type": "Microsoft.Network/publicIPAddresses",
      "apiVersion": "[variables('apiVersion')]",
      "name": "[variables('publicIPAddressName')]",
      "location": "[variables('location')]",
      "properties": {
        "publicIPAllocationMethod": "[variables('publicIPAddressType')]",
        "dnsSettings": {
          "domainNameLabel": "[variables('dnsLabelPrefix')]"
        }
      }
    }
  ]
}
DEPLOY
    deployment_mode = "Complete"
  }
`, rInt, location, rInt)
}

func testaccAzureRMTemplateDeployment_withParamsBody(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

output "test" {
  value = "${azurerm_template_deployment.test.outputs["testOutput"]}"
}

resource "azurerm_storage_container" "using-outputs" {
  name = "vhds"
  resource_group_name = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_template_deployment.test.outputs["accountName"]}"
  container_access_type = "private"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test-kv" {
  location = "%s"
  name = "vault%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  "sku" {
    name = "standard"
  }
  tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  enabled_for_template_deployment = true

  access_policy {
    key_permissions = []
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
    secret_permissions = [
      "get","list","set","purge"]
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  }
}

resource "azurerm_key_vault_secret" "test-secret" {
  name = "acctestsecret-%d"
  value = "terraform-test-%d"
  vault_uri = "${azurerm_key_vault.test-kv.vault_uri}"
}

locals {
	"templated-file" = <<TPL
{
"dnsLabelPrefix": {
    "reference": {
      "keyvault": {
        "id": "${azurerm_key_vault.test-kv.id}"
      },
      "secretName": "${azurerm_key_vault_secret.test-secret.name}"
    }
  },
"storageAccountType": {
   "value": "Standard_GRS"
  }
}
TPL
}

resource "azurerm_template_deployment" "test" {
  name = "acctesttemplate-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  template_body = <<DEPLOY
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "storageAccountType": {
      "type": "string",
      "defaultValue": "Standard_LRS",
      "allowedValues": [
        "Standard_LRS",
        "Standard_GRS",
        "Standard_ZRS"
      ],
      "metadata": {
        "description": "Storage Account type"
      }
    },
    "dnsLabelPrefix": {
      "type": "string",
      "metadata": {
        "description": "DNS Label for the Public IP. Must be lowercase. It should match with the following regular expression: ^[a-z][a-z0-9-]{1,61}[a-z0-9]$ or it will raise an error."
      }
    }
  },
  "variables": {
    "location": "[resourceGroup().location]",
    "storageAccountName": "[concat(uniquestring(resourceGroup().id), 'storage')]",
    "publicIPAddressName": "[concat('myPublicIp', uniquestring(resourceGroup().id))]",
    "publicIPAddressType": "Dynamic",
    "apiVersion": "2015-06-15"
  },
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts",
      "name": "[variables('storageAccountName')]",
      "apiVersion": "[variables('apiVersion')]",
      "location": "[variables('location')]",
      "properties": {
        "accountType": "[parameters('storageAccountType')]"
      }
    },
    {
      "type": "Microsoft.Network/publicIPAddresses",
      "apiVersion": "[variables('apiVersion')]",
      "name": "[variables('publicIPAddressName')]",
      "location": "[variables('location')]",
      "properties": {
        "publicIPAllocationMethod": "[variables('publicIPAddressType')]",
        "dnsSettings": {
          "domainNameLabel": "[parameters('dnsLabelPrefix')]"
        }
      }
    }
  ],
  "outputs": {
    "testOutput": {
      "type": "string",
      "value": "Output Value"
    },
    "accountName": {
      "type": "string",
      "value": "[variables('storageAccountName')]"
    }
  }
}
DEPLOY

  parameters_body = "${local.templated-file}"
  deployment_mode = "Complete"
  depends_on = ["azurerm_key_vault_secret.test-secret"]
}
`, rInt, location, location, rInt, rInt, rInt, rInt)

}

func testAccAzureRMTemplateDeployment_withParams(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
  }

  output "test" {
    value = "${azurerm_template_deployment.test.outputs["testOutput"]}"
  }

  resource "azurerm_storage_container" "using-outputs" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_template_deployment.test.outputs["accountName"]}"
    container_access_type = "private"
  }

  resource "azurerm_template_deployment" "test" {
    name = "acctesttemplate-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    template_body = <<DEPLOY
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "storageAccountType": {
      "type": "string",
      "defaultValue": "Standard_LRS",
      "allowedValues": [
        "Standard_LRS",
        "Standard_GRS",
        "Standard_ZRS"
      ],
      "metadata": {
        "description": "Storage Account type"
      }
    },
    "dnsLabelPrefix": {
      "type": "string",
      "metadata": {
        "description": "DNS Label for the Public IP. Must be lowercase. It should match with the following regular expression: ^[a-z][a-z0-9-]{1,61}[a-z0-9]$ or it will raise an error."
      }
    }
  },
  "variables": {
    "location": "[resourceGroup().location]",
    "storageAccountName": "[concat(uniquestring(resourceGroup().id), 'storage')]",
    "publicIPAddressName": "[concat('myPublicIp', uniquestring(resourceGroup().id))]",
    "publicIPAddressType": "Dynamic",
    "apiVersion": "2015-06-15"
  },
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts",
      "name": "[variables('storageAccountName')]",
      "apiVersion": "[variables('apiVersion')]",
      "location": "[variables('location')]",
      "properties": {
        "accountType": "[parameters('storageAccountType')]"
      }
    },
    {
      "type": "Microsoft.Network/publicIPAddresses",
      "apiVersion": "[variables('apiVersion')]",
      "name": "[variables('publicIPAddressName')]",
      "location": "[variables('location')]",
      "properties": {
        "publicIPAllocationMethod": "[variables('publicIPAddressType')]",
        "dnsSettings": {
          "domainNameLabel": "[parameters('dnsLabelPrefix')]"
        }
      }
    }
  ],
  "outputs": {
    "testOutput": {
      "type": "string",
      "value": "Output Value"
    },
    "accountName": {
      "type": "string",
      "value": "[variables('storageAccountName')]"
    }
  }
}
DEPLOY
    parameters {
	dnsLabelPrefix = "terraform-test-%d"
	storageAccountType = "Standard_GRS"
    }
    deployment_mode = "Complete"
  }
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTemplateDeployment_withOutputs(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
  }

  output "tfStringOutput" {
    value = "${lookup(azurerm_template_deployment.test.outputs, "stringOutput")}"
  }

  output "tfIntOutput" {
    value = "${lookup(azurerm_template_deployment.test.outputs, "intOutput")}"
  }

  output "tfFalseOutput" {
    value = "${lookup(azurerm_template_deployment.test.outputs, "falseOutput")}"
  }

  output "tfTrueOutput" {
    value = "${lookup(azurerm_template_deployment.test.outputs, "trueOutput")}"
  }

  resource "azurerm_template_deployment" "test" {
    name = "acctesttemplate-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    template_body = <<DEPLOY
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "storageAccountType": {
      "type": "string",
      "defaultValue": "Standard_LRS",
      "allowedValues": [
        "Standard_LRS",
        "Standard_GRS",
        "Standard_ZRS"
      ],
      "metadata": {
        "description": "Storage Account type"
      }
    },
    "dnsLabelPrefix": {
      "type": "string",
      "metadata": {
        "description": "DNS Label for the Public IP. Must be lowercase. It should match with the following regular expression: ^[a-z][a-z0-9-]{1,61}[a-z0-9]$ or it will raise an error."
      }
    },
    "intParameter": {
      "type": "int",
      "defaultValue": -123
    },
    "falseParameter": {
      "type": "bool",
      "defaultValue": false
    },
    "trueParameter": {
      "type": "bool",
      "defaultValue": true
    }
  },
  "variables": {
    "location": "[resourceGroup().location]",
    "storageAccountName": "[concat(uniquestring(resourceGroup().id), 'storage')]",
    "publicIPAddressName": "[concat('myPublicIp', uniquestring(resourceGroup().id))]",
    "publicIPAddressType": "Dynamic",
    "apiVersion": "2015-06-15"
  },
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts",
      "name": "[variables('storageAccountName')]",
      "apiVersion": "[variables('apiVersion')]",
      "location": "[variables('location')]",
      "properties": {
        "accountType": "[parameters('storageAccountType')]"
      }
    },
    {
      "type": "Microsoft.Network/publicIPAddresses",
      "apiVersion": "[variables('apiVersion')]",
      "name": "[variables('publicIPAddressName')]",
      "location": "[variables('location')]",
      "properties": {
        "publicIPAllocationMethod": "[variables('publicIPAddressType')]",
        "dnsSettings": {
          "domainNameLabel": "[parameters('dnsLabelPrefix')]"
        }
      }
    }
  ],
  "outputs": {
    "stringOutput": {
      "type": "string",
      "value": "[parameters('storageAccountType')]"
    },
    "intOutput": {
      "type": "int",
      "value": "[parameters('intParameter')]"
    },
    "falseOutput": {
      "type": "bool",
      "value": "[parameters('falseParameter')]"
    },
    "trueOutput": {
      "type": "bool",
      "value": "[parameters('trueParameter')]"
    }
  }
}
DEPLOY
    parameters {
      dnsLabelPrefix = "terraform-test-%d"
      storageAccountType = "Standard_GRS"
    }
    deployment_mode = "Incremental"
  }
`, rInt, location, rInt, rInt)
}

// StorageAccount name is too long, forces error
func testAccAzureRMTemplateDeployment_withError(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
  }

  output "test" {
    value = "${lookup(azurerm_template_deployment.test.outputs, "testOutput")}"
  }

  resource "azurerm_template_deployment" "test" {
    name = "acctesttemplate-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    template_body = <<DEPLOY
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "storageAccountType": {
      "type": "string",
      "defaultValue": "Standard_LRS",
      "allowedValues": [
        "Standard_LRS",
        "Standard_GRS",
        "Standard_ZRS"
      ],
      "metadata": {
        "description": "Storage Account type"
      }
    }
  },
  "variables": {
    "location": "[resourceGroup().location]",
    "storageAccountName": "badStorageAccountNameTooLong",
    "apiVersion": "2015-06-15"
  },
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts",
      "name": "[variables('storageAccountName')]",
      "apiVersion": "[variables('apiVersion')]",
      "location": "[variables('location')]",
      "properties": {
        "accountType": "[parameters('storageAccountType')]"
      }
    }
  ],
  "outputs": {
    "testOutput": {
      "type": "string",
      "value": "Output Value"
    }
  }
}
DEPLOY
    parameters {
        storageAccountType = "Standard_GRS"
    }
    deployment_mode = "Complete"
  }
`, rInt, location, rInt)
}
