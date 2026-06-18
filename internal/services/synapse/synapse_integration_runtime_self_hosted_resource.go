// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
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

func resourceSynapseIntegrationRuntimeSelfHosted() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseIntegrationRuntimeSelfHostedCreateUpdate,
		Read:   resourceSynapseIntegrationRuntimeSelfHostedRead,
		Update: resourceSynapseIntegrationRuntimeSelfHostedCreateUpdate,
		Delete: resourceSynapseIntegrationRuntimeSelfHostedDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationRuntimeID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SynapseIntegrationRuntimeSelfHostedV0ToV1{},
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
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid name for Self-Hosted Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
				),
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"authorization_key_primary": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"authorization_key_secondary": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSynapseIntegrationRuntimeSelfHostedCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
				return tf.ImportAsExistsError("azurerm_synapse_integration_runtime_self_hosted", id.ID())
			}
		}
	}

	integrationRuntime := integrationruntime.IntegrationRuntimeResource{
		Name: pointer.To(id.IntegrationRuntimeName),
		Properties: integrationruntime.SelfHostedIntegrationRuntime{
			Description: pointer.To(d.Get("description").(string)),
			Type:        integrationruntime.IntegrationRuntimeTypeSelfHosted,
		},
	}

	if err := client.CreateThenPoll(ctx, id, integrationRuntime, integrationruntime.DefaultCreateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceSynapseIntegrationRuntimeSelfHostedRead(d, meta)
}

func resourceSynapseIntegrationRuntimeSelfHostedRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.IntegrationRuntimesClient
	authKeysClient := meta.(*clients.Client).Synapse.IntegrationRuntimeAuthKeysClient
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
		selfHostedIntegrationRuntime, ok := model.Properties.(integrationruntime.SelfHostedIntegrationRuntime)
		if !ok {
			return fmt.Errorf("converting integration runtime to Self-Hosted integration runtime (%q)", id)
		}

		d.Set("description", selfHostedIntegrationRuntime.Description)
	}

	respKey, err := authKeysClient.AuthKeysList(ctx, *id)
	if err != nil {
		if response.WasNotFound(respKey.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving auth keys (%q): %+v", id, err)
	}

	if model := respKey.Model; model != nil {
		d.Set("authorization_key_primary", model.AuthKey1)
		d.Set("authorization_key_secondary", model.AuthKey2)
	}

	return nil
}

func resourceSynapseIntegrationRuntimeSelfHostedDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
