// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managementgroup_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managementgroups/2020-05-01/managementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagementGroupSubscriptionAssociation struct{}

// NOTE: this is a combined test rather than separate split out tests due to
// all testcases in this file share the same subscription instance so that
// these testcases have to be run sequentially.

func TestAccManagementGroupSubscriptionAssociation(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"Resource": {
			"basic":          testAccManagementGroupSubscriptionAssociation_basic,
			"requiresImport": testAccManagementGroupSubscriptionAssociation_requiresImport,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccManagementGroupSubscriptionAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_subscription_association", "test")

	r := ManagementGroupSubscriptionAssociation{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func testAccManagementGroupSubscriptionAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_subscription_association", "test")

	r := ManagementGroupSubscriptionAssociation{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r ManagementGroupSubscriptionAssociation) basic() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "test" {
  subscription_id = %q
}

resource "azurerm_management_group" "test" {
}

resource "azurerm_management_group_subscription_association" "test" {
  management_group_id = azurerm_management_group.test.id
  subscription_id     = data.azurerm_subscription.test.id
}
`, os.Getenv("ARM_SUBSCRIPTION_ID_ALT"))
}

func (r ManagementGroupSubscriptionAssociation) requiresImport(_ acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_management_group_subscription_association" "import" {
  management_group_id = azurerm_management_group_subscription_association.test.management_group_id
  subscription_id     = azurerm_management_group_subscription_association.test.subscription_id
}
`, r.basic())
}

func (r ManagementGroupSubscriptionAssociation) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managementgroups.ParseSubscriptionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ManagementGroups.GroupsClient.Get(ctx, commonids.NewManagementGroupID(id.GroupId), managementgroups.GetOperationOptions{
		CacheControl: pointer.To("no-cache"),
		Expand:       pointer.To(managementgroups.ExpandChildren),
		Recurse:      pointer.FromBool(false),
	})
	if err != nil {
		return nil, fmt.Errorf("retrieving Management Group to check for Subscription Association: %+v", err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Children == nil {
		return utils.Bool(false), nil
	}

	present := false
	for _, v := range *resp.Model.Properties.Children {
		if v.Type != nil && *v.Type == managementgroups.ManagementGroupChildTypeSubscriptions && v.Name != nil && strings.EqualFold(*v.Name, id.SubscriptionId) {
			present = true
		}
	}

	return utils.Bool(present), nil
}
