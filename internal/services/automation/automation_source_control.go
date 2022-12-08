package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
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
			ValidateFunc: validate.AutomationAccountID,
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
				string(automation.SourceTypeVsoGit),
				string(automation.SourceTypeVsoTfvc),
				string(automation.SourceTypeGitHub),
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
							string(automation.TokenTypeOauth),
							string(automation.TokenTypePersonalAccessToken),
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
			client := meta.Client.Automation.SourceControlClient

			var model SourceControlModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			accountID, _ := parse.AutomationAccountID(model.AutomationAccountID)
			id := parse.NewSourceControlID(subscriptionID, accountID.ResourceGroup, accountID.Name, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
			if !utils.ResponseWasNotFound(existing.Response) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			var param automation.SourceControlCreateOrUpdateParameters
			param.SourceControlCreateOrUpdateProperties = &automation.SourceControlCreateOrUpdateProperties{}
			param.RepoURL = utils.String(model.RepoURL)
			param.Branch = utils.String(model.Branch)
			param.FolderPath = utils.String(model.FolderPath)
			param.AutoSync = utils.Bool(model.AutoSync)
			param.PublishRunbook = utils.Bool(model.PublishRunbook)
			param.SourceType = automation.SourceType(model.SourceType)
			param.Description = utils.String(model.Description)

			param.SecurityToken = &automation.SourceControlSecurityTokenProperties{}
			if len(model.SecurityToken) > 0 {
				token := model.SecurityToken[0]
				param.SecurityToken.TokenType = automation.TokenType(token.TokenType)
				param.SecurityToken.AccessToken = utils.String(token.Token)
				if token.RefreshToken != "" {
					param.SecurityToken.RefreshToken = utils.String(token.RefreshToken)
				}
			}

			_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, param)
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
			id, err := parse.SourceControlID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			client := meta.Client.Automation.SourceControlClient
			result, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
			if err != nil {
				return err
			}
			if result.SourceControlProperties == nil {
				return fmt.Errorf("retrieving resource control with nil SourceControlProperties")
			}
			var output SourceControlModel
			if err := meta.Decode(&output); err != nil {
				return err
			}
			output.AutomationAccountID = parse.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroup, id.AutomationAccountName).ID()
			output.Name = id.Name
			output.RepoURL = utils.NormalizeNilableString(result.RepoURL)
			output.Branch = utils.NormalizeNilableString(result.Branch)
			output.FolderPath = utils.NormalizeNilableString(result.FolderPath)
			output.AutoSync = utils.NormaliseNilableBool(result.AutoSync)
			output.PublishRunbook = utils.NormaliseNilableBool(result.PublishRunbook)
			output.SourceType = string(result.SourceType)
			output.Description = utils.NormalizeNilableString(result.Description)
			return meta.Encode(&output)
		},
	}
}

func (m SourceControlResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.SourceControlClient

			id, err := parse.SourceControlID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SourceControlModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			var upd automation.SourceControlUpdateParameters
			prop := &automation.SourceControlUpdateProperties{}
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

			if meta.ResourceData.HasChange("security") {
				prop.SecurityToken = &automation.SourceControlSecurityTokenProperties{
					AccessToken:  utils.String(model.SecurityToken[0].TokenType),
					RefreshToken: utils.String(model.SecurityToken[0].RefreshToken),
					TokenType:    automation.TokenType(model.SecurityToken[0].TokenType),
				}
			}
			upd.SourceControlUpdateProperties = prop
			if _, err = client.Update(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m SourceControlResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.SourceControlID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Automation.SourceControlClient
			if _, err = client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m SourceControlResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SourceControlID
}
