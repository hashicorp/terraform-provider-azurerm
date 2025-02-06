// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceApplicationInsights() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmApplicationInsightsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"instrumentation_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"application_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"app_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"retention_in_days": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"workspace_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceArmApplicationInsightsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.ComponentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := components.NewComponentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

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

			workspaceId := ""
			if props.WorkspaceResourceId != nil {
				workspaceId = *props.WorkspaceResourceId
			}
			d.Set("workspace_id", workspaceId)
		}
	}
	return nil
}
