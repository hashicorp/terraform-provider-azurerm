package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualHubRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualHubRouteTableCreateUpdate,
		Read:   resourceArmVirtualHubRouteTableRead,
		Update: resourceArmVirtualHubRouteTableCreateUpdate,
		Delete: resourceArmVirtualHubRouteTableDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.HubRouteTableID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.VirtualHubRouteTableName,
			},

			"virtual_hub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.VirtualHubID,
			},

			"labels": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"route": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"destinations": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"destinations_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CIDR",
								"ResourceId",
								"Service",
							}, false),
						},

						"next_hop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"next_hop_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ResourceId",
							ValidateFunc: validation.StringInSlice([]string{
								"ResourceId",
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceArmVirtualHubRouteTableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubRouteTableClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.Name, virtualHubResourceName)
	defer locks.UnlockByName(id.Name, virtualHubResourceName)

	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", name, id.ResourceGroup, id.Name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_hub_route_table", *existing.ID)
		}
	}

	parameters := network.HubRouteTable{
		Name: utils.String(d.Get("name").(string)),
		HubRouteTableProperties: &network.HubRouteTableProperties{
			Labels: utils.ExpandStringSlice(d.Get("labels").(*schema.Set).List()),
			Routes: expandArmVirtualHubRouteTableHubRoutes(d.Get("route").(*schema.Set).List()),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", name, id.ResourceGroup, id.Name, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating future for HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", name, id.ResourceGroup, id.Name, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", name, id.ResourceGroup, id.Name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for HubRouteTable %q (Resource Group %q / Virtual Hub %q) ID", name, id.ResourceGroup, id.Name)
	}

	d.SetId(*resp.ID)

	return resourceArmVirtualHubRouteTableRead(d, meta)
}

func resourceArmVirtualHubRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubRouteTableClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HubRouteTableID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Virtual Hub Route Table %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", id.Name, id.ResourceGroup, id.VirtualHubName, err)
	}

	d.Set("name", id.Name)
	d.Set("virtual_hub_id", parse.NewVirtualHubID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName).ID(""))

	if props := resp.HubRouteTableProperties; props != nil {
		d.Set("labels", utils.FlattenStringSlice(props.Labels))

		if err := d.Set("route", flattenArmVirtualHubRouteTableHubRoutes(props.Routes)); err != nil {
			return fmt.Errorf("setting `route`: %+v", err)
		}
	}
	return nil
}

func resourceArmVirtualHubRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubRouteTableClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HubRouteTableID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", id.Name, id.ResourceGroup, id.VirtualHubName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", id.Name, id.ResourceGroup, id.VirtualHubName, err)
	}

	return nil
}

func expandArmVirtualHubRouteTableHubRoutes(input []interface{}) *[]network.HubRoute {
	results := make([]network.HubRoute, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := network.HubRoute{
			Name:            utils.String(v["name"].(string)),
			DestinationType: utils.String(v["destinations_type"].(string)),
			Destinations:    utils.ExpandStringSlice(v["destinations"].(*schema.Set).List()),
			NextHopType:     utils.String(v["next_hop_type"].(string)),
			NextHop:         utils.String(v["next_hop"].(string)),
		}

		results = append(results, result)
	}

	return &results
}

func flattenArmVirtualHubRouteTableHubRoutes(input *[]network.HubRoute) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		var destinationType string
		if item.DestinationType != nil {
			destinationType = *item.DestinationType
		}

		var nextHop string
		if item.NextHop != nil {
			nextHop = *item.NextHop
		}

		var nextHopType string
		if item.NextHopType != nil {
			nextHopType = *item.NextHopType
		}

		v := map[string]interface{}{
			"name":              name,
			"destinations":      utils.FlattenStringSlice(item.Destinations),
			"destinations_type": destinationType,
			"next_hop":          nextHop,
			"next_hop_type":     nextHopType,
		}

		results = append(results, v)
	}

	return results
}
