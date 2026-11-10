// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_runtime_environment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runtimeenvironment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.RuntimeEnvironment
			id, err := runtimeenvironment.ParseRuntimeEnvironmentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return meta.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AutomationRuntimeEnvironmentResourceModel{
				Name:                id.RuntimeEnvironmentName,
				AutomationAccountId: runtimeenvironment.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = model.Location

				if model.Properties != nil {
					if model.Properties.DefaultPackages != nil {
						state.RuntimeDefaultPackages = *model.Properties.DefaultPackages
					}

					state.Description = pointer.From(model.Properties.Description)

					if model.Properties.Runtime != nil {
						state.RuntimeLanguage = pointer.From(model.Properties.Runtime.Language)
						state.RuntimeVersion = pointer.From(model.Properties.Runtime.Version)
					}
				}

				if model.Tags != nil {
					state.Tags = *model.Tags
				}
			}

			return meta.Encode(&state)
		},
	}
}
