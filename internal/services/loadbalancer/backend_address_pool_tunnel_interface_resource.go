package loadbalancer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.Resource = BackendAddressPoolTunnelInterfaceResource{}
var _ sdk.ResourceWithUpdate = BackendAddressPoolTunnelInterfaceResource{}

type BackendAddressPoolTunnelInterfaceResource struct{}

type BackendAddressPoolTunnelInterfaceModel struct {
	Identifier           int    `tfschema:"identifier"`
	BackendAddressPoolId string `tfschema:"backend_address_pool_id"`
	Type                 string `tfschema:"type"`
	Protocol             string `tfschema:"protocol"`
	Port                 int    `tfschema:"port"`
}

func (r BackendAddressPoolTunnelInterfaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"identifier": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ForceNew: true,
		},

		"backend_address_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LoadBalancerBackendAddressPoolID,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(network.GatewayLoadBalancerTunnelInterfaceTypeNone),
				string(network.GatewayLoadBalancerTunnelInterfaceTypeInternal),
				string(network.GatewayLoadBalancerTunnelInterfaceTypeExternal),
			},
				false,
			),
		},

		"protocol": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(network.GatewayLoadBalancerTunnelProtocolNone),
				string(network.GatewayLoadBalancerTunnelProtocolNative),
				string(network.GatewayLoadBalancerTunnelProtocolVXLAN),
			},
				false,
			),
		},

		"port": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
	}
}

func (r BackendAddressPoolTunnelInterfaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r BackendAddressPoolTunnelInterfaceResource) ModelObject() interface{} {
	return &BackendAddressPoolTunnelInterfaceModel{}
}

func (r BackendAddressPoolTunnelInterfaceResource) ResourceType() string {
	return "azurerm_lb_backend_address_pool_tunnel_interface"
}

func (r BackendAddressPoolTunnelInterfaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model BackendAddressPoolTunnelInterfaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			poolId, err := parse.LoadBalancerBackendAddressPoolID(model.BackendAddressPoolId)
			if err != nil {
				return err
			}

			locks.ByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)

			// Tunnel Interface can only be created for Gateway sku, so we have to check
			lb, err := metadata.Client.LoadBalancers.LoadBalancersClient.Get(ctx, poolId.ResourceGroup, poolId.LoadBalancerName, "")
			if err != nil {
				return fmt.Errorf("retrieving Load Balancer %q (Resource Group %q): %+v", poolId.LoadBalancerName, poolId.ResourceGroup, err)
			}
			if lb.Sku == nil || lb.Sku.Name != network.LoadBalancerSkuNameGateway {
				return fmt.Errorf("Tunnel Interface is only supported on Gateway SKU Load Balancers")
			}

			id := parse.NewBackendAddressPoolTunnelInterfaceID(subscriptionId, poolId.ResourceGroup, poolId.LoadBalancerName, poolId.BackendAddressPoolName, strconv.Itoa(model.Identifier))
			pool, err := client.Get(ctx, poolId.ResourceGroup, poolId.LoadBalancerName, poolId.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *poolId, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *poolId)
			}

			tunnelInterfaces := make([]network.GatewayLoadBalancerTunnelInterface, 0)
			if pool.BackendAddressPoolPropertiesFormat.TunnelInterfaces != nil {
				tunnelInterfaces = *pool.BackendAddressPoolPropertiesFormat.TunnelInterfaces
			}

			metadata.Logger.Infof("checking for existing %s..", id)
			for _, tunnelInterface := range tunnelInterfaces {
				if tunnelInterface.Identifier == nil {
					continue
				}

				if strconv.Itoa(int(*tunnelInterface.Identifier)) == id.TunnelInterfaceName {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			tunnelInterfaces = append(tunnelInterfaces, network.GatewayLoadBalancerTunnelInterface{
				Identifier: utils.Int32(int32(model.Identifier)),
				Type:       network.GatewayLoadBalancerTunnelInterfaceType(model.Type),
				Protocol:   network.GatewayLoadBalancerTunnelProtocol(model.Protocol),
				Port:       utils.Int32(int32(model.Port)),
			})
			pool.BackendAddressPoolPropertiesFormat.TunnelInterfaces = &tunnelInterfaces

			metadata.Logger.Infof("adding %s..", id)
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			metadata.Logger.Infof("waiting for update %s..", id)
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
			}
			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r BackendAddressPoolTunnelInterfaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			id, err := parse.BackendAddressPoolTunnelInterfaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			var tunnelInterface *network.GatewayLoadBalancerTunnelInterface
			if pool.BackendAddressPoolPropertiesFormat.TunnelInterfaces != nil {
				for _, itf := range *pool.BackendAddressPoolPropertiesFormat.TunnelInterfaces {
					if itf.Identifier == nil {
						continue
					}

					if strconv.Itoa(int(*itf.Identifier)) == id.TunnelInterfaceName {
						tunnelInterface = &itf
						break
					}
				}
			}
			if tunnelInterface == nil {
				return metadata.MarkAsGone(id)
			}

			identifier, err := strconv.Atoi(id.TunnelInterfaceName)
			if err != nil {
				return fmt.Errorf("converting the Identifier of Tunnel Interface %q from string to int: %+v", id, err)
			}

			backendAddressPoolId := parse.NewLoadBalancerBackendAddressPoolID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)

			var port int
			if tunnelInterface.Port != nil {
				port = int(*tunnelInterface.Port)
			}

			model := BackendAddressPoolTunnelInterfaceModel{
				Identifier:           identifier,
				BackendAddressPoolId: backendAddressPoolId.ID(),
				Type:                 string(tunnelInterface.Type),
				Protocol:             string(tunnelInterface.Protocol),
				Port:                 port,
			}

			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r BackendAddressPoolTunnelInterfaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			id, err := parse.BackendAddressPoolTunnelInterfaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

			pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			tunnelInterfaces := make([]network.GatewayLoadBalancerTunnelInterface, 0)
			if pool.BackendAddressPoolPropertiesFormat.TunnelInterfaces != nil {
				tunnelInterfaces = *pool.BackendAddressPoolPropertiesFormat.TunnelInterfaces
			}

			newTunnelInterfaces := make([]network.GatewayLoadBalancerTunnelInterface, 0)
			for _, itf := range tunnelInterfaces {
				if itf.Identifier == nil {
					continue
				}

				if strconv.Itoa(int(*itf.Identifier)) != id.TunnelInterfaceName {
					newTunnelInterfaces = append(newTunnelInterfaces, itf)
				}
			}
			pool.BackendAddressPoolPropertiesFormat.TunnelInterfaces = &newTunnelInterfaces

			metadata.Logger.Infof("removing %s..", *id)
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
			if err != nil {
				return fmt.Errorf("removing %s: %+v", *id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for removal of %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r BackendAddressPoolTunnelInterfaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.BackendAddressPoolTunnelInterfaceID
}

func (r BackendAddressPoolTunnelInterfaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			id, err := parse.BackendAddressPoolTunnelInterfaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

			var model BackendAddressPoolTunnelInterfaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			tunnelInterfaces := make([]network.GatewayLoadBalancerTunnelInterface, 0)
			if pool.BackendAddressPoolPropertiesFormat.TunnelInterfaces != nil {
				tunnelInterfaces = *pool.BackendAddressPoolPropertiesFormat.TunnelInterfaces
			}
			index := -1
			for i, itf := range tunnelInterfaces {
				if itf.Identifier == nil {
					continue
				}

				if strconv.Itoa(int(*itf.Identifier)) == id.TunnelInterfaceName {
					index = i
					break
				}
			}
			if index == -1 {
				return fmt.Errorf("%s was not found", *id)
			}

			tunnelInterfaces[index] = network.GatewayLoadBalancerTunnelInterface{
				Identifier: utils.Int32(int32(model.Identifier)),
				Type:       network.GatewayLoadBalancerTunnelInterfaceType(model.Type),
				Protocol:   network.GatewayLoadBalancerTunnelProtocol(model.Protocol),
				Port:       utils.Int32(int32(model.Port)),
			}
			pool.BackendAddressPoolPropertiesFormat.TunnelInterfaces = &tunnelInterfaces

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
