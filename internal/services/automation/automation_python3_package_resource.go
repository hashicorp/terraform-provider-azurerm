package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/python3package"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type Python3PackageModel struct {
	Name                  string            `tfschema:"name"`
	ResourceGroupName     string            `tfschema:"resource_group_name"`
	AutomationAccountName string            `tfschema:"automation_account_name"`
	ContentUri            string            `tfschema:"content_uri"`
	ContentVersion        string            `tfschema:"content_version"`
	HashAlgorithm         string            `tfschema:"hash_algorithm"`
	HashValue             string            `tfschema:"hash_value"`
	Tags                  map[string]string `tfschema:"tags"`
}

type Python3PackageResource struct{}

var _ sdk.Resource = (*Python3PackageResource)(nil)

func (m Python3PackageResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"content_uri": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"content_version": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"hash_algorithm": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"hash_value"},
		},

		"hash_value": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"hash_algorithm"},
		},

		"tags": commonschema.Tags(),
	}
}

func (m Python3PackageResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (m Python3PackageResource) ModelObject() interface{} {
	return &Python3PackageModel{}
}

func (m Python3PackageResource) ResourceType() string {
	return "azurerm_automation_python3_package"
}

func (m Python3PackageResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.Python3Package

			var model Python3PackageModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := python3package.NewPython3PackageID(subscriptionID, model.ResourceGroupName, model.AutomationAccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			req := python3package.PythonPackageCreateParameters{}
			req.Properties = python3package.PythonPackageCreateProperties{
				ContentLink: python3package.ContentLink{
					Uri: &model.ContentUri,
				},
			}

			if model.ContentVersion != "" {
				req.Properties.ContentLink.Version = &model.ContentVersion
			}

			if model.HashAlgorithm != "" {
				req.Properties.ContentLink.ContentHash = &python3package.ContentHash{
					Algorithm: model.HashAlgorithm,
					Value:     model.HashValue,
				}
			}
			req.Tags = &model.Tags

			if err = client.CreateOrUpdateThenPoll(ctx, id, req); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m Python3PackageResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := python3package.ParsePython3PackageID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Automation.Python3Package
			result, err := client.Get(ctx, *id)
			if err != nil {
				return err
			}

			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}

			var stateModel Python3PackageModel
			if err = meta.Decode(&stateModel); err != nil {
				return err
			}

			output := Python3PackageModel{
				ResourceGroupName:     id.ResourceGroupName,
				AutomationAccountName: id.AutomationAccountName,
				Name:                  id.Python3PackageName,

				// the fields below don't return by the API, remove it when issue fixed
				// https://github.com/Azure/azure-rest-api-specs/issues/25538
				ContentVersion: stateModel.ContentVersion,
				ContentUri:     stateModel.ContentUri,
				HashValue:      stateModel.HashValue,
				HashAlgorithm:  stateModel.HashAlgorithm,
			}

			model := result.Model
			if model.Properties != nil {
				if content := model.Properties.ContentLink; content != nil {
					output.ContentUri = pointer.From(content.Uri)
					output.ContentVersion = pointer.From(content.Version)
					if hash := content.ContentHash; hash != nil {
						output.HashAlgorithm = hash.Algorithm
						output.HashValue = hash.Value
					}
				}
				output.Tags = pointer.From(model.Tags)
			}

			return meta.Encode(&output)
		},
	}
}

func (m Python3PackageResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.Python3Package
			id, err := python3package.ParsePython3PackageID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model Python3PackageModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding %s: %+v", id, err)
			}

			var upd python3package.PythonPackageUpdateParameters
			if meta.ResourceData.HasChange("tags") {
				upd.Tags = &model.Tags
			}

			if _, err = client.Update(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m Python3PackageResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := python3package.ParsePython3PackageID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Automation.Python3Package
			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m Python3PackageResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return python3package.ValidatePython3PackageID
}
