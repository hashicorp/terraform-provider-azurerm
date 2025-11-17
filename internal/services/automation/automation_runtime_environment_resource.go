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
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runtimeenvironment"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AutomationRuntimeEnvironmentResourceModel struct {
	Name                   string            `tfschema:"name"`
	ResourceGroupName      string            `tfschema:"resource_group_name"`
	AutomationAccountName  string            `tfschema:"automation_account_name"`
	RuntimeLanguage        string            `tfschema:"runtime_language"`
	RuntimeVersion         string            `tfschema:"runtime_version"`
	RuntimeDefaultPackages map[string]string `tfschema:"runtime_default_packages"`
	Location               string            `tfschema:"location"`
	Description            string            `tfschema:"description"`
	Tags                   map[string]string `tfschema:"tags"`
}

type AutomationRuntimeEnvironmentResource struct{}

var _ sdk.Resource = (*AutomationRuntimeEnvironmentResource)(nil)

func (m AutomationRuntimeEnvironmentResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if rd.HasChange("runtime_default_packages") {
				old, new := rd.GetChange("runtime_default_packages")
				oldMap, newMap := old.(map[string]interface{}), new.(map[string]interface{})

				// Packages can not be removed, only added
				if len(oldMap) > len(newMap) {
					rd.ForceNew("runtime_default_packages")
				}

				// Packages changes not the version can not be done in place
				for k := range oldMap {
					if _, ok := newMap[k]; !ok {
						rd.ForceNew("runtime_default_packages")
					}
				}
			}

			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (m AutomationRuntimeEnvironmentResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: validate.AutomationAccount(),
		},
		"runtime_language": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string("Python"),
				string("PowerShell"),
			}, false),
		},
		"runtime_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"runtime_default_packages": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"location": commonschema.Location(),
		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"tags": commonschema.Tags(),
	}
}

func (m AutomationRuntimeEnvironmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (m AutomationRuntimeEnvironmentResource) ModelObject() interface{} {
	return &AutomationRuntimeEnvironmentResourceModel{}
}

func (m AutomationRuntimeEnvironmentResource) ResourceType() string {
	return "azurerm_automation_runtime_environment"
}

func (m AutomationRuntimeEnvironmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.RuntimeEnvironmentClient

			var model AutomationRuntimeEnvironmentResourceModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := runtimeenvironment.NewRuntimeEnvironmentID(subscriptionID, model.ResourceGroupName, model.AutomationAccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			req := runtimeenvironment.RuntimeEnvironment{}

			req.Properties = &runtimeenvironment.RuntimeEnvironmentProperties{
				Runtime: &runtimeenvironment.RuntimeProperties{
					Language: pointer.FromString(model.RuntimeLanguage),
					Version:  pointer.FromString(model.RuntimeVersion),
				},
				DefaultPackages: &model.RuntimeDefaultPackages,
				Description:     pointer.FromString(model.Description),
			}

			req.Location = azure.NormalizeLocation(model.Location)
			req.Tags = &model.Tags

			if _, err = client.Create(ctx, id, req); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m AutomationRuntimeEnvironmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := runtimeenvironment.ParseRuntimeEnvironmentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Automation.RuntimeEnvironmentClient
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}

			var stateModel AutomationRuntimeEnvironmentResourceModel
			if err = meta.Decode(&stateModel); err != nil {
				return err
			}

			output := AutomationRuntimeEnvironmentResourceModel{
				Name:                   id.RuntimeEnvironmentName,
				ResourceGroupName:      stateModel.ResourceGroupName,
				AutomationAccountName:  id.AutomationAccountName,
				RuntimeLanguage:        *result.Model.Properties.Runtime.Language,
				RuntimeVersion:         *result.Model.Properties.Runtime.Version,
				RuntimeDefaultPackages: *result.Model.Properties.DefaultPackages,
				Description:            *result.Model.Properties.Description,
				Location:               result.Model.Location,
				Tags:                   *result.Model.Tags,
			}

			return meta.Encode(&output)
		},
	}
}

func (m AutomationRuntimeEnvironmentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.RuntimeEnvironmentClient
			id, err := runtimeenvironment.ParseRuntimeEnvironmentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var stateModel AutomationRuntimeEnvironmentResourceModel
			if err = meta.Decode(&stateModel); err != nil {
				return err
			}

			var upd runtimeenvironment.RuntimeEnvironmentUpdateParameters
			if meta.ResourceData.HasChange("runtime_default_packages") {
				upd.Properties = &runtimeenvironment.RuntimeEnvironmentUpdateProperties{
					DefaultPackages: &stateModel.RuntimeDefaultPackages,
				}
			}

			if _, err = client.Update(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m AutomationRuntimeEnvironmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := runtimeenvironment.ParseRuntimeEnvironmentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Automation.RuntimeEnvironmentClient
			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m AutomationRuntimeEnvironmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return runtimeenvironment.ValidateRuntimeEnvironmentID
}
