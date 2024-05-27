// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/module"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = PowerShell72ModuleResource{}

type ModuleLinkModel struct {
	Uri  string       `tfschema:"uri"`
	Hash []ModuleHash `tfschema:"hash"`
}

type ModuleHash struct {
	Algorithm string `tfschema:"algorithm"`
	Value     string `tfschema:"value"`
}

type AutomationPowerShell72ModuleModel struct {
	AutomationAccountID string                 `tfschema:"automation_account_id"`
	Name                string                 `tfschema:"name"`
	ModuleLink          []ModuleLinkModel      `tfschema:"module_link"`
	Tags                map[string]interface{} `tfschema:"tags"`
}

type PowerShell72ModuleResource struct{}

func (r PowerShell72ModuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"automation_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: module.ValidateAutomationAccountID,
		},

		"module_link": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"hash": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"algorithm": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"value": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		"tags": commonschema.Tags(),
	}
}

func (r PowerShell72ModuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PowerShell72ModuleResource) ModelObject() interface{} {
	return &AutomationPowerShell72ModuleModel{}
}

func (r PowerShell72ModuleResource) ResourceType() string {
	return "azurerm_automation_powershell72_module"
}

func (r PowerShell72ModuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return module.ValidatePowerShell72ModuleID
}

func (r PowerShell72ModuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.Module

			var model AutomationPowerShell72ModuleModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			accountID, _ := module.ParseAutomationAccountID(model.AutomationAccountID)
			name := metadata.ResourceData.Get("name").(string)

			id := module.NewPowerShell72ModuleID(subscriptionId, accountID.ResourceGroupName, accountID.AutomationAccountName, name)

			existing, err := client.PowerShell72ModuleGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			// for existing global module do update instead of raising ImportAsExistsError
			isGlobal := existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.IsGlobal != nil && *existing.Model.Properties.IsGlobal
			if !response.WasNotFound(existing.HttpResponse) && !isGlobal {
				return tf.ImportAsExistsError("azurerm_automation_powershell72_module", id.ID())
			}

			parameters := module.ModuleCreateOrUpdateParameters{
				Properties: module.ModuleCreateOrUpdateProperties{
					ContentLink: expandPowerShell72ModuleLink(model.ModuleLink),
				},
				Tags: tags.Expand(model.Tags),
			}

			if _, err := client.PowerShell72ModuleCreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			// the API returns 'done' but it's not actually finished provisioning yet
			// tracking issue: https://github.com/Azure/azure-rest-api-specs/pull/25435
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{
					string(module.ModuleProvisioningStateActivitiesStored),
					string(module.ModuleProvisioningStateConnectionTypeImported),
					string(module.ModuleProvisioningStateContentDownloaded),
					string(module.ModuleProvisioningStateContentRetrieved),
					string(module.ModuleProvisioningStateContentStored),
					string(module.ModuleProvisioningStateContentValidated),
					string(module.ModuleProvisioningStateCreated),
					string(module.ModuleProvisioningStateCreating),
					string(module.ModuleProvisioningStateModuleDataStored),
					string(module.ModuleProvisioningStateModuleImportRunbookComplete),
					string(module.ModuleProvisioningStateRunningImportModuleRunbook),
					string(module.ModuleProvisioningStateStartingImportModuleRunbook),
					string(module.ModuleProvisioningStateUpdating),
				},
				Target: []string{
					string(module.ModuleProvisioningStateSucceeded),
				},
				MinTimeout: 30 * time.Second,
				Refresh: func() (interface{}, string, error) {
					resp, err2 := client.PowerShell72ModuleGet(ctx, id)
					if err2 != nil {
						return resp, "Error", fmt.Errorf("retrieving %s: %+v", id, err2)
					}

					provisioningState := "Unknown"
					if model := resp.Model; model != nil {
						if props := model.Properties; props != nil {
							if props.ProvisioningState != nil {
								provisioningState = string(*props.ProvisioningState)
							}
							if props.Error != nil && props.Error.Message != nil && *props.Error.Message != "" {
								return resp, provisioningState, fmt.Errorf(*props.Error.Message)
							}
							return resp, provisioningState, nil
						}
					}
					return resp, provisioningState, nil
				},
				Timeout: time.Until(deadline),
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to finish provisioning: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r PowerShell72ModuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.Module

			id, err := module.ParsePowerShell72ModuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AutomationPowerShell72ModuleModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			parameters := module.ModuleCreateOrUpdateParameters{
				Properties: module.ModuleCreateOrUpdateProperties{
					ContentLink: expandPowerShell72ModuleLink(model.ModuleLink),
				},
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = tags.Expand(model.Tags)
			}

			if _, err := client.PowerShell72ModuleCreateOrUpdate(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			// the API returns 'done' but it's not actually finished provisioning yet
			// tracking issue: https://github.com/Azure/azure-rest-api-specs/pull/25435
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{
					string(module.ModuleProvisioningStateActivitiesStored),
					string(module.ModuleProvisioningStateConnectionTypeImported),
					string(module.ModuleProvisioningStateContentDownloaded),
					string(module.ModuleProvisioningStateContentRetrieved),
					string(module.ModuleProvisioningStateContentStored),
					string(module.ModuleProvisioningStateContentValidated),
					string(module.ModuleProvisioningStateCreated),
					string(module.ModuleProvisioningStateCreating),
					string(module.ModuleProvisioningStateModuleDataStored),
					string(module.ModuleProvisioningStateModuleImportRunbookComplete),
					string(module.ModuleProvisioningStateRunningImportModuleRunbook),
					string(module.ModuleProvisioningStateStartingImportModuleRunbook),
					string(module.ModuleProvisioningStateUpdating),
				},
				Target: []string{
					string(module.ModuleProvisioningStateSucceeded),
				},
				MinTimeout: 30 * time.Second,
				Refresh: func() (interface{}, string, error) {
					resp, err2 := client.PowerShell72ModuleGet(ctx, *id)
					if err2 != nil {
						return resp, "Error", fmt.Errorf("retrieving %s: %+v", id, err2)
					}

					provisioningState := "Unknown"
					if model := resp.Model; model != nil {
						if props := model.Properties; props != nil {
							if props.ProvisioningState != nil {
								provisioningState = string(*props.ProvisioningState)
							}
							if props.Error != nil && props.Error.Message != nil && *props.Error.Message != "" {
								return resp, provisioningState, fmt.Errorf(*props.Error.Message)
							}
							return resp, provisioningState, nil
						}
					}
					return resp, provisioningState, nil
				},
				Timeout: time.Until(deadline),
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r PowerShell72ModuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.Module
			id, err := module.ParsePowerShell72ModuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.PowerShell72ModuleGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var output AutomationPowerShell72ModuleModel
			if err := metadata.Decode(&output); err != nil {
				return err
			}

			output.Name = id.PowerShell72ModuleName
			output.AutomationAccountID = module.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName).ID()
			if resp.Model != nil {
				output.Tags = tags.Flatten(resp.Model.Tags)
			}

			return metadata.Encode(&output)
		},
	}
}

func (PowerShell72ModuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		// the Func returns a function which deletes the Resource Group
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.Module
			id, err := module.ParsePowerShell72ModuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.PowerShell72ModuleDelete(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}

				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandPowerShell72ModuleLink(m []ModuleLinkModel) module.ContentLink {
	moduleLink := m[0]
	if len(moduleLink.Hash) > 0 {
		hash := moduleLink.Hash[0]
		return module.ContentLink{
			Uri: &moduleLink.Uri,
			ContentHash: &module.ContentHash{
				Algorithm: hash.Algorithm,
				Value:     hash.Value,
			},
		}
	}

	return module.ContentLink{
		Uri: &moduleLink.Uri,
	}
}
