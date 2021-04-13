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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceExpressRouteCircuitConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceExpressRouteCircuitConnectionCreateUpdate,
		Read:   resourceExpressRouteCircuitConnectionRead,
		Update: resourceExpressRouteCircuitConnectionCreateUpdate,
		Delete: resourceExpressRouteCircuitConnectionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ExpressRouteCircuitConnectionID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"circuit_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"peering_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"peer_peering_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"address_prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
			},

			"authorization_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"ipv6circuit_connection_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_prefix": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsCIDR,
						},

						"circuit_connection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func resourceExpressRouteCircuitConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnectionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	circuitName := d.Get("circuit_name").(string)

	id := parse.NewExpressRouteCircuitConnectionID(subscriptionId, resourceGroup, circuitName, "AzurePrivatePeering", name).ID()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, circuitName, "AzurePrivatePeering", name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing ExpressRouteCircuitConnection %q (Resource Group %q / circuitName %q): %+v", name, resourceGroup, circuitName, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_express_route_circuit_connection", id)
		}
	}

	expressRouteCircuitConnectionParameters := network.ExpressRouteCircuitConnection{
		Name: utils.String(d.Get("name").(string)),
		ExpressRouteCircuitConnectionPropertiesFormat: &network.ExpressRouteCircuitConnectionPropertiesFormat{
			AddressPrefix: utils.String(d.Get("address_prefix").(string)),
			ExpressRouteCircuitPeering: &network.SubResource{
				ID: utils.String(d.Get("peering_id").(string)),
			},
			PeerExpressRouteCircuitPeering: &network.SubResource{
				ID: utils.String(d.Get("peer_peering_id").(string)),
			},
			Ipv6CircuitConnectionConfig: expandExpressRouteCircuitConnectionIpv6CircuitConnectionConfig(d.Get("ipv6circuit_connection_config").([]interface{})),
		},
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		expressRouteCircuitConnectionParameters.ExpressRouteCircuitConnectionPropertiesFormat.AuthorizationKey = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, circuitName, "AzurePrivatePeering", name, expressRouteCircuitConnectionParameters)
	if err != nil {
		return fmt.Errorf("creating/updating ExpressRouteCircuitConnection %q (Resource Group %q / circuitName %q): %+v", name, resourceGroup, circuitName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of the ExpressRouteCircuitConnection %q (Resource Group %q / circuitName %q): %+v", name, resourceGroup, circuitName, err)
	}

	d.SetId(id)
	return resourceExpressRouteCircuitConnectionRead(d, meta)
}

func resourceExpressRouteCircuitConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteCircuitConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.CircuitName, id.PeeringName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] express route circuit connection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving ExpressRouteCircuitConnection %q (Resource Group %q / circuitName %q / peeringName %q): %+v", id.Name, id.ResourceGroup, id.CircuitName, id.PeeringName, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("circuit_name", id.CircuitName)
	d.Set("peering_id", resp.ExpressRouteCircuitPeering.ID)
	d.Set("peer_peering_id", resp.PeerExpressRouteCircuitPeering.ID)
	if props := resp.ExpressRouteCircuitConnectionPropertiesFormat; props != nil {
		d.Set("address_prefix", props.AddressPrefix)
		d.Set("authorization_key", props.AuthorizationKey)
		if err := d.Set("ipv6circuit_connection_config", flattenExpressRouteCircuitConnectionIpv6CircuitConnectionConfig(props.Ipv6CircuitConnectionConfig)); err != nil {
			return fmt.Errorf("setting `ipv6circuit_connection_config`: %+v", err)
		}
	}
	return nil
}

func resourceExpressRouteCircuitConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteCircuitConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.CircuitName, id.PeeringName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting ExpressRouteCircuitConnection %q (Resource Group %q / circuitName %q / peeringName %q): %+v", id.Name, id.ResourceGroup, id.CircuitName, id.PeeringName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the ExpressRouteCircuitConnection %q (Resource Group %q / circuitName %q / peeringName %q): %+v", id.Name, id.ResourceGroup, id.CircuitName, id.PeeringName, err)
	}
	return nil
}

func expandExpressRouteCircuitConnectionIpv6CircuitConnectionConfig(input []interface{}) *network.Ipv6CircuitConnectionConfig {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &network.Ipv6CircuitConnectionConfig{
		AddressPrefix: utils.String(v["address_prefix"].(string)),
	}
}

func flattenExpressRouteCircuitConnectionIpv6CircuitConnectionConfig(input *network.Ipv6CircuitConnectionConfig) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var addressPrefix string
	if input.AddressPrefix != nil {
		addressPrefix = *input.AddressPrefix
	}
	var circuitConnectionStatus network.CircuitConnectionStatus
	if input.CircuitConnectionStatus != "" {
		circuitConnectionStatus = input.CircuitConnectionStatus
	}
	return []interface{}{
		map[string]interface{}{
			"address_prefix":            addressPrefix,
			"circuit_connection_status": circuitConnectionStatus,
		},
	}
}
