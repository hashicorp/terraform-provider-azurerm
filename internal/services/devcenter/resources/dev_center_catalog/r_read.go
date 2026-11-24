// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/catalogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Catalogs

			id, err := catalogs.ParseDevCenterCatalogID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing catalog id: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := DevCenterCatalogsResourceModel{
				Name:              id.CatalogName,
				ResourceGroupName: id.ResourceGroupName,
			}
			state.DevCenterID = catalogs.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID()

			if properties := model.Properties; properties != nil {
				if gitHub := properties.GitHub; gitHub != nil {
					state.CatalogGitHub = []CatalogPropertiesModel{
						{
							URI:            pointer.From(gitHub.Uri),
							Branch:         pointer.From(gitHub.Branch),
							KeyVaultKeyUrl: pointer.From(gitHub.SecretIdentifier),
							Path:           pointer.From(gitHub.Path),
						},
					}
				}

				if adoGit := properties.AdoGit; adoGit != nil {
					state.CatalogAdoGit = []CatalogPropertiesModel{
						{
							URI:            pointer.From(adoGit.Uri),
							Branch:         pointer.From(adoGit.Branch),
							KeyVaultKeyUrl: pointer.From(adoGit.SecretIdentifier),
							Path:           pointer.From(adoGit.Path),
						},
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func CatalogPropertiesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"uri": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"branch": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"key_vault_key_url": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"path": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
			},
		},
	}
}
