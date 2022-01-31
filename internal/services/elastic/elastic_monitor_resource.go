package elastic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/legacysdk/elastic/mgmt/2020-07-01/elastic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceElasticMonitor() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceElasticMonitorCreate,
		Read:   resourceElasticMonitorRead,
		Update: resourceElasticMonitorUpdate,
		Delete: resourceElasticMonitorDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ElasticMonitorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ElasticMonitorName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"user_info": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"email_address": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.ElasticEmailAddress,
						},
					},
				},
			},

			"monitoring_status": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"elastic_cloud_user": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"email_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"elastic_cloud_sso_default_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"elastic_cloud_deployment": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"deployment_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"azure_subscription_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"elasticsearch_region": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"elasticsearch_service_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"kibana_service_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"kibana_sso_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceElasticMonitorCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Elastic.MonitorClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewElasticMonitorID(subscriptionId, resourceGroup, name)

	existing, err := client.Get(ctx, id.ResourceGroup, id.MonitorName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %q: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_elastic_monitor", id.ID())
	}

	monitoringStatus := elastic.MonitoringStatusDisabled
	if d.Get("monitoring_status").(bool) {
		monitoringStatus = elastic.MonitoringStatusEnabled
	}

	body := elastic.MonitorResource{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Sku:      expandMonitorResourceSku(d.Get("sku").([]interface{})),
		Properties: &elastic.MonitorProperties{
			UserInfo:         expandMonitorUserInfo(d.Get("user_info").([]interface{})),
			MonitoringStatus: monitoringStatus,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	future, err := client.Create(ctx, id.ResourceGroup, id.MonitorName, &body)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceElasticMonitorRead(d, meta)
}

func resourceElasticMonitorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.MonitorClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ElasticMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MonitorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Elastic monitor %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}
	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		if err := d.Set("elastic_cloud_user", flattenElasticCloudUser(props.ElasticProperties)); err != nil {
			return fmt.Errorf("setting `elastic_cloud_user`: %+v", err)
		}
		if err := d.Set("elastic_cloud_deployment", flattenElasticCloudDeployment(props.ElasticProperties)); err != nil {
			return fmt.Errorf("setting `elastic_cloud_deployment`: %+v", err)
		}
		d.Set("user_info", flattenUserInfo(props.ElasticProperties))
		d.Set("monitoring_status", props.MonitoringStatus == elastic.MonitoringStatusEnabled)
	}
	if err := d.Set("sku", flattenMonitorResourceSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceElasticMonitorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.MonitorClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ElasticMonitorID(d.Id())
	if err != nil {
		return err
	}

	body := elastic.MonitorResourceUpdateParameters{}

	if d.HasChange("tags") {
		body.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.MonitorName, &body); err != nil {
		return fmt.Errorf("updating Elastic Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}
	return resourceElasticMonitorRead(d, meta)
}

func resourceElasticMonitorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.MonitorClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ElasticMonitorID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.MonitorName)
	if err != nil {
		return fmt.Errorf("deleting Elastic Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the Elastic Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}
	return nil
}

func expandMonitorResourceSku(input []interface{}) *elastic.ResourceSku {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &elastic.ResourceSku{
		Name: utils.String(v["name"].(string)),
	}
}

func expandMonitorUserInfo(input []interface{}) *elastic.UserInfo {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	return &elastic.UserInfo{
		EmailAddress: utils.String(v["email_address"].(string)),
	}
}

func flattenUserInfo(input *elastic.Properties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var email_address string
	if input.ElasticCloudUser != nil {
		if input.ElasticCloudUser.EmailAddress != nil {
			email_address = *input.ElasticCloudUser.EmailAddress
		}
	}
	return []interface{}{
		map[string]interface{}{
			"email_address": email_address,
		},
	}
}

func flattenElasticCloudUser(input *elastic.Properties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var elastic_cloud_user []interface{}

	if input.ElasticCloudUser != nil {
		var email_address string
		if input.ElasticCloudUser.EmailAddress != nil {
			email_address = *input.ElasticCloudUser.EmailAddress
		}
		var id string
		if input.ElasticCloudUser.ID != nil {
			id = *input.ElasticCloudUser.ID
		}
		var elastic_cloud_sso_default_url string
		if input.ElasticCloudUser.ElasticCloudSsoDefaultURL != nil {
			elastic_cloud_sso_default_url = *input.ElasticCloudUser.ElasticCloudSsoDefaultURL
		}
		elastic_cloud_user = []interface{}{
			map[string]interface{}{
				"email_address":                 email_address,
				"id":                            id,
				"elastic_cloud_sso_default_url": elastic_cloud_sso_default_url,
			},
		}

	} else {
		elastic_cloud_user = make([]interface{}, 0)
	}

	return elastic_cloud_user
}

func flattenElasticCloudDeployment(input *elastic.Properties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var elastic_cloud_deployment []interface{}

	if input.ElasticCloudDeployment != nil {
		var name string
		if input.ElasticCloudDeployment.Name != nil {
			name = *input.ElasticCloudDeployment.Name
		}
		var deployment_id string
		if input.ElasticCloudDeployment.DeploymentID != nil {
			deployment_id = *input.ElasticCloudDeployment.DeploymentID
		}
		var azure_subscription_id string
		if input.ElasticCloudDeployment.AzureSubscriptionID != nil {
			azure_subscription_id = *input.ElasticCloudDeployment.AzureSubscriptionID
		}
		var elasticsearch_region string
		if input.ElasticCloudDeployment.ElasticsearchRegion != nil {
			elasticsearch_region = *input.ElasticCloudDeployment.ElasticsearchRegion
		}
		var elasticsearch_service_url string
		if input.ElasticCloudDeployment.ElasticsearchServiceURL != nil {
			elasticsearch_service_url = *input.ElasticCloudDeployment.ElasticsearchServiceURL
		}
		var kibana_service_url string
		if input.ElasticCloudDeployment.KibanaServiceURL != nil {
			kibana_service_url = *input.ElasticCloudDeployment.KibanaServiceURL
		}
		var kibana_sso_url string
		if input.ElasticCloudDeployment.KibanaSsoURL != nil {
			kibana_sso_url = *input.ElasticCloudDeployment.KibanaSsoURL
		}
		elastic_cloud_deployment = []interface{}{
			map[string]interface{}{
				"name":                      name,
				"deployment_id":             deployment_id,
				"azure_subscription_id":     azure_subscription_id,
				"elasticsearch_region":      elasticsearch_region,
				"elasticsearch_service_url": elasticsearch_service_url,
				"kibana_service_url":        kibana_service_url,
				"kibana_sso_url":            kibana_sso_url,
			},
		}

	} else {
		elastic_cloud_deployment = make([]interface{}, 0)
	}

	return elastic_cloud_deployment
}

func flattenMonitorResourceSku(input *elastic.ResourceSku) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var name string
	if input.Name != nil {
		name = *input.Name
	}
	return []interface{}{
		map[string]interface{}{
			"name": name,
		},
	}
}
