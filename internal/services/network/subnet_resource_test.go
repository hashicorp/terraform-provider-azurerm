// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SubnetResource struct{}

func TestAccSubnet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

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

func TestAccSubnet_basic_addressPrefixes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_addressPrefixes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnet_complete_addressPrefixes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete_addressPrefixes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnet_update_addressPrefixes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_addressPrefixes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete_addressPrefixes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic_addressPrefixes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_subnet"),
		},
	})
}

func TestAccSubnet_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccSubnet_defaultOutbound(t *testing.T) {
	dataInternal := acceptance.BuildTestData(t, "azurerm_subnet", "internal")
	dataPublic := acceptance.BuildTestData(t, "azurerm_subnet", "public")
	r := SubnetResource{}

	dataInternal.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultOutbound(dataInternal),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(dataInternal.ResourceName).ExistsInAzure(r),
				check.That(dataInternal.ResourceName).Key("default_outbound_access_enabled").HasValue("false"),
				check.That(dataPublic.ResourceName).ExistsInAzure(r),
				check.That(dataPublic.ResourceName).Key("default_outbound_access_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccSubnet_delegation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.delegation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.delegationUpdated(data),
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
			Config: r.delegation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// TODO 4.0: Remove test
func TestAccSubnet_enablePrivateEndpointNetworkPolicies(t *testing.T) {
	if !features.FourPointOhBeta() {
		data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
		r := SubnetResource{}

		data.ResourceTest(t, r, []acceptance.TestStep{
			{
				Config: r.enablePrivateEndpointNetworkPolicies(data, true),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
			{
				Config: r.enablePrivateEndpointNetworkPolicies(data, false),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
			{
				Config: r.enablePrivateEndpointNetworkPolicies(data, true),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
		})
	} else {
		t.Skip("skipping due to deprecation of the 'private_endpoint_network_policies_enabled' fields in 4.0")
	}
}

func TestAccSubnet_privateEndpointNetworkPolicies(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateEndpointNetworkPolicies(data, "Enabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.privateEndpointNetworkPolicies(data, "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.privateEndpointNetworkPolicies(data, "NetworkSecurityGroupEnabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.privateEndpointNetworkPolicies(data, "RouteTableEnabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnet_enablePrivateLinkServiceNetworkPolicies(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enablePrivateLinkServiceNetworkPolicies(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enablePrivateLinkServiceNetworkPolicies(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enablePrivateLinkServiceNetworkPolicies(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// TODO 4.0: Remove test
func TestAccSubnet_enforcePrivateLinkEndpointNetworkPolicies(t *testing.T) {
	if !features.FourPointOhBeta() {
		data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
		r := SubnetResource{}

		data.ResourceTest(t, r, []acceptance.TestStep{
			{
				Config: r.enforcePrivateLinkEndpointNetworkPolicies(data, true),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
			{
				Config: r.enforcePrivateLinkEndpointNetworkPolicies(data, false),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
			{
				Config: r.enforcePrivateLinkEndpointNetworkPolicies(data, true),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
		})
	} else {
		t.Skip("@WodansSon: skipping due to deprecation of the 'enforce_private_link_endpoint_network_policies' field in 4.0")
	}
}

// TODO 4.0: Remove test
func TestAccSubnet_enforcePrivateLinkServiceNetworkPolicies(t *testing.T) {
	if !features.FourPointOhBeta() {
		data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
		r := SubnetResource{}

		data.ResourceTest(t, r, []acceptance.TestStep{
			{
				Config: r.enforcePrivateLinkServiceNetworkPolicies(data, true),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
			{
				Config: r.enforcePrivateLinkServiceNetworkPolicies(data, false),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
			{
				Config: r.enforcePrivateLinkServiceNetworkPolicies(data, true),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
		})
	} else {
		t.Skip("@WodansSon: skipping due to deprecation of the 'enforce_private_link_service_network_policies' field in 4.0")
	}
}

// TODO 4.0: Remove test
func TestAccSubnet_PrivateLinkPoliciesToggleWithEnforceFirst(t *testing.T) {
	if !features.FourPointOhBeta() {
		data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
		r := SubnetResource{}

		data.ResourceTest(t, r, []acceptance.TestStep{
			{
				Config: r.enforcePrivateLinkEndpointNetworkPolicies(data, true),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("true"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("false"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Disabled"),
				),
			},
			data.ImportStep(),
			{
				Config: r.enablePrivateEndpointNetworkPolicies(data, true),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Enabled"),
				),
			},
			data.ImportStep(),
			{
				Config: r.privateEndpointNetworkPolicies(data, "Enabled"),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Enabled"),
				),
			},
			data.ImportStep(),
			{
				Config: r.enforcePrivateLinkServiceNetworkPolicies(data, true),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Enabled"),
				),
			},
			data.ImportStep(),
			{
				Config: r.enablePrivateLinkServiceNetworkPolicies(data, true),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Enabled"),
				),
			},
			data.ImportStep(),
			{
				Config: r.privateEndpointNetworkPolicies(data, "Enabled"),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Enabled"),
				),
			},
			data.ImportStep(),
		})
	} else {
		t.Skip("@WodansSon: skipping due to deprecation of the 'enforce_private_link_endpoint_network_policies' and 'enforce_private_link_service_network_policies' fields in 4.0")
	}
}

// TODO 4.0: Remove test
func TestAccSubnet_PrivateLinkPoliciesToggleWithEnabledFirst(t *testing.T) {
	if !features.FourPointOhBeta() {
		data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
		r := SubnetResource{}

		data.ResourceTest(t, r, []acceptance.TestStep{
			{
				Config: r.privateEndpointNetworkPolicies(data, "Disabled"),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("true"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("false"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Disabled"),
				),
			},
			data.ImportStep(),
			{
				Config: r.enablePrivateEndpointNetworkPolicies(data, false),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("true"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("false"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Disabled"),
				),
			},
			data.ImportStep(),
			{
				Config: r.enforcePrivateLinkEndpointNetworkPolicies(data, false),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Enabled"),
				),
			},
			data.ImportStep(),
			{
				Config: r.enablePrivateLinkServiceNetworkPolicies(data, false),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Enabled"),
				),
			},
			data.ImportStep(),
			{
				Config: r.enforcePrivateLinkServiceNetworkPolicies(data, false),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Enabled"),
				),
			},
			data.ImportStep(),
		})
	} else {
		t.Skip("@WodansSon: skipping due to deprecation of the 'enforce_private_link_endpoint_network_policies' and 'enforce_private_link_service_network_policies' fields in 4.0")
	}
}

func TestAccSubnet_serviceEndpoints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceEndpoints(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceEndpointsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// remove them
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceEndpoints(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnet_serviceEndpointPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceEndpointPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceEndpointPolicyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceEndpointPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccSubnet_updateAddressPrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updatedAddressPrefix(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnet_privateLinkEndpointNetworkPoliciesValidateDefaultValues(t *testing.T) {
	if features.FourPointOhBeta() {
		data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
		r := SubnetResource{}

		data.ResourceTest(t, r, []acceptance.TestStep{
			{
				Config: r.privateLinkEndpointNetworkPoliciesDefaults(data),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Disabled"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("true"),
				),
			},
			data.ImportStep(),
		})
	} else {
		data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
		r := SubnetResource{}

		data.ResourceTest(t, r, []acceptance.TestStep{
			{
				Config: r.privateLinkEndpointNetworkPoliciesDefaults(data),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("enforce_private_link_endpoint_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("enforce_private_link_service_network_policies").HasValue("false"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_link_service_network_policies_enabled").HasValue("true"),
					check.That(data.ResourceName).Key("private_endpoint_network_policies").HasValue("Enabled"),
				),
			},
			data.ImportStep(),
		})
	}
}

func TestAccSubnet_updateServiceDelegation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")
	r := SubnetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateServiceDelegation(data, "NGINX.NGINXPLUS/nginxDeployments"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateServiceDelegation(data, "PaloAltoNetworks.Cloudngfw/firewalls"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateServiceDelegation(data, "Qumulo.Storage/fileSystems"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateServiceDelegationNetworkInterfaces(data, "Oracle.Database/networkAttachments"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t SubnetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseSubnetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.Client.Subnets.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("reading Subnet (%s): %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (SubnetResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseSubnetID(state.ID)
	if err != nil {
		return nil, err
	}

	if err := client.Network.Client.Subnets.DeleteThenPoll(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting Subnet %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (SubnetResource) hasNoNatGateway(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(15*time.Minute))
	defer cancel()

	id, err := commonids.ParseSubnetID(state.ID)
	if err != nil {
		return err
	}

	subnet, err := client.Network.Client.Subnets.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(subnet.HttpResponse) {
			return fmt.Errorf("%s does not exist", id)
		}
		return fmt.Errorf("Bad: Get on subnetClient: %+v", err)
	}

	model := subnet.Model
	if model == nil {
		return fmt.Errorf("model was nil for %s", id)
	}

	props := subnet.Model.Properties
	if props == nil {
		return fmt.Errorf("properties was nil for %s", id)
	}

	if props.NatGateway != nil && ((props.NatGateway.Id == nil) || (props.NatGateway.Id != nil && *props.NatGateway.Id == "")) {
		return fmt.Errorf("no Route Table should exist for %s but got %q", id, *props.RouteTable.Id)
	}
	return nil
}

func (SubnetResource) hasNoNetworkSecurityGroup(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(15*time.Minute))
	defer cancel()

	id, err := commonids.ParseSubnetID(state.ID)
	if err != nil {
		return err
	}

	resp, err := client.Network.Client.Subnets.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s does not exist", id)
		}

		return fmt.Errorf("Bad: Get on subnetClient: %+v", err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("model was nil for %s", id)
	}

	props := resp.Model.Properties
	if props == nil {
		return fmt.Errorf("properties was nil for %s", id)
	}

	if props.NetworkSecurityGroup != nil && ((props.NetworkSecurityGroup.Id == nil) || (props.NetworkSecurityGroup.Id != nil && *props.NetworkSecurityGroup.Id == "")) {
		return fmt.Errorf("no Network Security Group should exist for %s but got %q", id, *props.NetworkSecurityGroup.Id)
	}

	return nil
}

func (SubnetResource) hasNoRouteTable(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(15*time.Minute))
	defer cancel()

	id, err := commonids.ParseSubnetID(state.ID)
	if err != nil {
		return err
	}

	resp, err := client.Network.Client.Subnets.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s does not exist", id)
		}

		return fmt.Errorf("Bad: Get on subnetClient: %+v", err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("model was nil for %s", id)
	}

	props := resp.Model.Properties
	if props == nil {
		return fmt.Errorf("properties was nil for %s", id)
	}

	if props.RouteTable != nil && ((props.RouteTable.Id == nil) || (props.RouteTable.Id != nil && *props.RouteTable.Id == "")) {
		return fmt.Errorf("no Route Table should exist for %s but got %q", id, *props.RouteTable.Id)
	}

	return nil
}

func (r SubnetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_subnet" "test2" {
  name                 = "internal2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
}
`, r.template(data))
}

func (r SubnetResource) delegation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "first"

    service_delegation {
      name = "Microsoft.ContainerInstance/containerGroups"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/action",
      ]
    }
  }
}
`, r.template(data))
}

func (r SubnetResource) defaultOutbound(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_subnet" "internal" {
  name                            = "internal"
  resource_group_name             = azurerm_resource_group.test.name
  virtual_network_name            = azurerm_virtual_network.test.name
  address_prefixes                = ["10.0.2.0/24"]
  default_outbound_access_enabled = false
}
resource "azurerm_subnet" "public" {
  name                            = "public"
  resource_group_name             = azurerm_resource_group.test.name
  virtual_network_name            = azurerm_virtual_network.test.name
  address_prefixes                = ["10.0.3.0/24"]
  default_outbound_access_enabled = true
}
`, r.template(data))
}

func (r SubnetResource) delegationUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "first"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}
`, r.template(data))
}

// TODO 4.0: Remove test
func (r SubnetResource) enablePrivateEndpointNetworkPolicies(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  private_endpoint_network_policies = %t
}
`, r.template(data), enabled)
}

func (r SubnetResource) privateEndpointNetworkPolicies(data acceptance.TestData, enabled string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  private_endpoint_network_policies = "%s"
}
`, r.template(data), enabled)
}

func (r SubnetResource) enablePrivateLinkServiceNetworkPolicies(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  private_link_service_network_policies_enabled = %t
}
`, r.template(data), enabled)
}

// TODO 4.0: Remove test
func (r SubnetResource) enforcePrivateLinkEndpointNetworkPolicies(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  private_endpoint_network_policies = %t
}
`, r.template(data), enabled)
}

// TODO 4.0: Remove test
func (r SubnetResource) enforcePrivateLinkServiceNetworkPolicies(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  private_link_service_network_policies_enabled = %t
}
`, r.template(data), enabled)
}

func (r SubnetResource) privateLinkEndpointNetworkPoliciesDefaults(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}
`, r.template(data))
}

func (SubnetResource) basic_addressPrefixes(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-n-%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16", "ace:cab:deca::/48"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefixes     = ["10.0.0.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (SubnetResource) complete_addressPrefixes(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-n-%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16", "ace:cab:deca::/48"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefixes     = ["10.0.0.0/24", "ace:cab:deca:deed::/64"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r SubnetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "import" {
  name                 = azurerm_subnet.test.name
  resource_group_name  = azurerm_subnet.test.resource_group_name
  virtual_network_name = azurerm_subnet.test.virtual_network_name
  address_prefixes     = azurerm_subnet.test.address_prefixes
}
`, r.basic(data))
}

func (r SubnetResource) serviceEndpoints(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "internal2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
  service_endpoints    = ["Microsoft.Sql"]
}
`, r.template(data))
}

func (r SubnetResource) serviceEndpointsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Sql", "Microsoft.Storage"]
}

resource "azurerm_subnet" "test2" {
  name                 = "internal2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
  service_endpoints    = ["Microsoft.Sql", "Microsoft.Storage"]
}
`, r.template(data))
}

func (r SubnetResource) serviceEndpointPolicyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_service_endpoint_storage_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}
`, r.template(data), data.RandomInteger)
}

func (r SubnetResource) serviceEndpointPolicyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_service_endpoint_storage_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_subnet" "test" {
  name                        = "internal"
  resource_group_name         = azurerm_resource_group.test.name
  virtual_network_name        = azurerm_virtual_network.test.name
  address_prefixes            = ["10.0.2.0/24"]
  service_endpoints           = ["Microsoft.Sql"]
  service_endpoint_policy_ids = [azurerm_subnet_service_endpoint_storage_policy.test.id]
}
`, r.template(data), data.RandomInteger)
}

func (r SubnetResource) updatedAddressPrefix(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
}
`, r.template(data))
}

func (r SubnetResource) updateServiceDelegation(data acceptance.TestData, serviceName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "first"

    service_delegation {
      name = "%s"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}
`, r.template(data), serviceName)
}

func (r SubnetResource) updateServiceDelegationNetworkInterfaces(data acceptance.TestData, serviceName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "first"

    service_delegation {
      name = "%s"

      actions = [
        "Microsoft.Network/networkinterfaces/*",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}
`, r.template(data), serviceName)
}

func (SubnetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
