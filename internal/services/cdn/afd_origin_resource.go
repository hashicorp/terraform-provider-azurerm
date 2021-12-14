package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
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

	id := parse.NewAfdOriginsID(originGroup.SubscriptionId, originGroup.ResourceGroup, originGroup.ProfileName, originGroup.OriginGroupName, originname)

	afdOrigin := cdn.AFDOrigin{
		AFDOriginProperties: &cdn.AFDOriginProperties{
			HostName: &hostname,
			Priority: &priority,
			Weight:   &weight,
		},
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

	return nil
}

func resourceAfdOriginsUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.AfdOriginsID(d.Id())
	if err != nil {
		return err
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
