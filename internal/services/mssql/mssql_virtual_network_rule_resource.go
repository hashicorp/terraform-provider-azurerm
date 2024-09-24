// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/virtualnetworkrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlVirtualNetworkRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlVirtualNetworkRuleCreateUpdate,
		Read:   resourceMsSqlVirtualNetworkRuleRead,
		Update: resourceMsSqlVirtualNetworkRuleCreateUpdate,
		Delete: resourceMsSqlVirtualNetworkRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualNetworkRuleID(id)
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
				ValidateFunc: validate.VirtualNetworkRuleName,
			},

			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerID,
			},

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"ignore_missing_vnet_service_endpoint": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false, // When not provided, Azure defaults to false
			},
		},
	}
}

func resourceMsSqlVirtualNetworkRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverId, err := parse.ServerID(d.Get("server_id").(string))
	if err != nil {
		return fmt.Errorf("parsing server ID %q: %+v", d.Get("server_id"), err)
	}

	id := virtualnetworkrules.NewVirtualNetworkRuleID(serverId.SubscriptionId, serverId.ResourceGroup, serverId.Name, d.Get("name").(string))

	subnetId, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("parsing subnet ID %q: %+v", d.Get("subnet_id"), err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing MSSQL %s: %+v", id.String(), err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mssql_virtual_network_rule", id.ID())
		}
	}

	parameters := virtualnetworkrules.VirtualNetworkRule{
		Properties: &virtualnetworkrules.VirtualNetworkRuleProperties{
			VirtualNetworkSubnetId:           subnetId.ID(),
			IgnoreMissingVnetServiceEndpoint: utils.Bool(d.Get("ignore_missing_vnet_service_endpoint").(bool)),
		},
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating MSSQL %s: %+v", id.String(), err)
	}

	// Wait for the provisioning state to become ready
	log.Printf("[DEBUG] Waiting for MSSQL %s to become ready", id.String())

	d.SetId(id.ID())

	return resourceMsSqlVirtualNetworkRuleRead(d, meta)
}

func resourceMsSqlVirtualNetworkRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkrules.ParseVirtualNetworkRuleID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing ID %q: %+v", d.Id(), err)
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] MSSQL %s was not found - removing from state", id.String())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving MSSQL %s: %+v", id.String(), err)
	}

	d.Set("name", id.VirtualNetworkRuleName)

	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
	d.Set("server_id", serverId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ignore_missing_vnet_service_endpoint", props.IgnoreMissingVnetServiceEndpoint)

			subnetId, err := commonids.ParseSubnetIDInsensitively(props.VirtualNetworkSubnetId)
			if err != nil {
				return fmt.Errorf("parsing subnet ID returned by API %q: %+v", props.VirtualNetworkSubnetId, err)
			}

			d.Set("subnet_id", subnetId.ID())
		}
	}

	return nil
}

func resourceMsSqlVirtualNetworkRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkrules.ParseVirtualNetworkRuleID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing ID %q: %+v", d.Id(), err)
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting MSSQL %s: %+v", id.String(), err)
	}

	return nil
}
