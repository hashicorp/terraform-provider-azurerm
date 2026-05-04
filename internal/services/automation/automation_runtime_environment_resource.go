// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runtimeenvironment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AutomationRuntimeEnvironmentResourceModel struct {
	Name                   string            `tfschema:"name"`
	AutomationAccountId    string            `tfschema:"automation_account_id"`
	RuntimeLanguage        string            `tfschema:"runtime_language"`
	RuntimeVersion         string            `tfschema:"runtime_version"`
	RuntimeDefaultPackages map[string]string `tfschema:"runtime_default_packages"`
	Location               string            `tfschema:"location"`
	Description            string            `tfschema:"description"`
	Tags                   map[string]string `tfschema:"tags"`
}

type AutomationRuntimeEnvironmentResource struct{}

var (
	_ sdk.ResourceWithUpdate        = AutomationRuntimeEnvironmentResource{}
	_ sdk.ResourceWithCustomizeDiff = AutomationRuntimeEnvironmentResource{}
)

func (m AutomationRuntimeEnvironmentResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if rdp, ok := rd.GetOk("runtime_default_packages"); ok && rdp != nil {
				if rl, ok := rd.GetOk("runtime_language"); ok {
					if rl.(string) == "Python" {
						return errors.New("runtime_default_packages cannot be set for runtime_language `Python`")
					}
				}
			}

			if rd.HasChange("runtime_default_packages") {
				old, new := rd.GetChange("runtime_default_packages")
				oldMap, newMap := old.(map[string]interface{}), new.(map[string]interface{})

				// Azure API limitation: Runtime environment packages cannot be removed once added
				if len(oldMap) > len(newMap) {
					if err := rd.ForceNew("runtime_default_packages"); err != nil {
						return err
					}
				}

				// Azure API limitation: Package names cannot change, only versions can be updated
				for k := range oldMap {
					if _, ok := newMap[k]; !ok {
						if err := rd.ForceNew("runtime_default_packages"); err != nil {
							return err
						}
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

		"location": commonschema.Location(),

		"automation_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: runtimeenvironment.ValidateAutomationAccountID,
		},

		"runtime_language": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Python",
				"PowerShell",
			}, false),
		},

		"runtime_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},

		"runtime_default_packages": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"tags": {
			Type:         pluginsdk.TypeMap,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: tags.Validate,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
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
			client := meta.Client.Automation.RuntimeEnvironment

			var model AutomationRuntimeEnvironmentResourceModel
			if err := meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionID := meta.Client.Account.SubscriptionId

			autAccId, err := runtimeenvironment.ParseAutomationAccountIDInsensitively(model.AutomationAccountId)
			if err != nil {
				return err
			}

			id := runtimeenvironment.NewRuntimeEnvironmentID(subscriptionID, autAccId.ResourceGroupName, autAccId.AutomationAccountName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			req := runtimeenvironment.RuntimeEnvironment{
				Location: location.Normalize(model.Location),
				Name:     pointer.To(model.Name),
				Properties: &runtimeenvironment.RuntimeEnvironmentProperties{
					Runtime: &runtimeenvironment.RuntimeProperties{
						Language: pointer.To(model.RuntimeLanguage),
						Version:  pointer.To(model.RuntimeVersion),
					},
				},
			}

			if model.Tags != nil {
				req.Tags = &model.Tags
			}

			if model.RuntimeDefaultPackages != nil {
				req.Properties.DefaultPackages = pointer.To(model.RuntimeDefaultPackages)
			}

			if model.Description != "" {
				req.Properties.Description = pointer.To(model.Description)
			}

			if _, err = client.Create(ctx, id, req); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
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
			client := meta.Client.Automation.RuntimeEnvironment
			id, err := runtimeenvironment.ParseRuntimeEnvironmentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return meta.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AutomationRuntimeEnvironmentResourceModel{
				Name:                id.RuntimeEnvironmentName,
				AutomationAccountId: runtimeenvironment.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = model.Location

				if model.Properties != nil {
					if model.Properties.DefaultPackages != nil {
						state.RuntimeDefaultPackages = *model.Properties.DefaultPackages
					}

					state.Description = pointer.From(model.Properties.Description)

					if model.Properties.Runtime != nil {
						state.RuntimeLanguage = pointer.From(model.Properties.Runtime.Language)
						state.RuntimeVersion = pointer.From(model.Properties.Runtime.Version)
					}
				}

				if model.Tags != nil {
					state.Tags = *model.Tags
				}
			}

			return meta.Encode(&state)
		},
	}
}

func (m AutomationRuntimeEnvironmentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.RuntimeEnvironment
			id, err := runtimeenvironment.ParseRuntimeEnvironmentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var state AutomationRuntimeEnvironmentResourceModel
			if err = meta.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			param := runtimeenvironment.RuntimeEnvironmentUpdateParameters{
				Properties: &runtimeenvironment.RuntimeEnvironmentUpdateProperties{
					DefaultPackages: pointer.To(state.RuntimeDefaultPackages),
				},
			}

			if _, err = client.Update(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (m AutomationRuntimeEnvironmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.RuntimeEnvironment

			id, err := runtimeenvironment.ParseRuntimeEnvironmentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (m AutomationRuntimeEnvironmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return runtimeenvironment.ValidateRuntimeEnvironmentID
}
