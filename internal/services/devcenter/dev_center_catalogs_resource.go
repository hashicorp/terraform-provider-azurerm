package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/catalogs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DevCenterCatalogsResourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	DevCenterID       string `tfschema:"dev_center_id"`
	CatalogType       string `tfschema:"catalog_type"`
	URI               string `tfschema:"uri"`
	Branch            string `tfschema:"branch"`
	KeyVaultKeyUrl    string `tfschema:"key_vault_key_url"`
	Path              string `tfschema:"path"`
}

type DevCenterCatalogsResource struct{}

var _ sdk.Resource = DevCenterCatalogsResource{}

func (r DevCenterCatalogsResource) ModelObject() interface{} {
	return &DevCenterCatalogsResourceModel{}
}

func (r DevCenterCatalogsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return catalogs.ValidateCatalogID
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

		"catalog_type": {
			Required: true,
			Type:     pluginsdk.TypeString,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"gitHub",
				"adoGit",
			}, false),
		},

		"uri": {
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"branch": {
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"key_vault_key_url": {
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"path": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},
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

			client := metadata.Client.DevCenter.V20230401.Catalogs
			subscriptionId := metadata.Client.Account.SubscriptionId
			devCenterId, err := catalogs.ParseDevCenterID(model.DevCenterID)
			if err != nil {
				return fmt.Errorf("parsing dev center id: %+v", err)
			}
			devCenterName := devCenterId.DevCenterName
			id := catalogs.NewCatalogID(subscriptionId, model.ResourceGroupName, devCenterName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var properties catalogs.Catalog
			if model.CatalogType == "gitHub" {
				gitHubProps := catalogs.GitCatalog{
					Branch:           pointer.To(model.Branch),
					Uri:              pointer.To(model.URI),
					Path:             pointer.To(model.Path),
					SecretIdentifier: pointer.To(model.KeyVaultKeyUrl),
				}

				properties = catalogs.Catalog{
					Properties: &catalogs.CatalogProperties{
						GitHub: pointer.To(gitHubProps),
					},
				}
			} else {
				adoGitProps := catalogs.GitCatalog{
					Branch:           pointer.To(model.Branch),
					Uri:              pointer.To(model.URI),
					Path:             pointer.To(model.Path),
					SecretIdentifier: pointer.To(model.KeyVaultKeyUrl),
				}

				properties = catalogs.Catalog{
					Properties: &catalogs.CatalogProperties{
						AdoGit: pointer.To(adoGitProps),
					},
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
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

			client := metadata.Client.DevCenter.V20230401.Catalogs
			id, err := catalogs.ParseCatalogID(metadata.ResourceData.Id())
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

			if metadata.ResourceData.HasChange("uri") {
				if model.CatalogType == "gitHub" {
					properties.Properties.GitHub.Uri = pointer.To(model.URI)
				} else {
					properties.Properties.AdoGit.Uri = pointer.To(model.URI)
				}
			}

			if metadata.ResourceData.HasChange("branch") {
				if model.CatalogType == "gitHub" {
					properties.Properties.GitHub.Branch = pointer.To(model.Branch)
				} else {
					properties.Properties.AdoGit.Branch = pointer.To(model.Branch)
				}
			}

			if metadata.ResourceData.HasChange("path") {
				if model.CatalogType == "gitHub" {
					properties.Properties.GitHub.Path = pointer.To(model.Path)
				} else {
					properties.Properties.AdoGit.Path = pointer.To(model.Path)
				}
			}

			if metadata.ResourceData.HasChange("key_vault_key_url") {
				if model.CatalogType == "gitHub" {
					properties.Properties.GitHub.SecretIdentifier = pointer.To(model.KeyVaultKeyUrl)
				} else {
					properties.Properties.AdoGit.SecretIdentifier = pointer.To(model.KeyVaultKeyUrl)
				}
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
			client := metadata.Client.DevCenter.V20230401.Catalogs

			id, err := catalogs.ParseCatalogID(metadata.ResourceData.Id())
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
					state.CatalogType = "gitHub"
					state.URI = pointer.From(gitHub.Uri)
					state.Branch = pointer.From(gitHub.Branch)
					state.Path = pointer.From(gitHub.Path)
					state.KeyVaultKeyUrl = pointer.From(gitHub.SecretIdentifier)
				} else if adoGit := properties.AdoGit; adoGit != nil {
					state.CatalogType = "adoGit"
					state.URI = pointer.From(adoGit.Uri)
					state.Branch = pointer.From(adoGit.Branch)
					state.Path = pointer.From(adoGit.Path)
					state.KeyVaultKeyUrl = pointer.From(adoGit.SecretIdentifier)
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
			client := metadata.Client.DevCenter.V20230401.Catalogs

			id, err := catalogs.ParseCatalogID(metadata.ResourceData.Id())
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
