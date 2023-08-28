// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/hybridrunbookworker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type HybridRunbookWorkerModel struct {
	ResourceGroupName     string `tfschema:"resource_group_name"`
	AutomationAccountName string `tfschema:"automation_account_name"`
	WorkerGroupName       string `tfschema:"worker_group_name"`
	WorkerName            string `tfschema:"worker_name"`
	WorkerId              string `tfschema:"worker_id"`
	VmResourceId          string `tfschema:"vm_resource_id"`
	Ip                    string `tfschema:"ip"`
	RegisteredDateTime    string `tfschema:"registration_date_time"`
	LastSeenDateTime      string `tfschema:"last_seen_date_time"`
	WorkerType            string `tfschema:"worker_type"`
}

type HybridRunbookWorkerResource struct{}

var _ sdk.Resource = (*HybridRunbookWorkerResource)(nil)

func (m HybridRunbookWorkerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"worker_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"worker_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"vm_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (m HybridRunbookWorkerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"ip": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"registration_date_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"last_seen_date_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"worker_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"worker_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m HybridRunbookWorkerResource) ModelObject() interface{} {
	return &HybridRunbookWorkerModel{}
}

func (m HybridRunbookWorkerResource) ResourceType() string {
	return "azurerm_automation_hybrid_runbook_worker"
}

func (m HybridRunbookWorkerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.HybridRunbookWorker

			var model HybridRunbookWorkerModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := hybridrunbookworker.NewHybridRunbookWorkerID(subscriptionID, model.ResourceGroupName,
				model.AutomationAccountName, model.WorkerGroupName, model.WorkerId)
			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			req := hybridrunbookworker.HybridRunbookWorkerCreateParameters{}
			if model.VmResourceId != "" {
				req.Properties.VMResourceId = utils.String(model.VmResourceId)
			}

			future, err := client.Create(ctx, id, req)
			if err != nil {
				// Workaround swagger issue https://github.com/Azure/azure-rest-api-specs/issues/19741
				if !response.WasStatusCode(future.HttpResponse, http.StatusCreated) {
					return fmt.Errorf("creating %s: %v", id, err)
				}
			}
			_ = future

			meta.SetID(id)
			return nil
		},
	}
}

func (m HybridRunbookWorkerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := hybridrunbookworker.ParseHybridRunbookWorkerID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			client := meta.Client.Automation.HybridRunbookWorker
			result, err := client.Get(ctx, *id)
			if err != nil {
				return err
			}
			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}

			var output HybridRunbookWorkerModel

			// the name in response corresponding to work_id in request
			output.WorkerId = utils.NormalizeNilableString(result.Model.Name)
			output.AutomationAccountName = id.AutomationAccountName
			output.ResourceGroupName = id.ResourceGroupName
			output.WorkerGroupName = id.HybridRunbookWorkerGroupName
			if prop := result.Model.Properties; prop != nil {
				output.VmResourceId = utils.NormalizeNilableString(prop.VMResourceId)
				output.WorkerType = utils.NormalizeNilableString((*string)(prop.WorkerType))
				output.LastSeenDateTime = utils.NormalizeNilableString(prop.LastSeenDateTime)
				output.RegisteredDateTime = utils.NormalizeNilableString(prop.RegisteredDateTime)
				output.Ip = utils.NormalizeNilableString(prop.IP)
				output.WorkerName = utils.NormalizeNilableString(prop.WorkerName)
			}
			return meta.Encode(&output)
		},
	}
}

func (m HybridRunbookWorkerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := hybridrunbookworker.ParseHybridRunbookWorkerID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Automation.HybridRunbookWorker
			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m HybridRunbookWorkerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return hybridrunbookworker.ValidateHybridRunbookWorkerID
}
