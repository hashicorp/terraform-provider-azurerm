package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
			_, err := parse.LogAnalyticsClusterID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsClusterName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"identity": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(operationalinsights.SystemAssigned),
							}, false),
						},

						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			// Per the documentation cluster capacity must start at 1000 GB and can go above 3000 GB with an exception by Microsoft
			// so I am not limiting the upperbound here by design
			// https://docs.microsoft.com/en-us/azure/azure-monitor/platform/manage-cost-storage#log-analytics-dedicated-clusters
			"size_gb": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  1000,
				ValidateFunc: validation.All(
					validation.IntAtLeast(1000),
					validation.IntDivisibleBy(100),
				),
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

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewLogAnalyticsClusterID(subscriptionId, resourceGroup, name)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Log Analytics Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_log_analytics_cluster", *existing.ID)
	}

	parameters := operationalinsights.Cluster{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: expandLogAnalyticsClusterIdentity(d.Get("identity").([]interface{})),
		Sku: &operationalinsights.ClusterSku{
			Capacity: utils.Int64(int64(d.Get("size_gb").(int))),
			Name:     operationalinsights.CapacityReservation,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Log Analytics Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Log Analytics Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if _, err = client.Get(ctx, resourceGroup, name); err != nil {
		return fmt.Errorf("retrieving Log Analytics Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	createWait := logAnalyticsClusterWaitForState(ctx, meta, d.Timeout(pluginsdk.TimeoutCreate), id.ResourceGroup, id.ClusterName)

	if _, err := createWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Log Analytics Cluster to finish updating %q (Resource Group %q): %v", id.ClusterName, id.ResourceGroup, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsClusterRead(d, meta)
}

func resourceLogAnalyticsClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Log Analytics %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Log Analytics Cluster %q (Resource Group %q): %+v", id.ClusterName, id.ResourceGroup, err)
	}
	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if err := d.Set("identity", flattenLogAnalyticsIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	if props := resp.ClusterProperties; props != nil {
		d.Set("cluster_id", props.ClusterID)
	}

	capacity := 0
	if sku := resp.Sku; sku != nil {
		if sku.Capacity != nil {
			capacity = int(*sku.Capacity)
		}
	}
	d.Set("size_gb", capacity)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLogAnalyticsClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsClusterID(d.Id())
	if err != nil {
		return err
	}

	parameters := operationalinsights.ClusterPatch{}

	if d.HasChange("size_gb") {
		parameters.Sku = &operationalinsights.ClusterSku{
			Capacity: utils.Int64(int64(d.Get("size_gb").(int))),
			Name:     operationalinsights.CapacityReservation,
		}
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.ClusterName, parameters); err != nil {
		return fmt.Errorf("updating Log Analytics Cluster %q (Resource Group %q): %+v", id.ClusterName, id.ResourceGroup, err)
	}

	// Need to wait for the cluster to actually finish updating the resource before continuing
	// since the service returns a 200 instantly while it's still updating in the background
	log.Printf("[INFO] Checking for Log Analytics Cluster provisioning state")

	updateWait := logAnalyticsClusterWaitForState(ctx, meta, d.Timeout(pluginsdk.TimeoutUpdate), id.ResourceGroup, id.ClusterName)

	if _, err := updateWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Log Analytics Cluster to finish updating %q (Resource Group %q): %v", id.ClusterName, id.ResourceGroup, err)
	}

	return resourceLogAnalyticsClusterRead(d, meta)
}

func resourceLogAnalyticsClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ClusterName)
	if err != nil {
		return fmt.Errorf("deleting Log Analytics Cluster %q (Resource Group %q): %+v", id.ClusterName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Log Analytics Cluster %q (Resource Group %q): %+v", id.ClusterName, id.ResourceGroup, err)
	}

	return nil
}

func expandLogAnalyticsClusterIdentity(input []interface{}) *operationalinsights.Identity {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &operationalinsights.Identity{
		Type: operationalinsights.IdentityType(v["type"].(string)),
	}
}

func flattenLogAnalyticsIdentity(input *operationalinsights.Identity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var t operationalinsights.IdentityType
	if input.Type != "" {
		t = input.Type
	}
	var principalId string
	if input.PrincipalID != nil {
		principalId = *input.PrincipalID
	}
	var tenantId string
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}
	return []interface{}{
		map[string]interface{}{
			"type":         t,
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}
}
