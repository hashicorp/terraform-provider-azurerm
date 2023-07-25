// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
			ValidateFunc: watcher.ValidateAutomationAccountID,
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
			accountID, _ := watcher.ParseAutomationAccountID(model.AutomationAccountID)
			id := watcher.NewWatcherID(subscriptionID, accountID.ResourceGroupName, accountID.AutomationAccountName, model.Name)

			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			tags := expandStringInterfaceMap(model.Tags)
			scriptParameters := expandStringInterfaceMap(model.ScriptParameters)

			param := watcher.Watcher{
				Properties: &watcher.WatcherProperties{
					Description:                 utils.String(model.Description),
					ExecutionFrequencyInSeconds: utils.Int64(model.ExecutionFrequencyInSeconds),
					ScriptName:                  utils.String(model.ScriptName),
					ScriptParameters:            &scriptParameters,
					ScriptRunOn:                 utils.String(model.ScriptRunOn),
				},
				Etag:     utils.String(model.Etag),
				Location: utils.String(model.Location),
				Tags:     &tags,
			}

			_, err = client.CreateOrUpdate(ctx, id, param)
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
			id, err := watcher.ParseWatcherID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Automation.WatcherClient
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			var output WatcherModel
			if err := meta.Decode(&output); err != nil {
				return err
			}

			output.Name = id.WatcherName

			if model := resp.Model; model != nil {
				if props := resp.Model.Properties; props != nil {
					output.AutomationAccountID = watcher.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName).ID()
					output.ExecutionFrequencyInSeconds = utils.NormaliseNilableInt64(props.ExecutionFrequencyInSeconds)
					output.ScriptName = utils.NormalizeNilableString(props.ScriptName)
					output.ScriptRunOn = utils.NormalizeNilableString(props.ScriptRunOn)
					output.Description = utils.NormalizeNilableString(props.Description)
					output.Status = utils.NormalizeNilableString(props.Status)

					if props.ScriptParameters != nil {
						output.ScriptParameters = flattenMap(*props.ScriptParameters)
					}
				}
			}

			// tags, etag and location are not returned by response, so do NOT encode them

			return meta.Encode(&output)
		},
	}
}

func (m WatcherResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.WatcherClient

			id, err := watcher.ParseWatcherID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WatcherModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			var upd watcher.WatcherUpdateParameters
			upd.Properties = &watcher.WatcherUpdateProperties{}
			if meta.ResourceData.HasChange("execution_frequency_in_seconds") {
				upd.Properties.ExecutionFrequencyInSeconds = utils.Int64(model.ExecutionFrequencyInSeconds)
			}
			if _, err = client.Update(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", *id, err)
			}

			return nil
		},
	}
}

func (m WatcherResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := watcher.ParseWatcherID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Automation.WatcherClient
			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", *id, err)
			}
			return nil
		},
	}
}

func (m WatcherResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return watcher.ValidateWatcherID
}
