// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/raiblocklists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type cognitiveRaiBlocklistModel struct {
	Name               string `tfschema:"name"`
	CognitiveAccountId string `tfschema:"cognitive_account_id"`
	Description        string `tfschema:"description"`
}

type CognitiveRaiBlocklistResource struct{}

var _ sdk.Resource = CognitiveRaiBlocklistResource{}

func (c CognitiveRaiBlocklistResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cognitive_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: raiblocklists.ValidateAccountID,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (c CognitiveRaiBlocklistResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (c CognitiveRaiBlocklistResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiBlocklistsClient

			var model cognitiveRaiBlocklistModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId, err := raiblocklists.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id := raiblocklists.NewRaiBlocklistID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(c.ResourceType(), id)
			}

			payload := &raiblocklists.RaiBlocklist{
				Properties: &raiblocklists.RaiBlocklistProperties{},
			}

			if model.Description != "" {
				payload.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, *payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (c CognitiveRaiBlocklistResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiBlocklistsClient

			var model cognitiveRaiBlocklistModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId, err := raiblocklists.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id, err := raiblocklists.ParseRaiBlocklistID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (c CognitiveRaiBlocklistResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiBlocklistsClient

			id, err := raiblocklists.ParseRaiBlocklistID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			accountId := raiblocklists.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (c CognitiveRaiBlocklistResource) IDValidationFunc() func(interface{}, string) ([]string, []error) {
	return raiblocklists.ValidateRaiBlocklistID
}

func (c CognitiveRaiBlocklistResource) ModelObject() interface{} {
	return &cognitiveRaiBlocklistModel{}
}

func (c CognitiveRaiBlocklistResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiBlocklistsClient

			id, err := raiblocklists.ParseRaiBlocklistID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := cognitiveRaiBlocklistModel{
				Name:               id.RaiBlocklistName,
				CognitiveAccountId: raiblocklists.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
			}

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (c CognitiveRaiBlocklistResource) ResourceType() string {
	return "azurerm_cognitive_account_rai_blocklist"
}
