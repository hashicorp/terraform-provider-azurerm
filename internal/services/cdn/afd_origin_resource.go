package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAfdOrigin() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAfdOriginsCreate,
		Read:   resourceAfdOriginsRead,
		Update: resourceAfdOriginsUpdate,
		Delete: resourceAfdOriginsDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CdnEndpointV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.OriginGroupsID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"origin_group_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"profile_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"host_name": {
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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 1000),
			},
		},
	}
}

func resourceAfdOriginsCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	originname := d.Get("name").(string)
	originsgroup := d.Get("origins_group_name").(string)
	hostname := d.Get("hostname").(string)
	priority := int32(d.Get("priority").(int))
	weight := int32(d.Get("weight").(int))

	id := parse.NewEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))

	afdOrigin := cdn.AFDOrigin{
		AFDOriginProperties: &cdn.AFDOriginProperties{
			HostName: &hostname,
			Priority: &priority,
			Weight:   &weight,
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, originsgroup, originname, afdOrigin)
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
	return nil
}

func resourceAfdOriginsUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourceAfdOriginsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}
