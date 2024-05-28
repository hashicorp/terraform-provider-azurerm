// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IPGroupCidrResource struct{}

func TestAccIpGroupCidr_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group_cidr", "test")
	r := IPGroupCidrResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_ip_group_cidr.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIpGroupCidr_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group_cidr", "test")
	r := IPGroupCidrResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_ip_group_cidr.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_ip_group.test").Key("tags.env").HasValue("prod"),
				check.That("azurerm_ip_group_cidr.test").ExistsInAzure(r),
				check.That("azurerm_ip_group_cidr.multiple_1").ExistsInAzure(r),
				check.That("azurerm_ip_group_cidr.multiple_2").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIpGroupCidr_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group_cidr", "test")
	r := IPGroupCidrResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_ip_group_cidr.test").ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_ip_group_cidr"),
		},
	})
}

func (t IPGroupCidrResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IpGroupCidrID(state.ID)
	if err != nil {
		return nil, err
	}

	ipGroupId := ipgroups.NewIPGroupID(id.SubscriptionId, id.ResourceGroup, id.IpGroupName)

	resp, err := clients.Network.Client.IPGroups.Get(ctx, ipGroupId, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", ipGroupId)
	}
	if resp.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", ipGroupId)
	}

	if !utils.SliceContainsValue(*resp.Model.Properties.IPAddresses, state.Attributes["cidr"]) {
		return pointer.To(false), nil
	}

	return pointer.To(true), nil
}

func (IPGroupCidrResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_ip_group" "test" {
  name                = "acceptanceTestIpGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    env = "prod"
  }

  lifecycle {
    ignore_changes = ["cidrs"]
  }
}

resource "azurerm_ip_group_cidr" "test" {
  ip_group_id = azurerm_ip_group.test.id
  cidr        = "10.0.0.0/24"
}


`, data.RandomInteger, data.Locations.Primary)
}

func (r IPGroupCidrResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ip_group_cidr" "multiple_1" {
  ip_group_id = azurerm_ip_group.test.id
  cidr        = "10.10.0.0/24"
}

resource "azurerm_ip_group_cidr" "multiple_2" {
  ip_group_id = azurerm_ip_group.test.id
  cidr        = "10.20.0.0/24"
}
`, r.basic(data))
}

func (r IPGroupCidrResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ip_group_cidr" "import" {
  ip_group_id = azurerm_ip_group_cidr.test.ip_group_id
  cidr        = azurerm_ip_group_cidr.test.cidr
}
`, r.basic(data))
}
