package sourcecontrol

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-01-15/web"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppServiceGitHubTokenResource struct{}

type AppServiceGitHubTokenModel struct {
	Token string `tfschema:"token"`
}

var _ sdk.ResourceWithUpdate = AppServiceGitHubTokenResource{}
var _ sdk.Resource = AppServiceGitHubTokenResource{}

func (r AppServiceGitHubTokenResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"token": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r AppServiceGitHubTokenResource) Attributes() map[string]*pluginsdk.Schema {
	return nil
}

func (r AppServiceGitHubTokenResource) ModelObject() interface{} {
	return AppServiceSourceControlModel{}
}

func (r AppServiceGitHubTokenResource) ResourceType() string {
	return "azurerm_app_service_github_token"
}

func (r AppServiceGitHubTokenResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var appServiceGitHubToken AppServiceGitHubTokenModel

			if err := metadata.Decode(&appServiceGitHubToken); err != nil {
				return err
			}

			client := metadata.Client.AppService.BaseClient

			existing, err := client.GetSourceControl(ctx, "GitHub")
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for existing GitHub Token configuration")
				}
			}
			if existing.SourceControlProperties != nil && existing.SourceControlProperties.Token != nil && *existing.SourceControlProperties.Token != "" {
				return metadata.ResourceRequiresImport(r.ResourceType(), parse.AppServiceGitHubTokenId{})
			}

			sourceControlOAuth := web.SourceControl{
				SourceControlProperties: &web.SourceControlProperties{
					Token: utils.String(appServiceGitHubToken.Token),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, "GitHub", sourceControlOAuth); err != nil {
				return err
			}

			id := parse.AppServiceGitHubTokenId{}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r AppServiceGitHubTokenResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.BaseClient

			resp, err := client.GetSourceControl(ctx, "GitHub")
			if err != nil || resp.SourceControlProperties == nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(parse.AppServiceGitHubTokenId{})
				}
				return fmt.Errorf("reading App Service Source Control GitHub Token")
			}

			state := AppServiceGitHubTokenModel{}

			state.Token = utils.NormalizeNilableString(resp.Token)

			return metadata.Encode(&state)
		},
	}
}

func (r AppServiceGitHubTokenResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.BaseClient

			sourceControlOAuth := web.SourceControl{
				SourceControlProperties: &web.SourceControlProperties{
					Token: utils.String(""),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, "GitHub", sourceControlOAuth); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r AppServiceGitHubTokenResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SourceControlGitHubTokenID
}

func (r AppServiceGitHubTokenResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var appServiceGitHubToken AppServiceGitHubTokenModel

			if err := metadata.Decode(&appServiceGitHubToken); err != nil {
				return err
			}

			client := metadata.Client.AppService.BaseClient

			sourceControlOAuth := web.SourceControl{
				SourceControlProperties: &web.SourceControlProperties{
					Token: utils.String(appServiceGitHubToken.Token),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, "GitHub", sourceControlOAuth); err != nil {
				return err
			}

			return nil
		},
	}
}
