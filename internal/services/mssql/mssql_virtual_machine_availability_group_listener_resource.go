package mssql

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	lbParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	lbValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	sqlValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	networkParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlVirtualMachineAvailabilityGroupListenerResource struct{}

type MsSqlVirtualMachineAvailabilityGroupListenerModel struct {
	Name                     string `tfschema:"name"`
	SqlVirtualMachineGroupId string `tfschema:"sql_virtual_machine_group_id"`
	AvailabilityGroupName    string `tfschema:"availability_group_name"`

	Port                       int                                                                      `tfschema:"port"`
	LoadBalancerConfiguration  []LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener  `tfschema:"load_balancer_configuration"`
	MultiSubnetIpConfiguration []MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener `tfschema:"multi_subnet_ip_configuration"`
	Replica                    []ReplicaMsSqlVirtualMachineAvailabilityGroupListener                    `tfschema:"replica"`
}

type LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener struct {
	LoadBalancerId       string   `tfschema:"load_balancer_id"`
	PrivateIpAddress     string   `tfschema:"private_ip_address"`
	ProbePort            int      `tfschema:"probe_port"`
	SqlVirtualMachineIds []string `tfschema:"sql_virtual_machine_ids"`
	SubnetId             string   `tfschema:"subnet_id"`
}

type MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener struct {
	PrivateIpAddress    string `tfschema:"private_ip_address"`
	SqlVirtualMachineId string `tfschema:"sql_virtual_machine_id"`
	SubnetId            string `tfschema:"subnet_id"`
}

type ReplicaMsSqlVirtualMachineAvailabilityGroupListener struct {
	SqlVirtualMachineId string `tfschema:"sql_virtual_machine_id"`
	Role                string `tfschema:"role"`
	Commit              string `tfschema:"commit"`
	FailoverMode        string `tfschema:"failover_mode"`
	ReadableSecondary   string `tfschema:"readable_secondary"`
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

		"load_balancer_configuration": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ExactlyOneOf: []string{"load_balancer_configuration", "multi_subnet_ip_configuration"},
			ForceNew:     true,
			MaxItems:     1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"load_balancer_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: lbValidate.LoadBalancerID,
					},

					"private_ip_address": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsIPAddress,
					},

					"probe_port": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.PortNumber,
					},

					"sql_virtual_machine_ids": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: sqlvirtualmachines.ValidateSqlVirtualMachineID,
						},
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: networkValidate.SubnetID,
					},
				},
			},
		},

		"multi_subnet_ip_configuration": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			ExactlyOneOf: []string{"load_balancer_configuration", "multi_subnet_ip_configuration"},
			ForceNew:     true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"private_ip_address": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsIPAddress,
					},

					"sql_virtual_machine_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: sqlvirtualmachines.ValidateSqlVirtualMachineID,
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: networkValidate.SubnetID,
					},
				},
			},
		},

		"replica": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"sql_virtual_machine_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: sqlValidate.SqlVirtualMachineID,
					},

					"role": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice([]string{string(availabilitygrouplisteners.RolePrimary), string(availabilitygrouplisteners.RoleSecondary)}, false),
					},

					"commit": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice([]string{string(availabilitygrouplisteners.CommitSynchronousCommit), string(availabilitygrouplisteners.CommitAsynchronousCommit)}, false),
					},

					"failover_mode": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice([]string{string(availabilitygrouplisteners.FailoverManual), string(availabilitygrouplisteners.FailoverAutomatic)}, false),
					},

					"readable_secondary": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice([]string{string(availabilitygrouplisteners.ReadableSecondaryNo), string(availabilitygrouplisteners.ReadableSecondaryReadOnly), string(availabilitygrouplisteners.ReadableSecondaryAll)}, false),
					},
				},
			},
			Set: ReplicaSchemaMsSqlVirtualMachineAvailabilityGroupListenerHash,
		},
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
				parameters.Properties.MultiSubnetIPConfigurations = expandMsSqlVirtualMachineAvailabilityGroupListenerMultiSubnetIpConfiguration(model.MultiSubnetIpConfiguration)
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

func expandMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(lbConfigs []LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener) (*[]availabilitygrouplisteners.LoadBalancerConfiguration, error) {
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

		lbConfig.PrivateIPAddress = &availabilitygrouplisteners.PrivateIPAddress{
			IPAddress:        pointer.To(lb.PrivateIpAddress),
			SubnetResourceId: pointer.To(lb.SubnetId),
		}

		results = append(results, lbConfig)
	}
	return &results, nil
}

func expandMsSqlVirtualMachineAvailabilityGroupListenerMultiSubnetIpConfiguration(multiSubnetIpConfiguration []MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener) *[]availabilitygrouplisteners.MultiSubnetIPConfiguration {
	results := make([]availabilitygrouplisteners.MultiSubnetIPConfiguration, 0)

	for _, item := range multiSubnetIpConfiguration {

		config := availabilitygrouplisteners.MultiSubnetIPConfiguration{
			SqlVirtualMachineInstance: item.SqlVirtualMachineId,
		}

		config.PrivateIPAddress = availabilitygrouplisteners.PrivateIPAddress{
			IPAddress:        pointer.To(item.PrivateIpAddress),
			SubnetResourceId: pointer.To(item.SubnetId),
		}

		results = append(results, config)
	}

	return &results
}

func flattenMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(input *[]availabilitygrouplisteners.LoadBalancerConfiguration, subscriptionId string) ([]LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener, error) {
	results := make([]LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, lbConfig := range *input {
		privateIpAddress := ""
		subnetResourceId := ""
		if v := lbConfig.PrivateIPAddress; v != nil {
			privateIpAddress = pointer.From(v.IPAddress)

			parsedSubnetResourceId, err := networkParse.SubnetIDInsensitively(pointer.From(v.SubnetResourceId))
			if err != nil {
				return nil, err
			}
			subnetResourceId = parsedSubnetResourceId.ID()
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
				if err != nil {
					return nil, err
				}

				// get correct casing for subscription in id due to https://github.com/Azure/azure-rest-api-specs/issues/25211
				newId := sqlvirtualmachines.NewSqlVirtualMachineID(subscriptionId, parsedId.ResourceGroupName, parsedId.SqlVirtualMachineName)

				parsedIds = append(parsedIds, newId.ID())
			}
			sqlVirtualMachineIds = parsedIds
		}

		v := LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener{
			LoadBalancerId:       loadBalancerId,
			PrivateIpAddress:     privateIpAddress,
			ProbePort:            int(pointer.From(lbConfig.ProbePort)),
			SqlVirtualMachineIds: sqlVirtualMachineIds,
			SubnetId:             subnetResourceId,
		}

		results = append(results, v)
	}
	return results, nil
}

func flattenMsSqlVirtualMachineAvailabilityGroupListenerMultiSubnetIpConfiguration(input *[]availabilitygrouplisteners.MultiSubnetIPConfiguration, subscriptionId string) ([]MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener, error) {
	results := make([]MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, config := range *input {
		parsedSubnetResourceId, err := networkParse.SubnetIDInsensitively(pointer.From(config.PrivateIPAddress.SubnetResourceId))
		if err != nil {
			return nil, err
		}

		parsedSqlVirtualMachineId, err := sqlvirtualmachines.ParseSqlVirtualMachineIDInsensitively(config.SqlVirtualMachineInstance)
		if err != nil {
			return nil, err
		}

		// get correct casing for subscription in id due to https://github.com/Azure/azure-rest-api-specs/issues/25211
		newSqlVirtualMachineId := sqlvirtualmachines.NewSqlVirtualMachineID(subscriptionId, parsedSqlVirtualMachineId.ResourceGroupName, parsedSqlVirtualMachineId.SqlVirtualMachineName)

		v := MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener{
			PrivateIpAddress:    pointer.From(config.PrivateIPAddress.IPAddress),
			SqlVirtualMachineId: newSqlVirtualMachineId.ID(),
			SubnetId:            parsedSubnetResourceId.ID(),
		}

		results = append(results, v)
	}
	return results, nil
}

func expandMsSqlVirtualMachineAvailabilityGroupListenerReplicas(replicas []ReplicaMsSqlVirtualMachineAvailabilityGroupListener) (*[]availabilitygrouplisteners.AgReplica, error) {
	results := make([]availabilitygrouplisteners.AgReplica, 0)

	for _, rep := range replicas {
		replica := availabilitygrouplisteners.AgReplica{
			Role:              pointer.To(availabilitygrouplisteners.Role(rep.Role)),
			Commit:            pointer.To(availabilitygrouplisteners.Commit(rep.Commit)),
			Failover:          pointer.To(availabilitygrouplisteners.Failover(rep.FailoverMode)),
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

func flattenMsSqlVirtualMachineAvailabilityGroupListenerReplicas(input *[]availabilitygrouplisteners.AgReplica, subscriptionId string) ([]ReplicaMsSqlVirtualMachineAvailabilityGroupListener, error) {
	results := make([]ReplicaMsSqlVirtualMachineAvailabilityGroupListener, 0)
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

			// get correct casing for subscription in id due to https://github.com/Azure/azure-rest-api-specs/issues/25211
			newId := sqlvirtualmachines.NewSqlVirtualMachineID(subscriptionId, parsedId.ResourceGroupName, parsedId.SqlVirtualMachineName)

			sqlVirtualMachineInstanceId = newId.ID()
		}

		v := ReplicaMsSqlVirtualMachineAvailabilityGroupListener{
			SqlVirtualMachineId: sqlVirtualMachineInstanceId,
			Role:                string(pointer.From(replica.Role)),
			Commit:              string(pointer.From(replica.Commit)),
			FailoverMode:        string(pointer.From(replica.Failover)),
			ReadableSecondary:   string(pointer.From(replica.ReadableSecondary)),
		}

		results = append(results, v)
	}
	return results, nil
}

func ReplicaSchemaMsSqlVirtualMachineAvailabilityGroupListenerHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["sql_virtual_machine_id"].(string))))
		buf.WriteString(fmt.Sprintf("%s-", m["role"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["commit"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["failover_mode"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["readable_secondary"].(string)))
	}

	return pluginsdk.HashString(buf.String())
}
