package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceLogAnalyticsWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceLogAnalyticsWorkspaceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"retention_in_days": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"daily_quota_gb": {
				Type:     pluginsdk.TypeFloat,
				Computed: true,
			},

			"workspace_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_shared_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_shared_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceLogAnalyticsWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	sharedKeysClient := meta.(*clients.Client).LogAnalytics.SharedKeysClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	id := workspaces.NewWorkspaceID(subscriptionId, resGroup, name)

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("log analytics workspaces %q (Resource Group %q) was not found", name, resGroup)
		}
		return fmt.Errorf("making Read request on AzureRM Log Analytics workspaces '%s': %+v", name, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))

		if props := model.Properties; props != nil {
			d.Set("workspace_id", props.CustomerId)

			if sku := props.Sku; sku != nil {
				d.Set("sku", sku.Name)
			}
			d.Set("retention_in_days", props.RetentionInDays)

			if workspaceCapping := props.WorkspaceCapping; workspaceCapping != nil {
				d.Set("daily_quota_gb", workspaceCapping.DailyQuotaGb)
			} else {
				d.Set("daily_quota_gb", utils.Float(-1))
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	sharedKeys, err := sharedKeysClient.GetSharedKeys(ctx, resGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Shared keys for Log Analytics workspaces %s: %+v", name, err)
	} else {
		d.Set("primary_shared_key", sharedKeys.PrimarySharedKey)
		d.Set("secondary_shared_key", sharedKeys.SecondarySharedKey)
	}

	return nil
}
