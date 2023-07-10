// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

type AppServiceSourceControlTokenResource struct{}

type AppServiceSourceControlTokenModel struct {
	Token       string `tfschema:"token"`
	TokenSecret string `tfschema:"token_secret"`
	Type        string `tfschema:"type"`
}

var (
	_ sdk.ResourceWithUpdate = AppServiceSourceControlTokenResource{}
	_ sdk.Resource           = AppServiceSourceControlTokenResource{}
)

func (r AppServiceSourceControlTokenResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Bitbucket",
				"Dropbox",
				"GitHub",
				"OneDrive",
			}, false),
		},

		"token": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"token_secret": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r AppServiceSourceControlTokenResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AppServiceSourceControlTokenResource) ModelObject() interface{} {
	return &AppServiceSourceControlTokenModel{}
}

func (r AppServiceSourceControlTokenResource) ResourceType() string {
	return "azurerm_source_control_token"
}

func (r AppServiceSourceControlTokenResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var sourceControlToken AppServiceSourceControlTokenModel

			if err := metadata.Decode(&sourceControlToken); err != nil {
				return err
			}

			id := parse.NewAppServiceSourceControlTokenID(sourceControlToken.Type)

			client := metadata.Client.AppService.BaseClient

			existing, err := client.GetSourceControl(ctx, id.Type)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("%s not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.SourceControlProperties != nil && existing.SourceControlProperties.Token != nil && *existing.SourceControlProperties.Token != "" {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			sourceControlOAuth := web.SourceControl{
				SourceControlProperties: &web.SourceControlProperties{
					Token:       utils.String(sourceControlToken.Token),
					TokenSecret: utils.String(sourceControlToken.TokenSecret),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, id.Type, sourceControlOAuth); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r AppServiceSourceControlTokenResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.BaseClient

			id, err := parse.AppServiceSourceControlTokenID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetSourceControl(ctx, id.Type)
			if err != nil || resp.SourceControlProperties == nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			state := AppServiceSourceControlTokenModel{}

			state.Type = id.Type

			if resp.SourceControlProperties != nil {
				state.Token = utils.NormalizeNilableString(resp.Token)
				state.TokenSecret = utils.NormalizeNilableString(resp.TokenSecret)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AppServiceSourceControlTokenResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.BaseClient

			id, err := parse.AppServiceSourceControlTokenID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			sourceControlOAuth := web.SourceControl{
				SourceControlProperties: &web.SourceControlProperties{
					Token:       utils.String(""),
					TokenSecret: utils.String(""),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, id.Type, sourceControlOAuth); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r AppServiceSourceControlTokenResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AppServiceSourceControlTokenID
}

func (r AppServiceSourceControlTokenResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var sourceControlToken AppServiceSourceControlTokenModel

			id, err := parse.AppServiceSourceControlTokenID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := metadata.Decode(&sourceControlToken); err != nil {
				return err
			}

			client := metadata.Client.AppService.BaseClient

			sourceControlOAuth := web.SourceControl{
				SourceControlProperties: &web.SourceControlProperties{
					Token:       utils.String(sourceControlToken.Token),
					TokenSecret: utils.String(sourceControlToken.TokenSecret),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, id.Type, sourceControlOAuth); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
