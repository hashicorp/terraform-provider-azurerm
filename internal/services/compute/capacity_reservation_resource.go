package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCapacityReservation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCapacityReservationCreate,
		Read:   resourceCapacityReservationRead,
		Update: resourceCapacityReservationUpdate,
		Delete: resourceCapacityReservationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CapacityReservationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CapacityReservationName(),
			},

			"capacity_reservation_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CapacityReservationGroupID,
			},

			"zone": commonschema.ZoneSingleOptionalForceNew(),

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"capacity": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceCapacityReservationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	capacityReservationGroupId, err := parse.CapacityReservationGroupID(d.Get("capacity_reservation_group_id").(string))
	if err != nil {
		return err
	}

	capacityReservationGroupClient := meta.(*clients.Client).Compute.CapacityReservationGroupsClient
	capacityReservationGroup, err := capacityReservationGroupClient.Get(ctx, capacityReservationGroupId.ResourceGroup, capacityReservationGroupId.Name, "")
	if err != nil {
		return err
	}

	id := parse.NewCapacityReservationID(subscriptionId, capacityReservationGroupId.ResourceGroup, capacityReservationGroupId.Name, d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.CapacityReservationGroupName, id.Name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_capacity_reservation", id.ID())
	}

	parameters := compute.CapacityReservation{
		Location: capacityReservationGroup.Location,
		Sku:      expandCapacityReservationSku(d.Get("sku").([]interface{})),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("zone"); ok {
		parameters.Zones = &[]string{
			v.(string),
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.CapacityReservationGroupName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCapacityReservationRead(d, meta)
}

func resourceCapacityReservationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CapacityReservationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.CapacityReservationGroupName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	capacityReservationReservationGroupId := parse.NewCapacityReservationGroupID(id.SubscriptionId, id.ResourceGroup, id.CapacityReservationGroupName)
	d.Set("name", resp.Name)
	d.Set("capacity_reservation_group_id", capacityReservationReservationGroupId.ID())
	if err := d.Set("sku", flattenCapacityReservationSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
	}

	zone := ""
	if resp.Zones != nil && len(*resp.Zones) > 0 {
		z := *resp.Zones
		zone = z[0]
	}
	d.Set("zone", zone)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceCapacityReservationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CapacityReservationID(d.Id())
	if err != nil {
		return err
	}

	parameters := compute.CapacityReservationUpdate{}
	if d.HasChange("sku") {
		parameters.Sku = expandCapacityReservationSku(d.Get("sku").([]interface{}))
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.CapacityReservationGroupName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}
	return resourceCapacityReservationRead(d, meta)
}

func resourceCapacityReservationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CapacityReservationID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.CapacityReservationGroupName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandCapacityReservationSku(input []interface{}) *compute.Sku {
	if len(input) == 0 {
		return &compute.Sku{}
	}

	v := input[0].(map[string]interface{})
	return &compute.Sku{
		Name:     utils.String(v["name"].(string)),
		Capacity: utils.Int64(int64(v["capacity"].(int))),
	}
}

func flattenCapacityReservationSku(input *compute.Sku) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var name string
	if input.Name != nil {
		name = *input.Name
	}

	var capacity int64
	if input.Capacity != nil {
		capacity = *input.Capacity
	}

	return []interface{}{
		map[string]interface{}{
			"name":     name,
			"capacity": capacity,
		},
	}
}
