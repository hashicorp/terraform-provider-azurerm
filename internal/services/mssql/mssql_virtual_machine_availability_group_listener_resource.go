package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	lbParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	networkParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlVirtualMachineAvailabilityGroupListenerResource struct{}

type MsSqlVirtualMachineAvailabilityGroupListenerModel struct {
	Name                       string `tfschema:"name"`
	ResourceGroup              string `tfschema:"resource_group_name"`
	SqlVirtualMachineGroupName string `tfschema:"sql_virtual_machine_group_name"`
	AvailabilityGroupName      string `tfschema:"availability_group_name"`

	CreateDefaultAvailabilityGroup bool                                                                           `tfschema:"create_default_availability_group"`
	Port                           int                                                                            `tfschema:"port"`
	LoadBalancerConfiguration      []helper.LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener `tfschema:"load_balancer_configuration"`
	Replica                        []helper.ReplicaMsSqlVirtualMachineAvailabilityGroupListener                   `tfschema:"replica"`
}

var _ sdk.Resource = MsSqlVirtualMachineAvailabilityGroupListenerResource{}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) ModelObject() interface{} {
	return &MsSqlVirtualMachineAvailabilityGroupListenerModel{}
}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) ResourceType() string {
	return "azurerm_mssql_virtual_machine_availability_group_listener"
}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return availabilitygrouplisteners.ValidateAvailabilityGroupListenerID
}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"sql_virtual_machine_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 15),
		},

		"availability_group_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"create_default_availability_group": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},

		"port": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.PortNumber,
		},

		"load_balancer_configuration": helper.LoadBalancerConfigurationSchemaMsSqlVirtualMachineAvailabilityGroupListener(),

		"replica": helper.ReplicaSchemaMsSqlVirtualMachineAvailabilityGroupListener(),
	}
}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MsSqlVirtualMachineAvailabilityGroupListenerModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.MSSQL.VirtualMachinesAvailabilityGroupListenersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := availabilitygrouplisteners.NewAvailabilityGroupListenerID(subscriptionId, model.ResourceGroup, model.SqlVirtualMachineGroupName, model.Name)

			existing, err := client.Get(ctx, id, availabilitygrouplisteners.GetOperationOptions{Expand: utils.String("AvailabilityGroupConfiguration")})
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			lbConfigs, err := expandMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(model.LoadBalancerConfiguration)
			if err != nil {
				return err
			}

			replicas, err := expandMsSqlVirtualMachineAvailabilityGroupListenerReplicas(model.Replica)
			if err != nil {
				return err
			}

			parameters := availabilitygrouplisteners.AvailabilityGroupListener{
				Properties: &availabilitygrouplisteners.AvailabilityGroupListenerProperties{
					AvailabilityGroupName:                    utils.String(model.AvailabilityGroupName),
					LoadBalancerConfigurations:               lbConfigs,
					CreateDefaultAvailabilityGroupIfNotExist: utils.Bool(model.CreateDefaultAvailabilityGroup),
					Port:                                     utils.Int64(int64(model.Port)),
					AvailabilityGroupConfiguration: &availabilitygrouplisteners.AgConfiguration{
						Replicas: replicas,
					},
				},
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.MSSQL.VirtualMachinesAvailabilityGroupListenersClient

			id, err := availabilitygrouplisteners.ParseAvailabilityGroupListenerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, availabilitygrouplisteners.GetOperationOptions{Expand: utils.String("AvailabilityGroupConfiguration")})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := MsSqlVirtualMachineAvailabilityGroupListenerModel{
				Name:                       id.AvailabilityGroupListenerName,
				ResourceGroup:              id.ResourceGroupName,
				SqlVirtualMachineGroupName: id.SqlVirtualMachineGroupName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {

					avGroupName := ""
					if props.AvailabilityGroupName != nil {
						avGroupName = *props.AvailabilityGroupName
					}
					state.AvailabilityGroupName = avGroupName

					createDefaultAvailabilityGroup := true
					if props.CreateDefaultAvailabilityGroupIfNotExist != nil {
						createDefaultAvailabilityGroup = *props.CreateDefaultAvailabilityGroupIfNotExist
					}
					state.CreateDefaultAvailabilityGroup = createDefaultAvailabilityGroup

					var port int64
					if props.Port != nil {
						port = *props.Port
					}
					state.Port = int(port)

					avGroupListenerLbConfigs, err := flattenMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(props.LoadBalancerConfigurations, id.SubscriptionId)
					if err != nil {
						return fmt.Errorf("setting `load_balancer_configuration`: %+v", err)
					}
					state.LoadBalancerConfiguration = avGroupListenerLbConfigs

					if props.AvailabilityGroupConfiguration != nil {
						if props.AvailabilityGroupConfiguration.Replicas != nil {

							replicas, err := flattenMsSqlVirtualMachineAvailabilityGroupListenerReplicas(props.AvailabilityGroupConfiguration.Replicas)
							if err != nil {
								return fmt.Errorf("setting `replica`: %+v", err)
							}
							state.Replica = replicas
						}
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.VirtualMachinesAvailabilityGroupListenersClient

			id, err := availabilitygrouplisteners.ParseAvailabilityGroupListenerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(lbConfigs []helper.LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener) (*[]availabilitygrouplisteners.LoadBalancerConfiguration, error) {
	results := make([]availabilitygrouplisteners.LoadBalancerConfiguration, 0)

	for _, lb := range lbConfigs {

		lbConfig := availabilitygrouplisteners.LoadBalancerConfiguration{
			ProbePort: utils.Int64(int64(lb.ProbePort)),
		}

		parsedLbId := ""
		if lb.LoadBalancerId != "" {
			id, err := lbParse.LoadBalancerID(lb.LoadBalancerId)
			if err != nil {
				return nil, err
			}
			parsedLbId = id.ID()
		}
		lbConfig.LoadBalancerResourceId = utils.String(parsedLbId)

		var parsedIds []interface{}
		for _, sqlVmId := range lb.SqlVirtualMachineInstances {
			parsedId, err := parse.SqlVirtualMachineID(sqlVmId)
			if err != nil {
				return nil, err
			}
			parsedIds = append(parsedIds, parsedId.ID())
		}
		lbConfig.SqlVirtualMachineInstances = utils.ExpandStringSlice(parsedIds)

		if lb.PrivateIpAddress != nil {
			privateIpAddress, err := expandMsSqlVirtualMachinePrivateIpAddress(lb.PrivateIpAddress)
			if err != nil {
				return nil, err
			}
			lbConfig.PrivateIPAddress = privateIpAddress
		}

		if publicIp := lb.PublicIpAddressId; publicIp != "" {
			id, err := networkParse.PublicIpAddressID(publicIp)
			if err != nil {
				return nil, err
			}
			lbConfig.PublicIPAddressResourceId = utils.String(id.ID())
		}
		results = append(results, lbConfig)
	}
	return &results, nil
}

func expandMsSqlVirtualMachinePrivateIpAddress(input []helper.PrivateIpAddressMsSqlVirtualMachineAvailabilityGroupListener) (*availabilitygrouplisteners.PrivateIPAddress, error) {
	if len(input) == 0 {
		return nil, nil
	}

	var ipAddress, subnetId string
	if input[0].IpAddress != "" {
		ipAddress = input[0].IpAddress
	}

	if input[0].SubnetId != "" {
		id, err := networkParse.SubnetID(input[0].SubnetId)
		if err != nil {
			return nil, err
		}
		subnetId = id.ID()
	}

	return &availabilitygrouplisteners.PrivateIPAddress{
		IPAddress: utils.String(ipAddress),

		SubnetResourceId: utils.String(subnetId),
	}, nil
}

func flattenMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(input *[]availabilitygrouplisteners.LoadBalancerConfiguration, subscriptionId string) ([]helper.LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener, error) {
	results := make([]helper.LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, lbConfig := range *input {

		var privateIpAddress []helper.PrivateIpAddressMsSqlVirtualMachineAvailabilityGroupListener
		if lbConfig.PrivateIPAddress != nil {
			flattenedPrivateIp, err := flattenPrivateIpAddress(*lbConfig.PrivateIPAddress)
			if err != nil {
				return nil, err
			}
			privateIpAddress = flattenedPrivateIp
		}

		publicIpAddressId := ""
		if lbConfig.PublicIPAddressResourceId != nil {
			publicIpAddressId = *lbConfig.PublicIPAddressResourceId
		}

		loadBalancerId := ""
		if lbConfig.LoadBalancerResourceId != nil {
			id, err := lbParse.LoadBalancerID(*lbConfig.LoadBalancerResourceId)
			if err != nil {
				return nil, err
			}
			loadBalancerId = id.ID()
		}

		var probePort int64
		if lbConfig.ProbePort != nil {
			probePort = *lbConfig.ProbePort
		}

		var sqlVirtualMachineIds []string
		if lbConfig.SqlVirtualMachineInstances != nil {
			sqlVirtualMachineIds = *lbConfig.SqlVirtualMachineInstances
			var parsedIds []string
			for _, sqlVmId := range sqlVirtualMachineIds {
				parsedId, err := sqlvirtualmachines.ParseSqlVirtualMachineID(sqlVmId)
				// get correct casing for subscription in id
				newId := sqlvirtualmachines.NewSqlVirtualMachineID(subscriptionId, parsedId.ResourceGroupName, parsedId.SqlVirtualMachineName)
				if err != nil {
					return nil, err
				}
				parsedIds = append(parsedIds, newId.ID())
			}
			sqlVirtualMachineIds = parsedIds
		}

		v := helper.LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener{
			PrivateIpAddress:           privateIpAddress,
			PublicIpAddressId:          publicIpAddressId,
			LoadBalancerId:             loadBalancerId,
			ProbePort:                  int(probePort),
			SqlVirtualMachineInstances: sqlVirtualMachineIds,
		}

		results = append(results, v)
	}
	return results, nil
}

func flattenPrivateIpAddress(input availabilitygrouplisteners.PrivateIPAddress) ([]helper.PrivateIpAddressMsSqlVirtualMachineAvailabilityGroupListener, error) {

	privateIpAddress := helper.PrivateIpAddressMsSqlVirtualMachineAvailabilityGroupListener{}

	ipAddress := ""
	if input.IPAddress != nil {
		ipAddress = *input.IPAddress
	}
	privateIpAddress.IpAddress = ipAddress

	subnetId := ""
	if input.SubnetResourceId != nil {
		id, err := networkParse.SubnetID(*input.SubnetResourceId)
		if err != nil {
			return nil, err
		}
		subnetId = id.ID()
	}

	privateIpAddress.SubnetId = subnetId

	return []helper.PrivateIpAddressMsSqlVirtualMachineAvailabilityGroupListener{privateIpAddress}, nil
}

func expandMsSqlVirtualMachineAvailabilityGroupListenerReplicas(replicas []helper.ReplicaMsSqlVirtualMachineAvailabilityGroupListener) (*[]availabilitygrouplisteners.AgReplica, error) {
	results := make([]availabilitygrouplisteners.AgReplica, 0)

	for _, rep := range replicas {

		role := availabilitygrouplisteners.Role(rep.Role)
		commit := availabilitygrouplisteners.Commit(rep.Commit)
		failover := availabilitygrouplisteners.Failover(rep.Failover)
		readableSecondary := availabilitygrouplisteners.ReadableSecondary(rep.ReadableSecondary)

		replica := availabilitygrouplisteners.AgReplica{
			Role:              &role,
			Commit:            &commit,
			Failover:          &failover,
			ReadableSecondary: &readableSecondary,
		}

		sqlVirtualMachineId := rep.SqlVirtualMachineInstanceId
		if sqlVirtualMachineId != "" {
			id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(sqlVirtualMachineId)
			if err != nil {
				return nil, err
			}
			sqlVirtualMachineId = id.ID()
		}
		replica.SqlVirtualMachineInstanceId = utils.String(sqlVirtualMachineId)

		results = append(results, replica)
	}
	return &results, nil
}

func flattenMsSqlVirtualMachineAvailabilityGroupListenerReplicas(input *[]availabilitygrouplisteners.AgReplica) ([]helper.ReplicaMsSqlVirtualMachineAvailabilityGroupListener, error) {
	results := make([]helper.ReplicaMsSqlVirtualMachineAvailabilityGroupListener, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, replica := range *input {

		sqlVirtualMachineInstanceId := ""
		if replica.SqlVirtualMachineInstanceId != nil {
			parsedSqlVmId, err := sqlvirtualmachines.ParseSqlVirtualMachineID(*replica.SqlVirtualMachineInstanceId)
			if err != nil {
				return nil, err
			}

			sqlVirtualMachineInstanceId = parsedSqlVmId.ID()
		}

		v := helper.ReplicaMsSqlVirtualMachineAvailabilityGroupListener{
			SqlVirtualMachineInstanceId: sqlVirtualMachineInstanceId,
		}

		if replica.Role != nil {
			v.Role = string(*replica.Role)
		}

		if replica.Commit != nil {
			v.Commit = string(*replica.Commit)
		}

		if replica.Failover != nil {
			v.Failover = string(*replica.Failover)
		}

		if replica.ReadableSecondary != nil {
			v.ReadableSecondary = string(*replica.ReadableSecondary)
		}

		results = append(results, v)
	}
	return results, nil
}
