// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PrivateEndpointApplicationSecurityGroupAssociationResource struct{}

func TestAccPrivateEndpointApplicationSecurityGroupAssociationResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint_application_security_group_association", "test")
	r := PrivateEndpointApplicationSecurityGroupAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateEndpointApplicationSecurityGroupAssociationResource_updatePrivateEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint_application_security_group_association", "test")
	r := PrivateEndpointApplicationSecurityGroupAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// Ensure a subsequent update to the PrivateEndpoint does not affect the association
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccPrivateEndpointApplicationSecurityGroupAssociationResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint_application_security_group_association", "test")
	r := PrivateEndpointApplicationSecurityGroupAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_private_endpoint_application_security_group_association"),
		},
	})
}

func TestAccPrivateEndpointApplicationSecurityGroupAssociationResource_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint_application_security_group_association", "test")
	r := PrivateEndpointApplicationSecurityGroupAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentionally not using a DisappearsStep as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.destroy),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func (r PrivateEndpointApplicationSecurityGroupAssociationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	splitId := strings.Split(state.ID, "|")
	exists := false

	if len(splitId) != 2 {
		return &exists, fmt.Errorf("expected ID to be in the format {PrivateEndpointId}|{ApplicationSecurityGroupId} but got %q", state.ID)
	}

	endpointId, err := privateendpoints.ParsePrivateEndpointID(splitId[0])
	if err != nil {
		return &exists, err
	}

	securityGroupId, err := applicationsecuritygroups.ParseApplicationSecurityGroupID(splitId[1])
	if err != nil {
		return &exists, err
	}

	if endpointId == nil || securityGroupId == nil {
		return &exists, fmt.Errorf("parse error, both PrivateEndpointId and ApplicationSecurityGroupId should not be nil")
	}

	privateEndpointClient := client.Network.PrivateEndpoints
	existingPrivateEndpoint, err := privateEndpointClient.Get(ctx, *endpointId, privateendpoints.DefaultGetOperationOptions())
	if err != nil && !response.WasNotFound(existingPrivateEndpoint.HttpResponse) {
		return &exists, fmt.Errorf("checking for the presence of existing PrivateEndpoint %q: %+v", endpointId, err)
	}

	if response.WasNotFound(existingPrivateEndpoint.HttpResponse) {
		return &exists, fmt.Errorf("PrivateEndpoint %q does not exsits", endpointId)
	}

	input := existingPrivateEndpoint
	if input.Model != nil && input.Model.Properties != nil && input.Model.Properties.ApplicationSecurityGroups != nil {
		for _, value := range *input.Model.Properties.ApplicationSecurityGroups {
			if value.Id != nil && *value.Id == securityGroupId.ID() {
				exists = true
				break
			}
		}
	}
	return &exists, nil
}

func (r PrivateEndpointApplicationSecurityGroupAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = false
    private_connection_resource_id = azurerm_private_link_service.test.id
  }
}

resource "azurerm_application_security_group" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint_application_security_group_association" "test" {
  private_endpoint_id           = azurerm_private_endpoint.test.id
  application_security_group_id = azurerm_application_security_group.test.id
}
`, r.template(data, r.serviceAutoApprove(data)), data.RandomInteger, data.RandomInteger)
}

func (r PrivateEndpointApplicationSecurityGroupAssociationResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = false
    private_connection_resource_id = azurerm_private_link_service.test.id
  }

  tags = {
    "test" = "value1"
  }
}

resource "azurerm_application_security_group" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint_application_security_group_association" "test" {
  private_endpoint_id           = azurerm_private_endpoint.test.id
  application_security_group_id = azurerm_application_security_group.test.id
}
`, r.template(data, r.serviceAutoApprove(data)), data.RandomInteger, data.RandomInteger)
}

func (r PrivateEndpointApplicationSecurityGroupAssociationResource) serviceAutoApprove(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%d"
    primary   = true
    subnet_id = azurerm_subnet.service.id
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, data.RandomInteger, data.RandomInteger)
}

func (r PrivateEndpointApplicationSecurityGroupAssociationResource) template(data acceptance.TestData, seviceCfg string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-PEASGAsso-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name                 = "acctestsnetservice-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]

  private_link_service_network_policies_enabled = false
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  frontend_ip_configuration {
    name                 = azurerm_public_ip.test.name
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

%s
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, seviceCfg)
}

func (r PrivateEndpointApplicationSecurityGroupAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint_application_security_group_association" "import" {
  private_endpoint_id           = azurerm_private_endpoint.test.id
  application_security_group_id = azurerm_application_security_group.test.id
}
`, r.basic(data))
}

func (r PrivateEndpointApplicationSecurityGroupAssociationResource) destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(15*time.Minute))
	defer cancel()

	endpointId, err := privateendpoints.ParsePrivateEndpointID(state.Attributes["private_endpoint_id"])
	if err != nil {
		return err
	}

	securityGroupId, err := applicationsecuritygroups.ParseApplicationSecurityGroupID(state.Attributes["application_security_group_id"])
	if err != nil {
		return err
	}

	privateEndpointClient := client.Network.PrivateEndpoints

	existingPrivateEndpoint, err := privateEndpointClient.Get(ctx, *endpointId, privateendpoints.DefaultGetOperationOptions())
	if err != nil && !response.WasNotFound(existingPrivateEndpoint.HttpResponse) {
		return fmt.Errorf("checking for the presence of existing PrivateEndpoint %q: %+v", endpointId, err)
	}

	if response.WasNotFound(existingPrivateEndpoint.HttpResponse) {
		return fmt.Errorf("PrivateEndpoint %q does not exsits", endpointId)
	}

	if existingPrivateEndpoint.Model == nil || existingPrivateEndpoint.Model.Properties == nil || existingPrivateEndpoint.Model.Properties.ApplicationSecurityGroups == nil {
		return fmt.Errorf("model/properties/application security groups was missing for %s", endpointId)
	}

	// flag: application security group exists in private endpoint configuration
	ASGInPE := false

	input := existingPrivateEndpoint
	ASGList := *input.Model.Properties.ApplicationSecurityGroups
	newASGList := make([]privateendpoints.ApplicationSecurityGroup, 0)
	for idx, value := range ASGList {
		if value.Id != nil && *value.Id == securityGroupId.ID() {
			newASGList = append(newASGList, ASGList[:idx]...)
			newASGList = append(newASGList, ASGList[idx+1:]...)
			ASGInPE = true
			break
		}
	}
	if ASGInPE {
		input.Model.Properties.ApplicationSecurityGroups = &newASGList
	} else {
		return fmt.Errorf("deletion failed, ApplicationSecurityGroup %q does not linked with PrivateEndpoint %q", securityGroupId, endpointId)
	}

	if err = privateEndpointClient.CreateOrUpdateThenPoll(ctx, *endpointId, *input.Model); err != nil {
		return fmt.Errorf("creating %s: %+v", endpointId, err)
	}

	return nil
}
