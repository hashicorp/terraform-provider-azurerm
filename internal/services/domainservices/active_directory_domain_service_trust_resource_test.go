// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package domainservices_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DomainServiceTrustResource struct {
	// The name of the resource group contains the AADDS
	AADDSResourceGroupName string
	// The name the AADDS
	AADDSName string
	// The FQDN of the on-premise ADDS
	ADDSFqdn string
	// The primary IP of the DNS that can resolve the on-premise ADDS
	ADDSIp1 string
	// The secondary IP of the DNS that can resolve the on-premise ADDS
	ADDSIp2 string
	// The password of the inbound trust set in the on-premise ADDS
	TrustPassword string
}

func NewDomainServiceTrustResource() (*DomainServiceTrustResource, error) {
	aaddsRgName := os.Getenv("ARM_TEST_AADDS_RESOURCE_GROUP_NAME")
	if aaddsRgName == "" {
		return nil, fmt.Errorf("`ARM_TEST_AADDS_RESOURCE_GROUP_NAME` is not set")
	}

	aaddsName := os.Getenv("ARM_TEST_AADDS_NAME")
	if aaddsName == "" {
		return nil, fmt.Errorf("`ARM_TEST_AADDS_NAME` is not set")
	}

	addsFqdn := os.Getenv("ARM_TEST_ADDS_FQDN")
	if addsFqdn == "" {
		return nil, fmt.Errorf("`ARM_TEST_ADDS_FQDN` is not set")
	}

	addsIp1 := os.Getenv("ARM_TEST_ADDS_IP1")
	if addsIp1 == "" {
		return nil, fmt.Errorf("`ARM_TEST_ADDS_IP1` is not set")
	}

	addsIp2 := os.Getenv("ARM_TEST_ADDS_IP2")
	if addsIp2 == "" {
		return nil, fmt.Errorf("`ARM_TEST_ADDS_IP2` is not set")
	}

	password := os.Getenv("ARM_TEST_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("`ARM_TEST_PASSWORD` is not set")
	}

	return &DomainServiceTrustResource{
		AADDSResourceGroupName: aaddsRgName,
		AADDSName:              aaddsName,
		ADDSFqdn:               addsFqdn,
		ADDSIp1:                addsIp1,
		ADDSIp2:                addsIp2,
		TrustPassword:          password,
	}, nil
}

func TestAccDomainServiceTrust_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_active_directory_domain_service_trust", "test")
	r, err := NewDomainServiceTrustResource()
	if err != nil {
		t.Skipf("Skipping: %v", err)
	}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccDomainServiceTrust_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_active_directory_domain_service_trust", "test")
	r, err := NewDomainServiceTrustResource()
	if err != nil {
		t.Skipf("Skipping: %v", err)
	}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r DomainServiceTrustResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.DomainServices.DomainServicesClient

	id, err := parse.DomainServiceTrustID(state.ID)
	if err != nil {
		return nil, err
	}

	idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName)

	resp, err := client.Get(ctx, idsdk)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, err
	}

	model := resp.Model
	if model == nil {
		return nil, fmt.Errorf("reading %s: returned with null model", idsdk)
	}

	existingTrusts := []domainservices.ForestTrust{}
	if props := model.Properties; props != nil {
		if fsettings := props.ResourceForestSettings; fsettings != nil {
			if settings := fsettings.Settings; settings != nil {
				existingTrusts = *settings
			}
		}
	}
	for _, setting := range existingTrusts {
		if setting.FriendlyName != nil && *setting.FriendlyName == id.TrustName {
			return utils.Bool(true), nil
		}
	}

	return utils.Bool(false), nil
}

func (r DomainServiceTrustResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_active_directory_domain_service_trust" "test" {
  domain_service_id      = data.azurerm_active_directory_domain_service.test.id
  name                   = "acctest-trust-%s"
  trusted_domain_fqdn    = %q
  trusted_domain_dns_ips = [%q, %q]
  password               = %q
}
`, template, data.RandomString, r.ADDSFqdn, r.ADDSIp1, r.ADDSIp2, r.TrustPassword)
}

func (r DomainServiceTrustResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_active_directory_domain_service_trust" "import" {
  domain_service_id      = azurerm_active_directory_domain_service_trust.test.domain_service_id
  name                   = azurerm_active_directory_domain_service_trust.test.name
  trusted_domain_fqdn    = azurerm_active_directory_domain_service_trust.test.trusted_domain_fqdn
  trusted_domain_dns_ips = azurerm_active_directory_domain_service_trust.test.trusted_domain_dns_ips
  password               = azurerm_active_directory_domain_service_trust.test.password
}
`, template)
}

func (r DomainServiceTrustResource) template(_ acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}


data "azurerm_active_directory_domain_service" "test" {
  name                = %q
  resource_group_name = %q
}
`, r.AADDSName, r.AADDSResourceGroupName)
}
