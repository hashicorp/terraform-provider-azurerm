// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_variable_datetime_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/variable"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func testCheckAzureRMAutomationVariableExists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState, varType string) (*bool, error) {
	id, err := variable.ParseVariableID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Automation.Variable.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Automation %s Variable %q (Automation Account Name %q / Resource Group %q) does not exist", varType, id.VariableName, id.AutomationAccountName, id.ResourceGroupName)
	}

	return pointer.To(resp.Model != nil), nil
}
