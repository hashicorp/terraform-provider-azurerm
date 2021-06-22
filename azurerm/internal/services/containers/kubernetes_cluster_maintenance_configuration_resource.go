package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-03-01/containerservice"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/parse"
	containerValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKubernetesClusterMaintenanceConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKubernetesClusterMaintenanceConfigurationCreateUpdate,
		Read:   resourceKubernetesClusterMaintenanceConfigurationRead,
		Update: resourceKubernetesClusterMaintenanceConfigurationCreateUpdate,
		Delete: resourceKubernetesClusterMaintenanceConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.MaintenanceConfigurationID(id)
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"kubernetes_cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: containerValidate.ClusterID,
			},

			"maintenance_allowed": {
				Type:         pluginsdk.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"maintenance_allowed", "maintenance_not_allowed_window"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"day": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.WeekDaySunday),
								string(containerservice.WeekDayMonday),
								string(containerservice.WeekDayTuesday),
								string(containerservice.WeekDayWednesday),
								string(containerservice.WeekDayThursday),
								string(containerservice.WeekDayFriday),
								string(containerservice.WeekDaySaturday),
							}, false),
						},

						"hour_slots": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeInt,
								ValidateFunc: validation.IntBetween(0, 23),
							},
						},
					},
				},
			},

			"maintenance_not_allowed_window": {
				Type:         pluginsdk.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"maintenance_allowed", "maintenance_not_allowed_window"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"end": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.RFC3339Time,
							ValidateFunc:     validation.IsRFC3339Time,
						},

						"start": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.RFC3339Time,
							ValidateFunc:     validation.IsRFC3339Time,
						},
					},
				},
			},
		},
	}
}

func resourceKubernetesClusterMaintenanceConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	kubernetesClusterId, err := parse.ClusterID(d.Get("kubernetes_cluster_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewMaintenanceConfigurationID(subscriptionId, kubernetesClusterId.ResourceGroup, kubernetesClusterId.ManagedClusterName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ManagedClusterName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("check for existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_kubernetes_cluster_maintenance_configuration", id.ID())
		}
	}

	parameters := containerservice.MaintenanceConfiguration{
		MaintenanceConfigurationProperties: &containerservice.MaintenanceConfigurationProperties{
			NotAllowedTime: expandKubernetesClusterMaintenanceConfigurationTimeSpans(d.Get("maintenance_not_allowed_window").(*pluginsdk.Set).List()),
			TimeInWeek:     expandKubernetesClusterMaintenanceConfigurationTimeInWeeks(d.Get("maintenance_allowed").(*pluginsdk.Set).List()),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedClusterName, id.Name, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKubernetesClusterMaintenanceConfigurationRead(d, meta)
}

func resourceKubernetesClusterMaintenanceConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
	clustersClient := meta.(*clients.Client).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceConfigurationID(d.Id())
	if err != nil {
		return err
	}

	// kubernetes cluster does not use id generator, the id format is not standard. So fetch cluster info here
	cluster, err := clustersClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			log.Printf("[DEBUG] Managed Kubernetes Cluster %q was not found in Resource Group %q - removing from state!", id.ManagedClusterName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ManagedClusterName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Kubernetes Cluster Maintenance Configuration %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("kubernetes_cluster_id", cluster.ID)

	if props := resp.MaintenanceConfigurationProperties; props != nil {
		if err := d.Set("maintenance_not_allowed_window", flattenKubernetesClusterMaintenanceConfigurationTimeSpans(props.NotAllowedTime)); err != nil {
			return fmt.Errorf("setting `maintenance_not_allowed_window`: %+v", err)
		}
		if err := d.Set("maintenance_allowed", flattenKubernetesClusterMaintenanceConfigurationTimeInWeeks(props.TimeInWeek)); err != nil {
			return fmt.Errorf("setting `maintenance_allowed`: %+v", err)
		}
	}
	return nil
}

func resourceKubernetesClusterMaintenanceConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.ManagedClusterName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandKubernetesClusterMaintenanceConfigurationTimeSpans(input []interface{}) *[]containerservice.TimeSpan {
	results := make([]containerservice.TimeSpan, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		start, _ := time.Parse(time.RFC3339, v["start"].(string))
		end, _ := time.Parse(time.RFC3339, v["end"].(string))
		results = append(results, containerservice.TimeSpan{
			Start: &date.Time{Time: start},
			End:   &date.Time{Time: end},
		})
	}
	return &results
}

func expandKubernetesClusterMaintenanceConfigurationTimeInWeeks(input []interface{}) *[]containerservice.TimeInWeek {
	results := make([]containerservice.TimeInWeek, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, containerservice.TimeInWeek{
			Day:       containerservice.WeekDay(v["day"].(string)),
			HourSlots: utils.ExpandInt32Slice(v["hour_slots"].(*pluginsdk.Set).List()),
		})
	}
	return &results
}

func flattenKubernetesClusterMaintenanceConfigurationTimeSpans(input *[]containerservice.TimeSpan) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var end string
		if item.End != nil {
			end = item.End.Format(time.RFC3339)
		}
		var start string
		if item.Start != nil {
			start = item.Start.Format(time.RFC3339)
		}
		results = append(results, map[string]interface{}{
			"end":   end,
			"start": start,
		})
	}
	return results
}

func flattenKubernetesClusterMaintenanceConfigurationTimeInWeeks(input *[]containerservice.TimeInWeek) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"day":        string(item.Day),
			"hour_slots": utils.FlattenInt32Slice(item.HourSlots),
		})
	}
	return results
}
