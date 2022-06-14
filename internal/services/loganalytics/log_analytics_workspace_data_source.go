package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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

			"tags": tags.SchemaDataSource(),
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

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("log analytics workspaces %q (Resource Group %q) was not found", name, resGroup)
		}
		return fmt.Errorf("making Read request on AzureRM Log Analytics workspaces '%s': %+v", name, err)
	}

	id := parse.NewLogAnalyticsWorkspaceID(subscriptionId, resGroup, name)
	d.SetId(id.ID())

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	d.Set("workspace_id", resp.CustomerID)
	if sku := resp.Sku; sku != nil {
		d.Set("sku", sku.Name)
	}
	d.Set("retention_in_days", resp.RetentionInDays)

	if workspaceCapping := resp.WorkspaceCapping; workspaceCapping != nil {
		d.Set("daily_quota_gb", resp.WorkspaceCapping.DailyQuotaGb)
	} else {
		d.Set("daily_quota_gb", utils.Float(-1))
	}

	sharedKeys, err := sharedKeysClient.GetSharedKeys(ctx, resGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Shared keys for Log Analytics workspaces %s: %+v", name, err)
	} else {
		d.Set("primary_shared_key", sharedKeys.PrimarySharedKey)
		d.Set("secondary_shared_key", sharedKeys.SecondarySharedKey)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
