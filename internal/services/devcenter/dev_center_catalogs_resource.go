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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DevCenterCatalogsResourceModel struct {
	Name              string                   `tfschema:"name"`
	ResourceGroupName string                   `tfschema:"resource_group_name"`
	DevCenterID       string                   `tfschema:"dev_center_id"`
	CatalogGitHub     []CatalogPropertiesModel `tfschema:"catalog_github"`
	CatalogAdoGit     []CatalogPropertiesModel `tfschema:"catalog_adogit"`
}

type CatalogPropertiesModel struct {
	URI            string `tfschema:"uri"`
	Branch         string `tfschema:"branch"`
	KeyVaultKeyUrl string `tfschema:"key_vault_key_url"`
	Path           string `tfschema:"path"`
}

type DevCenterCatalogsResource struct{}

var _ sdk.Resource = DevCenterCatalogsResource{}

func (r DevCenterCatalogsResource) ModelObject() interface{} {
	return &DevCenterCatalogsResourceModel{}
}

func (r DevCenterCatalogsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return catalogs.ValidateDevCenterCatalogID
}

func (r DevCenterCatalogsResource) ResourceType() string {
	return "azurerm_dev_center_catalog"
}

func (r DevCenterCatalogsResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Required:     true,
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"dev_center_id": {
			Required:     true,
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			ValidateFunc: catalogs.ValidateDevCenterID,
		},

		"catalog_github": CatalogPropertiesSchema(),

		"catalog_adogit": CatalogPropertiesSchema(),
	}
}

func (r DevCenterCatalogsResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r DevCenterCatalogsResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DevCenterCatalogsResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.DevCenter.V20250201.Catalogs
			subscriptionId := metadata.Client.Account.SubscriptionId
			devCenterId, err := catalogs.ParseDevCenterID(model.DevCenterID)
			if err != nil {
				return fmt.Errorf("parsing dev center id: %+v", err)
			}
			devCenterName := devCenterId.DevCenterName
			id := catalogs.NewDevCenterCatalogID(subscriptionId, model.ResourceGroupName, devCenterName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			catalogProperties := catalogs.Catalog{
				Properties: &catalogs.CatalogProperties{
					AdoGit: expandCatalogProperties(model.CatalogAdoGit),
					GitHub: expandCatalogProperties(model.CatalogGitHub),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, catalogProperties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterCatalogsResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DevCenterCatalogsResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.DevCenter.V20250201.Catalogs
			id, err := catalogs.ParseDevCenterCatalogID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing catalog id: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if metadata.ResourceData.HasChange("catalog_github") {
				properties.Properties.GitHub = expandCatalogProperties(model.CatalogGitHub)
			}

			if metadata.ResourceData.HasChange("catalog_adogit") {
				properties.Properties.AdoGit = expandCatalogProperties(model.CatalogAdoGit)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DevCenterCatalogsResource) Read() sdk.ResourceFunc {
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

func (r DevCenterCatalogsResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Catalogs

			id, err := catalogs.ParseDevCenterCatalogID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing catalog id: %+v", err)
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
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

func expandCatalogProperties(input []CatalogPropertiesModel) *catalogs.GitCatalog {
	if len(input) == 0 {
		return nil
	}

	return &catalogs.GitCatalog{
		Uri:              pointer.To(input[0].URI),
		Branch:           pointer.To(input[0].Branch),
		SecretIdentifier: pointer.To(input[0].KeyVaultKeyUrl),
		Path:             pointer.To(input[0].Path),
	}
}
