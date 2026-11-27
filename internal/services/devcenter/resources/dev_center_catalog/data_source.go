// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dev_center_catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/catalogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	devcenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

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

func (DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Dev Center Catalog.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dev_center_id": commonschema.ResourceIDReferenceRequired(&devcenters.DevCenterId{}, "Dev Center"),
	}
}

func (DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"catalog_github": CatalogPropertiesSchemaForDataSource(),

		"catalog_adogit": CatalogPropertiesSchemaForDataSource(),
	}
}

func (DataSource) ModelObject() any {
	return &DevCenterCatalogDataSourceModel{}
}

func (DataSource) ResourceType() string {
	return "azurerm_dev_center_catalog"
}

func (DataSource) Description() string {
	return "Use this data source to access information about an existing Dev Center Catalog."
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := devcenterClient.From(metadata.Client).V20250201.Catalogs
			subscriptionId := account.From(metadata.Client).SubscriptionId

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
					Description: "The Git URI of the Dev Center Catalog.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"branch": {
					Description: "The Git branch of the Dev Center Catalog.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"key_vault_key_url": {
					Description: "A reference to the Key Vault Secret containing a security token to authenticate to a Git repository.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"path": {
					Description: "The folder where the catalog items can be found inside the repository.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},
			},
		},
	}
}
