package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestValidateMetricAlertRuleTags(t *testing.T) {
	cases := []struct {
		Name     string
		Value    map[string]interface{}
		ErrCount int
	}{
		{
			Name: "Single Valid",
			Value: map[string]interface{}{
				"hello": "world",
			},
			ErrCount: 0,
		},
		{
			Name: "Single Invalid",
			Value: map[string]interface{}{
				"$Type": "hello/world",
			},
			ErrCount: 1,
		},
		{
			Name: "Single Invalid lowercase",
			Value: map[string]interface{}{
				"$type": "hello/world",
			},
			ErrCount: 1,
		},
		{
			Name: "Multiple Valid",
			Value: map[string]interface{}{
				"hello": "world",
				"foo":   "bar",
			},
			ErrCount: 0,
		},
		{
			Name: "Multiple Invalid",
			Value: map[string]interface{}{
				"hello": "world",
				"$type": "Microsoft.Foo/Bar",
			},
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := monitor.ValidateMetricAlertRuleTags(tc.Value, "azurerm_metric_alert_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q to return %d errors but returned %d", tc.Name, tc.ErrCount, len(errors))
		}
	}
}

func TestAccAzureRMMetricAlertRule_virtualMachineCpu(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_metric_alertrule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMetricAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMetricAlertRule_virtualMachineCpu(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMetricAlertRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "tags.$type"),
				),
			},
			{
				Config: testAccAzureRMMetricAlertRule_virtualMachineCpu(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMetricAlertRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "tags.$type"),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMetricAlertRuleExists(data.ResourceName),
					resource.TestCheckNoResourceAttr(data.ResourceName, "tags.$type"),
				),
			},
		},
	})
}

func TestAccAzureRMMetricAlertRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_metric_alertrule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMetricAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMetricAlertRule_virtualMachineCpu(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMetricAlertRuleExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMetricAlertRule_requiresImport(data, true),
				ExpectError: acceptance.RequiresImportError("azurerm_metric_alertrule"),
			},
		},
	})
}

func TestAccAzureRMMetricAlertRule_sqlDatabaseStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_metric_alertrule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMetricAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMetricAlertRule_sqlDatabaseStorage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMetricAlertRuleExists(data.ResourceName),
					resource.TestCheckNoResourceAttr(data.ResourceName, "tags.$type"),
				),
			},
		},
	})
}

func testCheckAzureRMMetricAlertRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.AlertRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Alert Rule: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Alert Rule %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on monitorAlertRulesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMetricAlertRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.AlertRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_metric_alertrule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Alert Rule still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMMetricAlertRule_virtualMachineCpu(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "50"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}

resource "azurerm_metric_alertrule" "test" {
  name                = "${azurerm_virtual_machine.test.name}-cpu"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  description = "An alert rule to watch the metric Percentage CPU"

  enabled = %t

  resource_id = "${azurerm_virtual_machine.test.id}"
  metric_name = "Percentage CPU"
  operator    = "GreaterThan"
  threshold   = 75
  aggregation = "Average"
  period      = "PT5M"

  email_action {
    send_to_service_owners = false

    custom_emails = [
      "support@azure.microsoft.com",
    ]
  }

  webhook_action {
    service_uri = "https://requestb.in/18jamc41"

    properties = {
      severity        = "incredible"
      acceptance_test = "true"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, enabled)
}

func testAccAzureRMMetricAlertRule_requiresImport(data acceptance.TestData, enabled bool) string {
	template := testAccAzureRMMetricAlertRule_virtualMachineCpu(data, enabled)
	return fmt.Sprintf(`
%s

resource "azurerm_metric_alertrule" "import" {
  name                = "${azurerm_metric_alertrule.test.name}"
  resource_group_name = "${azurerm_metric_alertrule.test.resource_group_name}"
  location            = "${azurerm_metric_alertrule.test.location}"
  description         = "${azurerm_metric_alertrule.test.description}"
  enabled             = "${azurerm_metric_alertrule.test.enabled}"

  resource_id = "${azurerm_virtual_machine.test.id}"
  metric_name = "Percentage CPU"
  operator    = "GreaterThan"
  threshold   = 75
  aggregation = "Average"
  period      = "PT5M"

  email_action {
    send_to_service_owners = false

    custom_emails = [
      "support@azure.microsoft.com",
    ]
  }

  webhook_action {
    service_uri = "https://requestb.in/18jamc41"

    properties = {
      severity        = "incredible"
      acceptance_test = "true"
    }
  }
}
`, template)
}

func testAccAzureRMMetricAlertRule_sqlDatabaseStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}


resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%[1]d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%[1]d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%[1]d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "50"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hn%[1]d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}

resource "azurerm_metric_alertrule" "test" {
  name                = "${azurerm_sql_database.test.name}-storage"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  description = "An alert rule to watch the metric Storage"

  enabled = true

  resource_id = "${azurerm_sql_database.test.id}"
  metric_name = "storage"
  operator    = "GreaterThan"
  threshold   = 1073741824
  aggregation = "Maximum"
  period      = "PT10M"

  email_action {
    send_to_service_owners = false

    custom_emails = [
      "support@azure.microsoft.com",
    ]
  }

  webhook_action {
    service_uri = "https://requestb.in/18jamc41"

    properties = {
      severity        = "incredible"
      acceptance_test = "true"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
