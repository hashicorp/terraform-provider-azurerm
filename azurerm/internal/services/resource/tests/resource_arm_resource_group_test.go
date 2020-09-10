package tests

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/acctest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestMain(m *testing.M) {
	acctest.UseBinaryDriver("azurerm", provider.AzureProvider)
	acctest.UseBinaryDriver("azurerm2", provider.AzureProvider)
	resource.TestMain(m)
}

func TestAccAzureRMResourceGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.SupportedProviders,
		CheckDestroy:      testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMResourceGroup_multipleSubscriptions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.SupportedProviders,
		CheckDestroy:      testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_multipleSubscriptions(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(data.ResourceName),
				),
			},
			// data.ImportStep(),
		},
	})
}

func TestAccAzureRMResourceGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.SupportedProviders,
		CheckDestroy:      testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMResourceGroup_requiresImport),
		},
	})
}

func TestAccAzureRMResourceGroup_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.SupportedProviders,
		CheckDestroy:      testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			data.DisappearsStep(acceptance.DisappearsStepData{
				Config:      testAccAzureRMResourceGroup_basic,
				CheckExists: testCheckAzureRMResourceGroupExists,
				Destroy:     testCheckAzureRMResourceGroupDisappears,
			}),
		},
	})
}

func TestAccAzureRMResourceGroup_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.SupportedProviders,
		CheckDestroy:      testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMResourceGroup_withTagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMResourceGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.GroupsClient
		testClient, err := configureTestProvider()
		if err != nil {
			panic(err)
		}
		ctx := context.TODO()

		client := testClient.Resource.GroupsClient
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API

		resp, err := client.Get(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("Bad: Get on resourceGroupClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: resource group: %q does not exist", resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMResourceGroupDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient, err := configureTestProvider()
		if err != nil {
			panic(err)
		}
		ctx := context.TODO()

		client := testClient.Resource.GroupsClient

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API

		deleteFuture, err := client.Delete(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("Failed deleting Resource Group %q: %+v", resourceGroup, err)
		}

		err = deleteFuture.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Failed long polling for the deletion of Resource Group %q: %+v", resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMResourceGroupDestroy(s *terraform.State) error {
	testClient, err := configureTestProvider()
	if err != nil {
		panic(err)
	}
	ctx := context.TODO()

	client := testClient.Resource.GroupsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_resource_group" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Resource Group still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func configureTestProvider() (*clients.Client, error) {
	var auxTenants []string

	if v := os.Getenv("ARM_AUXILIARY_TENANT_IDS"); v != "" {
		auxTenants = strings.Split(v, ";")
	}

	if len(auxTenants) > 3 {
		return nil, fmt.Errorf("The provider only supports 3 auxiliary tenant IDs")
	}

	metadataHost := os.Getenv("ARM_METADATA_HOSTNAME")
	// TODO: remove in 3.0
	// note: this is inline to avoid calling out deprecations for users not setting this
	if v := os.Getenv("ARM_METADATA_URL"); v != "" {
		metadataHost = v
	}

	var useMSI bool
	if os.Getenv("ARM_USE_MSI") != "" {
		useMSI, _ = strconv.ParseBool(os.Getenv("ARM_USE_MSI"));
	}

	builder := &authentication.Builder{
		SubscriptionID:     os.Getenv("ARM_SUBSCRIPTION_ID"),
		ClientID:           os.Getenv("ARM_CLIENT_ID"),
		ClientSecret:       os.Getenv("ARM_CLIENT_SECRET"),
		TenantID:           os.Getenv("ARM_TENANT_ID"),
		AuxiliaryTenantIDs: auxTenants,
		Environment:        os.Getenv("ARM_ENVIRONMENT"),
		MetadataURL:        metadataHost, // TODO: rename this in Helpers too
		MsiEndpoint:        os.Getenv("ARM_MSI_ENDPOINT"),
		ClientCertPassword:   os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"),
		ClientCertPath:      os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"),

		// Feature Toggles
		SupportsClientCertAuth:         true,
		SupportsClientSecretAuth:       true,
		SupportsManagedServiceIdentity: useMSI,
		SupportsAzureCliToken:          true,
		SupportsAuxiliaryTenants:       len(auxTenants) > 0,

		// Doc Links
		ClientSecretDocsLink: "https://www.terraform.io/docs/providers/azurerm/guides/service_principal_client_secret.html",
	}

	config, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("Error building AzureRM Client: %s", err)
	}


	skipProviderRegistration := false
	clientBuilder := clients.ClientBuilder{
		AuthConfig:                  config,
		SkipProviderRegistration:    skipProviderRegistration,
		TerraformVersion:            "0.12.29",
		PartnerId:                  os.Getenv("ARM_PARTNER_ID"),
		//DisableCorrelationRequestID: d.Get("disable_correlation_request_id").(bool),
		//DisableTerraformPartnerID:   d.Get("disable_terraform_partner_id").(bool),
		// Features:                    expandFeatures(d.Get("features").([]interface{})),
		StorageUseAzureAD:           false,
	}
	client, err := clients.Build(context.TODO(), clientBuilder)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func testAccAzureRMResourceGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMResourceGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMResourceGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "import" {
  name     = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
}
`, template)
}

func testAccAzureRMResourceGroup_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMResourceGroup_withTagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMResourceGroup_multipleSubscriptions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azurerm2" {
  features {}
  subscription_id = "%s"
}

resource "azurerm_resource_group" "test" {
  provider = azurerm
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test1" {
  provider = azurerm2
  name     = "acctestRG-%d"
  location = "%s"
}
`, os.Getenv("ARM_SUBSCRIPTION_ID_ALT"), data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}
