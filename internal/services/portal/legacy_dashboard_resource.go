// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package portal

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLegacyDashboard() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLegacyDashboardCreateUpdate,
		Read:   resourceLegacyDashboardRead,
		Update: resourceLegacyDashboardCreateUpdate,
		Delete: resourceLegacyDashboardDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DashboardID(id)
			return err
		}),

		DeprecationMessage: "The `azurerm_dashboard` resource is deprecated and will be removed in v4.0 of the AzureRM Provider - the replacement is available as `azurerm_portal_dashboard`.",

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
				ValidateFunc: validate.DashboardName,
			},
			"resource_group_name": commonschema.ResourceGroupName(),
			"location":            commonschema.Location(),
			"tags":                commonschema.Tags(),
			"dashboard_properties": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: utils.NormalizeJson,
			},
		},
	}
}

func resourceLegacyDashboardCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dashboard.NewDashboardID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_dashboard", id.ID())
		}
	}

	payload := dashboard.Dashboard{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	var dashboardProperties dashboard.DashboardProperties

	dashboardPropsRaw := d.Get("dashboard_properties").(string)
	if err := json.Unmarshal([]byte(dashboardPropsRaw), &dashboardProperties); err != nil {
		return fmt.Errorf("parsing JSON: %+v", err)
	}
	payload.Properties = &dashboardProperties

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLegacyDashboardRead(d, meta)
}

func resourceLegacyDashboardRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dashboard.ParseDashboardID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DashboardName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		props, jsonErr := json.Marshal(model.Properties)
		if jsonErr != nil {
			return fmt.Errorf("parsing JSON for Dashboard Properties: %+v", jsonErr)
		}
		d.Set("dashboard_properties", string(props))

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceLegacyDashboardDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dashboard.ParseDashboardID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
