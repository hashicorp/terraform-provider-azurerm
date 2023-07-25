// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSynapseWorkspaceSecurityAlertPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseWorkspaceSecurityAlertPolicyCreateUpdate,
		Read:   resourceSynapseWorkspaceSecurityAlertPolicyRead,
		Update: resourceSynapseWorkspaceSecurityAlertPolicyCreateUpdate,
		Delete: resourceSynapseWorkspaceSecurityAlertPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspaceSecurityAlertPolicyID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"disabled_alerts": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Set:      pluginsdk.HashString,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"Sql_Injection",
						"Sql_Injection_Vulnerability",
						"Access_Anomaly",
						"Data_Exfiltration",
						"Unsafe_Action",
					}, false),
				},
			},

			"email_account_admins_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"email_addresses": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
				Set: pluginsdk.HashString,
			},

			"retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"policy_state": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(synapse.SecurityAlertPolicyStateDisabled),
					string(synapse.SecurityAlertPolicyStateEnabled),
					string(synapse.SecurityAlertPolicyStateNew),
				}, false),
			},

			"storage_account_access_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_endpoint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceSynapseWorkspaceSecurityAlertPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceSecurityAlertPolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewWorkspaceSecurityAlertPolicyID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, "Default")

	alertPolicy := expandServerSecurityAlertPolicy(d)

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, *alertPolicy)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSynapseWorkspaceSecurityAlertPolicyRead(d, meta)
}

func resourceSynapseWorkspaceSecurityAlertPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceSecurityAlertPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] synapse %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("synapse_workspace_id", workspaceId.ID())

	if props := resp.ServerSecurityAlertPolicyProperties; props != nil {
		d.Set("policy_state", string(props.State))

		if props.DisabledAlerts != nil {
			disabledAlerts := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
			for _, v := range *props.DisabledAlerts {
				if v != "" {
					disabledAlerts.Add(v)
				}
			}

			d.Set("disabled_alerts", disabledAlerts)
		}

		if props.EmailAccountAdmins != nil {
			d.Set("email_account_admins_enabled", props.EmailAccountAdmins)
		}

		if props.EmailAddresses != nil {
			emailAddresses := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
			for _, v := range *props.EmailAddresses {
				if v != "" {
					emailAddresses.Add(v)
				}
			}

			d.Set("email_addresses", emailAddresses)
		}

		if props.RetentionDays != nil {
			d.Set("retention_days", int(*props.RetentionDays))
		}

		if v, ok := d.GetOk("storage_account_access_key"); ok {
			d.Set("storage_account_access_key", v)
		}

		if props.StorageEndpoint != nil {
			d.Set("storage_endpoint", props.StorageEndpoint)
		}
	}

	return nil
}

func resourceSynapseWorkspaceSecurityAlertPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceSecurityAlertPolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	disabledPolicy := synapse.ServerSecurityAlertPolicy{
		ServerSecurityAlertPolicyProperties: &synapse.ServerSecurityAlertPolicyProperties{
			State: synapse.SecurityAlertPolicyStateDisabled,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, disabledPolicy)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}

func expandServerSecurityAlertPolicy(d *pluginsdk.ResourceData) *synapse.ServerSecurityAlertPolicy {
	policyState := synapse.SecurityAlertPolicyState(d.Get("policy_state").(string))

	policy := synapse.ServerSecurityAlertPolicy{
		ServerSecurityAlertPolicyProperties: &synapse.ServerSecurityAlertPolicyProperties{
			State: policyState,
		},
	}

	props := policy.ServerSecurityAlertPolicyProperties

	if v, ok := d.GetOk("disabled_alerts"); ok {
		disabledAlerts := make([]string, 0)
		for _, v := range v.(*pluginsdk.Set).List() {
			disabledAlerts = append(disabledAlerts, v.(string))
		}
		props.DisabledAlerts = &disabledAlerts
	}

	if v, ok := d.GetOk("email_addresses"); ok {
		emailAddresses := make([]string, 0)
		for _, v := range v.(*pluginsdk.Set).List() {
			emailAddresses = append(emailAddresses, v.(string))
		}
		props.EmailAddresses = &emailAddresses
	}

	if v, ok := d.GetOk("email_account_admins_enabled"); ok {
		props.EmailAccountAdmins = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("retention_days"); ok {
		props.RetentionDays = utils.Int32(int32(v.(int)))
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		props.StorageAccountAccessKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("storage_endpoint"); ok {
		props.StorageEndpoint = utils.String(v.(string))
	}

	return &policy
}
