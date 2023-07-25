// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
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

	id := parse.NewVirtualNetworkRuleID(serverId.SubscriptionId, serverId.ResourceGroup, serverId.Name, d.Get("name").(string))

	subnetId, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("parsing subnet ID %q: %+v", d.Get("subnet_id"), err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing MSSQL %s: %+v", id.String(), err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_mssql_virtual_network_rule", id.ID())
		}
	}

	parameters := sql.VirtualNetworkRule{
		VirtualNetworkRuleProperties: &sql.VirtualNetworkRuleProperties{
			VirtualNetworkSubnetID:           utils.String(subnetId.ID()),
			IgnoreMissingVnetServiceEndpoint: utils.Bool(d.Get("ignore_missing_vnet_service_endpoint").(bool)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating MSSQL %s: %+v", id.String(), err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
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

	id, err := parse.VirtualNetworkRuleID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing ID %q: %+v", d.Id(), err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] MSSQL %s was not found - removing from state", id.String())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving MSSQL %s: %+v", id.String(), err)
	}

	d.Set("name", id.Name)

	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)
	d.Set("server_id", serverId.ID())

	if props := resp.VirtualNetworkRuleProperties; props != nil {
		d.Set("ignore_missing_vnet_service_endpoint", props.IgnoreMissingVnetServiceEndpoint)

		subnetId := ""
		if sid := props.VirtualNetworkSubnetID; sid != nil {
			id, err := commonids.ParseSubnetIDInsensitively(*props.VirtualNetworkSubnetID)
			if err != nil {
				return fmt.Errorf("parsing subnet ID returned by API %q: %+v", *sid, err)
			}
			subnetId = id.ID()
		}
		d.Set("subnet_id", subnetId)
	}

	return nil
}

func resourceMsSqlVirtualNetworkRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkRuleID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing ID %q: %+v", d.Id(), err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting MSSQL %s: %+v", id.String(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of MSSQL %s: %+v", id.String(), err)
	}

	return nil
}
