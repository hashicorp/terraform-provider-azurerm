// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_powershell72_module

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/module"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.ModuleClientV2023
			id, err := module.ParsePowerShell72ModuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.PowerShell72ModuleGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var output AutomationPowerShell72ModuleModel
			if err := metadata.Decode(&output); err != nil {
				return err
			}

			output.Name = id.PowerShell72ModuleName
			output.AutomationAccountID = module.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName).ID()
			if resp.Model != nil {
				output.Tags = tags.Flatten(resp.Model.Tags)
			}

			return metadata.Encode(&output)
		},
	}
}

type ModuleHash struct {
	Algorithm string `tfschema:"algorithm"`
	Value     string `tfschema:"value"`
}
