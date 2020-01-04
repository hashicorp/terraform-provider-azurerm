package tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers"
)

func TestAccAzureRMContainerService_orchestrationPlatformValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{Value: "DCOS", ErrCount: 0},
		{Value: "Kubernetes", ErrCount: 0},
		{Value: "Swarm", ErrCount: 0},
		{Value: "Mesos", ErrCount: 1},
	}

	for _, tc := range cases {
		_, errors := containers.ValidateArmContainerServiceOrchestrationPlatform(tc.Value, "azurerm_container_service")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Container Service Orchestration Platform to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAccAzureRMContainerService_masterProfileCountValidation(t *testing.T) {
	cases := []struct {
		Value    int
		ErrCount int
	}{
		{Value: 0, ErrCount: 1},
		{Value: 1, ErrCount: 0},
		{Value: 2, ErrCount: 1},
		{Value: 3, ErrCount: 0},
		{Value: 4, ErrCount: 1},
		{Value: 5, ErrCount: 0},
		{Value: 6, ErrCount: 1},
	}

	for _, tc := range cases {
		_, errors := containers.ValidateArmContainerServiceMasterProfileCount(tc.Value, "azurerm_container_service")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Container Service Master Profile Count to trigger a validation error for '%d'", tc.Value)
		}
	}
}

func TestAccAzureRMContainerService_agentProfilePoolCountValidation(t *testing.T) {
	cases := []struct {
		Value    int
		ErrCount int
	}{
		{Value: 0, ErrCount: 1},
		{Value: 1, ErrCount: 0},
		{Value: 2, ErrCount: 0},
		{Value: 99, ErrCount: 0},
		{Value: 100, ErrCount: 0},
		{Value: 101, ErrCount: 1},
	}

	for _, tc := range cases {
		_, errors := containers.ValidateArmContainerServiceAgentPoolProfileCount(tc.Value, "azurerm_container_service")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Container Service Agent Pool Profile Count to trigger a validation error for '%d'", tc.Value)
		}
	}
}

func TestAccAzureRMContainerService_dcosBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerService_dcosBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerServiceExists("azurerm_container_service.test"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerService_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_container_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerService_dcosBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerServiceExists("azurerm_container_service.test"),
				),
			},
			{
				Config:      testAccAzureRMContainerService_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_container_service"),
			},
		},
	})
}

func TestAccAzureRMContainerService_kubernetesBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_service", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerService_kubernetesBasic(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerServiceExists("azurerm_container_service.test"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerService_kubernetesComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_service", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerService_kubernetesComplete(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerServiceExists("azurerm_container_service.test"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerService_swarmBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerService_swarmBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerServiceExists("azurerm_container_service.test"),
				),
			},
		},
	})
}

func testAccAzureRMContainerService_dcosBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_service" "test" {
  name                   = "acctestcontservice%d"
  location               = "${azurerm_resource_group.test.location}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  orchestration_platform = "DCOS"

  master_profile {
    count      = 1
    dns_prefix = "acctestmaster%d"
  }

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  agent_pool_profile {
    name       = "default"
    count      = 1
    dns_prefix = "acctestagent%d"
    vm_size    = "Standard_F2"
  }

  diagnostics_profile {
    enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMContainerService_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMContainerService_dcosBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_service" "import" {
  name                   = "${azurerm_container_service.test.name}"
  location               = "${azurerm_container_service.test.location}"
  resource_group_name    = "${azurerm_container_service.test.resource_group_name}"
  orchestration_platform = "DCOS"

  master_profile {
    count      = 1
    dns_prefix = "acctestmaster%d"
  }

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  agent_pool_profile {
    name       = "default"
    count      = 1
    dns_prefix = "acctestagent%d"
    vm_size    = "Standard_F2"
  }

  diagnostics_profile {
    enabled = false
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMContainerService_kubernetesBasic(data acceptance.TestData, clientId string, clientSecret string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_service" "test" {
  name                   = "acctestcontservice%d"
  location               = "${azurerm_resource_group.test.location}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  orchestration_platform = "Kubernetes"

  master_profile {
    count      = 1
    dns_prefix = "acctestmaster%d"
  }

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  agent_pool_profile {
    name       = "default"
    count      = 1
    dns_prefix = "acctestagent%d"
    vm_size    = "Standard_F2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  diagnostics_profile {
    enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMContainerService_kubernetesComplete(data acceptance.TestData, clientId string, clientSecret string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_service" "test" {
  name                   = "acctestcontservice%d"
  location               = "${azurerm_resource_group.test.location}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  orchestration_platform = "Kubernetes"

  master_profile {
    count      = 1
    dns_prefix = "acctestmaster%d"
  }

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  agent_pool_profile {
    name       = "default"
    count      = 1
    dns_prefix = "acctestagent%d"
    vm_size    = "Standard_F2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  diagnostics_profile {
    enabled = false
  }

  tags = {
    you = "me"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMContainerService_swarmBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_service" "test" {
  name                   = "acctestcontservice%d"
  location               = "${azurerm_resource_group.test.location}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  orchestration_platform = "Swarm"

  master_profile {
    count      = 1
    dns_prefix = "acctestmaster%d"
  }

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  agent_pool_profile {
    name       = "default"
    count      = 1
    dns_prefix = "acctestagent%d"
    vm_size    = "Standard_F2"
  }

  diagnostics_profile {
    enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testCheckAzureRMContainerServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Container Service Instance: %s", name)
		}

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Containers.ServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on containerServicesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Container Service Instance %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMContainerServiceDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Containers.ServicesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_container_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Container Service Instance still exists:\n%#v", resp)
		}
	}

	return nil
}
