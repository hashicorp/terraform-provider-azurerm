package deviceregistry_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assetendpointprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AssetEndpointProfileTestResource struct{}

const (
	ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME = "ARM_DEVICE_REGISTRY_CUSTOM_LOCATION"
	ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME  = "ARM_DEVICE_REGISTRY_RESOURCE_GROUP"
)

func TestAccAssetEndpointProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	if os.Getenv(ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME) == "" || os.Getenv(ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME, ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue(""),
				check.That(data.ResourceName).Key("authentication_method").HasValue(""),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetEndpointProfile_complete_certificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	if os.Getenv(ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME) == "" || os.Getenv(ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME, ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("Certificate"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue("myCertificateRef"),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetEndpointProfile_complete_usernamePassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	if os.Getenv(ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME) == "" || os.Getenv(ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME, ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeUsernamePassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("UsernamePassword"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue("myUsernameRef"),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue("myPasswordRef"),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetEndpointProfile_complete_anonymous(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	if os.Getenv(ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME) == "" || os.Getenv(ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME, ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeAnonymous(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("Anonymous"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetEndpointProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	if os.Getenv(ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME) == "" || os.Getenv(ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME, ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME)
	}

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

func TestAccAssetEndpointProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	if os.Getenv(ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME) == "" || os.Getenv(ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME, ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{ // first create the resource
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{ // update the authentication method to certificate
			Config: r.completeCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("Certificate"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue("myCertificateRef"),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
			),
		},
		data.ImportStep(),
		{ // update the authentication method to username/password
			Config: r.completeUsernamePassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("UsernamePassword"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue("myUsernameRef"),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue("myPasswordRef"),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
			),
		},
		data.ImportStep(),
		{ // update the authentication method to anonymous
			Config: r.completeAnonymous(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("Anonymous"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
			),
		},
		data.ImportStep(),
	})
}

func (AssetEndpointProfileTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := assetendpointprofiles.ParseAssetEndpointProfileID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DeviceRegistry.AssetEndpointProfileClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r AssetEndpointProfileTestResource) basic(data acceptance.TestData) string {
	template := r.template()
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
  name                                  = "acctest-assetendpointprofile-%[2]d"
  resource_group_name                   = local.resource_group_name
  extended_location_name                = local.custom_location_name
  extended_location_type                = "CustomLocation"
  target_address                        = "opc.tcp://foo"
  endpoint_profile_type                 = "OpcUa"
  discovered_asset_endpoint_profile_ref = "discoveredAssetEndpointProfile123"
  location                              = "%[3]s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) completeCertificate(data acceptance.TestData) string {
	template := r.template()
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
  name                                     = "acctest-assetendpointprofile-%[2]d"
  resource_group_name                      = local.resource_group_name
  extended_location_name                   = local.custom_location_name
  extended_location_type                   = "CustomLocation"
  target_address                           = "opc.tcp://foo"
  endpoint_profile_type                    = "OpcUa"
  discovered_asset_endpoint_profile_ref    = "discoveredAssetEndpointProfile123"
  additional_configuration                 = "{\"foo\": \"bar\"}"
  authentication_method                    = "Certificate"
  x509_credentials_certificate_secret_name = "myCertificateRef"
  tags = {
    "sensor" = "temperature,humidity"
  }
  location = "%[3]s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) completeUsernamePassword(data acceptance.TestData) string {
	template := r.template()
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
  name                                               = "acctest-assetendpointprofile-%[2]d"
  resource_group_name                                = local.resource_group_name
  extended_location_name                             = local.custom_location_name
  extended_location_type                             = "CustomLocation"
  target_address                                     = "opc.tcp://foo"
  endpoint_profile_type                              = "OpcUa"
  discovered_asset_endpoint_profile_ref              = "discoveredAssetEndpointProfile123"
  additional_configuration                           = "{\"foo\": \"bar\"}"
  authentication_method                              = "UsernamePassword"
  username_password_credentials_username_secret_name = "myUsernameRef"
  username_password_credentials_password_secret_name = "myPasswordRef"
  tags = {
    "sensor" = "temperature,humidity"
  }
  location = "%[3]s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) completeAnonymous(data acceptance.TestData) string {
	template := r.template()
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
  name                                  = "acctest-assetendpointprofile-%[2]d"
  resource_group_name                   = local.resource_group_name
  extended_location_name                = local.custom_location_name
  extended_location_type                = "CustomLocation"
  target_address                        = "opc.tcp://foo"
  endpoint_profile_type                 = "OpcUa"
  discovered_asset_endpoint_profile_ref = "discoveredAssetEndpointProfile123"
  additional_configuration              = "{\"foo\": \"bar\"}"
  authentication_method                 = "Anonymous"
  tags = {
    "sensor" = "temperature,humidity"
  }
  location = "%[3]s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "import" {
  name                                  = azurerm_device_registry_asset_endpoint_profile.test.name
  resource_group_name                   = azurerm_device_registry_asset_endpoint_profile.test.resource_group_name
  extended_location_name                = azurerm_device_registry_asset_endpoint_profile.test.extended_location_name
  extended_location_type                = azurerm_device_registry_asset_endpoint_profile.test.extended_location_type
  target_address                        = azurerm_device_registry_asset_endpoint_profile.test.target_address
  endpoint_profile_type                 = azurerm_device_registry_asset_endpoint_profile.test.endpoint_profile_type
  discovered_asset_endpoint_profile_ref = azurerm_device_registry_asset_endpoint_profile.test.discovered_asset_endpoint_profile_ref
  location                              = azurerm_device_registry_asset_endpoint_profile.test.location
}


`, template)
}

/*
Creates the terraform template for AzureRm provider and needed constants
*/
func (AssetEndpointProfileTestResource) template() string {
	customLocation := os.Getenv(ASSET_ENDPOINT_PROFILE_CUSTOM_LOCATION_NAME)
	resourceGroup := os.Getenv(ASSET_ENDPOINT_PROFILE_RESOURCE_GROUP_NAME)

	return fmt.Sprintf(`
locals {
  custom_location_name = "%[1]s"
  resource_group_name  = "%[2]s"
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}
`, customLocation, resourceGroup)
}
