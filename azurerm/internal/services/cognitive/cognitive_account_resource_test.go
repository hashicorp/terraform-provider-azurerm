package cognitive_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type CognitiveAccountResource struct {
}

func TestAccCognitiveAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("Face"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccount_speechServices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.speechServices(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("SpeechServices"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccount_speechServicesWithStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.speechServicesWithStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_cognitive_account"),
		},
	})
}

func TestAccCognitiveAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("Face"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Acceptance").HasValue("Test"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("Face"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("Face"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Acceptance").HasValue("Test"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
	})
}

func TestAccCognitiveAccount_qnaRuntimeEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.qnaRuntimeEndpoint(data, "https://localhost:8080/"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("QnAMaker"),
				check.That(data.ResourceName).Key("qna_runtime_endpoint").HasValue("https://localhost:8080/"),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.qnaRuntimeEndpoint(data, "https://localhost:9000/"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("QnAMaker"),
				check.That(data.ResourceName).Key("qna_runtime_endpoint").HasValue("https://localhost:9000/"),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccount_qnaRuntimeEndpointUnspecified(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.qnaRuntimeEndpointUnspecified(data),
			ExpectError: regexp.MustCompile("the QnAMaker runtime endpoint `qna_runtime_endpoint` is required when kind is set to `QnAMaker`"),
		},
	})
}

func TestAccCognitiveAccount_cognitiveServices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cognitiveServices(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccount_withMultipleCognitiveAccounts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withMultipleCognitiveAccounts(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccount_networkAclsVirtualNetworkRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkAclsVirtualNetworkRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkAclsVirtualNetworkRulesUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccount_networkAcls(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkAcls(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkAclsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccount_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").MatchesRegex(validate.UUIDRegExp),
				check.That(data.ResourceName).Key("identity.0.tenant_id").MatchesRegex(validate.UUIDRegExp),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").MatchesRegex(validate.UUIDRegExp),
				check.That(data.ResourceName).Key("identity.0.tenant_id").MatchesRegex(validate.UUIDRegExp),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccount_metricsAdvisor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")
	r := CognitiveAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.metricsAdvisor(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t CognitiveAccountResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cognitive.AccountsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Cognitive Services Account %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (CognitiveAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Face"
  sku_name            = "S0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CognitiveAccountResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Face"
  sku_name            = "S0"
  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CognitiveAccountResource) identityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Face"
  sku_name            = "S0"
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (CognitiveAccountResource) identitySystemAssignedUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Face"
  sku_name            = "S0"
  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (CognitiveAccountResource) speechServices(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "SpeechServices"
  sku_name            = "S0"
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

func (CognitiveAccountResource) speechServicesWithStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestrg%d"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-identity-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "SpeechServices"
  sku_name            = "S0"

  identity {
    type = "SystemAssigned"
  }

  storage {
    storage_account_id = azurerm_storage_account.test.id
    identity_client_id = azurerm_user_assigned_identity.test.client_id
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomIntOfLength(8), data.RandomInteger, data.RandomInteger)
}

func (CognitiveAccountResource) requiresImport(data acceptance.TestData) string {
	template := CognitiveAccountResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account" "import" {
  name                = azurerm_cognitive_account.test.name
  location            = azurerm_cognitive_account.test.location
  resource_group_name = azurerm_cognitive_account.test.resource_group_name
  kind                = azurerm_cognitive_account.test.kind
  sku_name            = "S0"
}
`, template)
}

func (CognitiveAccountResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Face"
  sku_name            = "S0"

  fqdns                             = ["foo.com", "bar.com"]
  public_network_access_enabled     = false
  outbound_network_access_restrited = true
  local_auth_enabled                = false

  tags = {
    Acceptance = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CognitiveAccountResource) qnaRuntimeEndpoint(data acceptance.TestData, url string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                 = "acctestcogacc-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  kind                 = "QnAMaker"
  qna_runtime_endpoint = "%s"
  sku_name             = "S0"
}
`, data.RandomInteger, "West US", data.RandomInteger, url) // QnAMaker only available in West US
}

func (CognitiveAccountResource) qnaRuntimeEndpointUnspecified(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "QnAMaker"
  sku_name            = "S0"
}
`, data.RandomInteger, "West US", data.RandomInteger) // QnAMaker only available in West US
}

func (CognitiveAccountResource) cognitiveServices(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "CognitiveServices"
  sku_name            = "S0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CognitiveAccountResource) metricsAdvisor(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}
resource "azurerm_cognitive_account" "test" {
  name                            = "acctestcogacc-%d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  kind                            = "MetricsAdvisor"
  sku_name                        = "S0"
  custom_subdomain_name           = "acctestcogacc-%d"
  metrics_advisor_aad_client_id   = "310d7b2e-d1d1-4b87-9807-5b885b290c00"
  metrics_advisor_aad_tenant_id   = "72f988bf-86f1-41af-91ab-2d7cd011db47"
  metrics_advisor_super_user_name = "mock_user1"
  metrics_advisor_website_name    = "mock_name2"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (CognitiveAccountResource) withMultipleCognitiveAccounts(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "CustomVision.Prediction"
  sku_name            = "S0"
}

resource "azurerm_cognitive_account" "test2" {
  name                = "acctestcogacc2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "CustomVision.Training"
  sku_name            = "S0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CognitiveAccountResource) networkAcls(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account" "test" {
  name                  = "acctestcogacc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  kind                  = "Face"
  sku_name              = "S0"
  custom_subdomain_name = "acctestcogacc-%d"

  network_acls {
    default_action             = "Deny"
    virtual_network_subnet_ids = [azurerm_subnet.test_a.id, azurerm_subnet.test_b.id]
  }
}
`, r.networkAclsTemplate(data), data.RandomInteger, data.RandomInteger)
}

func (r CognitiveAccountResource) networkAclsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_cognitive_account" "test" {
  name                  = "acctestcogacc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  kind                  = "Face"
  sku_name              = "S0"
  custom_subdomain_name = "acctestcogacc-%d"

  network_acls {
    default_action             = "Allow"
    ip_rules                   = ["123.0.0.101"]
    virtual_network_subnet_ids = [azurerm_subnet.test_a.id]
  }
}
`, r.networkAclsTemplate(data), data.RandomInteger, data.RandomInteger)
}

func (r CognitiveAccountResource) networkAclsVirtualNetworkRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account" "test" {
  name                  = "acctestcogacc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  kind                  = "Face"
  sku_name              = "S0"
  custom_subdomain_name = "acctestcogacc-%d"

  network_acls {
    default_action = "Deny"
    virtual_network_rules {
      subnet_id = azurerm_subnet.test_a.id
    }
    virtual_network_rules {
      subnet_id                            = azurerm_subnet.test_b.id
      ignore_missing_vnet_service_endpoint = true
    }

  }
}
`, r.networkAclsTemplate(data), data.RandomInteger, data.RandomInteger)
}

func (r CognitiveAccountResource) networkAclsVirtualNetworkRulesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_cognitive_account" "test" {
  name                  = "acctestcogacc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  kind                  = "Face"
  sku_name              = "S0"
  custom_subdomain_name = "acctestcogacc-%d"

  network_acls {
    default_action = "Allow"
    ip_rules       = ["123.0.0.101"]
    virtual_network_rules {
      subnet_id                            = azurerm_subnet.test_a.id
      ignore_missing_vnet_service_endpoint = true
    }
  }
}
`, r.networkAclsTemplate(data), data.RandomInteger, data.RandomInteger)
}

func (CognitiveAccountResource) networkAclsTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test_a" {
  name                 = "acctestsubneta%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.CognitiveServices"]
}

resource "azurerm_subnet" "test_b" {
  name                 = "acctestsubnetb%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.4.0/24"
  service_endpoints    = ["Microsoft.CognitiveServices"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
