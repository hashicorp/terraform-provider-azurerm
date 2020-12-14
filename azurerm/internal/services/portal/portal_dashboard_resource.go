package portal

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/portal/mgmt/2019-01-01-preview/portal"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/portal/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/portal/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceDashboardCreateUpdate,
		Read:   resourceDashboardRead,
		Update: resourceDashboardCreateUpdate,
		Delete: resourceDashboardDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DashboardID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DashboardName,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),
			"location":            azure.SchemaLocation(),
			"tags":                tags.Schema(),
			"dashboard_properties": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: utils.NormalizeJson,
			},
		},
	}
}

func resourceDashboardCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	t := d.Get("tags").(map[string]interface{})
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	dashboardProps := d.Get("dashboard_properties").(string)

	// TODO: requires import support

	dashboard := portal.Dashboard{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	var dashboardProperties portal.DashboardProperties

	if err := json.Unmarshal([]byte(dashboardProps), &dashboardProperties); err != nil {
		return fmt.Errorf("Error parsing JSON: %+v", err)
	}
	dashboard.DashboardProperties = &dashboardProperties

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, dashboard); err != nil {
		return fmt.Errorf("creating/updating Dashboard %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(parse.NewDashboardID(subscriptionId, resourceGroup, name).ID())
	return resourceDashboardRead(d, meta)
}

func resourceDashboardRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DashboardID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Dashboard %q was not found in Resource Group %q - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Dashboard %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*resp.Location))
	}

	props, jsonErr := json.Marshal(resp.DashboardProperties)
	if jsonErr != nil {
		return fmt.Errorf("parsing JSON for Dashboard Properties: %+v", jsonErr)
	}
	d.Set("dashboard_properties", string(props))

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DashboardID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Dashboard %q (Resource Group %q): %+v", id.Name, id.Name, err)
		}
	}

	return nil
}
