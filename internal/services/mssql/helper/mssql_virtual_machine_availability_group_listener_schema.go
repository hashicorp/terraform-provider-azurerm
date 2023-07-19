package helper

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	lbValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	sqlValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener struct {
	PrivateIpAddress     []PrivateIpAddressMsSqlVirtualMachineAvailabilityGroupListener `tfschema:"private_ip_address"`
	LoadBalancerId       string                                                         `tfschema:"load_balancer_id"`
	ProbePort            int                                                            `tfschema:"probe_port"`
	SqlVirtualMachineIds []string                                                       `tfschema:"sql_virtual_machine_ids"`
}

func LoadBalancerConfigurationSchemaMsSqlVirtualMachineAvailabilityGroupListener() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeList,
		Optional:     true,
		ExactlyOneOf: []string{"load_balancer_configuration", "multi_subnet_ip_configuration"},
		ForceNew:     true,
		MaxItems:     1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"private_ip_address": PrivateIpAddressSchemaMsSqlVirtualMachineAvailabilityGroupListener(),

				"load_balancer_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: lbValidate.LoadBalancerID,
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
			},
		},
	}
}

type MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener struct {
	PrivateIpAddress    []PrivateIpAddressMsSqlVirtualMachineAvailabilityGroupListener `tfschema:"private_ip_address"`
	SqlVirtualMachineId string                                                         `tfschema:"sql_virtual_machine_id"`
}

func MultiSubnetIpConfigurationSchemaMsSqlVirtualMachineAvailabilityGroupListener() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeSet,
		Optional:     true,
		ExactlyOneOf: []string{"load_balancer_configuration", "multi_subnet_ip_configuration"},
		ForceNew:     true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"private_ip_address": PrivateIpAddressSchemaMsSqlVirtualMachineAvailabilityGroupListener(),

				"sql_virtual_machine_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: sqlvirtualmachines.ValidateSqlVirtualMachineID,
				},
			},
		},
	}
}

type PrivateIpAddressMsSqlVirtualMachineAvailabilityGroupListener struct {
	IpAddress string `tfschema:"ip_address"`
	SubnetId  string `tfschema:"subnet_id"`
}

func PrivateIpAddressSchemaMsSqlVirtualMachineAvailabilityGroupListener() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"ip_address": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.IsIPAddress,
				},

				"subnet_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: networkValidate.SubnetID,
				},
			},
		},
	}
}

type ReplicaMsSqlVirtualMachineAvailabilityGroupListener struct {
	SqlVirtualMachineId string `tfschema:"sql_virtual_machine_id"`
	Role                string `tfschema:"role"`
	Commit              string `tfschema:"commit"`
	Failover            string `tfschema:"failover"`
	ReadableSecondary   string `tfschema:"readable_secondary"`
}

func ReplicaSchemaMsSqlVirtualMachineAvailabilityGroupListener() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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

				"failover": {
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
	}
}

func ReplicaSchemaMsSqlVirtualMachineAvailabilityGroupListenerHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["sql_virtual_machine_id"].(string))))
		buf.WriteString(fmt.Sprintf("%s-", m["role"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["commit"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["failover"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["readable_secondary"].(string)))
	}

	return pluginsdk.HashString(buf.String())
}
