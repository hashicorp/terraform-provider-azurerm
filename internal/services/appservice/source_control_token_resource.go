// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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

			id := resourceproviders.NewSourceControlID(sourceControlToken.Type)

			client := metadata.Client.AppService.ResourceProvidersClient

			existing, err := client.GetSourceControl(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model.Properties != nil && pointer.From(existing.Model.Properties.Token) != "" {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			sourceControlOAuth := resourceproviders.SourceControl{
				Properties: &resourceproviders.SourceControlProperties{
					Token:       pointer.To(sourceControlToken.Token),
					TokenSecret: pointer.To(sourceControlToken.TokenSecret),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, id, sourceControlOAuth); err != nil {
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
			client := metadata.Client.AppService.ResourceProvidersClient

			id, err := resourceproviders.ParseSourceControlID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetSourceControl(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := AppServiceSourceControlTokenModel{}

			state.Type = id.SourceControlName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Token = pointer.From(props.Token)
					state.TokenSecret = pointer.From(props.TokenSecret)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AppServiceSourceControlTokenResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ResourceProvidersClient

			id, err := resourceproviders.ParseSourceControlID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			sourceControlOAuth := resourceproviders.SourceControl{
				Properties: &resourceproviders.SourceControlProperties{
					Token:       pointer.To(""),
					TokenSecret: pointer.To(""),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, *id, sourceControlOAuth); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
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

			id, err := resourceproviders.ParseSourceControlID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := metadata.Decode(&sourceControlToken); err != nil {
				return err
			}

			client := metadata.Client.AppService.ResourceProvidersClient

			sourceControlOAuth := resourceproviders.SourceControl{
				Properties: &resourceproviders.SourceControlProperties{
					Token:       pointer.To(sourceControlToken.Token),
					TokenSecret: pointer.To(sourceControlToken.TokenSecret),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, *id, sourceControlOAuth); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
