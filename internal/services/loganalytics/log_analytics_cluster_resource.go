package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsClusterCreate,
		Read:   resourceLogAnalyticsClusterRead,
		Update: resourceLogAnalyticsClusterUpdate,
		Delete: resourceLogAnalyticsClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(6 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(6 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsClusterName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

			"size_gb": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default: func() int {
					if !features.FourPointOh() {
						return 1000
					}
					return 500
				}(),
				ValidateFunc: validation.IntInSlice([]int{500, 1000, 2000, 5000}),
			},

			"cluster_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceLogAnalyticsClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := clusters.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_log_analytics_cluster", id.ID())
	}

	expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	capacityReservation := clusters.ClusterSkuNameEnumCapacityReservation
	parameters := clusters.Cluster{
		Location: location.Normalize(d.Get("location").(string)),
		Identity: expandedIdentity,
		Sku: &clusters.ClusterSku{
			Capacity: utils.Int64(int64(d.Get("size_gb").(int))),
			Name:     &capacityReservation,
		},
		Tags: expandTags(d.Get("tags").(map[string]interface{})),
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if _, err = client.Get(ctx, id); err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	createWait, err := logAnalyticsClusterWaitForState(ctx, client, id)
	if err != nil {
		return err
	}
	if _, err := createWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Log Analytics Cluster to finish updating %q (Resource Group %q): %v", id.ClusterName, id.ResourceGroupName, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsClusterRead(d, meta)
}

func resourceLogAnalyticsClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Log Analytics %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Log Analytics Cluster %q (Resource Group %q): %+v", id.ClusterName, id.ResourceGroupName, err)
	}
	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))
		if err = d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
		if props := model.Properties; props != nil {
			d.Set("cluster_id", props.ClusterId)
		}
		capacity := 0
		if sku := model.Sku; sku != nil {
			if sku.Capacity != nil {
				capacity = int(*sku.Capacity)
			}
		}
		d.Set("size_gb", capacity)

		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}

func resourceLogAnalyticsClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	parameters := clusters.ClusterPatch{}
	capacityReservation := clusters.ClusterSkuNameEnumCapacityReservation
	if d.HasChange("size_gb") {
		parameters.Sku = &clusters.ClusterSku{
			Capacity: utils.Int64(int64(d.Get("size_gb").(int))),
			Name:     &capacityReservation,
		}
	}

	if d.HasChange("tags") {
		parameters.Tags = expandTags(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// Need to wait for the cluster to actually finish updating the resource before continuing
	// since the service returns a 200 instantly while it's still updating in the background
	log.Printf("[INFO] Checking for Log Analytics Cluster provisioning state")

	updateWait, err := logAnalyticsClusterWaitForState(ctx, client, *id)
	if err != nil {
		return err
	}

	if _, err := updateWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Log Analytics Cluster to finish updating %q (Resource Group %q): %v", id.ClusterName, id.ResourceGroupName, err)
	}

	return resourceLogAnalyticsClusterRead(d, meta)
}

func resourceLogAnalyticsClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting Log Analytics Cluster %q (Resource Group %q): %+v", id.ClusterName, id.ResourceGroupName, err)
	}

	return nil
}
