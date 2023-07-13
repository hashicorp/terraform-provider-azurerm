// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorOriginResource struct {
}

func TestAccCdnFrontDoorOrigin_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
	r := CdnFrontDoorOriginResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorOrigin_basicThreePointOh(t *testing.T) {
	if !features.FourPointOhBeta() {
		data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
		r := CdnFrontDoorOriginResource{}
		data.ResourceTest(t, r, []acceptance.TestStep{
			{
				Config: r.basicThreePointOh(data),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
		})
	} else {
		t.Skip("Test no longer valid due to deprecation of the 'health_probes_enabled' field in the 4.x version of the provider")
	}
}

func TestAccCdnFrontDoorOrigin_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
	r := CdnFrontDoorOriginResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCdnFrontDoorOrigin_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
	r := CdnFrontDoorOriginResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorOrigin_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
	r := CdnFrontDoorOriginResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorOrigin_privateLinkBlobPrimary(t *testing.T) {
	t.Skip("@tombuildsstuff: temporarily skipping until the private link is manually approved as part of the test step")

	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
	r := CdnFrontDoorOriginResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateLinkBlobPrimary(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// TODO: approve the connection by looking up and updating the private link
				// data.CheckWithClient(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
				//	clients.Network.PrivateLinkServiceClient.UpdatePrivateEndpointConnection()
				// }),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorOrigin_privateLinkStorageStaticWebSite(t *testing.T) {
	t.Skip("@tombuildsstuff: temporarily skipping until the private link is manually approved as part of the test step")

	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
	r := CdnFrontDoorOriginResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateLinkStaticWebSite(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// TODO: approve the connection by looking up and updating the private link
				// data.CheckWithClient(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
				//	clients.Network.PrivateLinkServiceClient.UpdatePrivateEndpointConnection()
				// }),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorOrigin_privateLinkAppServices(t *testing.T) {
	t.Skip("@tombuildsstuff: temporarily skipping until the private link is manually approved as part of the test step")

	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
	r := CdnFrontDoorOriginResource{}
	// NOTE: The Private Link will not be approved at this point but it will
	// be created. There is currently no way to automate the approval process.
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateLinkAppServices(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// TODO: approve the connection by looking up and updating the private link
				// data.CheckWithClient(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
				//	clients.Network.PrivateLinkServiceClient.UpdatePrivateEndpointConnection()
				// }),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorOrigin_privateLinkLoadBalancer(t *testing.T) {
	t.Skip("@tombuildsstuff: temporarily skipping until the private link is manually approved as part of the test step")

	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
	r := CdnFrontDoorOriginResource{}

	// NOTE: The Private Link will not be approved at this point but it will
	// be created. There is currently no way to automate the approval process.
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateLinkLoadBalancer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// TODO: approve the connection by looking up and updating the private link
				// data.CheckWithClient(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
				//	clients.Network.PrivateLinkServiceClient.UpdatePrivateEndpointConnection()
				// }),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorOrigin_removeOriginHostHeaderName(t *testing.T) {
	// regression test case for issue 20617
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
	r := CdnFrontDoorOriginResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeOriginHostHeader(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("origin_host_header").IsEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorOrigin_OriginHostHeaderRegression(t *testing.T) {
	// regression test case for issue 20866
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
	r := CdnFrontDoorOriginResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.OriginHostHeaderRegression(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("origin_host_header").HasValue("regression20866.australiaeast.cloudapp.azure.com"),
			),
		},
		data.ImportStep(),
		{
			Config: r.OriginHostHeaderRegressionUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("origin_host_header").HasValue("regression20866.australiaeast.cloudapp.azure.com"),
			),
		},
		data.ImportStep(),
		{
			Config: r.OriginHostHeaderRegression(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("origin_host_header").HasValue("regression20866.australiaeast.cloudapp.azure.com"),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeOriginHostHeader(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("origin_host_header").IsEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnFrontDoorOriginResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontDoorOriginID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorOriginsClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

//nolint:unused
func (r CdnFrontDoorOriginResource) templatePrivateLinkStorage(data acceptance.TestData) string {
	template := r.template(data, "Premium_AzureFrontDoor", false)
	return fmt.Sprintf(`

%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Premium"
  account_replication_type = "LRS"

  allow_nested_items_to_be_public = false

  network_rules {
    default_action = "Deny"
  }

  tags = {
    environment = "Test"
  }
}
`, template, data.RandomString)
}

// nolint: unused
func (r CdnFrontDoorOriginResource) templatePrivateLinkStorageStaticWebSite(data acceptance.TestData) string {
	template := r.template(data, "Premium_AzureFrontDoor", false)
	return fmt.Sprintf(`

%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"

  allow_nested_items_to_be_public = false

  network_rules {
    default_action = "Deny"
  }

  static_website {
    index_document     = "index.html"
    error_404_document = "404.html"
  }

  tags = {
    environment = "Test"
  }
}
`, template, data.RandomString)
}

// nolint: unused
func (r CdnFrontDoorOriginResource) templatePrivateLinkLoadBalancer(data acceptance.TestData) string {
	template := r.template(data, "Premium_AzureFrontDoor", true)
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                                          = "acctestsn-%[2]d"
  resource_group_name                           = azurerm_resource_group.test.name
  virtual_network_name                          = azurerm_virtual_network.test.name
  address_prefixes                              = ["10.5.1.0/24"]
  private_link_service_network_policies_enabled = false
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpi-%[2]d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[2]d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = azurerm_public_ip.test.name
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_private_link_service" "test" {
  name                = "acctestPLS-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  visibility_subscription_ids                 = [data.azurerm_client_config.current.subscription_id]
  load_balancer_frontend_ip_configuration_ids = [azurerm_lb.test.frontend_ip_configuration.0.id]

  nat_ip_configuration {
    name                       = "primary"
    private_ip_address         = "10.5.1.17"
    private_ip_address_version = "IPv4"
    subnet_id                  = azurerm_subnet.test.id
    primary                    = true
  }
}
`, template, data.RandomInteger)
}

//nolint:unused
func (r CdnFrontDoorOriginResource) templatePrivateLinkWebApp(data acceptance.TestData) string {
	template := r.template(data, "Premium_AzureFrontDoor", false)
	return fmt.Sprintf(`

%s

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "P1v3"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "testaccsc%[3]s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_share" "test" {
  name                 = "test"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 1
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2021-04-01"
  expiry = "2024-03-30"

  permissions {
    read    = false
    write   = true
    delete  = false
    list    = false
    add     = false
    create  = false
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}
`, template, data.RandomInteger, data.RandomString)
}

func (r CdnFrontDoorOriginResource) basic(data acceptance.TestData) string {
	template := r.template(data, "Standard_AzureFrontDoor", false)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-cdnfdorigin-%d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorOriginResource) basicThreePointOh(data acceptance.TestData) string {
	template := r.template(data, "Standard_AzureFrontDoor", false)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-cdnfdorigin-%d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id

  health_probes_enabled          = true
  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorOriginResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_origin" "import" {
  name                          = azurerm_cdn_frontdoor_origin.test.name
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}
`, config)
}

func (r CdnFrontDoorOriginResource) complete(data acceptance.TestData) string {
	template := r.template(data, "Standard_AzureFrontDoor", false)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-cdnfdorigin-%d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorOriginResource) update(data acceptance.TestData) string {
	template := r.template(data, "Standard_AzureFrontDoor", false)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-cdnfdorigin-%d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}
`, template, data.RandomInteger)
}

// nolint: unused
func (r CdnFrontDoorOriginResource) privateLinkBlobPrimary(data acceptance.TestData) string {
	template := r.templatePrivateLinkStorage(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-cdnfdorigin-%d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = true
  host_name                      = azurerm_storage_account.test.primary_blob_host
  origin_host_header             = azurerm_storage_account.test.primary_blob_host
  priority                       = 1
  weight                         = 500

  private_link {
    request_message        = "Request access for CDN Frontdoor Private Link Origin"
    target_type            = "blob"
    location               = azurerm_resource_group.test.location
    private_link_target_id = azurerm_storage_account.test.id
  }
}
`, template, data.RandomInteger)
}

// nolint: unused
func (r CdnFrontDoorOriginResource) privateLinkStaticWebSite(data acceptance.TestData) string {
	template := r.templatePrivateLinkStorageStaticWebSite(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-cdnfdorigin-%d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = true
  host_name                      = azurerm_storage_account.test.primary_web_host
  origin_host_header             = azurerm_storage_account.test.primary_web_host
  priority                       = 1
  weight                         = 500

  private_link {
    request_message        = "Request access for CDN Frontdoor Private Link Origin"
    target_type            = "web"
    location               = azurerm_resource_group.test.location
    private_link_target_id = azurerm_storage_account.test.id
  }
}
`, template, data.RandomInteger)
}

// nolint: unused
func (r CdnFrontDoorOriginResource) privateLinkAppServices(data acceptance.TestData) string {
	template := r.templatePrivateLinkWebApp(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-cdnfdorigin-%d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = true
  host_name                      = azurerm_linux_web_app.test.default_hostname
  origin_host_header             = azurerm_linux_web_app.test.default_hostname
  priority                       = 1
  weight                         = 500

  private_link {
    request_message        = "Request access for CDN Frontdoor Private Link Origin"
    target_type            = "sites"
    location               = azurerm_resource_group.test.location
    private_link_target_id = azurerm_linux_web_app.test.id
  }
}
`, template, data.RandomInteger)
}

// nolint: unused
func (r CdnFrontDoorOriginResource) privateLinkLoadBalancer(data acceptance.TestData) string {
	template := r.templatePrivateLinkLoadBalancer(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-cdnfdorigin-%d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = true
  host_name                      = azurerm_private_link_service.test.nat_ip_configuration.0.private_ip_address
  origin_host_header             = azurerm_private_link_service.test.nat_ip_configuration.0.private_ip_address
  priority                       = 1
  weight                         = 500

  private_link {
    request_message        = "Request access for CDN Frontdoor Private Link Origin"
    location               = azurerm_resource_group.test.location
    private_link_target_id = azurerm_private_link_service.test.id
  }
}
`, template, data.RandomInteger)
}

func (CdnFrontDoorOriginResource) template(data acceptance.TestData, profileSku string, isLoadBalancer bool) string {
	// NOTE: This is a hack (the private link service dependency in the profile resource) for what I believe is a bug in
	// the CDN Frontdoor API. I am currently speaking with the service team about how to correctly fix this issue,
	// but in the meantime this is what we need to do to get this scenario to work.
	var loadBalancerDependsOn string
	if isLoadBalancer {
		loadBalancerDependsOn = "depends_on = [azurerm_private_link_service.test]"
	}

	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-cdn-afdx-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  %s
  name                = "acctest-cdnfdprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = %q
}

resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name                     = "acctest-cdnfd-group-%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
}
`, data.RandomInteger, data.Locations.Primary, loadBalancerDependsOn, data.RandomInteger, profileSku, data.RandomInteger)
}

func (r CdnFrontDoorOriginResource) removeOriginHostHeader(data acceptance.TestData) string {
	template := r.template(data, "Standard_AzureFrontDoor", false)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-cdnfdorigin-%d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  priority                       = 1
  weight                         = 1
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorOriginResource) OriginHostHeaderRegression(data acceptance.TestData) string {
	template := r.template(data, "Standard_AzureFrontDoor", false)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-cdnfdorigin-%d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  origin_host_header             = "regression20866.australiaeast.cloudapp.azure.com"
  http_port                      = 80
  https_port                     = 443
  priority                       = 1
  weight                         = 1
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorOriginResource) OriginHostHeaderRegressionUpdate(data acceptance.TestData) string {
	template := r.template(data, "Standard_AzureFrontDoor", false)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-cdnfdorigin-%d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  origin_host_header             = "regression20866.australiaeast.cloudapp.azure.com"
  http_port                      = 80
  https_port                     = 443
  priority                       = 5
  weight                         = 1
}
`, template, data.RandomInteger)
}
