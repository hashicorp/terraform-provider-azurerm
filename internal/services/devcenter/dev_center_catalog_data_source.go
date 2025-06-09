// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/catalogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = DevCenterCatalogDataSource{}

type DevCenterCatalogDataSource struct{}

type DevCenterCatalogDataSourceModel struct {
	Name          string                             `tfschema:"name"`
	DevCenterID   string                             `tfschema:"dev_center_id"`
	CatalogGitHub []CatalogPropertiesDataSourceModel `tfschema:"catalog_github"`
	CatalogAdoGit []CatalogPropertiesDataSourceModel `tfschema:"catalog_adogit"`
}

type CatalogPropertiesDataSourceModel struct {
	URI            string `tfschema:"uri"`
	Branch         string `tfschema:"branch"`
	KeyVaultKeyUrl string `tfschema:"key_vault_key_url"`
	Path           string `tfschema:"path"`
}

func (DevCenterCatalogDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dev_center_id": commonschema.ResourceIDReferenceRequired(&devcenters.DevCenterId{}),
	}
}

func (DevCenterCatalogDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"catalog_github": CatalogPropertiesSchemaForDataSource(),

		"catalog_adogit": CatalogPropertiesSchemaForDataSource(),
	}
}

func (DevCenterCatalogDataSource) ModelObject() interface{} {
	return &DevCenterCatalogDataSourceModel{}
}

func (DevCenterCatalogDataSource) ResourceType() string {
	return "azurerm_dev_center_catalog"
}

func (r DevCenterCatalogDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Catalogs
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DevCenterCatalogDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterId, err := catalogs.ParseDevCenterID(state.DevCenterID)
			if err != nil {
				return err
			}

			id := catalogs.NewDevCenterCatalogID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			state.Name = id.CatalogName
			state.DevCenterID = catalogs.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID()

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if gitHub := props.GitHub; gitHub != nil {
						state.CatalogGitHub = []CatalogPropertiesDataSourceModel{
							{
								URI:            pointer.From(gitHub.Uri),
								Branch:         pointer.From(gitHub.Branch),
								KeyVaultKeyUrl: pointer.From(gitHub.SecretIdentifier),
								Path:           pointer.From(gitHub.Path),
							},
						}
					}

					if adoGit := props.AdoGit; adoGit != nil {
						state.CatalogAdoGit = []CatalogPropertiesDataSourceModel{
							{
								URI:            pointer.From(adoGit.Uri),
								Branch:         pointer.From(adoGit.Branch),
								KeyVaultKeyUrl: pointer.From(adoGit.SecretIdentifier),
								Path:           pointer.From(adoGit.Path),
							},
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func CatalogPropertiesSchemaForDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"uri": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"branch": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"key_vault_key_url": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}
