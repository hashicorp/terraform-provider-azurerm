// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
				Default:  string(synapse.DataFlowComputeTypeGeneral),
				ValidateFunc: validation.StringInSlice([]string{
					string(synapse.DataFlowComputeTypeGeneral),
					string(synapse.DataFlowComputeTypeComputeOptimized),
					string(synapse.DataFlowComputeTypeMemoryOptimized),
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

	id := parse.NewIntegrationRuntimeID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_integration_runtime_azure", id.ID())
		}
	}

	integrationRuntime := synapse.IntegrationRuntimeResource{
		Name: utils.String(id.Name),
		Properties: synapse.ManagedIntegrationRuntime{
			Description: utils.String(d.Get("description").(string)),
			Type:        synapse.TypeBasicIntegrationRuntimeTypeManaged,
			ManagedIntegrationRuntimeTypeProperties: &synapse.ManagedIntegrationRuntimeTypeProperties{
				ComputeProperties: &synapse.IntegrationRuntimeComputeProperties{
					Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
					DataFlowProperties: &synapse.IntegrationRuntimeDataFlowProperties{
						ComputeType: synapse.DataFlowComputeType(d.Get("compute_type").(string)),
						CoreCount:   utils.Int32(int32(d.Get("core_count").(int))),
						TimeToLive:  utils.Int32(int32(d.Get("time_to_live_min").(int))),
					},
				},
			},
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, integrationRuntime, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSynapseIntegrationRuntimeAzureRead(d, meta)
}

func resourceSynapseIntegrationRuntimeAzureRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("synapse_workspace_id", parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())

	managedIntegrationRuntime, convertSuccess := resp.Properties.AsManagedIntegrationRuntime()
	if !convertSuccess {
		return fmt.Errorf("converting managed integration runtime to Azure integration runtime (%q)", id)
	}

	d.Set("description", managedIntegrationRuntime.Description)
	if computeProps := managedIntegrationRuntime.ComputeProperties; computeProps != nil {
		d.Set("location", computeProps.Location)
		if dataFlowProps := computeProps.DataFlowProperties; dataFlowProps != nil {
			d.Set("compute_type", dataFlowProps.ComputeType)
			d.Set("core_count", dataFlowProps.CoreCount)
			d.Set("time_to_live_min", dataFlowProps.TimeToLive)
		}
	}

	return nil
}

func resourceSynapseIntegrationRuntimeAzureDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}
