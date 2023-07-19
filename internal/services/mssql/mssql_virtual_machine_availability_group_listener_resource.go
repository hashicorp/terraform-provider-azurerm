package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
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
	Name                     string `tfschema:"name"`
	SqlVirtualMachineGroupId string `tfschema:"sql_virtual_machine_group_id"`
	AvailabilityGroupName    string `tfschema:"availability_group_name"`

	Port                       int                                                                             `tfschema:"port"`
	LoadBalancerConfiguration  []helper.LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener  `tfschema:"load_balancer_configuration"`
	MultiSubnetIpConfiguration []helper.MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener `tfschema:"multi_subnet_ip_configuration"`
	Replica                    []helper.ReplicaMsSqlVirtualMachineAvailabilityGroupListener                    `tfschema:"replica"`
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
			ValidateFunc: validation.StringLenBetween(1, 15),
		},

		"sql_virtual_machine_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: sqlvirtualmachines.ValidateSqlVirtualMachineGroupID,
		},

		"availability_group_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"port": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.PortNumber,
		},

		"load_balancer_configuration": helper.LoadBalancerConfigurationSchemaMsSqlVirtualMachineAvailabilityGroupListener(),

		"multi_subnet_ip_configuration": helper.MultiSubnetIpConfigurationSchemaMsSqlVirtualMachineAvailabilityGroupListener(),

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

			sqlVirtualMachineGroupId, err := availabilitygrouplisteners.ParseSqlVirtualMachineGroupID(model.SqlVirtualMachineGroupId)
			if err != nil {
				return err
			}

			id := availabilitygrouplisteners.NewAvailabilityGroupListenerID(subscriptionId, sqlVirtualMachineGroupId.ResourceGroupName, sqlVirtualMachineGroupId.SqlVirtualMachineGroupName, model.Name)

			existing, err := client.Get(ctx, id, availabilitygrouplisteners.GetOperationOptions{Expand: pointer.To("AvailabilityGroupConfiguration")})
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			replicas, err := expandMsSqlVirtualMachineAvailabilityGroupListenerReplicas(model.Replica)
			if err != nil {
				return err
			}

			parameters := availabilitygrouplisteners.AvailabilityGroupListener{
				Properties: &availabilitygrouplisteners.AvailabilityGroupListenerProperties{
					AvailabilityGroupName:                    pointer.To(model.AvailabilityGroupName),
					CreateDefaultAvailabilityGroupIfNotExist: pointer.To(true),
					Port:                                     pointer.To(int64(model.Port)),
					AvailabilityGroupConfiguration: &availabilitygrouplisteners.AgConfiguration{
						Replicas: replicas,
					},
				},
			}

			if model.LoadBalancerConfiguration != nil && len(model.LoadBalancerConfiguration) != 0 {
				lbConfigs, err := expandMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(model.LoadBalancerConfiguration)
				if err != nil {
					return err
				}
				parameters.Properties.LoadBalancerConfigurations = lbConfigs
			}

			if model.MultiSubnetIpConfiguration != nil && len(model.MultiSubnetIpConfiguration) != 0 {
				multiSubnetIpConfiguration, err := expandMsSqlVirtualMachineAvailabilityGroupListenerMultiSubnetIpConfiguration(model.MultiSubnetIpConfiguration)
				if err != nil {
					return err
				}
				parameters.Properties.MultiSubnetIPConfigurations = multiSubnetIpConfiguration
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

			resp, err := client.Get(ctx, *id, availabilitygrouplisteners.GetOperationOptions{Expand: pointer.To("AvailabilityGroupConfiguration")})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := MsSqlVirtualMachineAvailabilityGroupListenerModel{
				Name:                     id.AvailabilityGroupListenerName,
				SqlVirtualMachineGroupId: availabilitygrouplisteners.NewSqlVirtualMachineGroupID(id.SubscriptionId, id.ResourceGroupName, id.SqlVirtualMachineGroupName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {

					state.AvailabilityGroupName = pointer.From(props.AvailabilityGroupName)
					state.Port = int(pointer.From(props.Port))

					avGroupListenerLbConfigs, err := flattenMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(props.LoadBalancerConfigurations, id.SubscriptionId)
					if err != nil {
						return fmt.Errorf("setting `load_balancer_configuration`: %+v", err)
					}
					state.LoadBalancerConfiguration = avGroupListenerLbConfigs

					multiSubnetIpConfiguration, err := flattenMsSqlVirtualMachineAvailabilityGroupListenerMultiSubnetIpConfiguration(props.MultiSubnetIPConfigurations, id.SubscriptionId)
					if err != nil {
						return fmt.Errorf("setting `multi_subnet_ip_configuration`: %+v", err)
					}
					state.MultiSubnetIpConfiguration = multiSubnetIpConfiguration

					if props.AvailabilityGroupConfiguration != nil {
						if props.AvailabilityGroupConfiguration.Replicas != nil {

							replicas, err := flattenMsSqlVirtualMachineAvailabilityGroupListenerReplicas(props.AvailabilityGroupConfiguration.Replicas, id.SubscriptionId)
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
			ProbePort: pointer.To(int64(lb.ProbePort)),
		}

		parsedLbId := ""
		if lb.LoadBalancerId != "" {
			id, err := lbParse.LoadBalancerID(lb.LoadBalancerId)
			if err != nil {
				return nil, err
			}
			parsedLbId = id.ID()
		}
		lbConfig.LoadBalancerResourceId = pointer.To(parsedLbId)

		var parsedIds []interface{}
		for _, sqlVmId := range lb.SqlVirtualMachineIds {
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
		results = append(results, lbConfig)
	}
	return &results, nil
}

func expandMsSqlVirtualMachineAvailabilityGroupListenerMultiSubnetIpConfiguration(multiSubnetIpConfiguration []helper.MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener) (*[]availabilitygrouplisteners.MultiSubnetIPConfiguration, error) {
	results := make([]availabilitygrouplisteners.MultiSubnetIPConfiguration, 0)

	for _, item := range multiSubnetIpConfiguration {

		config := availabilitygrouplisteners.MultiSubnetIPConfiguration{
			SqlVirtualMachineInstance: item.SqlVirtualMachineId,
		}

		if item.PrivateIpAddress != nil {
			privateIpAddress, err := expandMsSqlVirtualMachinePrivateIpAddress(item.PrivateIpAddress)
			if err != nil {
				return nil, err
			}
			config.PrivateIPAddress = *privateIpAddress
		}
		results = append(results, config)
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
		IPAddress: pointer.To(ipAddress),

		SubnetResourceId: pointer.To(subnetId),
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

		loadBalancerId := ""
		if lbConfig.LoadBalancerResourceId != nil {
			id, err := lbParse.LoadBalancerID(pointer.From(lbConfig.LoadBalancerResourceId))
			if err != nil {
				return nil, err
			}
			loadBalancerId = id.ID()
		}

		var sqlVirtualMachineIds []string
		if lbConfig.SqlVirtualMachineInstances != nil {
			sqlVirtualMachineIds = *lbConfig.SqlVirtualMachineInstances
			var parsedIds []string
			for _, sqlVmId := range sqlVirtualMachineIds {
				parsedId, err := sqlvirtualmachines.ParseSqlVirtualMachineIDInsensitively(sqlVmId)
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
			PrivateIpAddress:     privateIpAddress,
			LoadBalancerId:       loadBalancerId,
			ProbePort:            int(pointer.From(lbConfig.ProbePort)),
			SqlVirtualMachineIds: sqlVirtualMachineIds,
		}

		results = append(results, v)
	}
	return results, nil
}

func flattenMsSqlVirtualMachineAvailabilityGroupListenerMultiSubnetIpConfiguration(input *[]availabilitygrouplisteners.MultiSubnetIPConfiguration, subscriptionId string) ([]helper.MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener, error) {
	results := make([]helper.MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, config := range *input {
		var privateIpAddress []helper.PrivateIpAddressMsSqlVirtualMachineAvailabilityGroupListener
		flattenedPrivateIp, err := flattenPrivateIpAddress(config.PrivateIPAddress)
		if err != nil {
			return nil, err
		}
		privateIpAddress = flattenedPrivateIp

		parsedId, err := sqlvirtualmachines.ParseSqlVirtualMachineIDInsensitively(config.SqlVirtualMachineInstance)
		newId := sqlvirtualmachines.NewSqlVirtualMachineID(subscriptionId, parsedId.ResourceGroupName, parsedId.SqlVirtualMachineName)
		if err != nil {
			return nil, err
		}

		v := helper.MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener{
			PrivateIpAddress:    privateIpAddress,
			SqlVirtualMachineId: newId.ID(),
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
		replica := availabilitygrouplisteners.AgReplica{
			Role:              pointer.To(availabilitygrouplisteners.Role(rep.Role)),
			Commit:            pointer.To(availabilitygrouplisteners.Commit(rep.Commit)),
			Failover:          pointer.To(availabilitygrouplisteners.Failover(rep.Failover)),
			ReadableSecondary: pointer.To(availabilitygrouplisteners.ReadableSecondary(rep.ReadableSecondary)),
		}

		sqlVirtualMachineId := rep.SqlVirtualMachineId
		if sqlVirtualMachineId != "" {
			id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(sqlVirtualMachineId)
			if err != nil {
				return nil, err
			}
			sqlVirtualMachineId = id.ID()
		}
		replica.SqlVirtualMachineInstanceId = pointer.To(sqlVirtualMachineId)

		results = append(results, replica)
	}
	return &results, nil
}

func flattenMsSqlVirtualMachineAvailabilityGroupListenerReplicas(input *[]availabilitygrouplisteners.AgReplica, subscriptionId string) ([]helper.ReplicaMsSqlVirtualMachineAvailabilityGroupListener, error) {
	results := make([]helper.ReplicaMsSqlVirtualMachineAvailabilityGroupListener, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, replica := range *input {

		sqlVirtualMachineInstanceId := ""
		if replica.SqlVirtualMachineInstanceId != nil {
			parsedId, err := sqlvirtualmachines.ParseSqlVirtualMachineIDInsensitively(*replica.SqlVirtualMachineInstanceId)
			if err != nil {
				return nil, err
			}

			newId := sqlvirtualmachines.NewSqlVirtualMachineID(subscriptionId, parsedId.ResourceGroupName, parsedId.SqlVirtualMachineName)
			if err != nil {
				return nil, err
			}

			sqlVirtualMachineInstanceId = newId.ID()
		}

		v := helper.ReplicaMsSqlVirtualMachineAvailabilityGroupListener{
			SqlVirtualMachineId: sqlVirtualMachineInstanceId,
		}

		if replica.Role != nil {
			v.Role = string(pointer.From(replica.Role))
		}

		if replica.Commit != nil {
			v.Commit = string(pointer.From(replica.Commit))
		}

		if replica.Failover != nil {
			v.Failover = string(pointer.From(replica.Failover))
		}

		if replica.ReadableSecondary != nil {
			v.ReadableSecondary = string(pointer.From(replica.ReadableSecondary))
		}

		results = append(results, v)
	}
	return results, nil
}
