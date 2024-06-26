// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-07-01-preview/privatelinkscopesapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorPrivateLinkScope() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorPrivateLinkScopeCreateUpdate,
		Read:   resourceMonitorPrivateLinkScopeRead,
		Update: resourceMonitorPrivateLinkScopeCreateUpdate,
		Delete: resourceMonitorPrivateLinkScopeDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := privatelinkscopesapis.ParsePrivateLinkScopeID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PrivateLinkScopeName,
			},

			"ingestion_access_mode": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(privatelinkscopesapis.AccessModeOpen),
				ValidateFunc: validation.StringInSlice(privatelinkscopesapis.PossibleValuesForAccessMode(), false),
			},

			"query_access_mode": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(privatelinkscopesapis.AccessModeOpen),
				ValidateFunc: validation.StringInSlice(privatelinkscopesapis.PossibleValuesForAccessMode(), false),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorPrivateLinkScopeCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Monitor.PrivateLinkScopesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := privatelinkscopesapis.NewPrivateLinkScopeID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.PrivateLinkScopesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_private_link_scope", id.ID())
		}
	}

	ingestionAccessMode := privatelinkscopesapis.AccessMode(d.Get("ingestion_access_mode").(string))
	queryaccessMode := privatelinkscopesapis.AccessMode(d.Get("query_access_mode").(string))

	parameters := privatelinkscopesapis.AzureMonitorPrivateLinkScope{
		Name:     &name,
		Location: "Global",
		Tags:     utils.ExpandPtrMapStringString(d.Get("tags").(map[string]interface{})),
		Properties: privatelinkscopesapis.AzureMonitorPrivateLinkScopeProperties{
			AccessModeSettings: privatelinkscopesapis.AccessModeSettings{
				IngestionAccessMode: ingestionAccessMode,
				QueryAccessMode:     queryaccessMode,
			},
		},
	}

	if _, err := client.PrivateLinkScopesCreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorPrivateLinkScopeRead(d, meta)
}

func resourceMonitorPrivateLinkScopeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.PrivateLinkScopesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkscopesapis.ParsePrivateLinkScopeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.PrivateLinkScopesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PrivateLinkScopeName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if err = d.Set("tags", utils.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}

		props := model.Properties
		d.Set("ingestion_access_mode", string(props.AccessModeSettings.IngestionAccessMode))
		d.Set("query_access_mode", string(props.AccessModeSettings.QueryAccessMode))

	}

	return nil
}

func resourceMonitorPrivateLinkScopeDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.PrivateLinkScopesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkscopesapis.ParsePrivateLinkScopeID(d.Id())
	if err != nil {
		return err
	}

	err = client.PrivateLinkScopesDeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
