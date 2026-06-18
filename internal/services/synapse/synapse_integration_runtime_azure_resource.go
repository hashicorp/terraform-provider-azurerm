// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/integrationruntime"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseIntegrationRuntimeAzure() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseIntegrationRuntimeAzureCreateUpdate,
		Read:   resourceSynapseIntegrationRuntimeAzureRead,
		Update: resourceSynapseIntegrationRuntimeAzureCreateUpdate,
		Delete: resourceSynapseIntegrationRuntimeAzureDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationRuntimeID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SynapseIntegrationRuntimeAzureV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^([a-zA-Z0-9](-|-?[a-zA-Z0-9]+)+[a-zA-Z0-9])$`),
					`Invalid name for Azure Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
				),
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"location": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					location.EnhancedValidate,
					validation.StringInSlice([]string{"AutoResolve"}, false),
				),
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},

			"compute_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(integrationruntime.DataFlowComputeTypeGeneral),
				ValidateFunc: validation.StringInSlice([]string{
					string(integrationruntime.DataFlowComputeTypeGeneral),
					string(integrationruntime.DataFlowComputeTypeComputeOptimized),
					string(integrationruntime.DataFlowComputeTypeMemoryOptimized),
				}, false),
			},

			"core_count": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  8,
				ValidateFunc: validation.IntInSlice([]int{
					8, 16, 32, 48, 80, 144, 272,
				}),
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"time_to_live_min": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceSynapseIntegrationRuntimeAzureCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	id := integrationruntime.NewIntegrationRuntimeID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, integrationruntime.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_synapse_integration_runtime_azure", id.ID())
			}
		}
	}

	integrationRuntime := integrationruntime.IntegrationRuntimeResource{
		Name: pointer.To(id.IntegrationRuntimeName),
		Properties: integrationruntime.ManagedIntegrationRuntime{
			Description: pointer.To(d.Get("description").(string)),
			Type:        integrationruntime.IntegrationRuntimeTypeManaged,
			TypeProperties: integrationruntime.ManagedIntegrationRuntimeTypeProperties{
				ComputeProperties: &integrationruntime.IntegrationRuntimeComputeProperties{
					Location: pointer.To(location.Normalize(d.Get("location").(string))),
					DataFlowProperties: &integrationruntime.IntegrationRuntimeDataFlowProperties{
						ComputeType: pointer.To(integrationruntime.DataFlowComputeType(d.Get("compute_type").(string))),
						CoreCount:   pointer.To(int64(d.Get("core_count").(int))),
						TimeToLive:  pointer.To(int64(d.Get("time_to_live_min").(int))),
					},
				},
			},
		},
	}

	if err := client.CreateThenPoll(ctx, id, integrationRuntime, integrationruntime.DefaultCreateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceSynapseIntegrationRuntimeAzureRead(d, meta)
}

func resourceSynapseIntegrationRuntimeAzureRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntime.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, integrationruntime.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.IntegrationRuntimeName)
	d.Set("synapse_workspace_id", parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	if model := resp.Model; model != nil {
		managedIntegrationRuntime, ok := model.Properties.(integrationruntime.ManagedIntegrationRuntime)
		if !ok {
			return fmt.Errorf("converting managed integration runtime to Azure integration runtime (%q)", id)
		}

		d.Set("description", managedIntegrationRuntime.Description)
		if computeProps := managedIntegrationRuntime.TypeProperties.ComputeProperties; computeProps != nil {
			d.Set("location", location.NormalizeNilable(computeProps.Location))
			if dataFlowProps := computeProps.DataFlowProperties; dataFlowProps != nil {
				computeType := ""
				if dataFlowProps.ComputeType != nil {
					computeType = string(*dataFlowProps.ComputeType)
				}
				d.Set("compute_type", computeType)
				d.Set("core_count", pointer.From(dataFlowProps.CoreCount))
				d.Set("time_to_live_min", pointer.From(dataFlowProps.TimeToLive))
			}
		}
	}

	return nil
}

func resourceSynapseIntegrationRuntimeAzureDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntime.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
