// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/module"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationModule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationModuleCreate,
		Read:   resourceAutomationModuleRead,
		Update: resourceAutomationModuleUpdate,
		Delete: resourceAutomationModuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := module.ParseModuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

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
		},
	}
}

func resourceAutomationModuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Module
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Module creation.")

	id := module.NewModuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	// for existing global module do update instead of raising ImportAsExistsError
	isGlobal := existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.IsGlobal != nil && *existing.Model.Properties.IsGlobal
	if !response.WasNotFound(existing.HttpResponse) && !isGlobal {
		return tf.ImportAsExistsError("azurerm_automation_module", id.ID())
	}

	parameters := module.ModuleCreateOrUpdateParameters{
		Properties: module.ModuleCreateOrUpdateProperties{
			ContentLink: expandModuleLink(d),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := waitForModuleProvisioningCompletion(ctx, client, id, d.Timeout(pluginsdk.TimeoutCreate)); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceAutomationModuleRead(d, meta)
}

func resourceAutomationModuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Module
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := module.ParseModuleID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}

	// Start from the existing content link so that fields managed via
	// ignore_changes retain their server-side value.
	var contentLink *module.ContentLink
	if existing.Model != nil && existing.Model.Properties != nil {
		contentLink = existing.Model.Properties.ContentLink
	}

	if d.HasChange("module_link") {
		expanded := expandModuleLink(d)
		contentLink = &expanded
	}

	parameters := module.ModuleUpdateParameters{
		Name: &id.ModuleName,
		Properties: &module.ModuleUpdateProperties{
			ContentLink: contentLink,
		},
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err := waitForModuleProvisioningCompletion(ctx, client, *id, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return err
	}

	return resourceAutomationModuleRead(d, meta)
}

func resourceAutomationModuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Module
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := module.ParseModuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.ModuleName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	return nil
}

func resourceAutomationModuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Module
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := module.ParseModuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandModuleLink(d *pluginsdk.ResourceData) module.ContentLink {
	inputs := d.Get("module_link").([]interface{})
	input := inputs[0].(map[string]interface{})
	uri := input["uri"].(string)

	hashes := input["hash"].([]interface{})

	if len(hashes) > 0 {
		hash := hashes[0].(map[string]interface{})
		return module.ContentLink{
			Uri: &uri,
			ContentHash: &module.ContentHash{
				Algorithm: hash["algorithm"].(string),
				Value:     hash["value"].(string),
			},
		}
	}

	return module.ContentLink{
		Uri: &uri,
	}
}

// waitForModuleProvisioningCompletion polls the module until it reaches the Succeeded
// provisioning state. The API returns 'done' before provisioning is actually complete.
// Tracking issue: https://github.com/Azure/azure-rest-api-specs/pull/25435
func waitForModuleProvisioningCompletion(ctx context.Context, client *module.ModuleClient, id module.ModuleId, timeout time.Duration) error {
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
		Timeout:    timeout,
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, id)
			if err != nil {
				return resp, "Error", fmt.Errorf("retrieving %s: %+v", id, err)
			}

			provisioningState := "Unknown"
			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.ProvisioningState != nil {
						provisioningState = string(*props.ProvisioningState)
					}
					if props.Error != nil && props.Error.Message != nil && *props.Error.Message != "" {
						return resp, provisioningState, errors.New(*props.Error.Message)
					}
					return resp, provisioningState, nil
				}
			}
			return resp, provisioningState, nil
		},
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish provisioning: %+v", id, err)
	}

	return nil
}
