// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AppServiceSourceControlTokenDataSource struct{}

var _ sdk.DataSource = AppServiceSourceControlTokenDataSource{}

func (d AppServiceSourceControlTokenDataSource) Arguments() map[string]*pluginsdk.Schema {
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
	}
}

func (d AppServiceSourceControlTokenDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"token": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"token_secret": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},
	}
}

func (d AppServiceSourceControlTokenDataSource) ModelObject() interface{} {
	return &AppServiceSourceControlTokenModel{}
}

func (d AppServiceSourceControlTokenDataSource) ResourceType() string {
	return "azurerm_source_control_token"
}

func (d AppServiceSourceControlTokenDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ResourceProvidersClient

			var sourceControlToken AppServiceSourceControlTokenModel
			if err := metadata.Decode(&sourceControlToken); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			id := resourceproviders.NewSourceControlID(sourceControlToken.Type)

			resp, err := client.GetSourceControl(ctx, id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					sourceControlToken.Token = pointer.From(props.Token)
					sourceControlToken.TokenSecret = pointer.From(props.TokenSecret)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&sourceControlToken)
		},
	}
}
