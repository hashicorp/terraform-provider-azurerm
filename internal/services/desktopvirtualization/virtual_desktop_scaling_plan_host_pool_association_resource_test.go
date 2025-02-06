// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualDesktopScalingPlanAssociationResource struct{}

func TestAccVirtualDesktopScalingPlanAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_scaling_plan_host_pool_association", "test")
	r := VirtualDesktopScalingPlanAssociationResource{}
	roleAssignmentId := uuid.New().String()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, roleAssignmentId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccVirtualDesktopScalingPlanAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_scaling_plan_host_pool_association", "test")
	r := VirtualDesktopScalingPlanAssociationResource{}
	roleAssignmentId := uuid.New().String()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, roleAssignmentId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, roleAssignmentId),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_desktop_scaling_plan_host_pool_association"),
		},
	})
}

func (VirtualDesktopScalingPlanAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ScalingPlanHostPoolAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DesktopVirtualization.ScalingPlansClient.Get(ctx, id.ScalingPlan)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	found := false
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.HostPoolReferences != nil {
			for _, hostpool := range *props.HostPoolReferences {
				if strings.EqualFold(*hostpool.HostPoolArmPath, id.HostPool.ID()) {
					found = true
				}
			}
		}
	}

	return utils.Bool(found), nil
}

func (VirtualDesktopScalingPlanAssociationResource) basic(data acceptance.TestData, roleAssignmentId string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_scaling_plan_host_pool_association" "test" {
  host_pool_id    = azurerm_virtual_desktop_host_pool.test.id
  scaling_plan_id = azurerm_virtual_desktop_scaling_plan.test.id
  enabled         = true
  depends_on      = [azurerm_role_assignment.test]
}


`, VirtualDesktopScalingPlanResource{}.basic(data, roleAssignmentId))
}

func (r VirtualDesktopScalingPlanAssociationResource) requiresImport(data acceptance.TestData, roleAssignmentId string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_scaling_plan_host_pool_association" "import" {
  host_pool_id    = azurerm_virtual_desktop_host_pool.test.id
  scaling_plan_id = azurerm_virtual_desktop_scaling_plan.test.id
  enabled         = true
  depends_on      = [azurerm_role_assignment.test]
}


`, r.basic(data, roleAssignmentId))
}
