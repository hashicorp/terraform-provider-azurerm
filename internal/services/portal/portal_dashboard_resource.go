// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2026-04-01/dashboards"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePortalDashboard() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourcePortalDashboardCreateUpdate,
		Read:   resourcePortalDashboardRead,
		Update: resourcePortalDashboardCreateUpdate,
		Delete: resourcePortalDashboardDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dashboards.ParseDashboardID(id)
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
				ValidateFunc: validate.DashboardName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"tags": commonschema.Tags(),

			"dashboard_properties": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DashboardProperties,
				StateFunc:    utils.NormalizeJson,
			},
		},
	}

	if !features.FivePointOh() {
		resource.Schema["dashboard_properties"].DiffSuppressOnRefresh = true
		resource.Schema["dashboard_properties"].DiffSuppressFunc = func(k, old, new string, d *schema.ResourceData) bool {
			if old == new {
				return true
			}
			if parsedLegacy, ok := parse.LegacyDashboardProperties(new); ok {
				if converted, err := json.Marshal(parsedLegacy); err == nil {
					return utils.NormalizeJson(old) == utils.NormalizeJson(string(converted))
				}
			}
			return false
		}
	}
	return resource
}

func resourcePortalDashboardCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dashboards.NewDashboardID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_portal_dashboard", id.ID())
			}
		}
	}

	props := dashboards.Dashboard{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	dashboardPropsRaw := d.Get("dashboard_properties").(string)

	if !features.FivePointOh() {
		if dashboardProperties, ok := parse.LegacyDashboardProperties(dashboardPropsRaw); ok {
			props.Properties = dashboardProperties

			if _, err := client.CreateOrUpdate(ctx, id, props); err != nil {
				return fmt.Errorf("creating/updating %s %+v", id, err)
			}

			if d.IsNewResource() {
				d.SetId(id.ID())
			}
			return resourcePortalDashboardRead(d, meta)
		}
	}

	var dashboardProperties dashboards.DashboardPropertiesWithProvisioningState
	if err := json.Unmarshal([]byte(dashboardPropsRaw), &dashboardProperties); err != nil {
		return fmt.Errorf("parsing JSON: %+v", err)
	}
	parse.NormalizeDashboardPartMetadata(&dashboardProperties)

	props.Properties = &dashboardProperties

	if _, err := client.CreateOrUpdate(ctx, id, props); err != nil {
		return fmt.Errorf("creating/updating %s %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourcePortalDashboardRead(d, meta)
}

func resourcePortalDashboardRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dashboards.ParseDashboardID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DashboardName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			parse.NormalizeDashboardPartMetadata(props)
			v, err := json.Marshal(props)
			if err != nil {
				return fmt.Errorf("parsing JSON for Dashboard Properties: %+v", err)
			}
			d.Set("dashboard_properties", string(v))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourcePortalDashboardDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dashboards.ParseDashboardID(d.Id())
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
