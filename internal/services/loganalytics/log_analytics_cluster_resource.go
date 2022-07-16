package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
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

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

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
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_log_analytics_cluster", id.ID())
	}

	expandedIdentity, err := expandLogAnalyticsClusterIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	parameters := operationalinsights.Cluster{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: expandedIdentity,
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

	createWait, err := logAnalyticsClusterWaitForState(ctx, meta, id.ResourceGroup, id.ClusterName)
	if err != nil {
		return err
	}
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

	updateWait, err := logAnalyticsClusterWaitForState(ctx, meta, id.ResourceGroup, id.ClusterName)
	if err != nil {
		return err
	}

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

func expandLogAnalyticsClusterIdentity(input []interface{}) (*operationalinsights.Identity, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &operationalinsights.Identity{
		Type: operationalinsights.IdentityType(string(expanded.Type)),
	}, nil
}

func flattenLogAnalyticsIdentity(input *operationalinsights.Identity) []interface{} {
	var transform *identity.SystemAssigned

	if input != nil {
		transform = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
	}

	return identity.FlattenSystemAssigned(transform)
}
