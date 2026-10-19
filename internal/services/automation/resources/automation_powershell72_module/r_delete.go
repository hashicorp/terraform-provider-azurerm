// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_powershell72_module

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/module"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		// the Func returns a function which deletes the Resource Group
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.ModuleClientV2023
			id, err := module.ParsePowerShell72ModuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.PowerShell72ModuleDelete(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}

				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
