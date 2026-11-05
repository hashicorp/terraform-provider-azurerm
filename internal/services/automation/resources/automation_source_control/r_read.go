// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_source_control

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/sourcecontrol"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_source_control/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := sourcecontrol.ParseSourceControlID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Automation.SourceControl
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: got nil model", *id)
			}

			var output SourceControlModel
			if err := meta.Decode(&output); err != nil {
				return err
			}
			output.AutomationAccountID = sourcecontrol.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName).ID()
			output.Name = id.SourceControlName

			if props := resp.Model.Properties; props != nil {
				output.RepoURL = pointer.From(props.RepoURL)
				output.Branch = pointer.From(props.Branch)
				output.FolderPath = pointer.From(props.FolderPath)
				output.AutoSync = pointer.From(props.AutoSync)
				output.PublishRunbook = pointer.From(props.PublishRunbook)
				sourceType := ""
				if props.SourceType != nil {
					sourceType = string(*props.SourceType)
				}
				output.SourceType = sourceType
				output.Description = pointer.From(props.Description)
			}

			return meta.Encode(&output)
		},
	}
}

type Security struct {
	Token        string `tfschema:"token"`
	RefreshToken string `tfschema:"refresh_token"`
	TokenType    string `tfschema:"token_type"`
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.AutomationSourceControlV0ToV1{},
		},
	}
}
