package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type Content struct {
	Uri           string `tfschema:"uri"`
	HashAlgorithm string `tfschema:"hash_algorithm"`
	HashValue     string `tfschema:"hash_value"`
	Version       string `tfschema:"version"`
}

type Python2PackageModel struct {
	ResourceGroupName     string                 `tfschema:"resource_group_name"`
	AutomationAccountName string                 `tfschema:"automation_account_name"`
	Name                  string                 `tfschema:"name"`
	Tags                  map[string]interface{} `tfschema:"tags"`
	Content               []Content              `tfschema:"content"`
	IsGlobal              bool                   `tfschema:"is_global"`
	SizeInBytes           int64                  `tfschema:"size_in_bytes"`
	ActivityCount         int32                  `tfschema:"activity_count"`
	Location              string                 `tfschema:"location"`
	ErrorCode             string                 `tfschema:"error_code"`
	ErrorMeesage          string                 `tfschema:"error_meesage"`
	IsComposite           bool                   `tfschema:"is_composite"`
}

type Python2PackageResource struct{}

var _ sdk.Resource = (*Python2PackageResource)(nil)

func (m Python2PackageResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),
		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		// tags can be update by swagger definition, but no change after patch call, maybe change it to force_new?
		// see issue: https://github.com/Azure/azure-rest-api-specs/issues/19772
		"tags": commonschema.Tags(),
		"content": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"hash_algorithm": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"hash_value": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func (m Python2PackageResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"is_global": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"size_in_bytes": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"activity_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"error_code": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"error_meesage": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"is_composite": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m Python2PackageResource) ModelObject() interface{} {
	return &Python2PackageModel{}
}

func (m Python2PackageResource) ResourceType() string {
	return "azurerm_automation_python2_package"
}

func (m Python2PackageResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.Python2PackageClient

			var model Python2PackageModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := parse.NewPython2PackageID(subscriptionID, model.ResourceGroupName, model.AutomationAccountName, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
			if !utils.ResponseWasNotFound(existing.Response) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			var param automation.PythonPackageCreateParameters
			param.PythonPackageCreateProperties = &automation.PythonPackageCreateProperties{}
			param.Tags = tags.Expand(meta.ResourceData.Get("tags").(map[string]interface{}))
			content := model.Content[0]
			param.ContentLink = &automation.ContentLink{
				URI: utils.String(content.Uri),
				ContentHash: &automation.ContentHash{
					Algorithm: utils.String(content.HashAlgorithm),
					Value:     utils.String(content.HashValue),
				},
			}
			if content.Version != "" {
				param.ContentLink.Version = utils.String(content.Version)
			}
			_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, param)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}
			// wait provisioningState to be successful
			if err := m.waitDone(ctx, client, &id); err != nil {
				return err
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m Python2PackageResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.Python2PackageID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			client := meta.Client.Automation.Python2PackageClient
			result, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
			if err != nil {
				return err
			}

			var output Python2PackageModel
			// the server did not return the ContentLink properties, so read from current state
			// but for import there is current state
			if err := meta.Decode(&output); err != nil {
				return fmt.Errorf("encoding current state: %+v", err)
			}
			output.ResourceGroupName = id.ResourceGroup
			output.Name = utils.NormalizeNilableString(result.Name)
			// the tags always response an old value even after update
			output.Tags = tags.Flatten(result.Tags)
			output.AutomationAccountName = id.AutomationAccountName
			output.SizeInBytes = utils.NormaliseNilableInt64(result.SizeInBytes)
			output.ActivityCount = utils.NormaliseNilableInt32(result.ActivityCount)
			output.Location = utils.NormalizeNilableString(result.Location)
			output.IsComposite = utils.NormaliseNilableBool(result.IsComposite)
			output.IsGlobal = utils.NormaliseNilableBool(result.IsGlobal)
			if content := result.ContentLink; content != nil {
				c := Content{
					Uri:     utils.NormalizeNilableString(content.URI),
					Version: utils.NormalizeNilableString(content.Version),
				}
				if hash := content.ContentHash; hash != nil {
					c.HashAlgorithm = utils.NormalizeNilableString(hash.Algorithm)
					c.HashValue = utils.NormalizeNilableString(hash.Value)
				}
				output.Content = append([]Content{}, c)
			}
			if e := result.Error; e != nil {
				output.ErrorCode = utils.NormalizeNilableString(e.Code)
				output.ErrorMeesage = utils.NormalizeNilableString(e.Message)
			}

			return meta.Encode(&output)
		},
	}
}

func (m Python2PackageResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.Python2PackageClient

			id, err := parse.Python2PackageID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model Python2PackageModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			var upd automation.PythonPackageUpdateParameters
			if meta.ResourceData.HasChange("tags") {
				upd.Tags = tags.Expand(meta.ResourceData.Get("tags").(map[string]interface{}))
			}
			if _, err = client.Update(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			// wait provisioningState to be successful
			if err := m.waitDone(ctx, client, id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (m Python2PackageResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.Python2PackageID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Automation.Python2PackageClient
			if _, err = client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m Python2PackageResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.Python2PackageID
}

func (m Python2PackageResource) waitDone(ctx context.Context, client *automation.Python2PackageClient, id *parse.Python2PackageId) (err error) {
	targetState := automation.ModuleProvisioningStateSucceeded
	var pending []string
	for _, s := range automation.PossibleModuleProvisioningStateValues() {
		if s != targetState {
			pending = append(pending, string(s))
		}
	}
	conf := pluginsdk.StateChangeConf{
		Pending:    pending,
		Target:     []string{string(targetState)},
		MinTimeout: 20 * time.Second,
		Timeout:    10 * time.Minute,
		Refresh: func() (result interface{}, state string, err error) {
			resp, err2 := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
			if err2 != nil {
				return resp, "Error", fmt.Errorf("retrieving %s: %+v", id, err2)
			}

			if properties := resp.ModuleProperties; properties != nil {
				if properties.Error != nil && properties.Error.Message != nil && *properties.Error.Message != "" {
					return resp, string(properties.ProvisioningState), fmt.Errorf(*properties.Error.Message)
				}
				return resp, string(properties.ProvisioningState), nil
			}

			return resp, "Unknown", nil
		},
	}
	if _, err := conf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish provisioning: %+v", id, err)
	}
	return nil
}
