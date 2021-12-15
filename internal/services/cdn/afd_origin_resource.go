package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAfdOrigin() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAfdOriginsCreate,
		Read:   resourceAfdOriginsRead,
		Update: resourceAfdOriginsUpdate,
		Delete: resourceAfdOriginsDelete,

		SchemaVersion: 1,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AfdOriginGroupsID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AfdOriginGroupsID,
			},

			"host_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"origin_host_header": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"priority": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 5),
			},

			"weight": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"http_port": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      80,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"https_port": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      443,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
		},
	}
}

func resourceAfdOriginsCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// parse origin_group_id
	originGroupId := d.Get("origin_group_id").(string)
	originGroup, err := parse.AfdOriginGroupsID(originGroupId)
	if err != nil {
		return err
	}

	originname := d.Get("name").(string)
	hostname := d.Get("host_name").(string)
	priority := int32(d.Get("priority").(int))
	weight := int32(d.Get("weight").(int))
	httpPort := int32(d.Get("http_port").(int))
	httpsPort := int32(d.Get("https_port").(int))
	originHostHeader := d.Get("origin_host_header").(string)

	id := parse.NewAfdOriginsID(originGroup.SubscriptionId, originGroup.ResourceGroup, originGroup.ProfileName, originGroup.OriginGroupName, originname)

	afdOrigin := cdn.AFDOrigin{
		AFDOriginProperties: &cdn.AFDOriginProperties{
			HostName:  &hostname,
			Priority:  &priority,
			Weight:    &weight,
			HTTPPort:  &httpPort,
			HTTPSPort: &httpsPort,
		},
	}

	if originHostHeader != "" {
		afdOrigin.OriginHostHeader = &originHostHeader
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, originGroup.OriginGroupName, originname, afdOrigin)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAfdOriginsRead(d, meta)
}

func resourceAfdOriginsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdOriginsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.OriginName)
	d.Set("http_port", resp.HTTPPort)
	d.Set("https_port", resp.HTTPSPort)
	d.Set("priority", resp.Priority)
	d.Set("weight", resp.Weight)
	d.Set("host_name", resp.HostName)
	d.Set("origin_host_header", resp.OriginHostHeader)

	return nil
}

func resourceAfdOriginsUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdOriginsID(d.Id())
	if err != nil {
		return err
	}

	var originUpdateProperties cdn.AFDOriginUpdateParameters

	if d.HasChange("origin_host_header") {
		originHostHeader := d.Get("origin_host_header").(string)
		originUpdateProperties.OriginHostHeader = &originHostHeader
	}

	if d.HasChange("http_port") {
		httpPort := d.Get("http_port").(string)
		originUpdateProperties.OriginHostHeader = &httpPort
	}

	if d.HasChange("https_port") {
		httpsPort := d.Get("https_port").(string)
		originUpdateProperties.OriginHostHeader = &httpsPort
	}

	if d.HasChange("priority") {
		priority := d.Get("priority").(string)
		originUpdateProperties.OriginHostHeader = &priority
	}

	if d.HasChange("weight") {
		weight := d.Get("weight").(string)
		originUpdateProperties.OriginHostHeader = &weight
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName, originUpdateProperties)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceAfdOriginsRead(d, meta)
}

func resourceAfdOriginsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdOriginsID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}
