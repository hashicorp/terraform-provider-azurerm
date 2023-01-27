package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WatcherModel struct {
	AutomationAccountID         string                 `tfschema:"automation_account_id"`
	Name                        string                 `tfschema:"name"`
	Location                    string                 `tfschema:"location"`
	Tags                        map[string]interface{} `tfschema:"tags"`
	Etag                        string                 `tfschema:"etag"`
	ExecutionFrequencyInSeconds int64                  `tfschema:"execution_frequency_in_seconds"`
	ScriptName                  string                 `tfschema:"script_name"`
	ScriptParameters            map[string]interface{} `tfschema:"script_parameters"`
	ScriptRunOn                 string                 `tfschema:"script_run_on"`
	Description                 string                 `tfschema:"description"`
	Status                      string                 `tfschema:"status"`
}

type WatcherResource struct{}

var _ sdk.Resource = (*WatcherResource)(nil)

func (m WatcherResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"automation_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccountID,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),

		"etag": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"execution_frequency_in_seconds": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"script_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"script_parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"script_run_on": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (m WatcherResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m WatcherResource) ModelObject() interface{} {
	return &WatcherModel{}
}

func (m WatcherResource) ResourceType() string {
	return "azurerm_automation_watcher"
}

func (m WatcherResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.WatcherClient

			var model WatcherModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			accountID, _ := parse.AutomationAccountID(model.AutomationAccountID)
			id := parse.NewWatcherID(subscriptionID, accountID.ResourceGroup, accountID.Name, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
			if !utils.ResponseWasNotFound(existing.Response) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			var param automation.Watcher
			param.Tags = tags.Expand(model.Tags)
			param.Etag = utils.String(model.Etag)
			param.Location = utils.String(model.Location)
			param.WatcherProperties = &automation.WatcherProperties{}
			prop := param.WatcherProperties
			prop.ExecutionFrequencyInSeconds = utils.Int64(model.ExecutionFrequencyInSeconds)
			prop.ScriptName = utils.String(model.ScriptName)
			prop.ScriptParameters = tags.Expand(model.ScriptParameters)
			prop.ScriptRunOn = utils.String(model.ScriptRunOn)
			prop.Description = utils.String(model.Description)

			_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, param)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m WatcherResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.WatcherID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			client := meta.Client.Automation.WatcherClient
			result, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
			if err != nil {
				return err
			}

			var output WatcherModel
			if err := meta.Decode(&output); err != nil {
				return err
			}

			prop := result.WatcherProperties
			output.AutomationAccountID = parse.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroup, id.AutomationAccountName).ID()
			output.Name = id.Name
			output.ExecutionFrequencyInSeconds = utils.NormaliseNilableInt64(prop.ExecutionFrequencyInSeconds)
			output.ScriptName = utils.NormalizeNilableString(prop.ScriptName)
			output.ScriptParameters = tags.Flatten(prop.ScriptParameters)
			output.ScriptRunOn = utils.NormalizeNilableString(prop.ScriptRunOn)
			output.Description = utils.NormalizeNilableString(prop.Description)
			output.Status = utils.NormalizeNilableString(prop.Status)

			// tags, etag and location do not returned by response, so do NOT encode them

			return meta.Encode(&output)
		},
	}
}

func (m WatcherResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.WatcherClient

			id, err := parse.WatcherID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WatcherModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			var upd automation.WatcherUpdateParameters
			upd.WatcherUpdateProperties = &automation.WatcherUpdateProperties{}
			if meta.ResourceData.HasChange("execution_frequency_in_seconds") {
				upd.WatcherUpdateProperties.ExecutionFrequencyInSeconds = utils.Int64(model.ExecutionFrequencyInSeconds)
			}
			if _, err = client.Update(ctx, id.ResourceGroup, id.Name, id.AutomationAccountName, upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m WatcherResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.WatcherID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Automation.WatcherClient
			if _, err = client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m WatcherResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.WatcherID
}
