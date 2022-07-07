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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePortalDashboard() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePortalDashboardCreateUpdate,
		Read:   resourcePortalDashboardRead,
		Update: resourcePortalDashboardCreateUpdate,
		Delete: resourcePortalDashboardDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dashboard.ParseDashboardID(id)
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"tags": commonschema.Tags(),

			"dashboard_properties": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: utils.NormalizeJson,
			},
		},
	}
}

func resourcePortalDashboardCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return tf.ImportAsExistsError("azurerm_portal_dashboard", id.ID())
		}
	}

	props := dashboard.Dashboard{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	var dashboardProperties dashboard.DashboardProperties

	dashboardPropsRaw := d.Get("dashboard_properties").(string)
	if err := json.Unmarshal([]byte(dashboardPropsRaw), &dashboardProperties); err != nil {
		return fmt.Errorf("parsing JSON: %+v", err)
	}

	props.Properties = &dashboardProperties

	if _, err := client.CreateOrUpdate(ctx, id, props); err != nil {
		return fmt.Errorf("creating/updating %s %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePortalDashboardRead(d, meta)
}

func resourcePortalDashboardRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DashboardName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		if props := model.Properties; props != nil {
			v, err := json.Marshal(props)
			if err != nil {
				return fmt.Errorf("parsing JSON for Dashboard Properties: %+v", err)
			}
			d.Set("dashboard_properties", string(v))
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourcePortalDashboardDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
