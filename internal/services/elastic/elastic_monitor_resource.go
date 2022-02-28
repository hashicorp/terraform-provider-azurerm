package elastic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/sdk/2020-07-01/monitorsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
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
			_, err := monitorsresource.ParseMonitorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ElasticMonitorName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceElasticMonitorCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Elastic.MonitorClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := monitorsresource.NewMonitorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.MonitorsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %q: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_elastic_monitor", id.ID())
	}

	monitoringStatus := monitorsresource.MonitoringStatusDisabled
	if d.Get("monitoring_status").(bool) {
		monitoringStatus = monitorsresource.MonitoringStatusEnabled
	}

	body := monitorsresource.ElasticMonitorResource{
		Location: location.Normalize(d.Get("location").(string)),
		Sku:      expandMonitorResourceSku(d.Get("sku").([]interface{})),
		Properties: &monitorsresource.MonitorProperties{
			UserInfo:         expandMonitorUserInfo(d.Get("user_info").([]interface{})),
			MonitoringStatus: &monitoringStatus,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.MonitorsCreateThenPoll(ctx, id, body); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceElasticMonitorRead(d, meta)
}

func resourceElasticMonitorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.MonitorClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := monitorsresource.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("elastic_cloud_user", flattenElasticCloudUser(props.ElasticProperties)); err != nil {
				return fmt.Errorf("setting `elastic_cloud_user`: %+v", err)
			}
			if err := d.Set("elastic_cloud_deployment", flattenElasticCloudDeployment(props.ElasticProperties)); err != nil {
				return fmt.Errorf("setting `elastic_cloud_deployment`: %+v", err)
			}
			d.Set("user_info", flattenUserInfo(props.ElasticProperties))
			monitoringEnabled := false
			if props.MonitoringStatus != nil {
				monitoringEnabled = *props.MonitoringStatus == monitorsresource.MonitoringStatusEnabled
			}
			d.Set("monitoring_status", monitoringEnabled)
		}

		if err := d.Set("sku", flattenMonitorResourceSku(model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceElasticMonitorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.MonitorClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := monitorsresource.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	body := monitorsresource.ElasticMonitorResourceUpdateParameters{}

	if d.HasChange("tags") {
		body.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.MonitorsUpdate(ctx, *id, body); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	return resourceElasticMonitorRead(d, meta)
}

func resourceElasticMonitorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.MonitorClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := monitorsresource.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	if err := client.MonitorsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandMonitorResourceSku(input []interface{}) *monitorsresource.ResourceSku {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &monitorsresource.ResourceSku{
		Name: v["name"].(string),
	}
}

func expandMonitorUserInfo(input []interface{}) *monitorsresource.UserInfo {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	return &monitorsresource.UserInfo{
		EmailAddress: utils.String(v["email_address"].(string)),
	}
}

func flattenUserInfo(input *monitorsresource.ElasticProperties) []interface{} {
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

func flattenElasticCloudUser(input *monitorsresource.ElasticProperties) []interface{} {
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
		if input.ElasticCloudUser.Id != nil {
			id = *input.ElasticCloudUser.Id
		}
		var elastic_cloud_sso_default_url string
		if input.ElasticCloudUser.ElasticCloudSsoDefaultUrl != nil {
			elastic_cloud_sso_default_url = *input.ElasticCloudUser.ElasticCloudSsoDefaultUrl
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

func flattenElasticCloudDeployment(input *monitorsresource.ElasticProperties) []interface{} {
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
		if input.ElasticCloudDeployment.DeploymentId != nil {
			deployment_id = *input.ElasticCloudDeployment.DeploymentId
		}
		var azure_subscription_id string
		if input.ElasticCloudDeployment.AzureSubscriptionId != nil {
			azure_subscription_id = *input.ElasticCloudDeployment.AzureSubscriptionId
		}
		var elasticsearch_region string
		if input.ElasticCloudDeployment.ElasticsearchRegion != nil {
			elasticsearch_region = *input.ElasticCloudDeployment.ElasticsearchRegion
		}
		var elasticsearch_service_url string
		if input.ElasticCloudDeployment.ElasticsearchServiceUrl != nil {
			elasticsearch_service_url = *input.ElasticCloudDeployment.ElasticsearchServiceUrl
		}
		var kibana_service_url string
		if input.ElasticCloudDeployment.KibanaServiceUrl != nil {
			kibana_service_url = *input.ElasticCloudDeployment.KibanaServiceUrl
		}
		var kibana_sso_url string
		if input.ElasticCloudDeployment.KibanaSsoUrl != nil {
			kibana_sso_url = *input.ElasticCloudDeployment.KibanaSsoUrl
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

func flattenMonitorResourceSku(input *monitorsresource.ResourceSku) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"name": input.Name,
		},
	}
}
