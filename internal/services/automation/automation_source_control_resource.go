// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/sourcecontrol"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type Security struct {
	Token        string `tfschema:"token"`
	RefreshToken string `tfschema:"refresh_token"`
	TokenType    string `tfschema:"token_type"`
}

type SourceControlModel struct {
	AutomationAccountID string     `tfschema:"automation_account_id"`
	Name                string     `tfschema:"name"`
	RepoURL             string     `tfschema:"repository_url"`
	Branch              string     `tfschema:"branch"`
	FolderPath          string     `tfschema:"folder_path"`
	AutoSync            bool       `tfschema:"automatic_sync"`
	PublishRunbook      bool       `tfschema:"publish_runbook_enabled"`
	SourceType          string     `tfschema:"source_control_type"`
	Description         string     `tfschema:"description"`
	SecurityToken       []Security `tfschema:"security"`
}

type SourceControlResource struct{}

var _ sdk.Resource = (*SourceControlResource)(nil)

func (r SourceControlResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.AutomationSourceControlV0ToV1{},
		},
	}
}

func (m SourceControlResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"automation_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: sourcecontrol.ValidateAutomationAccountID,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"repository_url": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(0, 2000),
		},

		"branch": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 255),
		},

		"folder_path": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(0, 255),
		},

		"automatic_sync": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"publish_runbook_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"source_control_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sourcecontrol.SourceTypeVsoGit),
				string(sourcecontrol.SourceTypeVsoTfvc),
				string(sourcecontrol.SourceTypeGitHub),
			}, true),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"security": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"token": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(0, 1024),
					},

					"refresh_token": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(0, 1024),
					},

					"token_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(sourcecontrol.TokenTypeOauth),
							string(sourcecontrol.TokenTypePersonalAccessToken),
						}, false),
					},
				},
			},
		},
	}
}

func (m SourceControlResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (m SourceControlResource) ModelObject() interface{} {
	return &SourceControlModel{}
}

func (m SourceControlResource) ResourceType() string {
	return "azurerm_automation_source_control"
}

func (m SourceControlResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.SourceControl

			var model SourceControlModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			accountID, _ := sourcecontrol.ParseAutomationAccountID(model.AutomationAccountID)
			id := sourcecontrol.NewSourceControlID(subscriptionID, accountID.ResourceGroupName, accountID.AutomationAccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			sourceType := sourcecontrol.SourceType(model.SourceType)

			var param sourcecontrol.SourceControlCreateOrUpdateParameters
			param.Properties = sourcecontrol.SourceControlCreateOrUpdateProperties{
				AutoSync:       utils.Bool(model.AutoSync),
				Branch:         utils.String(model.Branch),
				Description:    utils.String(model.Description),
				FolderPath:     utils.String(model.FolderPath),
				PublishRunbook: utils.Bool(model.PublishRunbook),
				RepoUrl:        utils.String(model.RepoURL),
				SourceType:     &sourceType,
			}

			param.Properties.SecurityToken = &sourcecontrol.SourceControlSecurityTokenProperties{}
			if len(model.SecurityToken) > 0 {
				token := model.SecurityToken[0]
				tokenType := sourcecontrol.TokenType(token.TokenType)
				param.Properties.SecurityToken.TokenType = &tokenType
				param.Properties.SecurityToken.AccessToken = utils.String(token.Token)
				if token.RefreshToken != "" {
					param.Properties.SecurityToken.RefreshToken = utils.String(token.RefreshToken)
				}
			}

			_, err = client.CreateOrUpdate(ctx, id, param)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}
			meta.SetID(id)
			return nil
		},
	}
}

func (m SourceControlResource) Read() sdk.ResourceFunc {
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
				output.RepoURL = utils.NormalizeNilableString(props.RepoUrl)
				output.Branch = utils.NormalizeNilableString(props.Branch)
				output.FolderPath = utils.NormalizeNilableString(props.FolderPath)
				output.AutoSync = utils.NormaliseNilableBool(props.AutoSync)
				output.PublishRunbook = utils.NormaliseNilableBool(props.PublishRunbook)
				sourceType := ""
				if props.SourceType != nil {
					sourceType = string(*props.SourceType)
				}
				output.SourceType = sourceType
				output.Description = utils.NormalizeNilableString(props.Description)
			}

			return meta.Encode(&output)
		},
	}
}

func (m SourceControlResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.SourceControl

			id, err := sourcecontrol.ParseSourceControlID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SourceControlModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			var upd sourcecontrol.SourceControlUpdateParameters
			prop := &sourcecontrol.SourceControlUpdateProperties{}
			if meta.ResourceData.HasChange("branch") {
				prop.Branch = utils.String(model.Branch)
			}
			if meta.ResourceData.HasChange("folder_path") {
				prop.FolderPath = utils.String(model.FolderPath)
			}
			if meta.ResourceData.HasChange("automatic_sync") {
				prop.AutoSync = utils.Bool(model.AutoSync)
			}
			if meta.ResourceData.HasChange("folder_path") {
				prop.FolderPath = utils.String(model.FolderPath)
			}
			if meta.ResourceData.HasChange("publish_runbook_enabled") {
				prop.PublishRunbook = utils.Bool(model.PublishRunbook)
			}
			if meta.ResourceData.HasChange("description") {
				prop.Description = utils.String(model.Description)
			}

			tokenType := sourcecontrol.TokenType(model.SecurityToken[0].TokenType)
			if meta.ResourceData.HasChange("security") {
				prop.SecurityToken = &sourcecontrol.SourceControlSecurityTokenProperties{
					AccessToken:  utils.String(model.SecurityToken[0].TokenType),
					RefreshToken: utils.String(model.SecurityToken[0].RefreshToken),
					TokenType:    &tokenType,
				}
			}
			upd.Properties = prop
			if _, err = client.Update(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", *id, err)
			}

			return nil
		},
	}
}

func (m SourceControlResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := sourcecontrol.ParseSourceControlID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", *id)
			client := meta.Client.Automation.SourceControl
			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", *id, err)
			}
			return nil
		},
	}
}

func (m SourceControlResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sourcecontrol.ValidateSourceControlID
}
