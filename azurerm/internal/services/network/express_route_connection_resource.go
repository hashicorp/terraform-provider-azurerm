package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceExpressRouteConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteConnectionCreate,
		Read:   resourceExpressRouteConnectionRead,
		Update: resourceExpressRouteConnectionUpdate,
		Delete: resourceExpressRouteConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ExpressRouteConnectionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteConnectionName,
			},

			"express_route_circuit_peering_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteCircuitPeeringID,
			},

			"express_route_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteGatewayID,
			},

			"authorization_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},

			"enable_internet_security": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"routing": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"associated_route_table_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.HubRouteTableID,
							AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table"},
						},

						"propagated_route_table": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"labels": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
										AtLeastOneOf: []string{"routing.0.propagated_route_table.0.labels", "routing.0.propagated_route_table.0.route_table_ids"},
									},

									"route_table_ids": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validate.HubRouteTableID,
										},
										AtLeastOneOf: []string{"routing.0.propagated_route_table.0.labels", "routing.0.propagated_route_table.0.route_table_ids"},
									},
								},
							},
							AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table"},
						},
					},
				},
			},

			"routing_weight": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 32000),
			},
		},
	}
}

func resourceExpressRouteConnectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteConnectionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	expressRouteGatewayId, err := parse.ExpressRouteGatewayID(d.Get("express_route_gateway_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewExpressRouteConnectionID(expressRouteGatewayId.SubscriptionId, expressRouteGatewayId.ResourceGroup, expressRouteGatewayId.Name, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteGatewayName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_express_route_connection", id.ID())
	}

	parameters := network.ExpressRouteConnection{
		Name: utils.String(id.Name),
		ExpressRouteConnectionProperties: &network.ExpressRouteConnectionProperties{
			ExpressRouteCircuitPeering: &network.ExpressRouteCircuitPeeringID{
				ID: utils.String(d.Get("express_route_circuit_peering_id").(string)),
			},
			EnableInternetSecurity: utils.Bool(d.Get("enable_internet_security").(bool)),
			RoutingConfiguration:   expandExpressRouteConnectionRouting(d.Get("routing").([]interface{})),
			RoutingWeight:          utils.Int32(int32(d.Get("routing_weight").(int))),
		},
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		parameters.ExpressRouteConnectionProperties.AuthorizationKey = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ExpressRouteGatewayName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteConnectionRead(d, meta)
}

func resourceExpressRouteConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteGatewayName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("express_route_gateway_id", parse.NewExpressRouteGatewayID(id.SubscriptionId, id.ResourceGroup, id.ExpressRouteGatewayName).ID())

	if props := resp.ExpressRouteConnectionProperties; props != nil {
		d.Set("routing_weight", props.RoutingWeight)
		d.Set("authorization_key", props.AuthorizationKey)
		d.Set("enable_internet_security", props.EnableInternetSecurity)

		circuitPeeringID := ""
		if v := props.ExpressRouteCircuitPeering; v != nil {
			circuitPeeringID = *v.ID
		}
		peeringId, err := parse.ExpressRouteCircuitPeeringID(circuitPeeringID)
		if err != nil {
			return err
		}
		d.Set("express_route_circuit_peering_id", peeringId.ID())

		routing, err := flattenExpressRouteConnectionRouting(props.RoutingConfiguration)
		if err != nil {
			return err
		}
		if err := d.Set("routing", routing); err != nil {
			return fmt.Errorf("setting `routing`: %+v", err)
		}
	}

	return nil
}

func resourceExpressRouteConnectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteConnectionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteConnectionID(d.Id())
	if err != nil {
		return err
	}

	parameters := network.ExpressRouteConnection{
		Name: utils.String(id.Name),
		ExpressRouteConnectionProperties: &network.ExpressRouteConnectionProperties{
			ExpressRouteCircuitPeering: &network.ExpressRouteCircuitPeeringID{
				ID: utils.String(d.Get("express_route_circuit_peering_id").(string)),
			},
			EnableInternetSecurity: utils.Bool(d.Get("enable_internet_security").(bool)),
			RoutingConfiguration:   expandExpressRouteConnectionRouting(d.Get("routing").([]interface{})),
			RoutingWeight:          utils.Int32(int32(d.Get("routing_weight").(int))),
		},
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		parameters.ExpressRouteConnectionProperties.AuthorizationKey = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ExpressRouteGatewayName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	return resourceExpressRouteConnectionRead(d, meta)
}

func resourceExpressRouteConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ExpressRouteGatewayName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandExpressRouteConnectionRouting(input []interface{}) *network.RoutingConfiguration {
	if len(input) == 0 || input[0] == nil {
		return &network.RoutingConfiguration{}
	}

	v := input[0].(map[string]interface{})
	result := network.RoutingConfiguration{}

	if associatedRouteTableId := v["associated_route_table_id"].(string); associatedRouteTableId != "" {
		result.AssociatedRouteTable = &network.SubResource{
			ID: utils.String(associatedRouteTableId),
		}
	}

	if propagatedRouteTable := v["propagated_route_table"].([]interface{}); len(propagatedRouteTable) != 0 {
		result.PropagatedRouteTables = expandExpressRouteConnectionPropagatedRouteTable(propagatedRouteTable)
	}

	return &result
}

func expandExpressRouteConnectionPropagatedRouteTable(input []interface{}) *network.PropagatedRouteTable {
	if len(input) == 0 || input[0] == nil {
		return &network.PropagatedRouteTable{}
	}

	v := input[0].(map[string]interface{})

	result := network.PropagatedRouteTable{}

	if labels := v["labels"].(*pluginsdk.Set).List(); len(labels) != 0 {
		result.Labels = utils.ExpandStringSlice(labels)
	}

	if routeTableIds := v["route_table_ids"].([]interface{}); len(routeTableIds) != 0 {
		result.Ids = expandIDsToSubResources(routeTableIds)
	}

	return &result
}

func flattenExpressRouteConnectionRouting(input *network.RoutingConfiguration) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	associatedRouteTableId := ""
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.ID != nil {
		associatedRouteTableId = *input.AssociatedRouteTable.ID
	}
	routeTableId, err := parse.HubRouteTableID(associatedRouteTableId)
	if err != nil {
		return nil, err
	}

	return []interface{}{
		map[string]interface{}{
			"associated_route_table_id": routeTableId.ID(),
			"propagated_route_table":    flattenExpressRouteConnectionPropagatedRouteTable(input.PropagatedRouteTables),
		},
	}, nil
}

func flattenExpressRouteConnectionPropagatedRouteTable(input *network.PropagatedRouteTable) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	labels := make([]interface{}, 0)
	if input.Labels != nil {
		labels = utils.FlattenStringSlice(input.Labels)
	}

	routeTableIds := make([]interface{}, 0)
	if input.Ids != nil {
		routeTableIds = flattenSubResourcesToIDs(input.Ids)
	}

	return []interface{}{
		map[string]interface{}{
			"labels":          labels,
			"route_table_ids": routeTableIds,
		},
	}
}
