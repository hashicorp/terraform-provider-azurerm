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
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.SourceControl

			var model SourceControlModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionID := metadata.Client.Account.SubscriptionId
			accountID, _ := sourcecontrol.ParseAutomationAccountID(model.AutomationAccountID)
			id := sourcecontrol.NewSourceControlID(subscriptionID, accountID.ResourceGroupName, accountID.AutomationAccountName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			var param sourcecontrol.SourceControlCreateOrUpdateParameters
			param.Properties = sourcecontrol.SourceControlCreateOrUpdateProperties{
				AutoSync:       pointer.To(model.AutoSync),
				Branch:         pointer.To(model.Branch),
				Description:    pointer.To(model.Description),
				FolderPath:     pointer.To(model.FolderPath),
				PublishRunbook: pointer.To(model.PublishRunbook),
				RepoURL:        pointer.To(model.RepoURL),
				SourceType:     pointer.ToEnum[sourcecontrol.SourceType](model.SourceType),
			}

			param.Properties.SecurityToken = &sourcecontrol.SourceControlSecurityTokenProperties{}
			if len(model.SecurityToken) > 0 {
				token := model.SecurityToken[0]
				param.Properties.SecurityToken.TokenType = pointer.ToEnum[sourcecontrol.TokenType](token.TokenType)
				param.Properties.SecurityToken.AccessToken = pointer.To(token.Token)
				if token.RefreshToken != "" {
					param.Properties.SecurityToken.RefreshToken = pointer.To(token.RefreshToken)
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}
			metadata.SetID(id)
			return nil
		},
	}
}
