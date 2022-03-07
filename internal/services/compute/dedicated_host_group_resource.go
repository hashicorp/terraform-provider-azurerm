package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDedicatedHostGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDedicatedHostGroupCreate,
		Read:   resourceDedicatedHostGroupRead,
		Update: resourceDedicatedHostGroupUpdate,
		Delete: resourceDedicatedHostGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.HostGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: func() map[string]*pluginsdk.Schema {
			s := map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validate.DedicatedHostGroupName(),
				},

				"location": azure.SchemaLocation(),

				// There's a bug in the Azure API where this is returned in upper-case
				// BUG: https://github.com/Azure/azure-rest-api-specs/issues/8068
				"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

				"platform_fault_domain_count": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(1, 3),
				},

				"automatic_placement_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"tags": tags.Schema(),
			}

			if features.ThreePointOhBeta() {
				s["zone"] = commonschema.ZoneSingleOptionalForceNew()
			} else {
				s["zones"] = azure.SchemaSingleZone()
			}

			return s
		}(),
	}
}

func resourceDedicatedHostGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewHostGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_dedicated_host_group", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	platformFaultDomainCount := d.Get("platform_fault_domain_count").(int)
	t := d.Get("tags").(map[string]interface{})

	parameters := compute.DedicatedHostGroup{
		Location: utils.String(location),
		DedicatedHostGroupProperties: &compute.DedicatedHostGroupProperties{
			PlatformFaultDomainCount: utils.Int32(int32(platformFaultDomainCount)),
		},
		Tags: tags.Expand(t),
	}
	if features.ThreePointOhBeta() {
		if zone, ok := d.GetOk("zone"); ok {
			parameters.Zones = &[]string{
				zone.(string),
			}
		}
	} else {
		if zones, ok := d.GetOk("zones"); ok {
			parameters.Zones = utils.ExpandStringSlice(zones.([]interface{}))
		}
	}

	if v, ok := d.GetOk("automatic_placement_enabled"); ok {
		parameters.DedicatedHostGroupProperties.SupportAutomaticPlacement = utils.Bool(v.(bool))
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDedicatedHostGroupRead(d, meta)
}

func resourceDedicatedHostGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HostGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Dedicated Host Group %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Dedicated Host Group %q (: %+v", id.String(), err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))
	if features.ThreePointOhBeta() {
		zone := ""
		if resp.Zones != nil && len(*resp.Zones) > 0 {
			z := *resp.Zones
			zone = z[0]
		}
		d.Set("zone", zone)
	} else {
		d.Set("zones", utils.FlattenStringSlice(resp.Zones))
	}

	if props := resp.DedicatedHostGroupProperties; props != nil {
		platformFaultDomainCount := 0
		if props.PlatformFaultDomainCount != nil {
			platformFaultDomainCount = int(*props.PlatformFaultDomainCount)
		}
		d.Set("platform_fault_domain_count", platformFaultDomainCount)

		d.Set("automatic_placement_enabled", props.SupportAutomaticPlacement)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDedicatedHostGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := compute.DedicatedHostGroupUpdate{
		Tags: tags.Expand(t),
	}

	if _, err := client.Update(ctx, resourceGroupName, name, parameters); err != nil {
		return fmt.Errorf("updating Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	return resourceDedicatedHostGroupRead(d, meta)
}

func resourceDedicatedHostGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HostGroupID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Dedicated Host Group %q : %+v", id.String(), err)
	}

	return nil
}
