// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementResource struct{}

func TestAccApiManagement_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

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

func TestAccApiManagement_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

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

func TestAccApiManagement_skuUpgradeDowngrade(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.standardSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_customProps(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customProps(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.0.expiry").Exists(),
				check.That(data.ResourceName).Key("certificate.0.subject").Exists(),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").Exists(),
				check.That(data.ResourceName).Key("certificate.1.expiry").Exists(),
				check.That(data.ResourceName).Key("certificate.1.subject").Exists(),
				check.That(data.ResourceName).Key("certificate.1.thumbprint").Exists(),
				check.That(data.ResourceName).Key("certificate.2.expiry").Exists(),
				check.That(data.ResourceName).Key("certificate.2.subject").Exists(),
				check.That(data.ResourceName).Key("certificate.2.thumbprint").Exists(),
				check.That(data.ResourceName).Key("certificate.3.expiry").Exists(),
				check.That(data.ResourceName).Key("certificate.3.subject").Exists(),
				check.That(data.ResourceName).Key("certificate.3.thumbprint").Exists(),
				check.That(data.ResourceName).Key("hostname_configuration.0.developer_portal.0.expiry").Exists(),
				check.That(data.ResourceName).Key("hostname_configuration.0.developer_portal.0.subject").Exists(),
				check.That(data.ResourceName).Key("hostname_configuration.0.developer_portal.0.thumbprint").Exists(),
				check.That(data.ResourceName).Key("hostname_configuration.0.portal.0.expiry").Exists(),
				check.That(data.ResourceName).Key("hostname_configuration.0.portal.0.subject").Exists(),
				check.That(data.ResourceName).Key("hostname_configuration.0.portal.0.thumbprint").Exists(),
			),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"certificate.0.certificate_password",
				"certificate.0.encoded_certificate",
				"certificate.1.certificate_password",
				"certificate.1.encoded_certificate",
				"certificate.2.encoded_certificate",
				"certificate.3.encoded_certificate",
				"hostname_configuration.0.portal.0.certificate",                    // not returned from API, sensitive
				"hostname_configuration.0.portal.0.certificate_password",           // not returned from API, sensitive
				"hostname_configuration.0.developer_portal.0.certificate",          // not returned from API, sensitive
				"hostname_configuration.0.developer_portal.0.certificate_password", // not returned from API, sensitive
				"hostname_configuration.0.proxy.1.certificate",                     // not returned from API, sensitive
				"hostname_configuration.0.proxy.1.certificate_password",            // not returned from API, sensitive
				"hostname_configuration.0.proxy.2.certificate",                     // not returned from API, sensitive
				"hostname_configuration.0.proxy.2.certificate_password",            // not returned from API, sensitive
			},
		},
	})
}

func TestAccApiManagement_completeUpdateAdditionalLocations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate", // not returned from API, sensitive
			"hostname_configuration.0.portal.0.certificate",                    // not returned from API, sensitive
			"hostname_configuration.0.portal.0.certificate_password",           // not returned from API, sensitive
			"hostname_configuration.0.developer_portal.0.certificate",          // not returned from API, sensitive
			"hostname_configuration.0.developer_portal.0.certificate_password", // not returned from API, sensitive
			"hostname_configuration.0.proxy.1.certificate",                     // not returned from API, sensitive
			"hostname_configuration.0.proxy.1.certificate_password",            // not returned from API, sensitive
			"hostname_configuration.0.proxy.2.certificate",                     // not returned from API, sensitive
			"hostname_configuration.0.proxy.2.certificate_password",            // not returned from API, sensitive
		),
		{
			Config: r.completeUpdateAdditionalLocations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate", // not returned from API, sensitive
			"hostname_configuration.0.portal.0.certificate",                    // not returned from API, sensitive
			"hostname_configuration.0.portal.0.certificate_password",           // not returned from API, sensitive
			"hostname_configuration.0.developer_portal.0.certificate",          // not returned from API, sensitive
			"hostname_configuration.0.developer_portal.0.certificate_password", // not returned from API, sensitive
			"hostname_configuration.0.proxy.1.certificate",                     // not returned from API, sensitive
			"hostname_configuration.0.proxy.1.certificate_password",            // not returned from API, sensitive
			"hostname_configuration.0.proxy.2.certificate",                     // not returned from API, sensitive
			"hostname_configuration.0.proxy.2.certificate_password",            // not returned from API, sensitive
		),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate", // not returned from API, sensitive
			"hostname_configuration.0.portal.0.certificate",                    // not returned from API, sensitive
			"hostname_configuration.0.portal.0.certificate_password",           // not returned from API, sensitive
			"hostname_configuration.0.developer_portal.0.certificate",          // not returned from API, sensitive
			"hostname_configuration.0.developer_portal.0.certificate_password", // not returned from API, sensitive
			"hostname_configuration.0.proxy.1.certificate",                     // not returned from API, sensitive
			"hostname_configuration.0.proxy.1.certificate_password",            // not returned from API, sensitive
			"hostname_configuration.0.proxy.2.certificate",                     // not returned from API, sensitive
			"hostname_configuration.0.proxy.2.certificate_password",            // not returned from API, sensitive
		),
	})
}

func TestAccApiManagement_signInSignUpSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.signInSignUpSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_delegationSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.delegationSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.delegationSettingsDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.delegationSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_virtualNetworkInternal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetworkInternal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_type").HasValue("Internal"),
				check.That(data.ResourceName).Key("private_ip_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_virtualNetworkInternalUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.virtualNetworkInternal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_virtualNetworkInternalAdditionalLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetworkInternalAdditionalLocation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_type").HasValue("Internal"),
				check.That(data.ResourceName).Key("public_ip_address_id").Exists(),
				check.That(data.ResourceName).Key("private_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("additional_location.0.private_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("additional_location.0.public_ip_address_id").Exists(),
				check.That(data.ResourceName).Key("additional_location.0.capacity").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

// Api Management doesn't support hostname keyvault using UserAssigned Identity
// There will be a inevitable dependency cycle here when using SystemAssigned Identity
// 1. create SystemAssigned Identity, grant the identity certificate access
// 2. Update the hostname configuration of the keyvault certificate
func TestAccApiManagement_identitySystemAssignedUpdateHostnameConfigurationsVersionedKeyVaultId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssignedUpdateHostnameConfigurationsKeyVaultId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssignedUpdateHostnameConfigurationsVersionedKeyVaultIdUpdateCD(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_identitySystemAssignedUpdateHostnameConfigurationsVersionlessKeyVaultId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssignedUpdateHostnameConfigurationsKeyVaultId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssignedUpdateHostnameConfigurationsVersionlessKeyVaultIdUpdateCD(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_identityUserAssignedHostnameConfigurationsKeyVaultId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityUserAssignedHostnameConfigurationsKeyVaultId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_consumption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_consumptionWithTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.consumptionWithTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_clientCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.consumptionClientCertificateEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.consumptionClientCertificateDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_gatewayDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleLocations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.gatewayDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleLocations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_minApiVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.consumptionMinApiVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.consumptionMinApiVersionUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.consumption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_removeSamples(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.removeSamples(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				r.testCheckHasNoProductsOrApis(data.ResourceName),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_softDeleteRecovery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.softDelete(data),
		},
		{
			Config: r.consumption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_softDeleteRecoveryDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.softDelete(data),
		},
		{
			Config:      r.consumptionRecoveryDisabled(data),
			ExpectError: regexp.MustCompile(`An existing soft-deleted API Management exists with the Name "[^"]+" in the location "[^"]+"`),
		},
	})
}

func (ApiManagementResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apimanagementservice.ParseServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ServiceClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (ApiManagementResource) testCheckHasNoProductsOrApis(resourceName string) pluginsdk.TestCheckFunc {
	return func(state *pluginsdk.State) error {
		client, err := testclient.Build()
		if err != nil {
			return fmt.Errorf("building client: %+v", err)
		}
		ctx := client.StopContext

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%q was not found in the state", resourceName)
		}

		id, err := apimanagementservice.ParseServiceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		timeout, cancel := context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()

		apiServiceId := api.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
		listResp, err := client.ApiManagement.ApiClient.ListByService(timeout, apiServiceId, api.ListByServiceOperationOptions{})
		if err != nil {
			return fmt.Errorf("listing APIs after creation of %s: %+v", *id, err)
		}

		if model := listResp.Model; model != nil {
			if count := len(*model); count > 0 {
				return fmt.Errorf("%s has %d unexpected associated APIs", *id, count)
			}
		}

		produceServiceId := product.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
		proListResp, err := client.ApiManagement.ProductsClient.ListByService(timeout, produceServiceId, product.ListByServiceOperationOptions{})
		if err != nil {
			return fmt.Errorf("listing products after creation of %s: %+v", *id, err)
		}
		if model := proListResp.Model; model != nil {
			if count := len(*model); count > 0 {
				return fmt.Errorf("%s has %d unexpected associated Products", *id, count)
			}
		}

		return nil
	}
}

func TestAccApiManagement_identityUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_identityNoneUpdateUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityNone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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
	})
}

func TestAccApiManagement_identityUserAssignedUpdateNone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityNone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_identitySystemAssignedUpdateNone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityNone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_identityNoneUpdateSystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityNone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_identitySystemAssignedUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_identitySystemAssignedUserAssignedUpdateNone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityNone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_identityNoneUpdateSystemAssignedUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityNone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_identitySystemAssignedUserAssignedUpdateSystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_identitySystemAssignedUserAssignedUpdateUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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
	})
}

func TestAccApiManagement_tenantAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tenantAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tenant_access.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("tenant_access.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("tenant_access.0.primary_key").Exists(),
				check.That(data.ResourceName).Key("tenant_access.0.secondary_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_additionalLocationGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.additionalLocationGateway(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagement_additionalLocationGateway_DivergentZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")
	r := ApiManagementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.additionalLocationGateway_DivergentZones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zones.#").HasValue("2"),
				check.That(data.ResourceName).Key("additional_location.0.zones.#").HasValue("0"),
			),
		},
	})
}

func (ApiManagementResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) additionalLocationGateway(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG2-%[1]d"
  location = "%[3]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Premium_2"

  additional_location {
    location         = azurerm_resource_group.test2.location
    gateway_disabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (ApiManagementResource) additionalLocationGateway_DivergentZones(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
  tags = {
    owner = "Dom Routley"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNET-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["192.168.0.0/24"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-gateway"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = azurerm_virtual_network.test.address_space
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "Client_communication_to_API_Management"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "80"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "VirtualNetwork"
  }

  security_rule {
    name                       = "Secure_Client_communication_to_API_Management"
    priority                   = 110
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "VirtualNetwork"
  }

  security_rule {
    name                       = "Management_endpoint_for_Azure_portal_and_Powershell"
    priority                   = 120
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3443"
    source_address_prefix      = "ApiManagement"
    destination_address_prefix = "VirtualNetwork"
  }

  security_rule {
    name                       = "Authenticate_To_Azure_Active_Directory"
    priority                   = 200
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_ranges    = ["80", "443"]
    source_address_prefix      = "ApiManagement"
    destination_address_prefix = "VirtualNetwork"
  }
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG2-%[1]d"
  location = "%[3]s"
  tags = {
    owner = "Dom Routley"
  }
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestVNET2-%[1]d"
  resource_group_name = azurerm_resource_group.test2.name
  location            = azurerm_resource_group.test2.location
  address_space       = ["192.168.1.0/24"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest2-gateway"
  resource_group_name  = azurerm_resource_group.test2.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = azurerm_virtual_network.test2.address_space
}

resource "azurerm_network_security_group" "test2" {
  name                = "acctest-NSG2-%[1]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name

  security_rule {
    name                       = "Client_communication_to_API_Management"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "80"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "VirtualNetwork"
  }

  security_rule {
    name                       = "Secure_Client_communication_to_API_Management"
    priority                   = 110
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "VirtualNetwork"
  }

  security_rule {
    name                       = "Management_endpoint_for_Azure_portal_and_Powershell"
    priority                   = 120
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3443"
    source_address_prefix      = "ApiManagement"
    destination_address_prefix = "VirtualNetwork"
  }

  security_rule {
    name                       = "Authenticate_To_Azure_Active_Directory"
    priority                   = 200
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_ranges    = ["80", "443"]
    source_address_prefix      = "ApiManagement"
    destination_address_prefix = "VirtualNetwork"
  }
}

resource "azurerm_subnet_network_security_group_association" "test2" {
  subnet_id                 = azurerm_subnet.test2.id
  network_security_group_id = azurerm_network_security_group.test2.id
}

resource "azurerm_public_ip" "test2" {
  name                = "acctest2IP-%[1]d"
  resource_group_name = azurerm_resource_group.test2.name
  location            = azurerm_resource_group.test2.location
  sku                 = "Standard"
  allocation_method   = "Static"
  domain_name_label   = "acctest2-ip-%[1]d"
}

resource "azurerm_api_management" "test" {
  name                 = "acctestAM-%[1]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  publisher_name       = "pub1"
  publisher_email      = "pub1@email.com"
  sku_name             = "Premium_2"
  virtual_network_type = "Internal"
  zones                = ["1", "2"]

  virtual_network_configuration {
    subnet_id = azurerm_subnet.test.id
  }

  additional_location {
    location             = azurerm_resource_group.test2.location
    public_ip_address_id = azurerm_public_ip.test2.id
    virtual_network_configuration {
      subnet_id = azurerm_subnet.test2.id
    }
  }

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_network_security_group_association.test2,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (ApiManagementResource) standardSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Standard_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiManagementResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "import" {
  name                = azurerm_api_management.test.name
  location            = azurerm_api_management.test.location
  resource_group_name = azurerm_api_management.test.resource_group_name
  publisher_name      = azurerm_api_management.test.publisher_name
  publisher_email     = azurerm_api_management.test.publisher_email

  sku_name = "Developer_1"
}
`, r.basic(data))
}

func (ApiManagementResource) customProps(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  security {
    frontend_tls10_enabled     = true
    triple_des_ciphers_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

func (ApiManagementResource) signInSignUpSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  sign_in {
    enabled = true
  }

  sign_up {
    enabled = true

    terms_of_service {
      enabled          = true
      consent_required = false
      text             = "Lorem Ipsum Dolor Morty"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) delegationSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  delegation {
    url                       = "https://google.com"
    subscriptions_enabled     = true
    validation_key            = "aW50ZWdyYXRpb24mMjAyMzAzMTAxODMwJkxRaUxzcUVsaUpEaHJRK01YZkJYV3paUi9qdzZDSWMrazhjUXB0bVdyTGxKcVYrd0R4OXRqMGRzTWZXU3hmeGQ0a2V0WjcrcE44U0dJdDNsYUQ3Rk5BPT0="
    user_registration_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) delegationSettingsDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  delegation {
    subscriptions_enabled     = false
    user_registration_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) complete(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-api1-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-api2-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test3" {
  name     = "acctestRG-api3-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                      = "acctestAM-%d"
  publisher_name            = "pub1"
  publisher_email           = "pub1@email.com"
  notification_sender_email = "notification@email.com"

  sku_name = "Premium_2"

  additional_location {
    zones    = []
    capacity = 1
    location = azurerm_resource_group.test2.location
  }

  additional_location {
    zones    = []
    location = azurerm_resource_group.test3.location
  }

  certificate {
    encoded_certificate  = filebase64("testdata/api_management_api_test.pfx")
    certificate_password = "terraform"
    store_name           = "CertificateAuthority"
  }

  certificate {
    encoded_certificate  = filebase64("testdata/api_management_api_test.pfx")
    certificate_password = "terraform"
    store_name           = "Root"
  }

  certificate {
    encoded_certificate = filebase64("testdata/api_management_api_test.cer")
    store_name          = "Root"
  }

  certificate {
    encoded_certificate = filebase64("testdata/api_management_api_test.cer")
    store_name          = "CertificateAuthority"
  }

  protocols {
    enable_http2 = true
  }

  security {
    enable_backend_tls11                                = true
    enable_backend_ssl30                                = true
    enable_backend_tls10                                = true
    enable_frontend_ssl30                               = true
    enable_frontend_tls10                               = true
    enable_frontend_tls11                               = true
    tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled = true
    tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled = true
    tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled   = true
    tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled   = true
    tls_rsa_with_aes128_cbc_sha256_ciphers_enabled      = true
    tls_rsa_with_aes128_cbc_sha_ciphers_enabled         = true
    tls_rsa_with_aes128_gcm_sha256_ciphers_enabled      = true
    tls_rsa_with_aes256_cbc_sha256_ciphers_enabled      = true
    tls_rsa_with_aes256_cbc_sha_ciphers_enabled         = true
    triple_des_ciphers_enabled                          = true
  }

  hostname_configuration {
    proxy {
      host_name                    = "acctestAM-%d.azure-api.net"
      negotiate_client_certificate = true
    }

    proxy {
      host_name                    = "api.terraform.io"
      certificate                  = filebase64("testdata/api_management_api_test.pfx")
      certificate_password         = "terraform"
      default_ssl_binding          = true
      negotiate_client_certificate = false
    }

    proxy {
      host_name                    = "api2.terraform.io"
      certificate                  = filebase64("testdata/api_management_api2_test.pfx")
      certificate_password         = "terraform"
      negotiate_client_certificate = true
    }

    portal {
      host_name            = "portal.terraform.io"
      certificate          = filebase64("testdata/api_management_portal_test.pfx")
      certificate_password = "terraform"
    }

    developer_portal {
      host_name            = "developer-portal.terraform.io"
      certificate          = filebase64("testdata/api_management_developer_portal_test.pfx")
      certificate_password = "terraform"
    }
  }

  tags = {
    "Acceptance" = "Test"
  }

  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.Locations.Ternary, data.RandomInteger, data.RandomInteger)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-api1-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-api2-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test3" {
  name     = "acctestRG-api3-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                      = "acctestAM-%d"
  publisher_name            = "pub1"
  publisher_email           = "pub1@email.com"
  notification_sender_email = "notification@email.com"

  sku_name = "Premium_2"

  additional_location {
    zones    = []
    capacity = 1
    location = azurerm_resource_group.test2.location
  }

  additional_location {
    zones    = []
    location = azurerm_resource_group.test3.location
  }

  certificate {
    encoded_certificate  = filebase64("testdata/api_management_api_test.pfx")
    certificate_password = "terraform"
    store_name           = "CertificateAuthority"
  }

  certificate {
    encoded_certificate  = filebase64("testdata/api_management_api_test.pfx")
    certificate_password = "terraform"
    store_name           = "Root"
  }

  certificate {
    encoded_certificate = filebase64("testdata/api_management_api_test.cer")
    store_name          = "Root"
  }

  certificate {
    encoded_certificate = filebase64("testdata/api_management_api_test.cer")
    store_name          = "CertificateAuthority"
  }

  protocols {
    http2_enabled = true
  }

  security {
    backend_tls11_enabled                               = true
    backend_ssl30_enabled                               = true
    backend_tls10_enabled                               = true
    frontend_ssl30_enabled                              = true
    frontend_tls10_enabled                              = true
    frontend_tls11_enabled                              = true
    tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled = true
    tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled = true
    tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled   = true
    tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled   = true
    tls_rsa_with_aes128_cbc_sha256_ciphers_enabled      = true
    tls_rsa_with_aes128_cbc_sha_ciphers_enabled         = true
    tls_rsa_with_aes128_gcm_sha256_ciphers_enabled      = true
    tls_rsa_with_aes256_cbc_sha256_ciphers_enabled      = true
    tls_rsa_with_aes256_cbc_sha_ciphers_enabled         = true
    triple_des_ciphers_enabled                          = true
  }

  hostname_configuration {
    proxy {
      host_name                    = "acctestAM-%d.azure-api.net"
      negotiate_client_certificate = true
    }

    proxy {
      host_name                    = "api.terraform.io"
      certificate                  = filebase64("testdata/api_management_api_test.pfx")
      certificate_password         = "terraform"
      default_ssl_binding          = true
      negotiate_client_certificate = false
    }

    proxy {
      host_name                    = "api2.terraform.io"
      certificate                  = filebase64("testdata/api_management_api2_test.pfx")
      certificate_password         = "terraform"
      negotiate_client_certificate = true
    }

    portal {
      host_name            = "portal.terraform.io"
      certificate          = filebase64("testdata/api_management_portal_test.pfx")
      certificate_password = "terraform"
    }

    developer_portal {
      host_name            = "developer-portal.terraform.io"
      certificate          = filebase64("testdata/api_management_developer_portal_test.pfx")
      certificate_password = "terraform"
    }
  }

  tags = {
    "Acceptance" = "Test"
  }

  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.Locations.Ternary, data.RandomInteger, data.RandomInteger)
}

func (ApiManagementResource) completeUpdateAdditionalLocations(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-api1-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-api2-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test3" {
  name     = "acctestRG-api3-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                      = "acctestAM-%d"
  publisher_name            = "pub1"
  publisher_email           = "pub1@email.com"
  notification_sender_email = "notification@email.com"

  sku_name = "Premium_2"

  additional_location {
    zones    = [1]
    location = azurerm_resource_group.test2.location
    capacity = 2
  }

  certificate {
    encoded_certificate  = filebase64("testdata/api_management_api_test.pfx")
    certificate_password = "terraform"
    store_name           = "CertificateAuthority"
  }

  certificate {
    encoded_certificate  = filebase64("testdata/api_management_api_test.pfx")
    certificate_password = "terraform"
    store_name           = "Root"
  }

  certificate {
    encoded_certificate = filebase64("testdata/api_management_api_test.cer")
    store_name          = "Root"
  }

  certificate {
    encoded_certificate = filebase64("testdata/api_management_api_test.cer")
    store_name          = "CertificateAuthority"
  }

  protocols {
    http2_enabled = true
  }

  security {
    backend_tls11_enabled                               = true
    backend_ssl30_enabled                               = true
    backend_tls10_enabled                               = true
    frontend_ssl30_enabled                              = true
    frontend_tls10_enabled                              = true
    frontend_tls11_enabled                              = true
    tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled = true
    tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled = true
    tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled   = true
    tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled   = true
    tls_rsa_with_aes128_cbc_sha256_ciphers_enabled      = true
    tls_rsa_with_aes128_cbc_sha_ciphers_enabled         = true
    tls_rsa_with_aes128_gcm_sha256_ciphers_enabled      = true
    tls_rsa_with_aes256_cbc_sha256_ciphers_enabled      = true
    tls_rsa_with_aes256_cbc_sha_ciphers_enabled         = true
    triple_des_ciphers_enabled                          = true
  }

  hostname_configuration {
    proxy {
      host_name                    = "acctestAM-%d.azure-api.net"
      negotiate_client_certificate = true
    }

    proxy {
      host_name                    = "api.terraform.io"
      certificate                  = filebase64("testdata/api_management_api_test.pfx")
      certificate_password         = "terraform"
      default_ssl_binding          = true
      negotiate_client_certificate = false
    }

    proxy {
      host_name                    = "api2.terraform.io"
      certificate                  = filebase64("testdata/api_management_api2_test.pfx")
      certificate_password         = "terraform"
      negotiate_client_certificate = true
    }

    portal {
      host_name            = "portal.terraform.io"
      certificate          = filebase64("testdata/api_management_portal_test.pfx")
      certificate_password = "terraform"
    }

    developer_portal {
      host_name            = "developer-portal.terraform.io"
      certificate          = filebase64("testdata/api_management_developer_portal_test.pfx")
      certificate_password = "terraform"
    }
  }

  tags = {
    "Acceptance" = "Test"
  }

  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.Locations.Ternary, data.RandomInteger, data.RandomInteger)
}

func (ApiManagementResource) virtualNetworkTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNET-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestSNET-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_network_security_rule" "client" {
  name                        = "Client_communication_to_API_Management"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "80"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "secure_client" {
  name                        = "Secure_Client_communication_to_API_Management"
  priority                    = 110
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "443"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "endpoint" {
  name                        = "Management_endpoint_for_Azure_portal_and_Powershell"
  priority                    = 120
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "3443"
  source_address_prefix       = "ApiManagement"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "authenticate" {
  name                        = "Authenticate_To_Azure_Active_Directory"
  priority                    = 200
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges     = ["80", "443"]
  source_address_prefix       = "ApiManagement"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementResource) virtualNetworkInternal(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  virtual_network_type = "Internal"
  virtual_network_configuration {
    subnet_id = azurerm_subnet.test.id
  }
}
`, r.virtualNetworkTemplate(data), data.RandomInteger)
}

func (r ApiManagementResource) virtualNetworkInternalAdditionalLocation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-%[2]d-2"
  location = "%[3]s"
}

// subnet2 from the second location
resource "azurerm_virtual_network" "test2" {
  name                = "acctestVNET2-%[2]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestSNET2-%[2]d"
  resource_group_name  = azurerm_resource_group.test2.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "azurerm_network_security_group" "test2" {
  name                = "acctest-NSG2-%[2]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
}

resource "azurerm_subnet_network_security_group_association" "test2" {
  subnet_id                 = azurerm_subnet.test2.id
  network_security_group_id = azurerm_network_security_group.test2.id
}

resource "azurerm_network_security_rule" "client2" {
  name                        = "Client_communication_to_API_Management"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "80"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test2.name
  network_security_group_name = azurerm_network_security_group.test2.name
}

resource "azurerm_network_security_rule" "secure_client2" {
  name                        = "Secure_Client_communication_to_API_Management"
  priority                    = 110
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "443"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test2.name
  network_security_group_name = azurerm_network_security_group.test2.name
}

resource "azurerm_network_security_rule" "endpoint2" {
  name                        = "Management_endpoint_for_Azure_portal_and_Powershell"
  priority                    = 120
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "3443"
  source_address_prefix       = "ApiManagement"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test2.name
  network_security_group_name = azurerm_network_security_group.test2.name
}

resource "azurerm_network_security_rule" "authenticate2" {
  name                        = "Authenticate_To_Azure_Active_Directory"
  priority                    = 200
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges     = ["80", "443"]
  source_address_prefix       = "ApiManagement"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test2.name
  network_security_group_name = azurerm_network_security_group.test2.name
}

resource "azurerm_public_ip" "test1" {
  name                = "acctest-IP1-%[4]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
  allocation_method   = "Static"
  domain_name_label   = "acctest-ip1-%[4]s"
}

resource "azurerm_public_ip" "test2" {
  name                = "acctest-IP2-%[4]s"
  resource_group_name = azurerm_resource_group.test2.name
  location            = azurerm_resource_group.test2.location
  sku                 = "Standard"
  allocation_method   = "Static"
  domain_name_label   = "acctest-ip2-%[4]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Premium_1"

  additional_location {
    location = azurerm_resource_group.test2.location
    capacity = 1

    public_ip_address_id = azurerm_public_ip.test2.id
    virtual_network_configuration {
      subnet_id = azurerm_subnet.test2.id
    }
  }

  virtual_network_type = "Internal"
  public_ip_address_id = azurerm_public_ip.test1.id
  virtual_network_configuration {
    subnet_id = azurerm_subnet.test.id
  }

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_network_security_group_association.test2,
  ]
}
`, r.virtualNetworkTemplate(data), data.RandomInteger, data.Locations.Secondary, data.RandomString)
}

func (ApiManagementResource) identityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_api_management" "test" {
  depends_on          = [azurerm_user_assigned_identity.test]
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ApiManagementResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) identitySystemAssignedUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ApiManagementResource) identityNone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) identitySystemAssignedUpdateHostnameConfigurationsTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestKV-%[4]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  certificate_permissions = [
    "Create",
    "Delete",
    "DeleteIssuers",
    "Get",
    "GetIssuers",
    "Import",
    "List",
    "ListIssuers",
    "ManageContacts",
    "ManageIssuers",
    "SetIssuers",
    "Update",
    "Purge",
  ]
  secret_permissions = [
    "Delete",
    "Get",
    "List",
    "Purge",
  ]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_api_management.test.identity[0].tenant_id
  object_id    = azurerm_api_management.test.identity[0].principal_id
  secret_permissions = [
    "Get",
    "List",
  ]
}

resource "azurerm_key_vault_certificate" "test" {
  depends_on   = [azurerm_key_vault_access_policy.test]
  name         = "acctestKVCert-%[3]d"
  key_vault_id = azurerm_key_vault.test.id
  certificate_policy {
    issuer_parameters {
      name = "Self"
    }
    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }
    secret_properties {
      content_type = "application/x-pkcs12"
    }
    x509_certificate_properties {
      # Server Authentication = 1.3.6.1.5.5.7.3.1
      # Client Authentication = 1.3.6.1.5.5.7.3.2
      extended_key_usage = ["1.3.6.1.5.5.7.3.1"]
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]
      subject_alternative_names {
        dns_names = ["api.pluginsdk.io"]
      }
      subject            = "CN=api.pluginsdk.io"
      validity_in_months = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r ApiManagementResource) identitySystemAssignedUpdateHostnameConfigurationsKeyVaultId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  identity {
    type = "SystemAssigned"
  }
}
`, r.identitySystemAssignedUpdateHostnameConfigurationsTemplate(data), data.RandomInteger)
}

func (r ApiManagementResource) identitySystemAssignedUpdateHostnameConfigurationsVersionlessKeyVaultIdUpdateCD(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  identity {
    type = "SystemAssigned"
  }

  hostname_configuration {
    proxy {
      host_name                    = "acctestAM-%d.azure-api.net"
      negotiate_client_certificate = true
    }

    proxy {
      host_name                    = "api.pluginsdk.io"
      key_vault_certificate_id     = azurerm_key_vault_certificate.test.versionless_secret_id
      default_ssl_binding          = true
      negotiate_client_certificate = false
    }
  }
}
`, r.identitySystemAssignedUpdateHostnameConfigurationsTemplate(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementResource) identitySystemAssignedUpdateHostnameConfigurationsVersionedKeyVaultIdUpdateCD(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  identity {
    type = "SystemAssigned"
  }

  hostname_configuration {
    proxy {
      host_name                    = "acctestAM-%d.azure-api.net"
      negotiate_client_certificate = true
    }

    proxy {
      host_name                    = "api.pluginsdk.io"
      key_vault_certificate_id     = azurerm_key_vault_certificate.test.secret_id
      default_ssl_binding          = true
      negotiate_client_certificate = false
    }
  }
}
`, r.identitySystemAssignedUpdateHostnameConfigurationsTemplate(data), data.RandomInteger, data.RandomInteger)
}

func (ApiManagementResource) identityUserAssignedHostnameConfigurationsKeyVaultId(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestKV-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  certificate_permissions = [
    "Create",
    "Delete",
    "DeleteIssuers",
    "Get",
    "GetIssuers",
    "Import",
    "List",
    "ListIssuers",
    "ManageContacts",
    "ManageIssuers",
    "SetIssuers",
    "Update",
    "Purge",
  ]
  secret_permissions = [
    "Delete",
    "Get",
    "List",
    "Purge",
  ]
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id
  secret_permissions = [
    "Get",
    "List",
  ]
}

resource "azurerm_key_vault_certificate" "test" {
  depends_on   = [azurerm_key_vault_access_policy.test]
  name         = "acctestKVCert-%[1]d"
  key_vault_id = azurerm_key_vault.test.id
  certificate_policy {
    issuer_parameters {
      name = "Self"
    }
    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }
    secret_properties {
      content_type = "application/x-pkcs12"
    }
    x509_certificate_properties {
      # Server Authentication = 1.3.6.1.5.5.7.3.1
      # Client Authentication = 1.3.6.1.5.5.7.3.2
      extended_key_usage = ["1.3.6.1.5.5.7.3.1"]
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]
      subject_alternative_names {
        dns_names = ["api.terraform.io"]
      }
      subject            = "CN=api.terraform.io"
      validity_in_months = 1
    }
  }
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  hostname_configuration {
    proxy {
      host_name                    = "acctestAM-%[1]d.azure-api.net"
      negotiate_client_certificate = true
    }

    proxy {
      host_name                       = "api.terraform.io"
      key_vault_id                    = azurerm_key_vault_certificate.test.secret_id
      default_ssl_binding             = true
      negotiate_client_certificate    = false
      ssl_keyvault_identity_client_id = azurerm_user_assigned_identity.test.client_id
    }
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [azurerm_key_vault_access_policy.test2]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestKV-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  certificate_permissions = [
    "Create",
    "Delete",
    "DeleteIssuers",
    "Get",
    "GetIssuers",
    "Import",
    "List",
    "ListIssuers",
    "ManageContacts",
    "ManageIssuers",
    "SetIssuers",
    "Update",
    "Purge",
  ]
  secret_permissions = [
    "Delete",
    "Get",
    "List",
    "Purge",
  ]
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id
  secret_permissions = [
    "Get",
    "List",
  ]
}

resource "azurerm_key_vault_certificate" "test" {
  depends_on   = [azurerm_key_vault_access_policy.test]
  name         = "acctestKVCert-%[1]d"
  key_vault_id = azurerm_key_vault.test.id
  certificate_policy {
    issuer_parameters {
      name = "Self"
    }
    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }
    secret_properties {
      content_type = "application/x-pkcs12"
    }
    x509_certificate_properties {
      # Server Authentication = 1.3.6.1.5.5.7.3.1
      # Client Authentication = 1.3.6.1.5.5.7.3.2
      extended_key_usage = ["1.3.6.1.5.5.7.3.1"]
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]
      subject_alternative_names {
        dns_names = ["api.terraform.io"]
      }
      subject            = "CN=api.terraform.io"
      validity_in_months = 1
    }
  }
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  hostname_configuration {
    proxy {
      host_name                    = "acctestAM-%[1]d.azure-api.net"
      negotiate_client_certificate = true
    }

    proxy {
      host_name                       = "api.terraform.io"
      key_vault_certificate_id        = azurerm_key_vault_certificate.test.secret_id
      default_ssl_binding             = true
      negotiate_client_certificate    = false
      ssl_keyvault_identity_client_id = azurerm_user_assigned_identity.test.client_id
    }
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [azurerm_key_vault_access_policy.test2]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (ApiManagementResource) consumption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) consumptionWithTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) consumptionClientCertificateEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                       = "acctestAM-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  publisher_name             = "pub1"
  publisher_email            = "pub1@email.com"
  sku_name                   = "Consumption_0"
  client_certificate_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) consumptionClientCertificateDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                       = "acctestAM-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  publisher_name             = "pub1"
  publisher_email            = "pub1@email.com"
  sku_name                   = "Consumption_0"
  client_certificate_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) multipleLocations(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
  additional_location {
    location = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary)
}

func (ApiManagementResource) gatewayDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
  gateway_disabled    = true
  additional_location {
    location = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary)
}

func (ApiManagementResource) consumptionMinApiVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
  min_api_version     = "2019-12-01"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) consumptionMinApiVersionUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
  min_api_version     = "2020-12-01"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) consumptionRecoveryDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    api_management {
      recover_soft_deleted = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) softDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    api_management {
      purge_soft_delete_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ApiManagementResource) removeSamples(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementResource) tenantAccess(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  tenant_access {
    enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
