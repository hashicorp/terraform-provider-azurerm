// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicefabric_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabric/2021-06-01/cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServiceFabricClusterResource struct{}

func TestAccAzureRMServiceFabricCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_basicNodeTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.is_primary").HasValue("true"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config: r.basicNodeTypeUpdate(data, 3, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("2"),
				check.That(data.ResourceName).Key("node_type.0.is_primary").HasValue("true"),
				check.That(data.ResourceName).Key("node_type.1.is_primary").HasValue("false"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMServiceFabricCluster_manualClusterCodeVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	codeVersion := "10.1.1541.9590"
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualClusterCodeVersion(data, codeVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_mode").HasValue("Manual"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "cluster_code_version", codeVersion),
			),
		},
		{
			Config: r.manualClusterCodeVersion(data, codeVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_mode").HasValue("Manual"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "cluster_code_version", codeVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_manualLatest(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualClusterCodeVersion(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_mode").HasValue("Manual"),
				check.That(data.ResourceName).Key("cluster_code_version").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_addOnFeatures(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.addOnFeatures(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_certificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.certificates(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.reverseProxyCertificates(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
				check.That(data.ResourceName).Key("node_type.0.reverse_proxy_endpoint_port").HasValue("19081"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyNotSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("node_type.0.reverse_proxy_endpoint_port").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config: r.reverseProxyCertificates(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
				check.That(data.ResourceName).Key("node_type.0.reverse_proxy_endpoint_port").HasValue("19081"),
			),
		},
		{
			Config: r.basic(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("node_type.0.reverse_proxy_endpoint_port").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_clientCertificateThumbprint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clientCertificateThumbprint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("1"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.is_admin").HasValue("true"),
				check.That(data.ResourceName).Key("client_certificate_common_name.#").HasValue("0"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_withMultipleClientCertificateThumbprints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withMultipleClientCertificateThumbprints(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_clientCertificateCommonNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clientCertificateCommonNames(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_certificate_common_name.#").HasValue("2"),
				check.That(data.ResourceName).Key("client_certificate_common_name.0.common_name").HasValue("firstcertcommonname"),
				check.That(data.ResourceName).Key("client_certificate_common_name.0.is_admin").HasValue("true"),
				check.That(data.ResourceName).Key("client_certificate_common_name.0.issuer_thumbprint").HasValue("3341db6cf2af72c611df3be3721a653af1d43ecd50f584f828793dbe9103c3ee"),
				check.That(data.ResourceName).Key("client_certificate_common_name.1.common_name").HasValue("secondcertcommonname"),
				check.That(data.ResourceName).Key("client_certificate_common_name.1.is_admin").HasValue("false"),
				check.That(data.ResourceName).Key("client_certificate_common_name.1.issuer_thumbprint").IsEmpty(),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("1"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.is_admin").HasValue("true"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_readerAdminClientCertificateThumbprint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.readerAdminClientCertificateThumbprint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("2"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.is_admin").HasValue("true"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.1.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.1.is_admin").HasValue("false"),
				check.That(data.ResourceName).Key("client_certificate_common_name.#").HasValue("0"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_certificateCommonNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.certificateCommonNames(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate_common_names.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyCertificateCommonNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.reverseProxyCertificateCommonNames(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("reverse_proxy_certificate_common_names.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_azureActiveDirectory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureActiveDirectory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("azure_active_directory.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("azure_active_directory.0.cluster_application_id").Exists(),
				check.That(data.ResourceName).Key("azure_active_directory.0.client_application_id").Exists(),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:19080"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_azureActiveDirectoryDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureActiveDirectory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("azure_active_directory.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("azure_active_directory.0.cluster_application_id").Exists(),
				check.That(data.ResourceName).Key("azure_active_directory.0.client_application_id").Exists(),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:19080"),
			),
		},
		{
			Config: r.azureActiveDirectoryDelete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:19080"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_diagnosticsConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.diagnosticsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("diagnostics_config.0.storage_account_name").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.protected_account_key_name").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.blob_endpoint").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.queue_endpoint").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.table_endpoint").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_diagnosticsConfigDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.diagnosticsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("diagnostics_config.0.storage_account_name").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.protected_account_key_name").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.blob_endpoint").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.queue_endpoint").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.table_endpoint").Exists(),
			),
		},
		{
			Config: r.diagnosticsConfigDelete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_fabricSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fabricSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fabric_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.%").HasValue("1"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("None"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_fabricSettingsRemove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fabricSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fabric_settings.#").HasValue("1"),
			),
		},
		{
			Config: r.basic(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fabric_settings.#").HasValue("0"),
			),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypeCustomPorts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodeTypeCustomPorts(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.application_ports.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.application_ports.0.start_port").HasValue("20000"),
				check.That(data.ResourceName).Key("node_type.0.application_ports.0.end_port").HasValue("29999"),
				check.That(data.ResourceName).Key("node_type.0.ephemeral_ports.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.ephemeral_ports.0.start_port").HasValue("30000"),
				check.That(data.ResourceName).Key("node_type.0.ephemeral_ports.0.end_port").HasValue("39999"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypesMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodeTypeMultiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.#").HasValue("2"),
				check.That(data.ResourceName).Key("node_type.0.name").HasValue("first"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("node_type.0.is_primary").HasValue("true"),
				check.That(data.ResourceName).Key("node_type.1.name").HasValue("second"),
				check.That(data.ResourceName).Key("node_type.1.instance_count").HasValue("4"),
				check.That(data.ResourceName).Key("node_type.1.is_primary").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
			),
		},
		{
			Config: r.basic(data, 4),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("4"),
			),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypeProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodeTypeProperties(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.0.placement_properties.%").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.placement_properties.HasSSD").HasValue("true"),
				check.That(data.ResourceName).Key("node_type.0.capacities.%").HasValue("2"),
				check.That(data.ResourceName).Key("node_type.0.capacities.ClientConnections").HasValue("20000"),
				check.That(data.ResourceName).Key("node_type.0.capacities.MemoryGB").HasValue("8"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceFabricCluster_clusterUpgradePolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clusterUpgradePolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_policy.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.clusterUpgradePolicyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_policy.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.clusterUpgradePolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_policy.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Hello").HasValue("World"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypesStateless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodeTypeStateless(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.#").HasValue("2"),
				check.That(data.ResourceName).Key("node_type.0.name").HasValue("first"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("node_type.0.is_primary").HasValue("true"),
				check.That(data.ResourceName).Key("node_type.1.name").HasValue("second"),
				check.That(data.ResourceName).Key("node_type.1.instance_count").HasValue("4"),
				check.That(data.ResourceName).Key("node_type.1.is_primary").HasValue("false"),
				check.That(data.ResourceName).Key("node_type.1.is_stateless").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_zonalUpgradeMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zonalUpgradeMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("node_type.0.multiple_availability_zones").HasValue("true"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("service_fabric_zonal_upgrade_mode").HasValue("Hierarchical"),
				check.That(data.ResourceName).Key("vmss_zonal_upgrade_mode").HasValue("Parallel"),
			),
		},
		data.ImportStep(),
	})
}

func (r ServiceFabricClusterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := cluster.ParseClusterID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ServiceFabric.ClustersClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}
	return utils.Bool(true), nil
}

func (r ServiceFabricClusterResource) basic(data acceptance.TestData, count int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = %d
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, count)
}

func (r ServiceFabricClusterResource) basicNodeTypeUpdate(data acceptance.TestData, count int, secondary_count int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = %d
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }

  node_type {
    name                 = "second"
    instance_count       = %d
    is_primary           = false
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, count, secondary_count)
}

func (r ServiceFabricClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_service_fabric_cluster" "import" {
  name                = azurerm_service_fabric_cluster.test.name
  resource_group_name = azurerm_service_fabric_cluster.test.resource_group_name
  location            = azurerm_service_fabric_cluster.test.location
  reliability_level   = azurerm_service_fabric_cluster.test.reliability_level
  upgrade_mode        = azurerm_service_fabric_cluster.test.upgrade_mode
  vm_image            = azurerm_service_fabric_cluster.test.vm_image
  management_endpoint = azurerm_service_fabric_cluster.test.management_endpoint

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, r.basic(data, 3))
}

func (r ServiceFabricClusterResource) manualClusterCodeVersion(data acceptance.TestData, clusterCodeVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                 = "acctest-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  reliability_level    = "Bronze"
  upgrade_mode         = "Manual"
  cluster_code_version = "%[3]s"
  vm_image             = "Windows"
  management_endpoint  = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, clusterCodeVersion)
}

func (r ServiceFabricClusterResource) addOnFeatures(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"
  add_on_features     = ["DnsService", "RepairManager"]

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) certificates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) reverseProxyCertificates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  reverse_proxy_certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                        = "first"
    instance_count              = 3
    is_primary                  = true
    client_endpoint_port        = 2020
    http_endpoint_port          = 80
    reverse_proxy_endpoint_port = 19081
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) clientCertificateThumbprint(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  client_certificate_thumbprint {
    thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = true
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) withMultipleClientCertificateThumbprints(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cluster-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-cluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  client_certificate_thumbprint {
    thumbprint = "1341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = true
  }

  client_certificate_thumbprint {
    thumbprint = "2341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = false
  }

  client_certificate_thumbprint {
    thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = true
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) clientCertificateCommonNames(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  client_certificate_common_name {
    common_name       = "firstcertcommonname"
    issuer_thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin          = true
  }

  client_certificate_common_name {
    common_name = "secondcertcommonname"
    is_admin    = false
  }

  client_certificate_thumbprint {
    thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = true
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) readerAdminClientCertificateThumbprint(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  client_certificate_thumbprint {
    thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = true
  }

  client_certificate_thumbprint {
    thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = false
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) certificateCommonNames(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate_common_names {
    common_names {
      certificate_common_name = "example"
    }

    x509_store_name = "My"
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) reverseProxyCertificateCommonNames(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate_common_names {
    common_names {
      certificate_common_name = "example"
    }

    x509_store_name = "My"
  }

  reverse_proxy_certificate_common_names {
    common_names {
      certificate_common_name = "example"
    }

    x509_store_name = "My"
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) azureActiveDirectory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {
}

data "azuread_domains" "test" {
}

resource "azuread_application" "cluster_explorer" {
  display_name    = "${azurerm_resource_group.test.name}-explorer-AAD"
  identifier_uris = ["https://test-%s.${data.azuread_domains.test.domains[0].domain_name}:19080/Explorer/index.html"]
  web {
    homepage_url  = "https://example:19080/Explorer/index.html"
    redirect_uris = ["https://example:19080/Explorer/index.html"]

    implicit_grant {
      access_token_issuance_enabled = true
    }
  }
  sign_in_audience = "AzureADMyOrg"


  # https://blogs.msdn.microsoft.com/aaddevsup/2018/06/06/guid-table-for-windows-azure-active-directory-permissions/
  # https://shawntabrizi.com/aad/common-microsoft-resources-azure-active-directory/
  required_resource_access {
    resource_app_id = "00000002-0000-0000-c000-000000000000"

    resource_access {
      id   = "311a71cc-e848-46a1-bdf8-97ff7156d8e6"
      type = "Scope"
    }
  }
}

resource "azuread_service_principal" "cluster_explorer" {
  application_id = azuread_application.cluster_explorer.application_id
}

resource "azuread_application" "cluster_console" {
  display_name     = "${azurerm_resource_group.test.name}-console-AAD"
  sign_in_audience = "AzureADMyOrg"
  web {
    redirect_uris = ["urn:ietf:wg:oauth:2.0:oob"]

    implicit_grant {
      access_token_issuance_enabled = true
    }
  }

  # https://blogs.msdn.microsoft.com/aaddevsup/2018/06/06/guid-table-for-windows-azure-active-directory-permissions/
  # https://shawntabrizi.com/aad/common-microsoft-resources-azure-active-directory/
  required_resource_access {
    resource_app_id = "00000002-0000-0000-c000-000000000000"

    resource_access {
      id   = "311a71cc-e848-46a1-bdf8-97ff7156d8e6"
      type = "Scope"
    }
  }

  required_resource_access {
    resource_app_id = azuread_application.cluster_explorer.application_id

    resource_access {
      id   = "311a71cc-e848-46a1-bdf8-97ff7156d8e6" # sign in and user profile permission ctx https://github.com/Azure/azure-cli/issues/7925
      type = "Scope"
    }
  }
}

resource "azuread_service_principal" "cluster_console" {
  application_id = azuread_application.cluster_console.application_id
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:19080"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  azure_active_directory {
    tenant_id              = data.azurerm_client_config.current.tenant_id
    cluster_application_id = azuread_application.cluster_explorer.application_id
    client_application_id  = azuread_application.cluster_console.application_id
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "system"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 19000
    http_endpoint_port   = 19080
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r ServiceFabricClusterResource) azureActiveDirectoryDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:19080"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "system"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 19000
    http_endpoint_port   = 19080
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) diagnosticsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  diagnostics_config {
    storage_account_name       = azurerm_storage_account.test.name
    protected_account_key_name = "StorageAccountKey1"
    blob_endpoint              = azurerm_storage_account.test.primary_blob_endpoint
    queue_endpoint             = azurerm_storage_account.test.primary_queue_endpoint
    table_endpoint             = azurerm_storage_account.test.primary_table_endpoint
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r ServiceFabricClusterResource) diagnosticsConfigDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r ServiceFabricClusterResource) fabricSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "None"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) nodeTypeCustomPorts(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80

    application_ports {
      start_port = 20000
      end_port   = 29999
    }

    ephemeral_ports {
      start_port = 30000
      end_port   = 39999
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) nodeTypeMultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }

  node_type {
    name                 = "second"
    instance_count       = 4
    is_primary           = false
    client_endpoint_port = 2121
    http_endpoint_port   = 81
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) nodeTypeProperties(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name = "first"

    placement_properties = {
      "HasSSD" = "true"
    }

    capacities = {
      "ClientConnections" = "20000"
      "MemoryGB"          = "8"
    }

    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) clusterUpgradePolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"
  diagnostics_config {
    storage_account_name       = azurerm_storage_account.test.name
    protected_account_key_name = "StorageAccountKey1"
    blob_endpoint              = azurerm_storage_account.test.primary_blob_endpoint
    queue_endpoint             = azurerm_storage_account.test.primary_queue_endpoint
    table_endpoint             = azurerm_storage_account.test.primary_table_endpoint
  }
  upgrade_policy {
    force_restart_enabled             = true
    health_check_retry_timeout        = "00:00:02"
    health_check_stable_duration      = "00:00:04"
    health_check_wait_duration        = "00:00:06"
    upgrade_domain_timeout            = "00:00:20"
    upgrade_replica_set_check_timeout = "00:00:10"
    upgrade_timeout                   = "00:00:40"
    health_policy {
      max_unhealthy_nodes_percent        = 5
      max_unhealthy_applications_percent = 40
    }
    delta_health_policy {
      max_delta_unhealthy_applications_percent         = 20
      max_delta_unhealthy_nodes_percent                = 40
      max_upgrade_domain_delta_unhealthy_nodes_percent = 60
    }
  }
  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r ServiceFabricClusterResource) clusterUpgradePolicyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"
  diagnostics_config {
    storage_account_name       = azurerm_storage_account.test.name
    protected_account_key_name = "StorageAccountKey1"
    blob_endpoint              = azurerm_storage_account.test.primary_blob_endpoint
    queue_endpoint             = azurerm_storage_account.test.primary_queue_endpoint
    table_endpoint             = azurerm_storage_account.test.primary_table_endpoint
  }
  upgrade_policy {
    force_restart_enabled        = false
    health_check_retry_timeout   = "00:00:02"
    health_check_stable_duration = "00:00:04"
    health_check_wait_duration   = "00:00:06"
    health_policy {
      max_unhealthy_nodes_percent = 5
    }
  }
  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r ServiceFabricClusterResource) nodeTypeStateless(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }

  node_type {
    name                 = "second"
    instance_count       = 4
    is_primary           = false
    is_stateless         = true
    client_endpoint_port = 2121
    http_endpoint_port   = 81
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) zonalUpgradeMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                              = "acctest-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  reliability_level                 = "Bronze"
  upgrade_mode                      = "Automatic"
  vm_image                          = "Windows"
  management_endpoint               = "http://example:80"
  service_fabric_zonal_upgrade_mode = "Hierarchical"
  vmss_zonal_upgrade_mode           = "Parallel"

  node_type {
    name                        = "first"
    instance_count              = 3
    is_primary                  = true
    client_endpoint_port        = 2020
    http_endpoint_port          = 80
    multiple_availability_zones = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
