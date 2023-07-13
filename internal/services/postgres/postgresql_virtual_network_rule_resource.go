// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/virtualnetworkrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePostgreSQLVirtualNetworkRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgreSQLVirtualNetworkRuleCreateUpdate,
		Read:   resourcePostgreSQLVirtualNetworkRuleRead,
		Update: resourcePostgreSQLVirtualNetworkRuleCreateUpdate,
		Delete: resourcePostgreSQLVirtualNetworkRuleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualnetworkrules.ParseVirtualNetworkRuleID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				// TODO: this should be using a validation func within the Postgres package
				ValidateFunc: networkValidate.VirtualNetworkRuleName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"ignore_missing_vnet_service_endpoint": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourcePostgreSQLVirtualNetworkRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.VirtualNetworkRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualnetworkrules.NewVirtualNetworkRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_postgresql_virtual_network_rule", id.ID())
		}
	}

	parameters := virtualnetworkrules.VirtualNetworkRule{
		Properties: &virtualnetworkrules.VirtualNetworkRuleProperties{
			VirtualNetworkSubnetId:           d.Get("subnet_id").(string),
			IgnoreMissingVnetServiceEndpoint: utils.Bool(d.Get("ignore_missing_vnet_service_endpoint").(bool)),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Wait for the provisioning state to become ready
	log.Printf("[DEBUG] Waiting for %s to become ready", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Initializing", "InProgress", "Unknown", "ResponseNotFound"},
		Target:                    []string{"Ready"},
		Refresh:                   postgreSQLVirtualNetworkStateStatusCodeRefreshFunc(ctx, client, id),
		MinTimeout:                1 * time.Minute,
		ContinuousTargetOccurence: 5,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be created or updated: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePostgreSQLVirtualNetworkRuleRead(d, meta)
}

func resourcePostgreSQLVirtualNetworkRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkrules.ParseVirtualNetworkRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error reading PostgreSQL Virtual Network Rule %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.VirtualNetworkRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("server_name", id.ServerName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("subnet_id", props.VirtualNetworkSubnetId)
			d.Set("ignore_missing_vnet_service_endpoint", props.IgnoreMissingVnetServiceEndpoint)
		}
	}

	return nil
}

func resourcePostgreSQLVirtualNetworkRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkrules.ParseVirtualNetworkRuleID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func postgreSQLVirtualNetworkStateStatusCodeRefreshFunc(ctx context.Context, client *virtualnetworkrules.VirtualNetworkRulesClient, id virtualnetworkrules.VirtualNetworkRuleId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				log.Printf("[DEBUG] Retrieving %s returned 404.", id)
				return nil, "ResponseNotFound", nil
			}

			return nil, "", fmt.Errorf("polling for the state of the %s: %+v", id, err)
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil && props.State != nil {
				log.Printf("[DEBUG] Retrieving %s returned Status %s", id, string(*props.State))
				return resp, string(*props.State), nil
			}
		}

		// Valid response was returned but VirtualNetworkRuleProperties was nil. Basically the rule exists, but with no properties for some reason. Assume Unknown instead of returning error.
		log.Printf("[DEBUG] Retrieving %s returned empty `properties``", id)
		return resp, "Unknown", nil
	}
}
