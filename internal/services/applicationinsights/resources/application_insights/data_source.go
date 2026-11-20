// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_insights

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Application Insights component.",
		Read:        dataSourceArmApplicationInsightsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"name": {
				Description:  "Specifies the name of the Application Insights component.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"instrumentation_key": {
				Description: "The instrumentation key of the Application Insights component.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"connection_string": {
				Description: "The connection string of the Application Insights component. (Sensitive)",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"location": {
				Description: "The Azure location where the component exists.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"application_type": {
				Description: "The type of the component.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"app_id": {
				Description: "The App ID associated with this Application Insights component.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"retention_in_days": {
				Description: "The retention period in days.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"workspace_id": {
				Description: "The id of the associated Log Analytics workspace",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": {
				Description: "Tags applied to the component.",
				Type:        pluginsdk.TypeMap,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceArmApplicationInsightsRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.ComponentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := componentsapis.NewComponentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.ComponentsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}
		if props := model.Properties; props != nil {
			d.Set("app_id", props.AppId)
			d.Set("application_type", props.ApplicationType)
			d.Set("connection_string", props.ConnectionString)
			d.Set("instrumentation_key", props.InstrumentationKey)
			retentionInDays := 0
			if props.RetentionInDays != nil {
				retentionInDays = int(*props.RetentionInDays)
			}
			d.Set("retention_in_days", retentionInDays)
			d.Set("workspace_id", pointer.From(props.WorkspaceResourceId))
		}
	}
	return nil
}
