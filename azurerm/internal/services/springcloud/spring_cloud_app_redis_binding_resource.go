package springcloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2020-11-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	redisValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redis/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const springCloudAppRedisBindingKeySSL = "useSsl"

func resourceSpringCloudAppRedisBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpringCloudAppRedisBindingCreateUpdate,
		Read:   resourceSpringCloudAppRedisBindingRead,
		Update: resourceSpringCloudAppRedisBindingCreateUpdate,
		Delete: resourceSpringCloudAppRedisBindingDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.SpringCloudAppBindingID(id)
			return err
		}, importSpringCloudAppBinding(springCloudAppRedisBindingType)),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudAppBindingName,
			},

			"spring_cloud_app_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudAppID,
			},

			"redis_cache_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: redisValidate.CacheID,
			},

			"redis_access_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"ssl_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceSpringCloudAppRedisBindingCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appId, err := parse.SpringCloudAppID(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewSpringCloudAppBindingID(appId.SubscriptionId, appId.ResourceGroup, appId.SpringName, appId.AppName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_app_redis_binding", id.ID())
		}
	}

	bindingResource := appplatform.BindingResource{
		Properties: &appplatform.BindingResourceProperties{
			BindingParameters: map[string]interface{}{
				springCloudAppRedisBindingKeySSL: d.Get("ssl_enabled").(bool),
			},
			Key:        utils.String(d.Get("redis_access_key").(string)),
			ResourceID: utils.String(d.Get("redis_cache_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName, bindingResource); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpringCloudAppRedisBindingRead(d, meta)
}

func resourceSpringCloudAppRedisBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppBindingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App Binding %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.Set("name", id.BindingName)
	d.Set("spring_cloud_app_id", parse.NewSpringCloudAppID(id.SubscriptionId, id.ResourceGroup, id.SpringName, id.AppName).ID())
	if props := resp.Properties; props != nil {
		d.Set("redis_cache_id", props.ResourceID)

		enableSSL := "false"
		if v, ok := props.BindingParameters[springCloudAppRedisBindingKeySSL]; ok {
			enableSSL = v.(string)
		}
		d.Set("ssl_enabled", strings.EqualFold(enableSSL, "true"))
	}
	return nil
}

func resourceSpringCloudAppRedisBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppBindingID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
