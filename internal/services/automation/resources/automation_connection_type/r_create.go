// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_connection_type

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/automationaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connectiontype"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.ConnectionType
			connClient := metadata.Client.Automation.AutomationAccount

			var model AutomationConnectionTypeModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}
			subscriptionID := metadata.Client.Account.SubscriptionId

			accountID := automationaccount.NewAutomationAccountID(subscriptionID, model.ResourceGroup, model.AutomationAccountName)
			account, err := connClient.Get(ctx, accountID)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", accountID, err)
			}
			if response.WasNotFound(account.HttpResponse) {
				return fmt.Errorf("%s was not found", accountID)
			}

			id := connectiontype.NewConnectionTypeID(accountID.SubscriptionId, model.ResourceGroup, model.AutomationAccountName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}
			param := connectiontype.ConnectionTypeCreateOrUpdateParameters{
				Name: model.Name,
				Properties: connectiontype.ConnectionTypeCreateOrUpdateProperties{
					IsGlobal:         pointer.To(model.IsGlobal),
					FieldDefinitions: map[string]connectiontype.FieldDefinition{},
				},
			}
			for _, field := range model.Field {
				param.Properties.FieldDefinitions[field.Name] = connectiontype.FieldDefinition{
					IsEncrypted: pointer.To(field.IsEncrypted),
					IsOptional:  pointer.To(field.IsOptional),
					Type:        field.Type,
				}
			}
			if _, err = client.CreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
