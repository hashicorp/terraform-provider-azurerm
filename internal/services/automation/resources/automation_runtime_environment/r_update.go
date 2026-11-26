// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_runtime_environment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runtimeenvironment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.RuntimeEnvironment
			id, err := runtimeenvironment.ParseRuntimeEnvironmentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var state AutomationRuntimeEnvironmentResourceModel
			if err = meta.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			param := runtimeenvironment.RuntimeEnvironmentUpdateParameters{
				Properties: &runtimeenvironment.RuntimeEnvironmentUpdateProperties{
					DefaultPackages: pointer.To(state.RuntimeDefaultPackages),
				},
			}

			if _, err = client.Update(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
