// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/hybridrunbookworkergroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type HybridRunbookWorkerGroupModel struct {
	ResourceGroupName     string `tfschema:"resource_group_name"`
	AutomationAccountName string `tfschema:"automation_account_name"`
	Name                  string `tfschema:"name"`
	CredentialName        string `tfschema:"credential_name"`
}

type HybridRunbookWorkerGroupResource struct{}

var _ sdk.Resource = (*HybridRunbookWorkerGroupResource)(nil)

func (m HybridRunbookWorkerGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(), // end if common
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
		"credential_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (m HybridRunbookWorkerGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (m HybridRunbookWorkerGroupResource) ModelObject() interface{} {
	return &HybridRunbookWorkerGroupModel{}
}

func (m HybridRunbookWorkerGroupResource) ResourceType() string {
	return "azurerm_automation_hybrid_runbook_worker_group"
}

func (m HybridRunbookWorkerGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.HybridRunbookWorkerGroup

			var model HybridRunbookWorkerGroupModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := hybridrunbookworkergroup.NewHybridRunbookWorkerGroupID(subscriptionID, model.ResourceGroupName,
				model.AutomationAccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}
			req := hybridrunbookworkergroup.HybridRunbookWorkerGroupCreateOrUpdateParameters{
				Name: pointer.To(model.Name),
			}
			if model.CredentialName != "" {
				req.Properties = &hybridrunbookworkergroup.HybridRunbookWorkerGroupCreateOrUpdateProperties{}
				req.Properties.Credential = &hybridrunbookworkergroup.RunAsCredentialAssociationProperty{
					Name: utils.String(model.CredentialName),
				}
			}
			// return 201 cause err in autorest sdk
			_, err = client.Create(ctx, id, req)
			// Workaround swagger issue https://github.com/Azure/azure-rest-api-specs/issues/19741
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m HybridRunbookWorkerGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := hybridrunbookworkergroup.ParseHybridRunbookWorkerGroupID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			client := meta.Client.Automation.HybridRunbookWorkerGroup
			result, err := client.Get(ctx, *id)
			if err != nil {
				return err
			}
			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}
			var output HybridRunbookWorkerGroupModel

			output.Name = utils.NormalizeNilableString(result.Model.Name)
			output.AutomationAccountName = id.AutomationAccountName
			if model := result.Model; model != nil {
				if prop := model.Properties; prop != nil {
					if prop.Credential != nil {
						output.CredentialName = pointer.From(prop.Credential.Name)
					}
				}
			}
			output.ResourceGroupName = id.ResourceGroupName
			return meta.Encode(&output)
		},
	}
}

func (m HybridRunbookWorkerGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.HybridRunbookWorkerGroup
			id, err := hybridrunbookworkergroup.ParseHybridRunbookWorkerGroupID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model HybridRunbookWorkerGroupModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			var upd hybridrunbookworkergroup.HybridRunbookWorkerGroupCreateOrUpdateParameters
			if meta.ResourceData.HasChange("credential_name") {
				upd.Properties = &hybridrunbookworkergroup.HybridRunbookWorkerGroupCreateOrUpdateProperties{}
				upd.Properties.Credential = &hybridrunbookworkergroup.RunAsCredentialAssociationProperty{
					Name: utils.String(model.CredentialName),
				}
			}
			if _, err = client.Update(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m HybridRunbookWorkerGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := hybridrunbookworkergroup.ParseHybridRunbookWorkerGroupID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Automation.HybridRunbookWorkerGroup
			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m HybridRunbookWorkerGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return hybridrunbookworkergroup.ValidateHybridRunbookWorkerGroupID
}
