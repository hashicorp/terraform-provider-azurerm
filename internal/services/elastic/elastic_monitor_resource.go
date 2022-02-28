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

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"elastic_cloud_email_address": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ElasticEmailAddress,
			},

			"monitoring_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"elastic_cloud_deployment_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"elastic_cloud_sso_default_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"elastic_cloud_user_id": {
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
			"kibana_sso_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
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
	if d.Get("monitoring_enabled").(bool) {
		monitoringStatus = monitorsresource.MonitoringStatusEnabled
	}

	body := monitorsresource.ElasticMonitorResource{
		Location: location.Normalize(d.Get("location").(string)),
		Sku: &monitorsresource.ResourceSku{
			Name: d.Get("sku_name").(string),
		},
		Properties: &monitorsresource.MonitorProperties{
			UserInfo: &monitorsresource.UserInfo{
				EmailAddress: utils.String(d.Get("elastic_cloud_email_address").(string)),
			},
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
			monitoringEnabled := false
			if props.MonitoringStatus != nil {
				monitoringEnabled = *props.MonitoringStatus == monitorsresource.MonitoringStatusEnabled
			}
			d.Set("monitoring_enabled", monitoringEnabled)

			if elastic := props.ElasticProperties; elastic != nil {
				if elastic.ElasticCloudDeployment != nil {
					// AzureSubscriptionId is the same as the subscription deployed into, so no point exposing it
					// ElasticsearchRegion is `{Cloud}-{Region}` - so the same as location/not worth exposing for now?
					d.Set("elastic_cloud_deployment_id", elastic.ElasticCloudDeployment.DeploymentId)
					d.Set("elasticsearch_service_url", elastic.ElasticCloudDeployment.ElasticsearchServiceUrl)
					d.Set("kibana_service_url", elastic.ElasticCloudDeployment.KibanaServiceUrl)
					d.Set("kibana_sso_uri", elastic.ElasticCloudDeployment.KibanaSsoUrl)
				}
				if elastic.ElasticCloudUser != nil {
					d.Set("elastic_cloud_user_id", elastic.ElasticCloudUser.Id)
					d.Set("elastic_cloud_email_address", elastic.ElasticCloudUser.EmailAddress)
					d.Set("elastic_cloud_sso_default_url", elastic.ElasticCloudUser.ElasticCloudSsoDefaultUrl)
				}
			}
		}

		skuName := ""
		if model.Sku != nil {
			skuName = model.Sku.Name
		}
		d.Set("sku_name", skuName)

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
